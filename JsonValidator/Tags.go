package JsonValidator

import "strings"

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

	ruleDefinitions := strings.Split(strings.TrimSpace(tagline), "|")

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
