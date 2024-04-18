package JsonValidator

import "reflect"

type ValidationContext struct {
	Json            *JsonContext
	RootContext     *ValidationContext
	ParentContext   *ValidationContext
	Field           reflect.Value
	FieldName       string
	StructFieldName string
	ValidationTag   *ValidationTag
}

func (context *ValidationContext) GetNeighborField(name string) (*ValidationContext, bool) {
	parent := reflect.Indirect(context.ParentContext.Field)
	structField, ok := parent.Type().FieldByName(name)

	if !ok {
		return nil, false
	}

	return buildFieldContext(
		context.ParentContext,
		structField,
		parent.FieldByName(name),
	), true
}

func (context *ValidationContext) IsRoot() bool {
	return context.Json.Path == ""
}
