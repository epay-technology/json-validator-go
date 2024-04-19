package Structure

import (
	"errors"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_validate_valid_maps_of_structs(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Child": {"a": {"Data": "AA"}, "b": {"Data": "BB"}}}`)
	type testDataChild struct {
		Data string `validation:"len:2"`
	}

	type testDataParent struct {
		Child map[string]testDataChild `validation:"required"`
	}

	// Act
	var data testDataParent
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}

func Test_it_can_validate_invalid_maps_of_structs(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Child": {"a": {"Data": "AAC"}, "b": {"Data": "BBD"}}}`)
	type testDataChild struct {
		Data string `validation:"len:2"`
	}

	type testDataParent struct {
		Child map[string]testDataChild `validation:"required"`
	}

	// Act
	var data testDataParent
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.Equal(t, 2, errorBag.CountErrors())
	require.True(t, errorBag.HasFailedKeyAndRule("Child.a.Data", "len"))
	require.True(t, errorBag.HasFailedKeyAndRule("Child.b.Data", "len"))
}
