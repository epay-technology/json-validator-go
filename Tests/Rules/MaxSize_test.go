package Rules

import (
	"errors"
	"fmt"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_validate_using_max_size_rule(t *testing.T) {
	// Setup
	cases := []struct {
		jsonString []byte
		shouldFail bool
	}{
		{[]byte(`{"Data": ""}`), false},
		{[]byte(`{"Data": "123456"}`), false},
		{[]byte(`{"Data": {"a": 0}}`), false},
		{[]byte(`{"Data": false}`), false},
		{[]byte(`{"Data": true}`), false},
		{[]byte(`{"Data": null}`), false},
		{[]byte(`{"Data": 123}`), false},
		{[]byte(`{"Data": 123, "otherKey": "hello very long message does not matter"}`), false},
		{[]byte(`{"Data": {"key": "largeValue"}}`), true},
		{[]byte(`{"Data": {"LongKey": 0}}`), true},
		{[]byte(`{"Data": "123456789"}`), true},
		{[]byte(`{"Data": 1234567890123}`), true},
	}

	type testData struct {
		Data any `validation:"maxSize:8"`
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
				require.True(t, errorBag.HasFailedKeyAndRule("Data", "maxSize"))
				require.Equal(t, 1, errorBag.CountErrors())
			} else {
				require.True(t, errorBag == nil)
				require.Equal(t, 0, errorBag.CountErrors())
			}
		})
	}
}
