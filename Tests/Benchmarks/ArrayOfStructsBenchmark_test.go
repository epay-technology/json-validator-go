package Benchmarks

import (
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"testing"
)

func BenchmarkNestedValidArraysOfStructs(b *testing.B) {
	type ChildStruct struct {
		ID     int    `json:"id" validation:"required|integer|min:0"`
		Text   string `json:"text" validation:"present|nullable|string|lenBetween:5,10"`
		Enable bool   `json:"enable" validation:"boolean"`
	}

	type MySimpleStruct struct {
		Children []ChildStruct `json:"children" validation:"required|lenBetween:2,10"`
	}

	var targetStruct MySimpleStruct
	JsonInput := []byte(`{"children": [
	{"id": 1, "text": "hello", "enable": true}, 
	{"id": 1, "text": "hello", "enable": true}, 
	{"id": 1, "text": "hello", "enable": true}, 
	{"id": 1, "text": "hello", "enable": true}, 
	{"id": 1, "text": "hello", "enable": true}, 
	{"id": 1, "text": "hello", "enable": true}
]}`)
	validator := JsonValidator.New()

	for i := 0; i < b.N; i++ {
		if err := validator.Validate(JsonInput, &targetStruct); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkNestedInvalidArraysOfStructs(b *testing.B) {
	type ChildStruct struct {
		ID     int    `json:"id" validation:"required|integer|min:0"`
		Text   string `json:"text" validation:"present|nullable|string|lenBetween:5,10"`
		Enable bool   `json:"enable" validation:"boolean"`
	}

	type MySimpleStruct struct {
		Children []ChildStruct `json:"children" validation:"required|lenBetween:2,10"`
	}

	var targetStruct MySimpleStruct
	JsonInput := []byte(`{"children": [
	{"id": -1, "text": "hello", "enable": true}, 
	{"id": 1, "text": false, "enable": true}, 
	{"id": 1, "text": "hello", "enable": null}, 
	{"text": "hello", "enable": true}, 
	{"id": 1, "enable": true}, 
	{"id": 1, "text": "hello", "enable": "world"}
]}`)
	validator := JsonValidator.New()

	for i := 0; i < b.N; i++ {
		err := validator.Validate(JsonInput, &targetStruct)

		if _, ok := err.(*JsonValidator.ErrorBag); err == nil || !ok {
			b.Error(err)
		}
	}
}
