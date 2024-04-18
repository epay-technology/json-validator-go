package Benchmarks

import (
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"testing"
)

func BenchmarkSimpleValidStructs(b *testing.B) {
	type MySimpleStruct struct {
		ID     int    `json:"id" validation:"required|integer|min:0"`
		Text   string `json:"text" validation:"present|nullable|string|maxLen:20"`
		Enable bool   `json:"enable" validation:"boolean"`
	}

	var targetStruct MySimpleStruct
	JsonInput := []byte(`{"id": 1532, "text": null}`)
	validator := JsonValidator.New()

	for i := 0; i < b.N; i++ {
		if err := validator.Validate(JsonInput, &targetStruct); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkSimpleInvalidStructs(b *testing.B) {
	type MySimpleStruct struct {
		ID     int    `json:"id" validation:"required|integer|min:0"`
		Text   string `json:"text" validation:"present|nullable|string|max:20"`
		Enable bool   `json:"enable" validation:"boolean"`
	}

	var targetStruct MySimpleStruct
	JsonInput := []byte(`{"id": "Not an integer", "text": 1234, "enable": "no"}`)
	validator := JsonValidator.New()

	for i := 0; i < b.N; i++ {
		if err := validator.Validate(JsonInput, &targetStruct); err == nil {
			b.Error("Validator expected errors, but none found")
		}
	}
}
