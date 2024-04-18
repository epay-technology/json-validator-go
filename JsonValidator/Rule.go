package JsonValidator

import (
	"strconv"
)

type Rule struct {
	Name           string
	IsPresenceRule bool
	IsNullableRule bool
	Function       RuleFunction
}

type RuleContext struct {
	Rule
	Params []string
}

func (context *RuleContext) GetStringParam(index int) string {
	return context.Params[index]
}

func (context *RuleContext) GetIntParam(index int) int {
	value, err := strconv.Atoi(context.Params[index])

	if err != nil {
		panic(err)
	}

	return value
}

func (context *RuleContext) GetFloatParam(index int) float64 {
	value, err := strconv.ParseFloat(context.Params[index], 64)

	if err != nil {
		panic(err)
	}

	return value
}

type rule struct {
	name   string
	params []string
}
