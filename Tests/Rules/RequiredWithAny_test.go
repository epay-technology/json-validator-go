package Rules

import (
	"errors"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_fails_requireWithAny_for_not_present_fields_when_one_sibling_is_there(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling1": null}`)
	type testData struct {
		Data     int `validation:"requiredWithAny:Sibling1,Sibling2"`
		Sibling1 any
		Sibling2 any
	}

	// Act
	_, err := JsonValidator.Validate[testData](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data", "requiredWithAny"))
}

func Test_it_fails_requireWithAny_for_not_present_fields_when_all_siblings_are_there(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling1": true, "Sibling2": true}`)
	type testData struct {
		Data     int `validation:"requiredWithAny:Sibling1,Sibling2"`
		Sibling1 any
		Sibling2 any
	}

	// Act
	_, err := JsonValidator.Validate[testData](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data", "requiredWithAny"))
}

func Test_it_does_not_fail_requireWithAny_for_not_present_fields_no_siblings_are_there(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{}`)
	type testData struct {
		Data     int `validation:"requiredWithAny:Sibling1,Sibling2"`
		Sibling1 any
		Sibling2 any
	}

	// Act
	_, err := JsonValidator.Validate[testData](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
}

func Test_it_does_not_fails_requireWithAny_for_present_fields_when_one_sibling_is_there(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling1": null, "Data": 0}`)
	type testData struct {
		Data     int `validation:"requiredWithAny:Sibling1,Sibling2"`
		Sibling1 any
		Sibling2 any
	}

	// Act
	_, err := JsonValidator.Validate[testData](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
}

func Test_it_fails_requireWithAny_for_null_fields_when_one_sibling_is_present(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling1": null, "Data": null}`)
	type testData struct {
		Data     int `validation:"requiredWithAny:Sibling1,Sibling2"`
		Sibling1 any
		Sibling2 any
	}

	// Act
	_, err := JsonValidator.Validate[testData](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data", "requiredWithAny"))
}

func Test_it_does_not_fail_requireWithAny_for_zero_value_for_int(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling1": null, "Data": 0}`)
	type testData struct {
		Data     int `validation:"requiredWithAny:Sibling1,Sibling2"`
		Sibling1 any
		Sibling2 any
	}

	// Act
	_, err := JsonValidator.Validate[testData](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
}

func Test_it_does_not_fail_requireWithAny_for_zero_value_for_bool(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling1": null, "Data": false}`)
	type testData struct {
		Data     bool `validation:"requiredWithAny:Sibling1,Sibling2"`
		Sibling1 any
		Sibling2 any
	}

	// Act
	_, err := JsonValidator.Validate[testData](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
}

func Test_it_does_not_fail_requireWithAny_for_zero_value_for_strings(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling1": null, "Data": ""}`)
	type testData struct {
		Data     string `validation:"requiredWithAny:Sibling1,Sibling2"`
		Sibling1 any
		Sibling2 any
	}

	// Act
	_, err := JsonValidator.Validate[testData](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
}

func Test_it_does_not_fail_requireWithAny_for_zero_value_for_arrays(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling1": null, "Data": []}`)
	type testData struct {
		Data     []any `validation:"requiredWithAny:Sibling1,Sibling2"`
		Sibling1 any
		Sibling2 any
	}

	// Act
	_, err := JsonValidator.Validate[testData](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
}

func Test_it_does_not_fail_requireWithAny_for_zero_value_for_objects(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Sibling1": null, "Data": {}}`)
	type testData struct {
		Data     any `validation:"requiredWithAny:Sibling1,Sibling2"`
		Sibling1 any
		Sibling2 any
	}

	// Act
	_, err := JsonValidator.Validate[testData](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
}
