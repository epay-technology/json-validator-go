package Tests

import (
	"errors"
	"github.com/stretchr/testify/require"
	"json-validator-go/JsonValidator"
	"testing"
)

func Test_it_can_register_custom_rules(t *testing.T) {
	// Arrange
	JsonValidator.RegisterRule("MyRule", func(context *JsonValidator.FieldValidationContext) (string, bool) {
		return "My Rule Ran", false
	})

	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{}`)
	type testData struct {
		Data []any `validation:"MyRule:2"`
	}

	// Act
	_, err := JsonValidator.Validate[testData](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data", "MyRule"))
}
