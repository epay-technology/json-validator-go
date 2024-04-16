package Tests

import (
	"errors"
	"github.com/stretchr/testify/require"
	"json-validator-go/JsonValidator"
	"testing"
)

func Test_it_can_validated_valid_nested_structs(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Child": {"Data": [1,2]}}`)
	type testDataChild struct {
		Data []any `validation:"len:2"`
	}

	type testDataParent struct {
		Child testDataChild `validation:"required"`
	}

	// Act
	_, err := JsonValidator.Validate[testDataParent](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
}

func Test_it_can_validated_errors_in_nested_structures(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Child": {"Data": [1]}}`)
	type testDataChild struct {
		Data []any `validation:"len:2"`
	}

	type testDataParent struct {
		Child testDataChild `validation:"required"`
	}

	// Act
	_, err := JsonValidator.Validate[testDataParent](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Child.Data", "len"))
}

func Test_it_does_not_run_validation_on_nested_structures_if_parent_is_not_present(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{}`)
	type testDataChild struct {
		Data []any `validation:"len:2"`
	}

	type testDataParent struct {
		Child testDataChild `validation:"required"`
	}

	// Act
	_, err := JsonValidator.Validate[testDataParent](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Child", "required"))
	require.Len(t, errorBag.Errors, 1)
}

func Test_it_does_not_run_validation_on_nested_structures_if_parent_is_null(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Child": null}`)
	type testDataChild struct {
		Data []any `validation:"len:2"`
	}

	type testDataParent struct {
		Child testDataChild `validation:"required"`
	}

	// Act
	_, err := JsonValidator.Validate[testDataParent](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Child", "required"))
	require.Len(t, errorBag.Errors, 1)
}

func Test_it_does_not_run_validation_on_nested_structures_if_parent_has_errors(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Child": [1,2,3]}`)
	type testDataChild struct {
		Data []any `validation:"len:2"`
	}

	type testDataParent struct {
		Child testDataChild `validation:"object|len:999"`
	}

	// Act
	_, err := JsonValidator.Validate[testDataParent](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.NotNil(t, errorBag)
	require.True(t, errorBag.HasFailedKeyAndRule("Child", "object"))
	require.True(t, errorBag.HasFailedKeyAndRule("Child", "len"))
	require.Len(t, errorBag.Errors, 1)
}
