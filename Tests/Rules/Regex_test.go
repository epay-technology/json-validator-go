package Rules

import (
	"errors"
	"fmt"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_validate_using_regex_rule(t *testing.T) {
	// Setup
	cases := []struct {
		jsonString []byte
		shouldFail bool
	}{
		{[]byte(`{"Data": "123456"}`), false},
		{[]byte(`{"Data": "abcdef"}`), false},
		{[]byte(`{"Data": "1a2b3c"}`), false},
		{[]byte(`{"Data": "ABC321"}`), false},
		{[]byte(`{"Data": "a123-b456"}`), true},
		{[]byte(`{"Data": "dsada56dsadas-dsadas561dsada"}`), true},
		{[]byte(`{"Data": {}}`), true},
		{[]byte(`{"Data": "hello world"}`), true},
		{[]byte(`{"Data": 123}`), true},
		{[]byte(`{"Data": null}`), true},
		{[]byte(`{"Data": true}`), true},
		{[]byte(`{"Data": []}`), true},
	}

	type testData struct {
		Data any `validation:"regex:^[a-zA-Z0-9]{1,10}$"`
	}

	for i, testCase := range cases {
		t.Run(fmt.Sprintf("Test: #%d", i), func(t *testing.T) {
			// Arrange
			var errorBag *JsonValidator.ErrorBag

			// Act
			var data testData
			err := JsonValidator.New().Validate(testCase.jsonString, &data)
			_ = errors.As(err, &errorBag)

			// Assert
			if testCase.shouldFail {
				require.True(t, errorBag != nil)
				require.True(t, errorBag.HasFailedKeyAndRule("Data", "regex"))
				require.Equal(t, 1, errorBag.CountErrors())
			} else {
				require.True(t, errorBag == nil)
				require.Equal(t, 0, errorBag.CountErrors())
			}
		})
	}
}
