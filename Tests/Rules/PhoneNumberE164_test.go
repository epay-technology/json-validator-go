package Rules

import (
	"errors"
	"fmt"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_validate_using_phone_number_e164_rule(t *testing.T) {
	// Setup
	cases := []struct {
		jsonString []byte
		shouldFail bool
	}{
		{[]byte(`{"Data": "+45 88888888"}`), false},
		{[]byte(`{"Data": "+1 2125554567"}`), false},
		{[]byte(`{"Data": "+354 212555456713"}`), false},
		{[]byte(`{"Data": "+1 8"}`), false}, // While this is technically valid e.164, it is not a valid phone number.
		{[]byte(`{"Data": "+1 1234564894981"}`), true},
		{[]byte(`{"Data": "+1 165645a"}`), true},
		{[]byte(`{"Data": "+1 165b645"}`), true},
		{[]byte(`{"Data": "+1 165-645"}`), true},
		{[]byte(`{"Data": "+1 165 645"}`), true},
		{[]byte(`{"Data": "niki@epay.technology"}`), true},
		{[]byte(`{"Data": ""}`), true},
		{[]byte(`{"Data": [1,2,3]}`), true},
		{[]byte(`{"Data": {}}`), true},
		{[]byte(`{"Data": "hello world"}`), true},
		{[]byte(`{"Data": 123}`), true},
		{[]byte(`{"Data": null}`), true},
		{[]byte(`{"Data": true}`), true},
	}

	type testData struct {
		Data any `validation:"phoneNumberE164"`
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
				require.True(t, errorBag.HasFailedKeyAndRule("Data", "phoneNumberE164"))
				require.Equal(t, 1, errorBag.CountErrors())
			} else {
				require.True(t, errorBag == nil)
				require.Equal(t, 0, errorBag.CountErrors())
			}
		})
	}
}
