package Rules

import (
	"errors"
	"fmt"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_validate_using_lenMin_rule(t *testing.T) {
	// Setup
	cases := []struct {
		jsonString []byte
		shouldFail bool
	}{
		// Valid values
		{[]byte(`{"Data": "bc"}`), false},
		{[]byte(`{"Data": "abcasdasdsadsa"}`), false},
		{[]byte(`{"Data": "13"}`), false},
		{[]byte(`{"Data": "8694541681"}`), false},
		{[]byte(`{"Data": "  "}`), false},
		{[]byte(`{"Data": "            "}`), false},
		{[]byte(`{"Data": [1,2]}`), false},
		{[]byte(`{"Data": [1,2,5,6,4,7,1,2]}`), false},
		{[]byte(`{"Data": ["a",3]}`), false},
		{[]byte(`{"Data": ["a",3,true,false,null,"adsa"]}`), false},
		{[]byte(`{"Data": {"hello": "world", "foo": "bar"}}`), false},
		{[]byte(`{"Data": {"hello": "world", "foo": "bar", "fizz": "buzz"}}`), false},

		// Invalid values
		{[]byte(`{"Data": "h"}`), true},
		{[]byte(`{"Data": 12}`), true},
		{[]byte(`{"Data": {}}`), true},
		{[]byte(`{"Data": {"foo": "bar"}}`), true},
		{[]byte(`{"Data": null}`), true},
		{[]byte(`{"Data": true}`), true},
	}

	type testData struct {
		Data any `validation:"lenMin:2"`
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
				require.True(t, errorBag.HasFailedKeyAndRule("Data", "lenMin"))
				require.Equal(t, 1, errorBag.CountErrors())
			} else {
				require.True(t, errorBag == nil)
				require.Equal(t, 0, errorBag.CountErrors())
			}
		})
	}
}
