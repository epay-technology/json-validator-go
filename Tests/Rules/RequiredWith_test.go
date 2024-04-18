package Rules

import (
	"errors"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_fails_requireWith_for_not_present_fields_when_sibling_is_there(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling": null}`)
	type testData struct {
		Data    int `validation:"requiredWith:Sibling"`
		Sibling any
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data", "requiredWith"))
	require.Equal(t, 1, errorBag.CountErrors())
}

func Test_it_does_not_fails_requireWith_for_present_fields_when_sibling_is_there(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling": null, "Data": 0}`)
	type testData struct {
		Data    int `validation:"requiredWith:Sibling"`
		Sibling any
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}

func Test_it_fails_requireWith_for_null_fields_when_sibling_is_present(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling": null, "Data": null}`)
	type testData struct {
		Data    int `validation:"requiredWith:Sibling"`
		Sibling any
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data", "requiredWith"))
	require.Equal(t, 1, errorBag.CountErrors())
}

func Test_it_does_not_fail_requireWith_for_zero_value_for_int(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling": null, "Data": 0}`)
	type testData struct {
		Data    int `validation:"requiredWith:Sibling"`
		Sibling any
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}

func Test_it_does_not_fail_requireWith_for_zero_value_for_bool(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling": null, "Data": false}`)
	type testData struct {
		Data    bool `validation:"requiredWith:Sibling"`
		Sibling any
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}

func Test_it_does_not_fail_requireWith_for_zero_value_for_strings(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling": null, "Data": ""}`)
	type testData struct {
		Data    string `validation:"requiredWith:Sibling"`
		Sibling any
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}

func Test_it_does_not_fail_requireWith_for_zero_value_for_arrays(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling": null, "Data": []}`)
	type testData struct {
		Data    []any `validation:"requiredWith:Sibling"`
		Sibling any
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}

func Test_it_does_not_fail_requireWith_for_zero_value_for_objects(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling": null, "Data": {}}`)
	type testData struct {
		Data    any `validation:"requiredWith:Sibling"`
		Sibling any
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}
