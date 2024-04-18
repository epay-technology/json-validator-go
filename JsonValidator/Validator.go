package JsonValidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

type Validator struct {
	*Rulebook
}

func New() *Validator {
	return &Validator{Rulebook: newRulebook(rules, nullableRules, presenceRules, aliases)}
}

type JsonContext struct {
	Path       string // The JSON path to the key under validation
	KeyPresent bool   // True if the key was present in the json
	EmptyValue bool   // True if the value is either null, 0, "", {}, [], or false
	IsNull     bool   // True if the value is null
	Value      any    // The raw json parsed value for the key. Will be nil if KeyPresent=false
}

func (validator *Validator) Validate(jsonData []byte, dataTarget any) error {
	var jsonRaw map[string]any

	// This also verifies the integrity of the payload being valid json
	if err := json.Unmarshal(jsonData, &jsonRaw); err != nil {
		return errors.New("invalid json cannot be parsed")
	}

	// Any errors at this point is not explicitly important.
	// It is only important if the validation does not also find any problems
	unmarshalErrors := json.Unmarshal(jsonData, dataTarget)

	targetReflection := reflect.ValueOf(dataTarget)
	validation := &ErrorBag{Errors: map[string][]string{}}
	context := &ValidationContext{
		Json: &JsonContext{
			KeyPresent: true,
			Value:      jsonRaw,
		},
		Field:     targetReflection,
		Validator: validator,
	}
	context.RootContext = context
	context.ParentContext = context

	// Runs the actual validation against the json
	validator.validateStructSubFields(context, validation)

	// Validation errors has priority over any unmarshal errors
	// Since the json validation should also discover such errors by itself
	if validation.IsInvalid() {
		return validation
	}

	// If there was no validation errors, but still unmarshal errors
	// Then our validation rules do not fully cover our API,
	// and we fall back to returning the unmarshal errors
	if unmarshalErrors != nil {
		return unmarshalErrors
	}

	return nil
}

// traverseField is responsible for continuing the traversal from a specific field.
// It does NOT validate the specific field, but traverses any sub-fields or slice entries
// and call functions which then perform the actual validation.
func (validator *Validator) traverseField(context *ValidationContext, validation *ErrorBag) {
	switch reflect.Indirect(context.Field).Kind() {
	case reflect.Struct:
		validator.validateStructSubFields(context, validation)
	case reflect.Slice, reflect.Array:
		validator.validateSliceEntries(context, validation)
	default:
		// Do nothing
	}
}

func (validator *Validator) validateSliceEntries(context *ValidationContext, validation *ErrorBag) {
	jsonReflection := reflect.ValueOf(context.Json.Value)

	// If the json value is not an array, then we cannot continue the traversal.
	// Since there is no data left to validate.
	// Other validation rules before this point should ensure that the type was an array and give an appropriate error
	if !validator.isReflectionOfArray(jsonReflection) {
		return
	}

	jsonArrayLen := jsonReflection.Len()
	sliceLen := reflect.Indirect(context.Field).Len()
	canUseSliceValues := jsonArrayLen == sliceLen

	sliceElem := reflect.Indirect(context.Field)
	sliceSubtype := sliceElem.Type().Elem()
	entryValue := reflect.New(sliceSubtype)

	// TODO: Support diving
	for i := 0; i < jsonArrayLen; i++ {
		// Whenever possible, we prefer to point to the actual data structure.
		// But if the unmarshalling failed, then that is not possible.
		// We must still run validation, so we can detect the reason for the unmarshalling failing.
		if canUseSliceValues {
			entryValue = sliceElem.Index(i)
		}

		// TODO: Support diving, so we can validate the entry itself, and not just the entry sub fields/entries
		// This will validate the individual entries by ensuring any of its subfields has correct values.
		validator.traverseField(validator.buildSliceEntryContext(context, entryValue, i), validation)
	}
}

func (validator *Validator) isReflectionOfArray(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		return true
	default:
		return false
	}
}

func (validator *Validator) buildSliceEntryContext(parentContext *ValidationContext, entryValue reflect.Value, index int) *ValidationContext {
	jsonTag := validator.getJsonTagForSliceEntry(index)

	return &ValidationContext{
		Json:            validator.getJsonContext(parentContext, jsonTag),
		RootContext:     parentContext.RootContext,
		ParentContext:   parentContext,
		Field:           entryValue,
		FieldName:       jsonTag.JsonKey,
		StructFieldName: strconv.Itoa(index),
		// The Validation tag does not really matter, since we are never validation this exact context
		// We are only using it as a parent context when validating an entry within a slice.
		ValidationTag: &ValidationTag{
			Rules:              []*Rule{},
			ExplicitlyNullable: false,
			PresenceRules:      []*Rule{},
		},
		Validator: parentContext.Validator,
	}
}

func (validator *Validator) validateStructSubFields(context *ValidationContext, validation *ErrorBag) {
	field := reflect.Indirect(context.Field)
	numFields := field.NumField()
	dataType := field.Type()

	for i := 0; i < numFields; i++ {
		structField := dataType.Field(i)
		fieldContext := validator.buildFieldContext(context, structField, field.Field(i))

		validator.validateField(fieldContext, validation)
	}
}

