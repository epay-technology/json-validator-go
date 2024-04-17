package Rules

import (
	"errors"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_allows_null_for_nullable_fields_when_pointer_type(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data": null}`)
	type testData struct {
		Data *int `validation:"nullable"`
	}

	// Act
	_, err := JsonValidator.Validate[testData](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}

func Test_it_allows_null_for_nullable_fields_even_when_not_a_pointer_type(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data": null}`)
	type testData struct {
		Data int `validation:"nullable"`
	}

	// Act
	_, err := JsonValidator.Validate[testData](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}

func Test_it_does_not_run_other_validation_rules_when_value_is_null(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data": null}`)
	type testData struct {
		Data int `validation:"nullable|integer"`
	}

	// Act
	_, err := JsonValidator.Validate[testData](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}
