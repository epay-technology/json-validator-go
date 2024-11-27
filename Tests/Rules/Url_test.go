package Rules

import (
	"errors"
	"fmt"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_validate_using_url_rule(t *testing.T) {
	// Setup
	cases := []struct {
		jsonString []byte
		shouldFail bool
	}{
		{[]byte(`{"Data": "https://www.google.dk/"}`), false},
		{[]byte(`{"Data": "http://www.google.dk"}`), false},
		{[]byte(`{"Data": "http://google.dk"}`), false},
		{[]byte(`{"Data": "https://google.dk"}`), false},
		{[]byte(`{"Data": "https://1234.dk"}`), false},
		{[]byte(`{"Data": "http://www.google.dk:8000"}`), true},
		{[]byte(`{"Data": "http://google"}`), true},
		{[]byte(`{"Data": "http://127.125"}`), true},
		{[]byte(`{"Data": "www.google.dk"}`), true},
		{[]byte(`{"Data": "google.dk"}`), true},
		{[]byte(`{"Data": [1,2,3]}`), true},
		{[]byte(`{"Data": {}}`), true},
		{[]byte(`{"Data": "hello world"}`), true},
		{[]byte(`{"Data": 123}`), true},
		{[]byte(`{"Data": null}`), true},
		{[]byte(`{"Data": true}`), true},
		{[]byte(`{"Data": "ftp://ftp.example.com"}`), true},
		{[]byte(`{"Data": "tcp://example.com:54321"}`), true},
		{[]byte(`{"Data": "udp://example.com:6000"}`), true},
		{[]byte(`{"Data": "wss://example.com"}`), true},
		{[]byte(`{"Data": "https://127.0.0.1"}`), true},
		{[]byte(`{"Data": "http://127.0.0.1"}`), true},
		{[]byte(`{"Data": "127.0.0.1"}`), true},
		{[]byte(`{"Data": "localhost"}`), true},
		{[]byte(`{"Data": "http://localhost"}`), true},
		{[]byte(`{"Data": "http://localhost"}`), true},
		{[]byte(`{"Data": "https://localhost"}`), true},
		{[]byte(`{"Data": "http://localhost:8000"}`), true},
		{[]byte(`{"Data": "https://localhost:8000"}`), true},
	}

	type testData struct {
		Data any `validation:"url"`
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
				require.True(t, errorBag.HasFailedKeyAndRule("Data", "url"))
				require.Equal(t, 1, errorBag.CountErrors())
			} else {
				require.True(t, errorBag == nil)
				require.Equal(t, 0, errorBag.CountErrors())
			}
		})
	}
}
