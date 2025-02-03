package JsonValidator

import (
	"fmt"
	"slices"
	"strings"
)

type RuleFunction func(*FieldValidationContext) (string, bool)

type ruleFunctionList map[string]RuleFunction

type Rulebook struct {
	rules      map[string]Rule
	composites map[string]string
}

func newRulebook(rules ruleFunctionList, nullableRules []string, presenceRules []string, aliases map[string]string) *Rulebook {
	rulebook := &Rulebook{
		rules:      make(map[string]Rule),
		composites: make(map[string]string),
	}

	for name, rule := range rules {
		rulebook.RegisterRule(Rule{
			Name:           name,
			Function:       rule,
			IsPresenceRule: slices.Contains(presenceRules, name),
			IsNullableRule: slices.Contains(nullableRules, name),
		})
	}

	for alias, name := range aliases {
		rulebook.RegisterAlias(alias, name)
	}

	return rulebook
}

func (rulebook Rulebook) RegisterRule(rule Rule) Rulebook {
	rulebook.rules[rule.Name] = rule

	return rulebook
}

func (rulebook Rulebook) RegisterAlias(alias string, name string) Rulebook {
	rule := rulebook.getRuleDefinition(name)

	rulebook.RegisterRule(Rule{
		Name:           alias,
		IsPresenceRule: rule.IsPresenceRule,
		IsNullableRule: rule.IsNullableRule,
		Function:       rule.Function,
	})

	return rulebook
}

func (rulebook Rulebook) RegisterComposite(name string, rules string) Rulebook {
	rulebook.composites[name] = rules

	return rulebook
}

func (rulebook Rulebook) IsComposite(ruleDefinition string) bool {
	name, _ := rulebook.parseRuleDefinition(ruleDefinition)
	_, ok := rulebook.composites[name]

	return ok
}

func (rulebook Rulebook) GetComposite(ruleDefinition string) string {
	name, params := rulebook.parseRuleDefinition(ruleDefinition)
	compositeRule, exists := rulebook.composites[name]
	if !exists {
		panic(fmt.Sprintf("No registered composite for name %s", name))
	}

	// Replace placeholders with provided params
	for i, param := range params {
		placeholder := fmt.Sprintf("$%d", i)
		compositeRule = strings.ReplaceAll(compositeRule, placeholder, param)
	}

	return compositeRule
}

func (rulebook Rulebook) GetRule(ruleDefinition string) *RuleContext {
	name, params := rulebook.parseRuleDefinition(ruleDefinition)
	definition := rulebook.getRuleDefinition(name)

	return &RuleContext{
		Rule:   definition,
		Params: params,
	}
}

func (rulebook Rulebook) getRuleDefinition(name string) Rule {
	if rule, ok := rulebook.rules[name]; ok {
		return rule
	}

	panic(fmt.Sprintf("No registered rule for name [%s]", name))
}

func (rulebook Rulebook) parseRuleDefinition(ruleDefinition string) (string, []string) {
	var params []string

	split := strings.Split(ruleDefinition, ":")
	name := split[0]

	if len(split) > 1 {
		params = strings.Split(split[1], ",")
	}

	return name, params
}
