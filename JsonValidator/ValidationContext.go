package JsonValidator

type ValidationContext struct {
	Json            *JsonContext
	RootContext     *ValidationContext
	ParentContext   *ValidationContext
	Field           *FieldCache
	FieldName       string
	StructFieldName string
	ValidationTag   *ValidationTag
	Validator       *Validator
}

func (context *ValidationContext) GetNeighborField(name string) (*ValidationContext, bool) {
	neighbor := context.Field.Parent.GetChildByName(name)

	if neighbor == nil {
		return nil, false
	}

	return context.Validator.buildFieldContext(context.ParentContext, neighbor), true
}

func (context *ValidationContext) IsRoot() bool {
	return context.Json.Path == ""
}
