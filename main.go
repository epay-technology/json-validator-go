package main

import (
	"errors"
	"json-validator-go/JsonValidator"
	"log"
)

type Parent struct {
	Child *Child `json:"child"`
}

type Child struct {
	MyInts []int `json:"myInts" validation:""`
	MyBool bool  `json:"myBool" validation:"requiredWith:MyInts"`
}

func main() {
	jsonString := `{"child": {"myInts": [1,2,3]}}`
	data, err := JsonValidator.Validate[Parent]([]byte(jsonString))

	var validationErrors *JsonValidator.ErrorBag
	errors.As(err, &validationErrors)

	log.Fatal(data.Child.MyInts, validationErrors.IsValid(), validationErrors) // false true
}
