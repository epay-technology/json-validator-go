package JsonValidator

type JsonTag struct {
	JsonKey string
}

type ValidationTag struct {
	Rules              []*rule
	PresenceRules      []*rule
	ExplicitlyNullable bool
}

func (tag *ValidationTag) HasRules() bool {
	return 0 < (len(tag.Rules) + len(tag.PresenceRules))
}
