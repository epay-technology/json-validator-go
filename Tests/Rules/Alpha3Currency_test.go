package Rules

import (
	"errors"
	"fmt"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_validate_using_alpha3Currency_rule(t *testing.T) {
	// Setup
	cases := []struct {
		jsonString []byte
		shouldFail bool
	}{
		{[]byte(`{"Data": "DKK"}`), false},
		{[]byte(`{"Data": "SEK"}`), false},
		{[]byte(`{"Data": "NOK"}`), false},
		{[]byte(`{"Data": "DK"}`), true},
		{[]byte(`{"Data": "DKKK"}`), true},
		{[]byte(`{"Data": "SE"}`), true},
		{[]byte(`{"Data": ""}`), true},
		{[]byte(`{"Data": null}`), true},
		{[]byte(`{"Data": true}`), true},
		{[]byte(`{"Data": {}}`), true},
		{[]byte(`{"Data": []}`), true},
		{[]byte(`{"Data": 123}`), true},
	}

	type testData struct {
		Data any `validation:"alpha3Currency"`
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
				require.True(t, errorBag.HasFailedKeyAndRule("Data", "alpha3Currency"))
				require.Equal(t, 1, errorBag.CountErrors())
			} else {
				require.True(t, errorBag == nil)
				require.Equal(t, 0, errorBag.CountErrors())
			}
		})
	}
}
