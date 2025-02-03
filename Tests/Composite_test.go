package Tests

import (
	"errors"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_use_composite_rules_for_valid_cases(t *testing.T) {
	// Arrange
	validator := JsonValidator.New()
	validator.RegisterComposite("MyComposite", "required|array|minLen:2")

	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data": ["a", "b"]}`)
	type testData struct {
		Data []any `validation:"MyComposite"`
	}

	// Act
	var data testData
	err := validator.Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
}

func Test_it_can_use_composite_rules_and_fail_with_actual_underlying_rule(t *testing.T) {
	// Arrange
	validator := JsonValidator.New()
	validator.RegisterComposite("MyComposite", "required|array|minLen:2")

	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data": ["a"]}`)
	type testData struct {
		Data []any `validation:"MyComposite"`
	}

	// Act
	var data testData
	err := validator.Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data", "minLen"))
	require.Equal(t, 1, errorBag.CountErrors())
}

func Test_it_can_use_composite_rules_and_fail_with_presence_rules(t *testing.T) {
	// Arrange
	validator := JsonValidator.New()
	validator.RegisterComposite("MyComposite", "required|array|minLen:2")

	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{}`)
	type testData struct {
		Data []any `validation:"MyComposite"`
	}

	// Act
	var data testData
	err := validator.Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data", "required"))
	require.Equal(t, 1, errorBag.CountErrors())
}

func Test_it_can_use_composite_rules_with_arguments_no_error(t *testing.T) {
	// Arrange
	validator := JsonValidator.New()
	validator.RegisterComposite("MyComposite", "required|minLen:$0|maxLen:$1|array")

	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data": ["a"]}`)
	type testData struct {
		Data []any `validation:"MyComposite:1,2"`
	}

	// Act
	var data testData
	err := validator.Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
}

func Test_it_can_use_composite_rules_with_arguments_with_error(t *testing.T) {
	// Arrange
	validator := JsonValidator.New()
	validator.RegisterComposite("MyComposite", "required|minLen:$0|maxLen:$1|array")

	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data": []}`)
	type testData struct {
		Data []any `validation:"MyComposite:1,2"`
	}

	// Act
	var data testData
	err := validator.Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data", "minLen"))
	require.Equal(t, 1, errorBag.CountErrors())
}

func Test_it_can_use_composite_rules_with_arguments_with_another_error(t *testing.T) {
	// Arrange
	validator := JsonValidator.New()
	validator.RegisterComposite("MyComposite", "required|minLen:$0|maxLen:$1|array")

	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data": ["a", "b", "c"]}`)
	type testData struct {
		Data []any `validation:"MyComposite:1,2"`
	}

	// Act
	var data testData
	err := validator.Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data", "maxLen"))
	require.Equal(t, 1, errorBag.CountErrors())
}
