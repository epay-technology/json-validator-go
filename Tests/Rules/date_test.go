package Rules

import (
	"errors"
	"fmt"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_validate_using_date_rule(t *testing.T) {
	// Setup
	cases := []struct {
		jsonString []byte
		shouldFail bool
	}{
		{[]byte(`{"Data": 123}`), true},
		{[]byte(`{"Data": 45.12}`), true},
		{[]byte(`{"Data": "abc"}`), true},
		{[]byte(`{"Data": false}`), true},
		{[]byte(`{"Data": ""}`), true},
		{[]byte(`{"Data": "hello world"}`), true},
		{[]byte(`{"Data": "ab"}`), true},
		{[]byte(`{"Data": "bc"}`), true},
		{[]byte(`{"Data": "ABC"}`), true},
		{[]byte(`{"Data": "aBc"}`), true},
		{[]byte(`{"Data": "AbC"}`), true},
		{[]byte(`{"Data": 12}`), true},
		{[]byte(`{"Data": 23}`), true},
		{[]byte(`{"Data": 13}`), true},
		{[]byte(`{"Data": 123.45}`), true},
		{[]byte(`{"Data": {}}`), true},
		{[]byte(`{"Data": [1,2,3]}`), true},
		{[]byte(`{"Data": null}`), true},
		{[]byte(`{"Data": true}`), true},
		{[]byte(`{"Data": "2024-00-01"}`), true},
		{[]byte(`{"Data": "2024-01-00"}`), true},
		{[]byte(`{"Data": "2024-01-32"}`), true},
		{[]byte(`{"Data": "2024-13-01"}`), true},
		{[]byte(`{"Data": " 2024-01-01"}`), true},
		{[]byte(`{"Data": "2024-01-01 "}`), true},
		{[]byte(`{"Data": "a2024-01-01"}`), true},
		{[]byte(`{"Data": "2024-01-01a"}`), true},
		{[]byte(`{"Data": "12024-01-01"}`), true},
		{[]byte(`{"Data": "2024-011-01"}`), true},
		{[]byte(`{"Data": "2024-01-011"}`), true},
		{[]byte(`{"Data": "0000-01-01"}`), false},
		{[]byte(`{"Data": "9999-12-31"}`), false},
		{[]byte(`{"Data": "2024-01-01"}`), false},
	}

	type testData struct {
		Data any `validation:"date"`
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
				require.True(t, errorBag.HasFailedKeyAndRule("Data", "date"))
				require.Equal(t, 1, errorBag.CountErrors())
			} else {
				require.True(t, errorBag == nil)
				require.Equal(t, 0, errorBag.CountErrors())
			}
		})
	}
}
