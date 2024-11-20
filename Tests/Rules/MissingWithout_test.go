package Rules

import (
	"errors"
	"fmt"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_validate_using_missingWithout_rule(t *testing.T) {
	// Setup
	cases := []struct {
		jsonString []byte
		shouldFail bool
	}{
		{[]byte(`{"Sibling": null}`), false},
		{[]byte(`{"Sibling": true}`), false},
		{[]byte(`{"Sibling": null, "Data": 0}`), false},
		{[]byte(`{"Sibling": null, "Data": null}`), false},
		{[]byte(`{"Sibling": false, "Data": null}`), false},
		{[]byte(`{"Sibling": null, "Data": false}`), false},
		{[]byte(`{"Sibling": [], "Data": false}`), false},
		{[]byte(`{"Sibling": "hello", "Data": false}`), false},
		{[]byte(`{"Sibling": {}, "Data": false}`), false},
		{[]byte(`{"Sibling": null, "Data": []}`), false},
		{[]byte(`{"Sibling": null, "Data": {}}`), false},
		{[]byte(`{"Data": {}}`), true},
		{[]byte(`{"Data": ""}`), true},
		{[]byte(`{"Data": 123}`), true},
		{[]byte(`{"Data": []}`), true},
	}

	type testData struct {
		Data    any `validation:"missingWithout:Sibling"`
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
				require.True(t, errorBag.HasFailedKeyAndRule("Data", "missingWithout"))
				require.Equal(t, 1, errorBag.CountErrors())
			} else {
				require.True(t, errorBag == nil)
				require.Equal(t, 0, errorBag.CountErrors())
			}
		})
	}
}
