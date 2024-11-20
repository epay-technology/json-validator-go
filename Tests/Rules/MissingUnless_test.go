package Rules

import (
	"errors"
	"fmt"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_validate_using_missingUnless_rule_with_number(t *testing.T) {
	// Setup
	cases := []struct {
		jsonString []byte
		shouldFail bool
	}{
		{[]byte(`{"Sibling": null}`), false},
		{[]byte(`{"Sibling": true}`), false},
		{[]byte(`{"Sibling": 123}`), false},
		{[]byte(`{"Sibling": "123"`), false},
		{[]byte(`{"Sibling": 123, "Data": true}`), false},
		{[]byte(`{"Sibling": "123", "Data": null}`), false},
		{[]byte(`{"Sibling": null, "Data": 0}`), true},
		{[]byte(`{"Sibling": null, "Data": null}`), true},
		{[]byte(`{"Sibling": 123.4, "Data": null}`), true},
		{[]byte(`{"Sibling": "", "Data": false}`), true},
		{[]byte(`{"Sibling": {}, "Data": []}`), true},
		{[]byte(`{"Sibling": [], "Data": {}}`), true},
	}

	type testData struct {
		Data    any `validation:"missingUnless:Sibling,123"`
		Sibling any
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
				require.Equal(t, 1, errorBag.CountErrors())
				require.True(t, errorBag.HasFailedKeyAndRule("Data", "missingUnless"))
			} else {
				require.True(t, errorBag == nil)
				require.Equal(t, 0, errorBag.CountErrors())
			}
		})
	}
}
