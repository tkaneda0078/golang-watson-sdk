package speechtotext

import (
	"testing"
)

func TestNewRequestHasRightInformation(t *testing.T) {
	r := NewRequest("bob", "1234")
	if r.username != "bob" {
		t.Fatalf("Expected username 'bob' got %v", r.username)
	}
	if r.password != "1234" {
		t.Fatalf("Expected password '1234' got %v", r.password)
	}
	if r.URL != API_URL {
		t.Fatalf("Expected URL %v got %v", API_URL, r.URL)
	}
	if r.EncodingModel != DEFAULT_ENCODING {
		t.Fatalf("Expected encoding %v got %v", DEFAULT_ENCODING, r.EncodingModel)
	}

}

func TestGoodResponseConverts(t *testing.T) {
	resp := []byte(`{
   "results": [
      {
         "alternatives": [
            {
               "confidence": 0.891,
               "transcript": "several tornadoes touch down as a line of severe thunderstorms swept through Colorado on Sunday "
            }
         ],
         "final": true
      }
   ],
   "result_index": 0
}`)
	c, err := toSpeechToTextStruct(resp)
	if err != nil {
		t.Fatalf("Unexpected error converting struct: %v", err)
	}
	if c.ResultIndex != 0 {
		t.Fatalf("Expected 'result_index' 0 but got %v", c.ResultIndex)
	}
	if len(c.Results) != 1 {
		t.Fatalf("Expected 1 result got %v", c.Results)
	}
	result := c.Results[0]
	if !result.Final {
		t.Fatal("Expected result final:true got false")
	}
	if len(result.Alternatives) != 1 {
		t.Fatalf("Expected 1 alternative got %v", len(result.Alternatives))
	}
	if result.Alternatives[0].Confidence != 0.891 {
		t.Fatalf("Expected confidence 0.891 got %v", result.Alternatives[0].Confidence)
	}
	if result.Alternatives[0].Transcript != "several tornadoes touch down as a line of severe thunderstorms swept through Colorado on Sunday " {
		t.Fatalf("Unexpected transcript: %v", result.Alternatives[0].Transcript)
	}
}

func TestErrorResponseConverts(t *testing.T) {
	resp := []byte(`{
  "error": "Model en-US_Broadband not found",
  "code": 404,
  "code_description": "No Such Resource"
}`)
	c, err := toSpeechToTextStruct(resp)
	if err != nil {
		t.Fatalf("Unexpected error converting struct: %v", err)
	}
	if c.Error != "Model en-US_Broadband not found" {
		t.Fatalf("Unexpected error, expected \"%v\" got %v", "Model en-US_Broadband not found", c.Error)
	}
	if c.ErrCode != 404 {
		t.Fatalf("Expected error 404 got %v", c.ErrCode)
	}
}

func TestJSONToModels(t *testing.T) {
	resp := []byte(`{
  "models": [
    {
      "name": "fr-FR_BroadbandModel",
      "language": "fr-FR",
      "url": "https://stream.watsonplatform.net/speech-to-text/api/v1/models/fr-FR_BroadbandModel",
      "rate": 16000,
      "supported_features": {
        "custom_language_model": false,
        "speaker_labels": false
      },
      "description": "French broadband model."
    },
    {
      "name": "en-US_NarrowbandModel",
      "language": "en-US",
      "url": "https://stream.watsonplatform.net/speech-to-text/api/v1/models/en-US_NarrowbandModel",
      "rate": 8000,
      "supported_features": {
        "custom_language_model": true,
        "speaker_labels": true
      },
      "description": "US English narrowband model."
    }
  ]
}`)
	models, err := toModels(resp)
	if err != nil {
		t.Fatal(err)
	}

	if len(models) != 2 {
		t.Fatalf("Expected 2 models, got %v", len(models))
	}
	if models[0].Name != "fr-FR_BroadbandModel" {
		t.Fatalf("Expected first model named fr-FR_BroadbandModel got %v", models[0].Name)
	}
	if models[1].Name != "en-US_NarrowbandModel" {
		t.Fatalf("Expected first model named en-US_NarrowbandModel got %v", models[1].Name)
	}
}
