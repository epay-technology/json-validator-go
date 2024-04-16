package Rules

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"json-validator-go/JsonValidator"
	"testing"
)

func Test_it_can_validate_using_intBetween_rule(t *testing.T) {
	// Setup
	cases := []struct {
		jsonString []byte
		shouldFail bool
	}{
		{[]byte(`{"Data": 5}`), false},
		{[]byte(`{"Data": 8}`), false},
		{[]byte(`{"Data": 10}`), false},
		{[]byte(`{"Data": [1,2,3]}`), true},
		{[]byte(`{}`), true},
		{[]byte(`{"Data": {}}`), true},
		{[]byte(`{"Data": "hello world"}`), true},
		{[]byte(`{"Data": 123}`), true},
		{[]byte(`{"Data": null}`), true},
		{[]byte(`{"Data": true}`), true},
	}

	type testData struct {
		Data any `validation:"intBetween:5,10"`
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
				require.True(t, errorBag.HasFailedKeyAndRule("Data", "intBetween"))
			} else {
				require.True(t, errorBag == nil)
			}
		})
	}
}
