package Rules

import (
	"errors"
	"fmt"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_validate_using_lenMax_rule(t *testing.T) {
	// Setup
	cases := []struct {
		jsonString []byte
		shouldFail bool
	}{
		// Valid values
		{[]byte(`{"Data": "bc"}`), false},
		{[]byte(`{"Data": "b"}`), false},
		{[]byte(`{"Data": "13"}`), false},
		{[]byte(`{"Data": "3"}`), false},
		{[]byte(`{"Data": ""}`), false},
		{[]byte(`{"Data": "  "}`), false},
		{[]byte(`{"Data": " "}`), false},
		{[]byte(`{"Data": ""}`), false},
		{[]byte(`{"Data": [1,2]}`), false},
		{[]byte(`{"Data": []}`), false},
		{[]byte(`{"Data": ["a",3]}`), false},
		{[]byte(`{"Data": ["a"]}`), false},
		{[]byte(`{"Data": {"hello": "world", "foo": "bar"}}`), false},
		{[]byte(`{"Data": {"foo": "bar"}}`), false},
		{[]byte(`{"Data": {}}`), false},

		// Invalid values
		{[]byte(`{"Data": 123}`), true},
		{[]byte(`{"Data": "hello world"}`), true},
		{[]byte(`{"Data": [1,2,3]}`), true},
		{[]byte(`{"Data": null}`), true},
		{[]byte(`{"Data": true}`), true},
	}

	type testData struct {
		Data any `validation:"lenMax:2"`
	}

	for i, testCase := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			// Arrange
			var errorBag *JsonValidator.ErrorBag

			// Act
			_, err := JsonValidator.Validate[testData](testCase.jsonString)
			_ = errors.As(err, &errorBag)

			// Assert
			if testCase.shouldFail {
				require.True(t, errorBag != nil)
				require.True(t, errorBag.HasFailedKeyAndRule("Data", "lenMax"))
				require.Equal(t, 1, errorBag.CountErrors())
			} else {
				require.True(t, errorBag == nil)
				require.Equal(t, 0, errorBag.CountErrors())
			}
		})
	}
}