func (validator *Validator) validateField(context *ValidationContext, validation *ErrorBag) {
	// We first execute any presence rules.
	// This is to handle null, and keys not existing separate from value/type assertions
	// If a presence error occurred, then no other validation rules should execute
	// This makes sure we do not run validation rules on fields that were not present or has a not allowed null value
	if presenceErrors := validator.runRules(context, validation, context.ValidationTag.PresenceRules); presenceErrors {
		return
	}

	// If the key is not present, then no additional validation should run.
	// Clients must define a presenceRule to require a field to be present or not
	if !context.Json.KeyPresent {
		return
	}

	// If the value is null and the field allows null explicitly, then stop the validation here.
	// This is due to the fact that nullable/required are special case rules.
	// And avoids rules regarding type assertions to fail on nullable fields.
	if context.Json.IsNull && context.ValidationTag.ExplicitlyNullable {
		return
	}

	// Then, run all non-presence rules.
	errorsFound := validator.runRules(context, validation, context.ValidationTag.Rules)

	if context.Json.KeyPresent && !context.Json.IsNull && !errorsFound {
		validator.traverseField(context, validation)
	}
}

func (validator *Validator) runRules(context *ValidationContext, validation *ErrorBag, rules []*Rule) bool {
	errorsFound := false

	for _, rule := range rules {
		if errorText, success := rule.Function(&FieldValidationContext{Validation: context, Params: rule.Params}); !success {
			errorsFound = true
			validation.AddError(context.Json.Path, fmt.Sprintf("[%s]: %s", rule.Name, errorText))
		}
	}

	return errorsFound
}

func (validator *Validator) buildFieldContext(parentContext *ValidationContext, fieldType reflect.StructField, fieldValue reflect.Value) *ValidationContext {
	jsonTag := validator.getJsonTagForStructField(fieldType)

	return &ValidationContext{
		Json:            validator.getJsonContext(parentContext, jsonTag),
		RootContext:     parentContext.RootContext,
		ParentContext:   parentContext,
		Field:           fieldValue,
		FieldName:       jsonTag.JsonKey,
		StructFieldName: fieldType.Name,
		ValidationTag:   validator.getValidationTag(fieldType),
		Validator:       parentContext.Validator,
	}
}

func (validator *Validator) getJsonContext(parentContext *ValidationContext, jsonTag *JsonTag) *JsonContext {
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

	// Handle array json values
	jsonRawArray, validArrayJson := parentContext.Json.Value.([]any)
	index, err := strconv.Atoi(jsonTag.JsonKey)

	if validArrayJson {
		if err == nil && len(jsonRawArray) > index {
			return validator.buildJsonContextForValue(path, true, jsonRawArray[index])
		}

		return validator.getEmptyJsonContext(path)
	}

	// Handle object json values
	jsonRawObject, validStructJson := parentContext.Json.Value.(map[string]any)

	if validStructJson {
		jsonValue, present := jsonRawObject[jsonTag.JsonKey]

		return validator.buildJsonContextForValue(path, present, jsonValue)
	}

	// Every other value type
	return validator.getEmptyJsonContext(path)
}

func (validator *Validator) buildJsonContextForValue(path string, present bool, jsonValue any) *JsonContext {
	return &JsonContext{
		Path:       path,
		KeyPresent: present,
		EmptyValue: !present || validator.isEmptyValue(jsonValue),
		IsNull:     present && jsonValue == nil,
		Value:      jsonValue,
	}
}

func (validator *Validator) getEmptyJsonContext(path string) *JsonContext {
	return validator.buildJsonContextForValue(path, false, nil)
}

func (validator *Validator) isEmptyValue(value any) bool {
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

func (validator *Validator) getJsonTagForStructField(field reflect.StructField) *JsonTag {
	tagline, ok := field.Tag.Lookup("json")

	if !ok {
		return &JsonTag{JsonKey: field.Name}
	}

	return &JsonTag{JsonKey: strings.Split(tagline, ",")[0]}
}

func (validator *Validator) getJsonTagForSliceEntry(index int) *JsonTag {
	return &JsonTag{JsonKey: strconv.Itoa(index)}
}

func (validator *Validator) getValidationTag(field reflect.StructField) *ValidationTag {
	tagline, ok := field.Tag.Lookup("validation")

	if !ok {
		return newValidationTag(validator.Rulebook, "")
	}

	return newValidationTag(validator.Rulebook, tagline)
}

// extractPresenceRules The first value in non-presence rules the second value is the presence rules
func (validator *Validator) extractPresenceRules(rules []string) (nonPresenceRulesList []*rule, presenceRulesList []*rule) {
	for _, ruleString := range rules {
		rule := validator.extractRule(ruleString)

		if slices.Contains(presenceRules, rule.name) {
			presenceRulesList = append(presenceRulesList, rule)
		} else {
			nonPresenceRulesList = append(nonPresenceRulesList, rule)
		}
	}

	return nonPresenceRulesList, presenceRulesList
}

func (validator *Validator) extractRule(ruleString string) *rule {
	var params []string
	split := strings.Split(ruleString, ":")
	name := split[0]

	if len(split) > 1 {
		params = strings.Split(split[1], ",")
	}

	return &rule{
		name:   name,
		params: params,
	}
}

func (validator *Validator) containsNullableRules(slice []*rule, values []string) bool {
	for _, value := range slice {
		if slices.Contains(values, value.name) {
			return true
		}
	}

	return false
}
