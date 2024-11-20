package JsonValidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Validator struct {
	*Rulebook
	structCache *StructCache
}

func New() *Validator {
	return &Validator{
		Rulebook:    newRulebook(rules, nullableRules, presenceRules, aliases),
		structCache: newStructCache(),
	}
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

	fieldCache, err := validator.structCache.Analyze(validator.Rulebook, reflect.TypeOf(dataTarget))

	if err != nil {
		return err
	}

	validation := newErrorBag()
	context := &ValidationContext{
		Json: &JsonContext{
			KeyPresent: true,
			Value:      jsonRaw,
		},
		Field:     fieldCache,
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
	if err := json.Unmarshal(jsonData, dataTarget); err != nil {
		return err
	}

	return nil
}

func (validator *Validator) Analyze(dataTarget any) (*FieldCache, error) {
	return validator.structCache.Analyze(validator.Rulebook, reflect.TypeOf(dataTarget))
}

// traverseField is responsible for continuing the traversal from a specific field.
// It does NOT validate the specific field, but traverses any sub-fields or slice entries
// and call functions which then perform the actual validation.
func (validator *Validator) traverseField(context *ValidationContext, validation *ErrorBag) {
	if context.Field.IsStruct {
		validator.validateStructSubFields(context, validation)
	} else if context.Field.IsSlice {
		validator.validateSliceEntries(context, validation)
	} else if context.Field.IsMap {
		validator.validateMapEntries(context, validation)
	}

	// Do nothing
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
	sliceSubtype := context.Field.Children[0]

	// TODO: Support diving
	for i := 0; i < jsonArrayLen; i++ {
		// TODO: Support diving, so we can validate the entry itself, and not just the entry sub fields/entries
		// This will validate the individual entries by ensuring any of its subfields has correct values.
		validator.traverseField(validator.buildSliceEntryContext(context, sliceSubtype, i), validation)
	}
}

func (validator *Validator) validateMapEntries(context *ValidationContext, validation *ErrorBag) {
	jsonReflection := reflect.ValueOf(context.Json.Value)

	// If the json value is not a map, then we cannot continue the traversal.
	// Since there is no data left to validate.
	// Other validation rules before this point should ensure that the type was an array and give an appropriate error
	if jsonReflection.Kind() != reflect.Map {
		return
	}

	mapKeys := jsonReflection.MapKeys()
	sliceSubtype := context.Field.Children[0]

	// TODO: Support diving
	for _, key := range mapKeys {
		// TODO: Support diving, so we can validate the entry itself, and not just the entry sub fields/entries
		// This will validate the individual entries by ensuring any of its subfields has correct values.
		validator.traverseField(validator.buildMapEntryContext(context, sliceSubtype, key.String()), validation)
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

func (validator *Validator) buildSliceEntryContext(parentContext *ValidationContext, fieldCache *FieldCache, index int) *ValidationContext {
	stringKey := strconv.Itoa(index)

	return &ValidationContext{
		Json:            validator.getJsonContextForIntegerKey(parentContext, index),
		RootContext:     parentContext.RootContext,
		ParentContext:   parentContext,
		Field:           fieldCache,
		FieldName:       stringKey,
		StructFieldName: stringKey,
		// The Validation tag does not really matter, since we are never validation this exact context
		// We are only using it as a parent context when validating an entry within a slice.
		ValidationTag: &ValidationTag{
			Rules:              []*RuleContext{},
			ExplicitlyNullable: false,
			PresenceRules:      []*RuleContext{},
		},
		Validator: parentContext.Validator,
	}
}

func (validator *Validator) buildMapEntryContext(parentContext *ValidationContext, fieldCache *FieldCache, key string) *ValidationContext {
	return &ValidationContext{
		Json:            validator.getJsonContextForStringKey(parentContext, key),
		RootContext:     parentContext.RootContext,
		ParentContext:   parentContext,
		Field:           fieldCache,
		FieldName:       key,
		StructFieldName: key,
		// The Validation tag does not really matter, since we are never validation this exact context
		// We are only using it as a parent context when validating an entry within a slice.
		ValidationTag: &ValidationTag{
			Rules:              []*RuleContext{},
			ExplicitlyNullable: false,
			PresenceRules:      []*RuleContext{},
		},
		Validator: parentContext.Validator,
	}
}

func (validator *Validator) validateStructSubFields(context *ValidationContext, validation *ErrorBag) {
	for _, subField := range context.Field.Children {
		fieldContext := validator.buildFieldContext(context, subField)

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

func (validator *Validator) runRules(context *ValidationContext, validation *ErrorBag, rules []*RuleContext) bool {
	errorsFound := false

	for _, rule := range rules {
		if errorText, success := rule.Function(&FieldValidationContext{Validation: context, Params: rule.Params, RuleName: rule.Name}); !success {
			errorsFound = true
			validation.AddError(context.Json.Path, fmt.Sprintf("[%s]: %s", rule.Name, errorText))
		}
	}

	return errorsFound
}

func (validator *Validator) buildFieldContext(parentContext *ValidationContext, fieldCache *FieldCache) *ValidationContext {
	return &ValidationContext{
		Json:            validator.getJsonContextForStringKey(parentContext, fieldCache.JsonKey),
		RootContext:     parentContext.RootContext,
		ParentContext:   parentContext,
		Field:           fieldCache,
		FieldName:       fieldCache.JsonKey,
		StructFieldName: fieldCache.StructKey,
		ValidationTag:   fieldCache.ValidationTag,
		Validator:       parentContext.Validator,
	}
}

func (validator *Validator) getJsonContextForIntegerKey(parentContext *ValidationContext, key int) *JsonContext {
	path := fmt.Sprintf("%s.%d", parentContext.Json.Path, key)

	if !parentContext.IsRoot() && !parentContext.Json.KeyPresent {
		return validator.getEmptyJsonContext(path)
	}

	// Handle array json values
	jsonRawArray, validArrayJson := parentContext.Json.Value.([]any)

	if !validArrayJson {
		return validator.getEmptyJsonContext(path)
	}

	if len(jsonRawArray) > key {
		return validator.buildJsonContextForValue(path, true, jsonRawArray[key])
	}

	return validator.getEmptyJsonContext(path)
}

func (validator *Validator) getJsonContextForStringKey(parentContext *ValidationContext, key string) *JsonContext {
	path := strings.TrimLeft(parentContext.Json.Path+"."+key, ".")

	if !parentContext.IsRoot() && !parentContext.Json.KeyPresent {
		return validator.getEmptyJsonContext(path)
	}

	// Handle object json values
	jsonRawObject, validStructJson := parentContext.Json.Value.(map[string]any)

	if !validStructJson {
		return validator.getEmptyJsonContext(path)
	}

	jsonValue, present := jsonRawObject[key]

	return validator.buildJsonContextForValue(path, present, jsonValue)
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
