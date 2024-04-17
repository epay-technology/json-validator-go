package Rules

import (
	"errors"
	"fmt"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_validate_using_integer_rule(t *testing.T) {
	// Setup
	cases := []struct {
		jsonString []byte
		shouldFail bool
	}{
		{[]byte(`{"Data": 123}`), false},
		{[]byte(`{"Data": true}`), true},
		{[]byte(`{"Data": "hello world"}`), true},
		{[]byte(`{"Data": {}}`), true},
		{[]byte(`{}`), true},
		{[]byte(`{"Data": [1,2,3]}`), true},
		{[]byte(`{"Data": null}`), true},
	}

	type testData struct {
		Data any `validation:"int"`
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
				require.True(t, errorBag.HasFailedKeyAndRule("Data", "int"))
			} else {
				require.True(t, errorBag == nil, string(testCase.jsonString))
			}
		})
	}
}
