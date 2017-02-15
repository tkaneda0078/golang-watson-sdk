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
	URL            string
	username       string
	password       string
	EncodingModel  string
	UseWholeSample bool
	WatsonOptOut   bool
	Timeout        time.Duration
}

func NewRequest(username, password string) *requestAudioProperties {
	return &requestAudioProperties{
		URL:            API_URL,
		Alternatives:   1,
		EncodingModel:  DEFAULT_ENCODING,
		UseWholeSample: true,
		username:       username,
		password:       password,
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

type SupportedFeatures struct {
	CustomLanguageModel bool `json:"custom_language_model"`
	SpeakerLabels       bool `json:"speaker_labels"`
}

type Model struct {
	Name              string             `json:"name"`
	Language          string             `json:"language"`
	URL               string             `json:"url"`
	Rate              int                `json:"rate"`
	SupportedFeatures *SupportedFeatures `json:"supported_features"`
	Description       string             `json:"description"`
}
type listModelsResponse struct {
	Models []*Model `json:"models"`
}

func (r *requestAudioProperties) ListModels() ([]*Model, error) {
	url := fmt.Sprintf("%s/v1/models", API_URL)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(r.username, r.password)

	c := &http.Client{}
	res, err := c.Do(request)
	if err != nil {
		return nil, err
	}

	all, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return toModels(all)
}

func toModels(in []byte) ([]*Model, error) {
	var l listModelsResponse
	err := json.Unmarshal(in, &l)
	return l.Models, err
}

// Convert sets up the request to the speech-to-text service and returns an object with
// the results.
func (r *requestAudioProperties) ToText(reader io.Reader, audioFormat string) (*SpeechToText, error) {
	if reader == nil {
		return nil, fmt.Errorf("No reader supplied")
	}
	if audioFormat == "" {
		return nil, fmt.Errorf("No audio format supplied, expected e.g. 'flac'")
	}

	url := fmt.Sprintf(AUDIO_RECOGNISE_SIGNATURE, r.URL, r.UseWholeSample, r.EncodingModel, r.Alternatives)

	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", fmt.Sprintf("audio/%v", audioFormat))
	if r.WatsonOptOut {
		request.Header.Add("X-Watson-Learning-Opt-Out", "true")
	}

	request.SetBasicAuth(r.username, r.password)

	// Safe to overwrite as if Timeout is 0, it defaults to 60 seconds.
	c := &http.Client{
		Timeout: r.Timeout,
	}
	res, err := c.Do(request)
	if err != nil {
		return nil, err
	}

	all, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return toSpeechToTextStruct(all)
}

func toSpeechToTextStruct(b []byte) (*SpeechToText, error) {
	var s SpeechToText
	err := json.Unmarshal(b, &s)
	return &s, err
}
