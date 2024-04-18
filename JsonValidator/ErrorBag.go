package JsonValidator

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ErrorBag struct {
	Errors map[string][]string
}

func newErrorBag() *ErrorBag {
	return &ErrorBag{Errors: map[string][]string{}}
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
