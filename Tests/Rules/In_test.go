package Rules

import (
	"errors"
	"fmt"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func Test_it_can_validate_using_in_rule(t *testing.T) {
	// Setup
	cases := []struct {
		jsonString []byte
		shouldFail bool
	}{
		{[]byte(`{"Data": 123}`), false},
		{[]byte(`{"Data": 45.12}`), false},
		{[]byte(`{"Data": "abc"}`), false},
		{[]byte(`{"Data": false}`), false},
		{[]byte(`{"Data": "hello world"}`), true},
		{[]byte(`{"Data": "ab"}`), true},
		{[]byte(`{"Data": "bc"}`), true},
		{[]byte(`{"Data": "ABC"}`), true},
		{[]byte(`{"Data": "aBc"}`), true},
		{[]byte(`{"Data": "AbC"}`), true},
		{[]byte(`{"Data": 12}`), true},
		{[]byte(`{"Data": 23}`), true},
		{[]byte(`{"Data": 13}`), true},
		{[]byte(`{"Data": 123.45}`), true},
		{[]byte(`{"Data": {}}`), true},
		{[]byte(`{"Data": [1,2,3]}`), true},
		{[]byte(`{"Data": null}`), true},
		{[]byte(`{"Data": true}`), true},
	}

	type testData struct {
		Data any `validation:"in:abc,123,45.12,false"`
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
				require.True(t, errorBag.HasFailedKeyAndRule("Data", "in"))
				require.Equal(t, 1, errorBag.CountErrors())
			} else {
				require.True(t, errorBag == nil)
				require.Equal(t, 0, errorBag.CountErrors())
			}
		})
	}
}

func Test_it_prints_the_actual_value_for_the_in_rule(t *testing.T) {
	// Setup
	cases := []struct {
		jsonString  []byte
		expectedMsg string
	}{
		{[]byte(`{"Data": "hello world"}`), "[hello world] given"},
		{[]byte(`{"Data": 12}`), "[12] given"},
		{[]byte(`{"Data": 12.45}`), "[12.45] given"},
		{[]byte(`{"Data": {}}`), "Object given"},
		{[]byte(`{"Data": [1,2,3]}`), "Array given"},
		{[]byte(`{"Data": null}`), "[NULL] given"},
		{[]byte(`{"Data": true}`), "[true] given"},
	}

	type testData struct {
		Data any `validation:"in:abc,123"`
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
			require.True(t, errorBag != nil)
			require.True(t, errorBag.HasFailedKeyAndRule("Data", "in"))
			require.Equal(t, 1, errorBag.CountErrors())
			require.True(t, strings.Contains(errorBag.Errors["Data"][0], testCase.expectedMsg), fmt.Sprintf("Expected substring: %s\nGot: %s", testCase.expectedMsg, errorBag.Errors["Data"][0]))
		})
	}
}
