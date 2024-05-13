package JsonValidator

import "strconv"

type FieldValidationContext struct {
	Validation *ValidationContext
	Params     []string
	RuleName   string
}

func (context *FieldValidationContext) GetParam(index int) string {
	return context.Params[index]
}

func (context *FieldValidationContext) GetIntParam(index int) int {
	value, err := strconv.Atoi(context.Params[index])

	if err != nil {
		panic(err)
	}

	return value
}

func (context *FieldValidationContext) GetFloatParam(index int) float64 {
	value, err := strconv.ParseFloat(context.Params[index], 64)

	if err != nil {
		panic(err)
	}

	return value
}
