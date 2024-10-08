package Rules

import (
	"errors"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_fails_requiredWithout_for_not_present_fields_when_sibling_is_there_and_null(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling": null}`)
	type testData struct {
		Data    int `validation:"requiredWithout:Sibling"`
		Sibling any
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data", "requiredWithout"))
	require.Equal(t, 1, errorBag.CountErrors())
}

func Test_it_does_not_fail_requiredWithout_for_not_present_fields_when_sibling_is_there_and_not_null(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling": true}`)
	type testData struct {
		Data    int `validation:"requiredWithout:Sibling"`
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

func Test_it_fails_requiredWithout_for_not_present_fields_when_sibling_is_there_and_null_with_json_alias(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"sibling": null}`)
	type testData struct {
		Data    int `json:"data" validation:"requiredWithout:Sibling"`
		Sibling any `json:"sibling"`
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("data", "requiredWithout"))
	require.Equal(t, 1, errorBag.CountErrors())
}

func Test_it_does_not_fail_requiredWithout_for_not_present_fields_when_sibling_is_there_and_not_null_with_json_alias(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"sibling": true}`)
	type testData struct {
		Data    int `json:"data" validation:"requiredWithout:Sibling"`
		Sibling any `json:"sibling"`
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}

func Test_it_fails_requiredWithout_for_not_present_fields_when_sibling_is_not_there(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{}`)
	type testData struct {
		Data    int `validation:"requiredWithout:Sibling"`
		Sibling any
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data", "requiredWithout"))
	require.Equal(t, 1, errorBag.CountErrors())
}

func Test_it_fails_requiredWithout_for_not_present_fields_when_sibling_is_not_there_with_json_aliases(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{}`)
	type testData struct {
		Data    int `json:"data" validation:"requiredWithout:Sibling"`
		Sibling any `json:"sibling"`
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("data", "requiredWithout"))
	require.Equal(t, 1, errorBag.CountErrors())
}

func Test_it_does_not_fail_requiredWithout_for_not_present_fields_when_sibling_is_there_in_nested_struct(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"child": {"sibling": true}}`)
	type nestedData struct {
		Data    int `json:"data" validation:"requiredWithout:Sibling"`
		Sibling any `json:"sibling" validation:"requiredWithout:Data"`
	}

	type testData struct {
		Child nestedData `json:"child"`
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}

func Test_it_fails_requiredWithout_for_not_present_fields_when_no_fields_are_present(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{}`)
	type testData struct {
		Data    int `validation:"requiredWithout:Sibling"`
		Sibling any `validation:"requiredWithout:Data"`
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data", "requiredWithout"))
	require.True(t, errorBag.HasFailedKeyAndRule("Sibling", "requiredWithout"))
	require.Equal(t, 2, errorBag.CountErrors())
}

func Test_it_does_not_fails_requiredWithout_for_present_fields_when_sibling_is_there(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling": null, "Data": 0}`)
	type testData struct {
		Data    int `validation:"requiredWithout:Sibling"`
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

func Test_it_fails_requiredWithout_for_null_fields_when_sibling_is_present_and_null(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling": null, "Data": null}`)
	type testData struct {
		Data    int `validation:"requiredWithout:Sibling"`
		Sibling any
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data", "requiredWithout"))
	require.Equal(t, 1, errorBag.CountErrors())
}

func Test_it_does_not_fail_requiredWithout_for_null_fields_when_sibling_is_present_and_not_null(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling": true, "Data": null}`)
	type testData struct {
		Data    int `validation:"requiredWithout:Sibling"`
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

func Test_it_does_not_fail_requiredWithout_for_zero_value_for_int(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling": null, "Data": 0}`)
	type testData struct {
		Data    int `validation:"requiredWithout:Sibling"`
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

func Test_it_does_not_fail_requiredWithout_for_zero_value_for_bool(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling": null, "Data": false}`)
	type testData struct {
		Data    bool `validation:"requiredWithout:Sibling"`
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

func Test_it_does_not_fail_requiredWithout_for_zero_value_for_strings(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling": null, "Data": ""}`)
	type testData struct {
		Data    string `validation:"requiredWithout:Sibling"`
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

func Test_it_does_not_fail_requiredWithout_for_zero_value_for_arrays(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling": null, "Data": []}`)
	type testData struct {
		Data    []any `validation:"requiredWithout:Sibling"`
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

func Test_it_does_not_fail_requiredWithout_for_zero_value_for_objects(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling": null, "Data": {}}`)
	type testData struct {
		Data    any `validation:"requiredWithout:Sibling"`
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
