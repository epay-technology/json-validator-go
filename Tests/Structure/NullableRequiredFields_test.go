package Structure

import (
	"errors"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_validate_valid_required_fields_with_the_nullable_rule(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"A":"Hello", "B": "World"}`)
	type testDataParent struct {
		A *string `validation:"requiredWithout:B|nullable|string"`
		B *string `validation:"requiredWithout:A|nullable|string"`
	}

	// Act
	_, err := JsonValidator.Validate[testDataParent](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}

func Test_it_can_validate_valid_required_fields_with_the_nullable_rule_when_one_field_is_not_present(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"B": "World"}`)
	type testDataParent struct {
		A *string `validation:"requiredWithout:B|nullable|string"`
		B *string `validation:"requiredWithout:A|nullable|string"`
	}

	// Act
	_, err := JsonValidator.Validate[testDataParent](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}

func Test_it_can_validate_invalid_required_fields_with_the_nullable_rule_when_all_fields_are_not_present(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{}`)
	type testDataParent struct {
		A *string `validation:"requiredWithout:B|nullable|string"`
		B *string `validation:"requiredWithout:A|nullable|string"`
	}

	// Act
	_, err := JsonValidator.Validate[testDataParent](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("A", "requiredWithout"))
	require.True(t, errorBag.HasFailedKeyAndRule("B", "requiredWithout"))
	require.Equal(t, 2, errorBag.CountErrors())
}
