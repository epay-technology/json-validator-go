package JsonValidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

type ErrorBag struct {
	Errors map[string][]string
}

func (v *ErrorBag) Error() string {
	return fmt.Sprintf("Validation Errors: %#v", v.Errors)
}

func (v *ErrorBag) GetErrorsForKey(key string) []string {
	if v.Errors == nil {
		return []string{}
	}

	errs, ok := v.Errors[key]

	if !ok {
		return []string{}
	}

	return errs
}

// HasFailedKeyAndRule Is used for testing. So performance is not critical.
func (v *ErrorBag) HasFailedKeyAndRule(key string, rule string) bool {
	if v.Errors == nil {
		return false
	}

	errorList, failedKey := v.Errors[key]

	if !failedKey {
		return false
	}

	for _, errorText := range errorList {
		if strings.Index(errorText, fmt.Sprintf("[%s]: ", rule)) == 0 {
			return true
		}
	}

	return false
}

func (v *ErrorBag) AddError(path string, description string) {
	if v.Errors == nil {
		v.Errors = map[string][]string{}
	}

	targetBucket, ok := v.Errors[path]

	if !ok {
		targetBucket = []string{}
	}

	v.Errors[path] = append(targetBucket, description)
}

func (v *ErrorBag) IsValid() bool {
	return len(v.Errors) == 0
}

func (v *ErrorBag) IsInvalid() bool {
	return !v.IsValid()
}

type ValidationError struct {
	Path      string
	ErrorText string
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("Validation Errors: %s %s", v.Path, v.ErrorText)
}

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

type JsonContext struct {
	Path       string // The JSON path to the key under validation
	KeyPresent bool   // True if the key was present in the json
	EmptyValue bool   // True if the value is either null, 0, "", {}, [], or false
	IsNull     bool   // True if the value is null
	Value      any    // The raw json parsed value for the key. Will be nil if KeyPresent=false
}

func (context *JsonContext) HasEmptyValue() bool {
	// TODO: Implement
	return false
}

func Validate[T any](jsonData []byte) (*T, error) {
	var dataTarget T
	var jsonRaw map[string]any

	// This also verifies the integrity of the payload being valid json
	if err := json.Unmarshal(jsonData, &jsonRaw); err != nil {
		return nil, errors.New("invalid json cannot be parsed")
	}

	// Any errors at this point is not explicitly important.
	// It is only important if the validation does not also find any problems
	unmarshalErrors := json.Unmarshal(jsonData, &dataTarget)

	targetReflection := reflect.ValueOf(dataTarget)
	validation := &ErrorBag{Errors: map[string][]string{}}
	context := &ValidationContext{
		Json: &JsonContext{
			KeyPresent: true,
			Value:      jsonRaw,
		},
		Field: targetReflection,
	}
	context.RootContext = context
	context.ParentContext = context

	// Runs the actual validation against the json
	validateStruct(context, validation)

	// Validation errors has priority over any unmarshal errors
	// Since the json validation should also discover such errors by itself
	if validation.IsInvalid() {
		return nil, validation
	}

	// If there was no validation errors, but still unmarshal errors
	// Then our validation rules do not fully cover our API,
	// and we fall back to returning the unmarshal errors
	if unmarshalErrors != nil {
		return nil, unmarshalErrors
	}

	return &dataTarget, nil
}

func traverseField(context *ValidationContext, validation *ErrorBag) {
	switch reflect.Indirect(context.Field).Kind() {
	case reflect.Struct:
		validateStruct(context, validation)
	}
}

func validateStruct(context *ValidationContext, validation *ErrorBag) {
	field := reflect.Indirect(context.Field)
	numFields := field.NumField()
	dataType := field.Type()

	for i := 0; i < numFields; i++ {
		structField := dataType.Field(i)
		fieldContext := buildFieldContext(context, structField, field.Field(i))

		validateField(fieldContext, validation)
	}
}

