package Rules

import (
	"errors"
	"fmt"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_validate_using_missingWith_rule(t *testing.T) {
	// Setup
	cases := []struct {
		jsonString []byte
		shouldFail bool
	}{
		{[]byte(`{"Sibling": null}`), false},
		{[]byte(`{"Sibling": true}`), false},
		{[]byte(`{"Sibling": null, "Data": 0}`), true},
		{[]byte(`{"Sibling": null, "Data": null}`), true},
		{[]byte(`{"Sibling": false, "Data": null}`), true},
		{[]byte(`{"Sibling": null, "Data": false}`), true},
		{[]byte(`{"Sibling": [], "Data": false}`), true},
		{[]byte(`{"Sibling": "hello", "Data": false}`), true},
		{[]byte(`{"Sibling": {}, "Data": false}`), true},
		{[]byte(`{"Sibling": null, "Data": []}`), true},
		{[]byte(`{"Sibling": null, "Data": {}}`), true},
	}

	type testData struct {
		Data    any `validation:"missingWith:Sibling"`
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
				require.True(t, errorBag.HasFailedKeyAndRule("Data", "missingWith"))
				require.Equal(t, 1, errorBag.CountErrors())
			} else {
				require.True(t, errorBag == nil)
				require.Equal(t, 0, errorBag.CountErrors())
			}
		})
	}
}
