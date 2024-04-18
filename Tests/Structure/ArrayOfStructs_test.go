package Structure

import (
	"errors"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_validate_valid_arrays_of_structs(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Child": [{"Data": "ab"}, {"Data": "bc"}]}`)
	type testDataChild struct {
		Data string `validation:"len:2"`
	}

	type testDataParent struct {
		Child []testDataChild `validation:"required"`
	}

	// Act
	var data testDataParent
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}

func Test_it_can_validate_invalid_arrays_of_structs(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Child": [{"Data": "ab"}, {"Data": "bcd"}]}`)
	type testDataChild struct {
		Data string `validation:"len:2"`
	}

	type testDataParent struct {
		Child []testDataChild `validation:"required"`
	}

	// Act
	var data testDataParent
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Child.1.Data", "len"))
	require.Equal(t, 1, errorBag.CountErrors())
}

func Test_it_can_validate_valid_arrays_of_arrays_of_structs(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Child": [[{"Data": "aa"}, {"Data": "bb"}], [{"Data": "cc"}]]}`)
	type testDataChild struct {
		Data string `validation:"len:2"`
	}

	type testDataParent struct {
		Child [][]testDataChild `validation:"required"`
	}

	// Act
	var data testDataParent
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}

func Test_it_can_validate_invalid_arrays_of_arrays_of_structs(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Child": [[{"Data": "ab"}, {"Data": "bcd"}], [{"Data": "abc"}, {"Data": "bc"}]]}`)
	type testDataChild struct {
		Data string `validation:"len:2"`
	}

	type testDataParent struct {
		Child [][]testDataChild `validation:"required"`
	}

	// Act
	var data testDataParent
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Child.0.1.Data", "len"))
	require.True(t, errorBag.HasFailedKeyAndRule("Child.1.0.Data", "len"))
	require.Equal(t, 2, errorBag.CountErrors())
}
