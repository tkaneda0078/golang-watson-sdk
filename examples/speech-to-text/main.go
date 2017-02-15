package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bluemixgaragelondon/golang-sdk/speech-to-text"
)

var DEFAULT_ENCODING = "en-US_BroadbandModel"
var N_ALTERNATIVES = 1
var AUDIO_TYPE = "flac"
var TIMEOUT = 60

func main() {
	pPassword := flag.String("pass", "", "Password for Watson speech-to-text API")
	pFilename := flag.String("filename", "", "File to convert to text")
	pUsername := flag.String("user", "", "Username for API")
	pEncodingType := flag.String("model", DEFAULT_ENCODING, "Model to use for input file")
	pAudioType := flag.String("audio-type", AUDIO_TYPE, "Type of input file")
	pAlternatives := flag.Int("alternatives", N_ALTERNATIVES, "Number of alternative texts to provide")
	pWholeAudio := flag.Bool("entire-sample", true, "Attempt to get transcription for entire sample")
	pWatsonOptOut := flag.Bool("no-watson-learning", false, "Set to opt out of teaching Watson with this input")
	pListModels := flag.Bool("list-models", false, "List the models instead of converting audio")
	pTimeout := flag.Int("timeout", TIMEOUT, "Time (seconds) to wait for speech-to-text to return")

	flag.Parse()

	if (!*pListModels && *pFilename == "") || *pPassword == "" || *pUsername == "" {
		flag.PrintDefaults()
		log.Fatalf("Need to supply an audio file and a username/password")
	}

	req := speechtotext.NewRequest(*pUsername, *pPassword)

	if *pListModels {
		models, err := req.ListModels()
		if err != nil {
			log.Fatal(err)
		}
		for _, m := range models {
			fmt.Printf("%v\n", m.Name)
		}
		return
	}

	// Override some defaults (Alternatives : 1, UseWholeSample: true, EncodingModel: speechtotext.DEFAULT_ENCODING)
	req.Alternatives = *pAlternatives
	req.UseWholeSample = *pWholeAudio
	req.EncodingModel = *pEncodingType
	req.WatsonOptOut = *pWatsonOptOut
	req.Timeout = time.Duration(*pTimeout) * time.Second

	file, err := os.Open(*pFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	s, err := req.ToText(file, *pAudioType)
	if err != nil {
		log.Fatal(err)
	}

	if s.Error != "" {
		log.Fatal(s.Error, s.ErrCode)
	}

	for _, r := range s.Results {
		for _, a := range r.Alternatives {
			log.Printf("Found line: \"%v\" (%.3v%% confident)", a.Transcript, 100*a.Confidence)
		}
	}
}
