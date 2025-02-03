package JsonValidator

import (
	"strings"
)

type JsonTag struct {
	JsonKey string
}

type ValidationTag struct {
	Rules              []*RuleContext
	PresenceRules      []*RuleContext
	ExplicitlyNullable bool
}

func newValidationTag(rulebook *Rulebook, tagline string) *ValidationTag {
	if tagline == "" {
		return &ValidationTag{
			Rules:              []*RuleContext{},
			PresenceRules:      []*RuleContext{},
			ExplicitlyNullable: false,
		}
	}

	var rules []*RuleContext
	var presenceRules []*RuleContext
	explicitNullable := false

	ruleParser := func(definition string) []string {
		return strings.Split(strings.TrimSpace(definition), "|")
	}

	ruleDefinitions := unwrapCompositeRules(rulebook, ruleParser(tagline), ruleParser)

	for _, ruleDefinition := range ruleDefinitions {
		rule := rulebook.GetRule(ruleDefinition)

		if rule.IsPresenceRule {
			presenceRules = append(presenceRules, rule)
		} else {
			rules = append(rules, rule)
		}

		if rule.IsNullableRule {
			explicitNullable = true
		}
	}

	return &ValidationTag{
		Rules:              rules,
		PresenceRules:      presenceRules,
		ExplicitlyNullable: explicitNullable,
	}
}

func unwrapCompositeRules(rulebook *Rulebook, ruleDefinitions []string, ruleParser func(tagLine string) []string) []string {
	rules := make([]string, 0, len(ruleDefinitions))

	// Unwrap all composite rules and replace them with their underlying rules.
	for _, ruleDefinition := range ruleDefinitions {
		if rulebook.IsComposite(ruleDefinition) {
			composite := rulebook.GetComposite(ruleDefinition)
			rules = append(rules, ruleParser(composite)...)
		} else {
			rules = append(rules, ruleDefinition)
		}
	}

	return rules
}

func (tag *ValidationTag) GetRules(name string) []*RuleContext {
	var rules []*RuleContext

	for _, ruleInstance := range tag.Rules {
		if ruleInstance.Name == name {
			rules = append(rules, ruleInstance)
		}
	}

	for _, ruleInstance := range tag.PresenceRules {
		if ruleInstance.Name == name {
			rules = append(rules, ruleInstance)
		}
	}

	return rules
}
