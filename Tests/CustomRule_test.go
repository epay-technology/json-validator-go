package Tests

import (
	"errors"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_register_custom_rules(t *testing.T) {
	// Arrange
	validator := JsonValidator.New()
	validator.RegisterRule(JsonValidator.Rule{
		Name:           "MyRule",
		IsPresenceRule: false,
		IsNullableRule: false,
		Function: func(context *JsonValidator.FieldValidationContext) (string, bool) {
			return "My RuleContext Ran", false
		},
	})

	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data": false}`)
	type testData struct {
		Data []any `validation:"MyRule:2"`
	}

	// Act
	var data testData
	err := validator.Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data", "MyRule"))
	require.Equal(t, 1, errorBag.CountErrors())
}