type FieldValidationContext struct {
	Validation *ValidationContext
	Params     []string
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

func validateField(context *ValidationContext, validation *ErrorBag) {
	// We first execute any presence rules.
	// This is to handle null, and keys not existing separate from value/type assertions
	// If a presence error occurred, then no other validation rules should execute
	// This makes sure we do not run validation rules on fields that were not present or has a not allowed null value
	if presenceErrors := runRules(context, validation, context.ValidationTag.PresenceRules); presenceErrors {
		return
	}

	// If the value is null and the field allows null explicitly, then stop the validation here.
	// This is due to the fact that nullable/required are special case rules.
	// And avoids rules regarding type assertions to fail on nullable fields.
	if context.Json.IsNull && context.ValidationTag.ExplicitlyNullable {
		return
	}

	// Then, run all non-presence rules.
	errorsFound := runRules(context, validation, context.ValidationTag.Rules)

	if context.Json.KeyPresent && !context.Json.IsNull && !errorsFound {
		traverseField(context, validation)
	}
}

func runRules(context *ValidationContext, validation *ErrorBag, rules []string) bool {
	errorsFound := false

	for _, rule := range rules {
		var params []string
		split := strings.Split(rule, ":")
		ruleName := split[0]

		if len(split) > 1 {
			params = strings.Split(split[1], ",")
		}

		ruleFunction := getRuleByName(ruleName)

		if ruleFunction == nil {
			log.Fatal("Could not locate Rule: " + ruleName)
		}

		if errorText, success := (*ruleFunction)(&FieldValidationContext{Validation: context, Params: params}); !success {
			errorsFound = true
			validation.AddError(context.Json.Path, fmt.Sprintf("[%s]: %s", ruleName, errorText))
		}
	}

	return errorsFound
}

func buildFieldContext(parentContext *ValidationContext, fieldType reflect.StructField, fieldValue reflect.Value) *ValidationContext {
	jsonTag := getJsonTag(fieldType)

	return &ValidationContext{
		Json:            getJsonContext(parentContext, jsonTag),
		RootContext:     parentContext.RootContext,
		ParentContext:   parentContext,
		Field:           fieldValue,
		FieldName:       jsonTag.JsonKey,
		StructFieldName: fieldType.Name,
		ValidationTag:   getValidationTag(fieldType),
	}
}

func getJsonContext(parentContext *ValidationContext, jsonTag *JsonTag) *JsonContext {
	path := strings.TrimLeft(parentContext.Json.Path+"."+jsonTag.JsonKey, ".")

	if !parentContext.IsRoot() && !parentContext.Json.KeyPresent {
		return &JsonContext{
			Path:       path,
			KeyPresent: false,
			EmptyValue: true,
			IsNull:     true,
			Value:      nil,
		}
	}

	var jsonContext JsonContext
	jsonRaw, validStructJson := parentContext.Json.Value.(map[string]any)

	if !validStructJson {
		jsonContext = JsonContext{
			Path:       path,
			KeyPresent: false,
			EmptyValue: true,
			IsNull:     true,
			Value:      nil,
		}
	} else {
		jsonValue, present := jsonRaw[jsonTag.JsonKey]
		jsonContext = JsonContext{
			Path:       path,
			KeyPresent: present,
			EmptyValue: !present || isEmptyValue(jsonValue),
			IsNull:     present && jsonValue == nil,
			Value:      jsonValue,
		}
	}

	return &jsonContext
}

func isEmptyValue(value any) bool {
	switch value := reflect.ValueOf(value); value.Kind() {
	case reflect.Map:
		return len(value.MapKeys()) == 0
	case reflect.Pointer:
		return value.IsNil()
	case reflect.Slice:
		return value.Len() == 0
	case reflect.Bool:
		return value.Bool() == true
	case reflect.String:
		return value.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return value.Int() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	}

	return false
}

type JsonTag struct {
	JsonKey string
}

type ValidationTag struct {
	Rules              []string
	PresenceRules      []string
	ExplicitlyNullable bool
}

func (tag *ValidationTag) HasRules() bool {
	return 0 < (len(tag.Rules) + len(tag.PresenceRules))
}

func getJsonTag(field reflect.StructField) *JsonTag {
	tagline, ok := field.Tag.Lookup("json")

	if !ok {
		return &JsonTag{JsonKey: field.Name}
	}

	return &JsonTag{JsonKey: strings.Split(tagline, ",")[0]}
}

func getValidationTag(field reflect.StructField) *ValidationTag {
	tagline, ok := field.Tag.Lookup("validation")

	if !ok || strings.TrimSpace(tagline) == "" {
		return &ValidationTag{
			Rules:              []string{},
			ExplicitlyNullable: false,
			PresenceRules:      []string{},
		}
	}

	rules, presenceRules := extractPresenceRules(strings.Split(strings.TrimSpace(tagline), "|"))

	return &ValidationTag{
		Rules:              rules,
		ExplicitlyNullable: sliceContainsAny(rules, nullableRules...),
		PresenceRules:      presenceRules,
	}
}

// extractPresenceRules The first value in non-presence rules the second value is the presence rules
func extractPresenceRules(rules []string) (nonPresenceRulesList []string, presenceRulesList []string) {
	for _, rule := range rules {
		if slices.Contains(presenceRules, rule) {
			presenceRulesList = append(presenceRulesList, rule)
		} else {
			nonPresenceRulesList = append(nonPresenceRulesList, rule)
		}
	}

	return nonPresenceRulesList, presenceRulesList
}

func sliceContainsAny[S ~[]E, E comparable](slice S, values ...E) bool {
	for _, value := range slice {
		if slices.Contains(values, value) {
			return true
		}
	}

	return false
}
