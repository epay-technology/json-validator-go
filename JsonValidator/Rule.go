package JsonValidator

import (
	"strconv"
)

type Rule struct {
	Function       RuleFunction
	Params         []string
	Name           string
	IsPresenceRule bool
	IsNullableRule bool
}

func (context *Rule) GetStringParam(index int) string {
	return context.Params[index]
}

func (context *Rule) GetIntParam(index int) int {
	value, err := strconv.Atoi(context.Params[index])

	if err != nil {
		panic(err)
	}

	return value
}

func (context *Rule) GetFloatParam(index int) float64 {
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
