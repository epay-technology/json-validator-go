package Rules

import (
	"errors"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_fails_present_for_not_present_fields(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{}`)
	type testData struct {
		Data int `validation:"present"`
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data", "present"))
	require.Equal(t, 1, errorBag.CountErrors())
}

func Test_it_does_not_fail_present_for_null_fields(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data": null}`)
	type testData struct {
		Data int `validation:"present|nullable"`
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}

func Test_it_does_not_fail_present_for_zero_value_for_int(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data": 0}`)
	type testData struct {
		Data int `validation:"present"`
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}

func Test_it_does_not_fail_present_for_zero_value_for_bool(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data": false}`)
	type testData struct {
		Data bool `validation:"present"`
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}

func Test_it_does_not_fail_present_for_zero_value_for_strings(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data": ""}`)
	type testData struct {
		Data string `validation:"present"`
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}

func Test_it_does_not_fail_present_for_zero_value_for_arrays(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data": []}`)
	type testData struct {
		Data []any `validation:"present"`
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}

func Test_it_does_not_fail_present_for_zero_value_for_objects(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data": {}}`)
	type testData struct {
		Data any `validation:"present"`
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}

func Test_it_does_not_run_non_presence_rules_for_a_field_that_was_not_present(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{}`)
	type testData struct {
		Data any `validation:"string"`
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}
