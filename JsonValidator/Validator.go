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
	if v == nil || v.Errors == nil {
		return "No validation errors"
	}

	jsonBytes, _ := json.MarshalIndent(v.Errors, "", "  ")

	return fmt.Sprintf("Validation Errors: \n%s", string(jsonBytes))
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

func (v *ErrorBag) CountErrors() int {
	if v == nil || v.Errors == nil {
		return 0
	}

	count := 0

	for _, list := range v.Errors {
		count += len(list)
	}

	return count
}

// HasFailedKeyAndRule Is used for testing. So performance is not critical.
func (v *ErrorBag) HasFailedKeyAndRule(key string, rule string) bool {
	if v == nil {
		return false
	}

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
	validateStructSubFields(context, validation)

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

// traverseField is responsible for continuing the traversal from a specific field.
// It does NOT validate the specific field, but traverses any sub-fields or slice entries
// and call functions which then perform the actual validation.
func traverseField(context *ValidationContext, validation *ErrorBag) {
	switch reflect.Indirect(context.Field).Kind() {
	case reflect.Struct:
		validateStructSubFields(context, validation)
	case reflect.Slice, reflect.Array:
		validateSliceEntries(context, validation)
	}
}

func validateSliceEntries(context *ValidationContext, validation *ErrorBag) {
	jsonReflection := reflect.ValueOf(context.Json.Value)

	// If the json value is not an array, then we cannot continue the traversal.
	// Since there is no data left to validate.
	// Other validation rules before this point should ensure that the type was an array and give an appropriate error
	if !isReflectionOfArray(jsonReflection) {
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
		traverseField(buildSliceEntryContext(context, entryValue, i), validation)
	}
}

func isReflectionOfArray(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		return true
	default:
		return false
	}
}

func buildSliceEntryContext(parentContext *ValidationContext, entryValue reflect.Value, index int) *ValidationContext {
	jsonTag := getJsonTagForSliceEntry(index)

	return &ValidationContext{
		Json:            getJsonContext(parentContext, jsonTag),
		RootContext:     parentContext.RootContext,
		ParentContext:   parentContext,
		Field:           entryValue,
		FieldName:       jsonTag.JsonKey,
		StructFieldName: strconv.Itoa(index),
		// The Validation tag does not really matter, since we are never validation this exact context
		// We are only using it as a parent context when validating an entry within a slice.
		ValidationTag: &ValidationTag{
			Rules:              []*rule{},
			ExplicitlyNullable: false,
			PresenceRules:      []*rule{},
		},
	}
}

func validateStructSubFields(context *ValidationContext, validation *ErrorBag) {
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
	errorsFound := runRules(context, validation, context.ValidationTag.Rules)

	if context.Json.KeyPresent && !context.Json.IsNull && !errorsFound {
		traverseField(context, validation)
	}
}

func runRules(context *ValidationContext, validation *ErrorBag, rules []*rule) bool {
	errorsFound := false

	for _, rule := range rules {
		ruleFunction := getRuleByName(rule.name)

		if ruleFunction == nil {
			log.Fatal("Could not locate Rule: " + rule.name)
		}

		if errorText, success := (*ruleFunction)(&FieldValidationContext{Validation: context, Params: rule.params}); !success {
			errorsFound = true
			validation.AddError(context.Json.Path, fmt.Sprintf("[%s]: %s", rule.name, errorText))
		}
	}

	return errorsFound
}

func buildFieldContext(parentContext *ValidationContext, fieldType reflect.StructField, fieldValue reflect.Value) *ValidationContext {
	jsonTag := getJsonTagForStructField(fieldType)

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

	// Handle array json values
	jsonRawArray, validArrayJson := parentContext.Json.Value.([]any)
	index, err := strconv.Atoi(jsonTag.JsonKey)

	if validArrayJson {
		if err == nil && len(jsonRawArray) > index {
			return buildJsonContextForValue(path, true, jsonRawArray[index])
		}

		return getEmptyJsonContext(path)
	}

	// Handle object json values
	jsonRawObject, validStructJson := parentContext.Json.Value.(map[string]any)

	if validStructJson {
		jsonValue, present := jsonRawObject[jsonTag.JsonKey]

		return buildJsonContextForValue(path, present, jsonValue)
	}

	// Every other value type
	return getEmptyJsonContext(path)
}

func buildJsonContextForValue(path string, present bool, jsonValue any) *JsonContext {
	return &JsonContext{
		Path:       path,
		KeyPresent: present,
		EmptyValue: !present || isEmptyValue(jsonValue),
		IsNull:     present && jsonValue == nil,
		Value:      jsonValue,
	}
}

func getEmptyJsonContext(path string) *JsonContext {
	return buildJsonContextForValue(path, false, nil)
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
	Rules              []*rule
	PresenceRules      []*rule
	ExplicitlyNullable bool
}

func (tag *ValidationTag) HasRules() bool {
	return 0 < (len(tag.Rules) + len(tag.PresenceRules))
}

func getJsonTagForStructField(field reflect.StructField) *JsonTag {
	tagline, ok := field.Tag.Lookup("json")

	if !ok {
		return &JsonTag{JsonKey: field.Name}
	}

	return &JsonTag{JsonKey: strings.Split(tagline, ",")[0]}
}

func getJsonTagForSliceEntry(index int) *JsonTag {
	return &JsonTag{JsonKey: strconv.Itoa(index)}
}

type rule struct {
	name   string
	params []string
}

func getValidationTag(field reflect.StructField) *ValidationTag {
	tagline, ok := field.Tag.Lookup("validation")

	if !ok || strings.TrimSpace(tagline) == "" {
		return &ValidationTag{
			Rules:              []*rule{},
			ExplicitlyNullable: false,
			PresenceRules:      []*rule{},
		}
	}

	rules, presenceRules := extractPresenceRules(strings.Split(strings.TrimSpace(tagline), "|"))

	return &ValidationTag{
		Rules:              rules,
		PresenceRules:      presenceRules,
		ExplicitlyNullable: containsNullableRules(rules, nullableRules),
	}
}

// extractPresenceRules The first value in non-presence rules the second value is the presence rules
func extractPresenceRules(rules []string) (nonPresenceRulesList []*rule, presenceRulesList []*rule) {
	for _, ruleString := range rules {
		rule := extractRule(ruleString)

		if slices.Contains(presenceRules, rule.name) {
			presenceRulesList = append(presenceRulesList, rule)
		} else {
			nonPresenceRulesList = append(nonPresenceRulesList, rule)
		}
	}

	return nonPresenceRulesList, presenceRulesList
}

func extractRule(ruleString string) *rule {
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

func containsNullableRules(slice []*rule, values []string) bool {
	for _, value := range slice {
		if slices.Contains(values, value.name) {
			return true
		}
	}

	return false
}
