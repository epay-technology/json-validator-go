package Structure

import (
	"errors"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_validated_valid_embedded_structs(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data": [1,2]}`)
	type testDataChild struct {
		Data []any `validation:"len:2"`
	}

	type testDataParent struct {
		testDataChild
	}

	// Act
	var data testDataParent
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}

func Test_it_can_validate_valid_deeply_embedded_structs(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data": [1,2]}`)
	type a struct {
		Data []any `validation:"len:2"`
	}

	type b struct {
		a
	}

	type c struct {
		b
	}

	type d struct {
		c
	}

	// Act
	var data d
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}

func Test_it_can_validated_errors_in_embedded_structures(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data": [1]}`)
	type testDataChild struct {
		Data []any `validation:"len:2"`
	}

	type testDataParent struct {
		testDataChild
	}

	// Act
	var data testDataParent
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data", "len"))
	require.Equal(t, 1, errorBag.CountErrors())
}

func Test_it_can_validate_errors_in_deeply_embedded_structs(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data": [1,2,3]}`)
	type a struct {
		Data []any `validation:"len:2"`
	}

	type b struct {
		a
	}

	type c struct {
		b
	}

	type d struct {
		c
	}

	// Act
	var data d
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data", "len"))
	require.Equal(t, 1, errorBag.CountErrors())
}

func Test_it_can_validate_errors_in_embedded_structs_with_parent_keys(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data": [1,2,3], "Key": "a"}`)
	type a struct {
		Data []any `validation:"len:2"`
	}

	type b struct {
		a
		Key string `validation:"len:2"`
	}

	// Act
	var data b
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data", "len"))
	require.True(t, errorBag.HasFailedKeyAndRule("Key", "len"))
	require.Equal(t, 2, errorBag.CountErrors())
}
