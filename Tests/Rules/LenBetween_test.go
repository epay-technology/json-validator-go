package Rules

import (
	"errors"
	"fmt"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_validate_using_lenBetween_rule(t *testing.T) {
	// Setup
	cases := []struct {
		jsonString []byte
		shouldFail bool
	}{
		// Valid values
		{[]byte(`{"Data": "bc"}`), false},
		{[]byte(`{"Data": "bc13"}`), false},
		{[]byte(`{"Data": "bac"}`), false},
		{[]byte(`{"Data": "13"}`), false},
		{[]byte(`{"Data": "1356"}`), false},
		{[]byte(`{"Data": "135"}`), false},
		{[]byte(`{"Data": "  "}`), false},
		{[]byte(`{"Data": "   "}`), false},
		{[]byte(`{"Data": [1,2]}`), false},
		{[]byte(`{"Data": [1,2,4]}`), false},
		{[]byte(`{"Data": [1,2,4,5]}`), false},
		{[]byte(`{"Data": ["a",3]}`), false},
		{[]byte(`{"Data": ["a",3, null]}`), false},
		{[]byte(`{"Data": {"hello": "world", "foo": "bar"}}`), false},
		{[]byte(`{"Data": {"hello": "world", "foo": "bar", "bar": "biz"}}`), false},

		// Invalid values
		{[]byte(`{"Data": "hello world"}`), true},
		{[]byte(`{"Data": "h"}`), true},
		{[]byte(`{"Data": {}}`), true},
		{[]byte(`{"Data": {"a": "a"}}`), true},
		{[]byte(`{"Data": {"a": "a", "b": "b", "c": "c", "d": "d", "e": "e"}}`), true},
		{[]byte(`{"Data": [1,2,3,4,5]}`), true},
		{[]byte(`{"Data": null}`), true},
		{[]byte(`{"Data": true}`), true},
	}

	type testData struct {
		Data any `validation:"lenBetween:2,4"`
	}

	for i, testCase := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			// Arrange
			var errorBag *JsonValidator.ErrorBag

			// Act
			var data testData
			err := JsonValidator.NewValidator().Validate(testCase.jsonString, &data)
			_ = errors.As(err, &errorBag)

			// Assert
			if testCase.shouldFail {
				require.True(t, errorBag != nil)
				require.True(t, errorBag.HasFailedKeyAndRule("Data", "lenBetween"))
				require.Equal(t, 1, errorBag.CountErrors())
			} else {
				require.True(t, errorBag == nil)
				require.Equal(t, 0, errorBag.CountErrors())
			}
		})
	}
}
