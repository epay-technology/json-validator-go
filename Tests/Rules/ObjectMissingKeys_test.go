package Rules

import (
	"errors"
	"fmt"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_validate_using_object_missing_keys_rule_on_any(t *testing.T) {
	// Setup
	cases := []struct {
		jsonString []byte
		shouldFail bool
	}{
		{[]byte(`{"Data": 0}`), true},
		{[]byte(`{"Data": null}`), true},
		{[]byte(`{"Data": false}`), true},
		{[]byte(`{"Data": ""}`), true},
		{[]byte(`{"Data": []}`), true},
		{[]byte(`{"Data": [1, 2, 3]}`), true},
		{[]byte(`{"Data": {"key1": true}}`), true},
		{[]byte(`{"Data": {"key1": "hello"}}`), true},
		{[]byte(`{"Data": {"key1": 123}}`), true},
		{[]byte(`{"Data": {"key1": []}}`), true},
		{[]byte(`{"Data": {"key1": {}}}`), true},
		{[]byte(`{"Data": {"key2": true}}`), true},
		{[]byte(`{"Data": {"3": true}}`), true},
		{[]byte(`{"Data": {"key1": true, "key2": true}}`), true},
		{[]byte(`{"Data": {"key3": true}}`), false},
		{[]byte(`{"Data": {"key3": {"key1": true}}}`), false},
		{[]byte(`{"Data": {}}`), false},
		{[]byte(`{"key1": true}`), false},
	}

	type testData struct {
		Data any `validation:"objectMissingKeys:key1,key2,3"`
	}

	for i, testCase := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			// Arrange
			var errorBag *JsonValidator.ErrorBag

			// Act
			var data testData
			err := JsonValidator.New().Validate(testCase.jsonString, &data)
			_ = errors.As(err, &errorBag)

			// Assert
			if testCase.shouldFail {
				require.True(t, errorBag != nil)
				require.True(t, errorBag.HasFailedKeyAndRule("Data", "objectMissingKeys"))
				require.Equal(t, 1, errorBag.CountErrors())
			} else {
				require.True(t, errorBag == nil)
				require.Equal(t, 0, errorBag.CountErrors())
			}
		})
	}
}

func Test_it_can_validate_using_object_missing_keys_rule_on_struct(t *testing.T) {
	// Setup
	cases := []struct {
		jsonString []byte
		shouldFail bool
	}{
		{[]byte(`{"Data": 0}`), true},
		{[]byte(`{"Data": null}`), true},
		{[]byte(`{"Data": false}`), true},
		{[]byte(`{"Data": ""}`), true},
		{[]byte(`{"Data": []}`), true},
		{[]byte(`{"Data": [1, 2, 3]}`), true},
		{[]byte(`{"Data": {"key1": true}}`), true},
		{[]byte(`{"Data": {"key1": "hello"}}`), true},
		{[]byte(`{"Data": {"key1": 123}}`), true},
		{[]byte(`{"Data": {"key1": []}}`), true},
		{[]byte(`{"Data": {"key1": {}}}`), true},
		{[]byte(`{"Data": {"key2": true}}`), true},
		{[]byte(`{"Data": {"3": true}}`), true},
		{[]byte(`{"Data": {"key1": true, "key2": true}}`), true},
		{[]byte(`{"Data": {"key3": true}}`), false},
		{[]byte(`{"Data": {"key3": {"key1": true}}}`), false},
		{[]byte(`{"Data": {}}`), false},
		{[]byte(`{"key1": true}`), false},
	}

	type subType struct {
		Key3 any `json:"key3"`
	}

	type testData struct {
		Data subType `validation:"objectMissingKeys:key1,key2,3"`
	}

	for i, testCase := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			// Arrange
			var errorBag *JsonValidator.ErrorBag

			// Act
			var data testData
			err := JsonValidator.New().Validate(testCase.jsonString, &data)
			_ = errors.As(err, &errorBag)

			// Assert
			if testCase.shouldFail {
				require.True(t, errorBag != nil)
				require.True(t, errorBag.HasFailedKeyAndRule("Data", "objectMissingKeys"))
				require.Equal(t, 1, errorBag.CountErrors())
			} else {
				require.True(t, errorBag == nil)
				require.Equal(t, 0, errorBag.CountErrors())
			}
		})
	}
}

func Test_it_can_validate_using_object_missing_keys_rule_on_map(t *testing.T) {
	// Setup
	cases := []struct {
		jsonString []byte
		shouldFail bool
	}{
		{[]byte(`{"Data": 0}`), true},
		{[]byte(`{"Data": null}`), true},
		{[]byte(`{"Data": false}`), true},
		{[]byte(`{"Data": ""}`), true},
		{[]byte(`{"Data": []}`), true},
		{[]byte(`{"Data": [1, 2, 3]}`), true},
		{[]byte(`{"Data": {"key1": true}}`), true},
		{[]byte(`{"Data": {"key1": "hello"}}`), true},
		{[]byte(`{"Data": {"key1": 123}}`), true},
		{[]byte(`{"Data": {"key1": []}}`), true},
		{[]byte(`{"Data": {"key1": {}}}`), true},
		{[]byte(`{"Data": {"key2": true}}`), true},
		{[]byte(`{"Data": {"3": true}}`), true},
		{[]byte(`{"Data": {"key1": true, "key2": true}}`), true},
		{[]byte(`{"Data": {"key3": true}}`), false},
		{[]byte(`{"Data": {"key3": {"key1": true}}}`), false},
		{[]byte(`{"Data": {}}`), false},
		{[]byte(`{"key1": true}`), false},
	}

	type testData struct {
		Data map[string]any `validation:"objectMissingKeys:key1,key2,3"`
	}

	for i, testCase := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			// Arrange
			var errorBag *JsonValidator.ErrorBag

			// Act
			var data testData
			err := JsonValidator.New().Validate(testCase.jsonString, &data)
			_ = errors.As(err, &errorBag)

			// Assert
			if testCase.shouldFail {
				require.True(t, errorBag != nil)
				require.True(t, errorBag.HasFailedKeyAndRule("Data", "objectMissingKeys"))
				require.Equal(t, 1, errorBag.CountErrors())
			} else {
				require.True(t, errorBag == nil)
				require.Equal(t, 0, errorBag.CountErrors())
			}
		})
	}
}
