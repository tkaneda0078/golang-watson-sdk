package speechtotext

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var API_URL = "https://stream.watsonplatform.net/speech-to-text/api"
var AUDIO_RECOGNISE_SIGNATURE = "%s/v1/recognize?continuous=%v&model=%s&max_alternatives=%v"
var DEFAULT_ENCODING = "en-US_BroadbandModel"

type requestAudioProperties struct {
	Alternatives   int
	fullURL        string
	URL            string
	Username       string
	Password       string
	EncodingModel  string
	UseWholeSample bool
	WatsonOptOut   bool
}

func NewRequest(u, p string) *requestAudioProperties {
	return &requestAudioProperties{
		URL:            API_URL,
		Alternatives:   1,
		EncodingModel:  DEFAULT_ENCODING,
		UseWholeSample: true,
		Username:       u,
		Password:       p,
	}
}

type SpeechToText struct {
	Results     []*ResultText `json:"results"`
	ResultIndex int           `json:"result_index"`
	Error       string        `json:"error"`
	ErrCode     int           `json:"code"`
}

type ResultText struct {
	Alternatives []*Alternatives `json:"alternatives"`
	Final        bool            `json:"final"`
}

type Alternatives struct {
	Confidence float64 `json:"confidence"`
	Transcript string  `json:"transcript"`
}

// Convert sets up the request to the speech-to-text service and returns an object with
// the results.
func (r *requestAudioProperties) ToText(reader io.Reader, audioType string) (*SpeechToText, error) {
	if reader == nil {
		return nil, fmt.Errorf("No reader supplied")
	}

	url := fmt.Sprintf(AUDIO_RECOGNISE_SIGNATURE, r.URL, r.UseWholeSample, r.EncodingModel, r.Alternatives)

	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", fmt.Sprintf("audio/%v", audioType))
	if r.WatsonOptOut {
		request.Header.Add("X-Watson-Learning-Opt-Out", "true")
	}

	request.SetBasicAuth(r.Username, r.Password)

	c := &http.Client{
		Timeout: 5 * time.Minute,
	}
	res, err := c.Do(request)
	if err != nil {
		return nil, err
	}

	all, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return convertToStruct(all)
}

func convertToStruct(b []byte) (*SpeechToText, error) {
	var s SpeechToText
	err := json.Unmarshal(b, &s)
	return &s, err
}
