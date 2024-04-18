package main

import (
	"errors"
	"github.com/epay-technology/json-validator-go/JsonValidator"
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

	var data Parent
	err := JsonValidator.NewValidator().Validate([]byte(jsonString), &data)

	var validationErrors *JsonValidator.ErrorBag
	errors.As(err, &validationErrors)

	log.Fatal(data.Child.MyInts, validationErrors.IsValid(), validationErrors) // false true
}
