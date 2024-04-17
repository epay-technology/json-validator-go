package Rules

import (
	"errors"
	"fmt"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_validate_using_uuid_rule(t *testing.T) {
	// Setup
	cases := []struct {
		jsonString []byte
		shouldFail bool
	}{
		{[]byte(`{"Data": "0d0f123e-41f7-0f94-9c53-275ce5a32c16"}`), false},
		{[]byte(`{"Data": "ad9d09c2-615f-12fa-ab28-11f0c2056c20"}`), false},
		{[]byte(`{"Data": "76f4201d-6b9e-2a54-9df2-04295158ebc1"}`), false},
		{[]byte(`{"Data": "4ef7dd5b-e10a-306b-ae28-38730d303236"}`), false},
		{[]byte(`{"Data": "0d0f123e-41f7-4f94-9c53-275ce5a32c16"}`), false},
		{[]byte(`{"Data": "4ef7dd5b-e10a-506b-ae28-38730d303236"}`), false},
		{[]byte(`{"Data": "31283319-d2e9-6b81-a02a-2f9e1129ad45"}`), false},
		{[]byte(`{"Data": "1146b96b-165c-720a-a684-5715cc398edc"}`), false},
		{[]byte(`{"Data": "d3ff5fc7-dc51-81ec-9c94-2cb093d74090"}`), false},
		{[]byte(`{"Data": "2bd93de8-1989-9af8-9518-a378acbdc9c2"}`), false},
		{[]byte(`{"Data": "00000000-0000-0000-0000-000000000000"}`), false},
		{[]byte(`{"Data": "d-00000000-0000-0000-0000-000000000000"}`), true},
		{[]byte(`{"Data": "00000000-0000-0000-0000-000000000000-a"}`), true},
		{[]byte(`{"Data": "d-00000000-0000-0000-0000-000000000000-a"}`), true},
		{[]byte(`{"Data": " 00000000-0000-0000-0000-000000000000 "}`), true},
		{[]byte(`{"Data": " 00000000-0000-0000-0000-000000000000"}`), true},
		{[]byte(`{"Data": "00000000-0000-0000-0000-000000000000 "}`), true},
		{[]byte(`{}`), true},
		{[]byte(`{"Data": {}}`), true},
		{[]byte(`{"Data": "hello world"}`), true},
		{[]byte(`{"Data": 123}`), true},
		{[]byte(`{"Data": null}`), true},
		{[]byte(`{"Data": true}`), true},
		{[]byte(`{"Data": []}`), true},
	}

	type testData struct {
		Data any `validation:"uuid"`
	}

	for i, testCase := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			// Arrange
			var errorBag *JsonValidator.ErrorBag

			// Act
			_, err := JsonValidator.Validate[testData](testCase.jsonString)
			_ = errors.As(err, &errorBag)

			// Assert
			if testCase.shouldFail {
				require.True(t, errorBag != nil)
				require.True(t, errorBag.HasFailedKeyAndRule("Data", "uuid"))
				require.Equal(t, 1, errorBag.CountErrors())
			} else {
				require.True(t, errorBag == nil)
				require.Equal(t, 0, errorBag.CountErrors())
			}
		})
	}
}
