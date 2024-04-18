package Benchmarks

import (
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"testing"
)

func BenchmarkNestedValidStructs(b *testing.B) {
	type ChildStruct struct {
		ID     int    `json:"id" validation:"required|integer|min:0"`
		Text   string `json:"text" validation:"present|nullable|string|lenBetween:5,10"`
		Enable bool   `json:"enable" validation:"boolean"`
	}

	type MySimpleStruct struct {
		Child1 ChildStruct `json:"child1" validation:"required"`
		Child2 ChildStruct `json:"child2" validation:"required"`
		Child3 ChildStruct `json:"child3" validation:"required"`
	}

	var targetStruct MySimpleStruct
	JsonInput := []byte(`{"child1": {"id": 1, "text": "hello", "enable": true}, "child2": {"id": 1, "text": "hello", "enable": true}, "child3": {"id": 1, "text": "hello", "enable": true}}`)
	validator := JsonValidator.New()

	for i := 0; i < b.N; i++ {
		if err := validator.Validate(JsonInput, &targetStruct); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkNestedInvalidStructs(b *testing.B) {
	type ChildStruct struct {
		ID     int    `json:"id" validation:"required|integer|min:0"`
		Text   string `json:"text" validation:"present|nullable|string|lenBetween:5,10"`
		Enable bool   `json:"enable" validation:"boolean"`
	}

	type MySimpleStruct struct {
		Child1 ChildStruct `json:"child1" validation:"required"`
		Child2 ChildStruct `json:"child2" validation:"required"`
		Child3 ChildStruct `json:"child3" validation:"required"`
	}

	var targetStruct MySimpleStruct
	JsonInput := []byte(`{"child1": {"id": -10, "text": "hello", "enable": true}, "child2": {"id": 1, "text": true, "enable": true}, "child3": {"id": 1, "text": "hello", "enable": null}}`)
	validator := JsonValidator.New()

	for i := 0; i < b.N; i++ {
		if err := validator.Validate(JsonInput, &targetStruct); err == nil {
			b.Error("Validator expected error. None returned")
		}
	}
}
