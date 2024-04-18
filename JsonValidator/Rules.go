package JsonValidator

import (
	"fmt"
	"net"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var rules = map[string]RuleFunction{
	"nullable":           nullable,
	"required":           required,
	"present":            present,
	"len":                lenRule,
	"lenMin":             lenMin,
	"lenMax":             lenMax,
	"lenBetween":         lenBetween,
	"requiredWith":       requiredWith,
	"requiredWithout":    requiredWithout,
	"requiredWithAny":    requiredWithAny,
	"requiredWithoutAny": requiredWithoutAny,
	"requiredWithAll":    requiredWithAll,
	"requiredWithoutAll": requiredWithoutAll,
	"array":              isArray,
	"object":             isObject,
	"string":             isString,
	"int":                isInteger,
	"float":              isFloat,
	"bool":               isBool,
	"in":                 isIn,
	"uuid":               isUuid,
	"regex":              matchesRegex,
	"between":            isBetween,
	"min":                isMin,
	"max":                isMax,
	"url":                isUrl,
	"ip":                 isIp,
	"email":              isEmail,
}

var aliases = map[string]string{
	"minLen":  "lenMin",
	"maxLen":  "lenMax",
	"nilable": "nullable",
	"boolean": "bool",
	"integer": "int",
}

var nullableRules = []string{
	"nullable",
	"nilable",
}

var presenceRules = []string{
	"present",
	"required",
	"requiredWith",
	"requiredWithAny",
	"requiredWithAll",
	"requiredWithout",
	"requiredWithoutAny",
	"requiredWithoutAll",
}

func requiredWith(context *FieldValidationContext) (string, bool) {
	siblingNames := []string{context.Params[0]}
	neighborExists := isNeighborsPresent(context, siblingNames, true)

	if !neighborExists {
		return "", true
	}

	if _, requiredOk := required(context); requiredOk {
		return "", true
	}

	siblingJsonKeys := getNeighborJsonKeys(context, siblingNames)

	return fmt.Sprintf("Is required when %s is present", siblingJsonKeys[0]), false
}

func nullable(context *FieldValidationContext) (string, bool) {
	return "", true
}

func requiredWithAny(context *FieldValidationContext) (string, bool) {
	siblingNames := context.Params
	neighborExists := isNeighborsPresent(context, siblingNames, false)

	if !neighborExists {
		return "", true
	}

	if _, requiredOk := required(context); requiredOk {
		return "", true
	}

	siblingJsonKeys := getNeighborJsonKeys(context, siblingNames)

	return fmt.Sprintf("Is required when any of [%s] is present", strings.Join(siblingJsonKeys, ",")), false
}

func requiredWithAll(context *FieldValidationContext) (string, bool) {
	siblingNames := context.Params
	allNeighborsPresent := isNeighborsPresent(context, siblingNames, true)

	if !allNeighborsPresent {
		return "", true
	}

	if _, requiredOk := required(context); requiredOk {
		return "", true
	}

	siblingJsonKeys := getNeighborJsonKeys(context, siblingNames)

	return fmt.Sprintf("Is required when all of [%s] is present", strings.Join(siblingJsonKeys, ",")), false
}

func requiredWithout(context *FieldValidationContext) (string, bool) {
	siblingNames := []string{context.Params[0]}
	neighborExists := isNeighborsPresent(context, siblingNames, true)

	if neighborExists {
		return "", true
	}

	if _, requiredOk := required(context); requiredOk {
		return "", true
	}

	siblingJsonKeys := getNeighborJsonKeys(context, siblingNames)

	return fmt.Sprintf("Is required when %s is not present", siblingJsonKeys[0]), false
}

func requiredWithoutAny(context *FieldValidationContext) (string, bool) {
	siblingNames := context.Params
	allNeighborsArePresent := isNeighborsPresent(context, siblingNames, true)

	if allNeighborsArePresent {
		return "", true
	}

	if _, requiredOk := required(context); requiredOk {
		return "", true
	}

	siblingJsonKeys := getNeighborJsonKeys(context, siblingNames)

	return fmt.Sprintf("Is required when any of [%s] is not present", strings.Join(siblingJsonKeys, ",")), false
}

func requiredWithoutAll(context *FieldValidationContext) (string, bool) {
	siblingNames := context.Params
	anyNeighborIsPresent := isNeighborsPresent(context, siblingNames, false)

	if anyNeighborIsPresent {
		return "", true
	}

	if _, requiredOk := required(context); requiredOk {
		return "", true
	}

	siblingJsonKeys := getNeighborJsonKeys(context, siblingNames)

	return fmt.Sprintf("Is required when all of [%s] is not present", strings.Join(siblingJsonKeys, ",")), false
}

func isNeighborsPresent(context *FieldValidationContext, fields []string, all bool) bool {
	isAllPresent := true

	for _, field := range fields {
		neighbor, ok := context.Validation.GetNeighborField(field)

		if !ok {
			panic(fmt.Sprintf("No such field within struct: %s - Remember: Cross field references must use the struct name, and not the json name", field))
		}

		fieldPresent := ok && neighbor.Json.KeyPresent
		isAllPresent = isAllPresent && fieldPresent

		// When we want to know all fields are present but one is not
		if all && !fieldPresent {
			return false
		}

		// When we want to know all fields are present and the current field is
		if all && fieldPresent {
			continue
		}

		// When we just want to know that any field is present
		if !all && fieldPresent {
			return true
		}
	}

	return isAllPresent
}

func getNeighborJsonKeys(context *FieldValidationContext, fields []string) []string {
	var names []string

	for _, field := range fields {
		neighbor, ok := context.Validation.GetNeighborField(field)

		if !ok {
			names = append(names, field)
		} else {
			names = append(names, neighbor.FieldName)
		}
	}

	return names
}

func required(context *FieldValidationContext) (string, bool) {
	_, isPresent := present(context)

	if isPresent && !context.Validation.Json.IsNull {
		return "", true
	}

	return "Is a required non-nullable field", false
}

func present(context *FieldValidationContext) (string, bool) {
	if context.Validation.Json.KeyPresent {
		return "", true
	}

	return "Key must be present", false
}

func lenRule(context *FieldValidationContext) (string, bool) {
	expectedLen, err := strconv.Atoi(context.Params[0])

	if err != nil {
		return "validation failed", false
	}

	return lenHelper(context, &expectedLen, &expectedLen)
}

func lenMin(context *FieldValidationContext) (string, bool) {
	expectedLen, err := strconv.Atoi(context.Params[0])

	if err != nil {
		return "validation failed", false
	}

	return lenHelper(context, &expectedLen, nil)
}

func lenMax(context *FieldValidationContext) (string, bool) {
	expectedLen, err := strconv.Atoi(context.Params[0])

	if err != nil {
		return "validation failed", false
	}

	return lenHelper(context, nil, &expectedLen)
}

func lenBetween(context *FieldValidationContext) (string, bool) {
	expectedMinLen, errMin := strconv.Atoi(context.Params[0])
	expectedMaxLen, errMax := strconv.Atoi(context.Params[1])

	if errMin != nil || errMax != nil {
		return "validation failed", false
	}

	return lenHelper(context, &expectedMinLen, &expectedMaxLen)
}

func isArray(context *FieldValidationContext) (string, bool) {
	errorMessage := "Must be an array"

	if context.Validation.Json.IsNull {
		return errorMessage, false
	}

	valueKind := reflect.Indirect(reflect.ValueOf(context.Validation.Json.Value)).Kind()
	isValidKind := slices.Contains([]reflect.Kind{reflect.Slice, reflect.Array}, valueKind)

	return errorMessage, isValidKind
}

func isObject(context *FieldValidationContext) (string, bool) {
	errorMessage := "Must be an object"

	if context.Validation.Json.IsNull {
		return errorMessage, false
	}

	valueKind := reflect.Indirect(reflect.ValueOf(context.Validation.Json.Value)).Kind()
	isValidKind := slices.Contains([]reflect.Kind{reflect.Struct, reflect.Map}, valueKind)

	return errorMessage, isValidKind
}

func isString(context *FieldValidationContext) (string, bool) {
	errorMessage := "Must be a string"

	if context.Validation.Json.IsNull {
		return errorMessage, false
	}

	valueKind := reflect.Indirect(reflect.ValueOf(context.Validation.Json.Value)).Kind()
	isValidKind := slices.Contains([]reflect.Kind{reflect.String}, valueKind)

	return errorMessage, isValidKind
}

func isInteger(context *FieldValidationContext) (string, bool) {
	errorMessage := "Must be an integer"

	if context.Validation.Json.IsNull || !context.Validation.Json.KeyPresent {
		return errorMessage, false
	}

	valueReflection := reflect.ValueOf(context.Validation.Json.Value)
	valueKind := valueReflection.Kind()
	isValidKind := slices.Contains([]reflect.Kind{
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
	}, valueKind)

	if isValidKind {
		return "", true
	}

	floatVal, ok := valueReflection.Interface().(float64)

	if !ok {
		return errorMessage, false
	}

	return errorMessage, floatVal == float64(int(floatVal))
}

func isFloat(context *FieldValidationContext) (string, bool) {
	errorMessage := "Must be a float"

	if context.Validation.Json.IsNull {
		return errorMessage, false
	}

	valueKind := reflect.Indirect(reflect.ValueOf(context.Validation.Json.Value)).Kind()
	isValidKind := slices.Contains([]reflect.Kind{
		// Int is here considered a subset of the float value space
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
	}, valueKind)

	return errorMessage, isValidKind
}

func isIn(context *FieldValidationContext) (string, bool) {
	var actualValue string
	errorMessage := fmt.Sprintf("Value must be in set: [%s]", strings.Join(context.Params, ", "))

	if context.Validation.Json.IsNull || !context.Validation.Json.KeyPresent {
		return errorMessage, false
	}

	if value, isInt := context.Validation.Json.Value.(int); isInt {
		actualValue = strconv.Itoa(value)
	} else if value, isFloat := context.Validation.Json.Value.(float64); isFloat {
		actualValue = fmt.Sprintf("%d", int(value))
	} else if value, isString := context.Validation.Json.Value.(string); isString {
		actualValue = value
	} else {
		return errorMessage, false
	}

	return errorMessage, slices.Contains(context.Params, actualValue)
}

func isUuid(context *FieldValidationContext) (string, bool) {
	errorMessage := "Must be a valid uuid string"

	return errorMessage, verifyRegex(context, "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$")
}

func isUrl(context *FieldValidationContext) (string, bool) {
	errorMessage := "Must be a valid url string"

	return errorMessage, verifyRegex(context, "^(?:ftp|tcp|udp|wss?|https?):\\/\\/[\\w\\.\\/#=?&-_%]+$")
}

func isIp(context *FieldValidationContext) (string, bool) {
	errorMessage := "Must be a valid ip string"

	fieldValue, isString := context.Validation.Json.Value.(string)

	if !isString {
		return errorMessage, false
	}

	return errorMessage, fieldValue != "" && net.ParseIP(fieldValue) != nil
}

func isEmail(context *FieldValidationContext) (string, bool) {
	errorMessage := "Must be a valid email string"

	return errorMessage, verifyRegex(context, "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$")
}

func matchesRegex(context *FieldValidationContext) (string, bool) {
	regexString := context.Params[0]
	errorMessage := fmt.Sprintf("Must be a string matching regex: %s", regexString)

	return errorMessage, verifyRegex(context, regexString)
}

func isBetween(context *FieldValidationContext) (string, bool) {
	minValue := context.GetFloatParam(0)
	maxValue := context.GetFloatParam(1)

	errorMessage := fmt.Sprintf("Must be a number between %s and %s", context.GetParam(0), context.GetParam(1))

	value, isNumber := convertJsonValueToNumber(context)

	if !isNumber {
		return errorMessage, false
	}

	return errorMessage, minValue <= value && value <= maxValue
}

func isMin(context *FieldValidationContext) (string, bool) {
	minValue := context.GetFloatParam(0)
	errorMessage := fmt.Sprintf("Must be a number greater than or equal to %s", context.GetParam(0))

	value, isNumber := convertJsonValueToNumber(context)

	if !isNumber {
		return errorMessage, false
	}

	return errorMessage, minValue <= value
}

func isMax(context *FieldValidationContext) (string, bool) {
	maxValue := context.GetFloatParam(0)
	errorMessage := fmt.Sprintf("Must be a number less than or equal to %s", context.GetParam(0))

	value, isNumber := convertJsonValueToNumber(context)

	if !isNumber {
		return errorMessage, false
	}

	return errorMessage, value <= maxValue
}

func convertJsonValueToNumber(context *FieldValidationContext) (float64, bool) {
	reflection := reflect.ValueOf(context.Validation.Json.Value)

	isFloat := slices.Contains([]reflect.Kind{reflect.Float32, reflect.Float64}, reflection.Kind())

	if isFloat {
		return reflection.Float(), true
	}

	isInt := slices.Contains([]reflect.Kind{
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
	}, reflection.Kind())

	if isInt {
		return float64(reflection.Int()), true
	}

	return 0, false
}

func verifyRegex(context *FieldValidationContext, regexString string) bool {
	fieldValue, isString := context.Validation.Json.Value.(string)

	if !isString {
		return false
	}

	regex, err := regexp.Compile(regexString)

	if err != nil {
		return false
	}

	return regex.MatchString(fieldValue)
}

func isBool(context *FieldValidationContext) (string, bool) {
	errorMessage := "Must be a boolean"

	if context.Validation.Json.IsNull {
		return errorMessage, false
	}

	valueKind := reflect.Indirect(reflect.ValueOf(context.Validation.Json.Value)).Kind()
	isValidKind := slices.Contains([]reflect.Kind{reflect.Bool}, valueKind)

	return errorMessage, isValidKind
}

func lenHelper(context *FieldValidationContext, min *int, max *int) (string, bool) {
	var errorText string

	if min != nil && max != nil {
		if min == max {
			errorText = fmt.Sprintf("Lenght must be exactly %d", *min)
		} else {
			errorText = fmt.Sprintf("Lenght must be between %d and %d", *min, *max)
		}
	} else if min != nil {
		errorText = fmt.Sprintf("Lenght must be longer than %d", *min)
	} else if max != nil {
		errorText = fmt.Sprintf("Lenght must be longer than %d", *max)
	}

	valueType := reflect.ValueOf(context.Validation.Json.Value)

	if !slices.Contains([]reflect.Kind{reflect.Slice, reflect.Map, reflect.Array, reflect.String}, valueType.Kind()) {
		return errorText, false
	}

	actualLen := valueType.Len()

	if min != nil && *min > actualLen {
		return errorText, false
	}

	if max != nil && *max < actualLen {
		return errorText, false
	}

	return "", true
}
