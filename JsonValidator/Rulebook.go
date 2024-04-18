package JsonValidator

import (
	"fmt"
	"slices"
	"strings"
)

type RuleFunction func(*FieldValidationContext) (string, bool)

type ruleFunctionList map[string]RuleFunction

type Rulebook struct {
	rules         ruleFunctionList
	nullableRules ruleFunctionList
	presenceRules ruleFunctionList
	aliases       map[string]string
}

func newRulebook(rules ruleFunctionList, nullableRules []string, presenceRules []string, aliases map[string]string) *Rulebook {
	rulebook := &Rulebook{
		rules:         ruleFunctionList{},
		nullableRules: ruleFunctionList{},
		presenceRules: ruleFunctionList{},
		aliases:       aliases,
	}

	for name, rule := range rules {
		rulebook.RegisterRule(
			name,
			slices.Contains(presenceRules, name),
			slices.Contains(nullableRules, name),
			rule,
		)
	}

	return rulebook
}

func (rulebook *Rulebook) RegisterRule(name string, isPresenceRule bool, isNullableRule bool, function RuleFunction) *Rulebook {
	if isPresenceRule {
		rulebook.presenceRules[name] = function
	}

	if isNullableRule {
		rulebook.nullableRules[name] = function
	}

	rulebook.rules[name] = function

	return rulebook
}

func (rulebook *Rulebook) RegisterAlias(alias string, name string) *Rulebook {
	rulebook.aliases[alias] = name

	return rulebook
}

func (rulebook *Rulebook) GetRule(ruleDefinition string) *Rule {
	name, params := rulebook.parseRuleDefinition(ruleDefinition)

	return &Rule{
		Function:       rulebook.getRuleFunction(name),
		Params:         params,
		Name:           name,
		IsNullableRule: inMap(rulebook.nullableRules, name),
		IsPresenceRule: inMap(rulebook.presenceRules, name),
	}
}

func (rulebook *Rulebook) getRuleFunction(name string) RuleFunction {
	if function, ok := rulebook.rules[name]; ok {
		return function
	}

	if aliases, ok := rulebook.aliases[name]; ok {
		return rulebook.getRuleFunction(aliases)
	}

	panic(fmt.Sprintf("No registered rule for name %s", name))
}

func (rulebook *Rulebook) parseRuleDefinition(ruleDefinition string) (string, []string) {
	var params []string

	split := strings.Split(ruleDefinition, ":")
	name := split[0]

	if len(split) > 1 {
		params = strings.Split(split[1], ",")
	}

	return name, params
}

func inMap[T comparable, V any](dictionary map[T]V, key T) bool {
	_, present := dictionary[key]
	return present
}
