package speechtotext

import (
	"testing"
)

func TestNewRequestHasRightInformation(t *testing.T) {
	r := NewRequest("bob", "1234")
	if r.Username != "bob" {
		t.Fatalf("Expected username 'bob' got %v", r.Username)
	}
	if r.Password != "1234" {
		t.Fatalf("Expected password '1234' got %v", r.Password)
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
	c, err := convertToStruct(resp)
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
	c, err := convertToStruct(resp)
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
