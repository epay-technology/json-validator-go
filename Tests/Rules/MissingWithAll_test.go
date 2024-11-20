package Rules

import (
	"errors"
	"fmt"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_validate_using_missingWithAll_rule(t *testing.T) {
	// Setup
	cases := []struct {
		jsonString []byte
		shouldFail bool
	}{
		{[]byte(`{}`), false},
		{[]byte(`{"Sibling1": true}`), false},
		{[]byte(`{"Sibling2": true}`), false},
		{[]byte(`{"Data": true}`), false},
		{[]byte(`{"Sibling1": true, "Sibling2": true}`), false},
		{[]byte(`{"Sibling1": true, "Data": true}`), false},
		{[]byte(`{"Sibling1": true, "Data": null}`), false},
		{[]byte(`{"Sibling1": null, "Data": null}`), false},
		{[]byte(`{"Sibling2": true, "Data": true}`), false},
		{[]byte(`{"Sibling2": true, "Data": null}`), false},
		{[]byte(`{"Sibling2": null, "Data": null}`), false},
		{[]byte(`{"Sibling1": true, "Sibling2": true, "Data": true}`), true},
		{[]byte(`{"Sibling1": true, "Sibling2": true, "Data": null}`), true},
		{[]byte(`{"Sibling1": null, "Sibling2": null, "Data": null}`), true},
	}

	type testData struct {
		Data     any `validation:"missingWithAll:Sibling1,Sibling2"`
		Sibling1 any
		Sibling2 any
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
				require.True(t, errorBag.HasFailedKeyAndRule("Data", "missingWithAll"))
				require.Equal(t, 1, errorBag.CountErrors())
			} else {
				require.True(t, errorBag == nil)
				require.Equal(t, 0, errorBag.CountErrors())
			}
		})
	}
}
