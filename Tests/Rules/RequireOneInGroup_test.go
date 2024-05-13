package Rules

import (
	"errors"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_fails_require_one_in_group_if_no_fields_are_present(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{}`)
	type testData struct {
		Data1 any `validation:"requireOneInGroup:Group1"`
		Data2 any `validation:"requireOneInGroup:Group1"`
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data1", "requireOneInGroup"))
	require.True(t, errorBag.HasFailedKeyAndRule("Data2", "requireOneInGroup"))
	require.Equal(t, 2, errorBag.CountErrors())
}

func Test_it_fails_require_one_in_group_if_two_fields_are_present(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data1": true, "Data2": false}`)
	type testData struct {
		Data1 any `validation:"requireOneInGroup:Group1"`
		Data2 any `validation:"requireOneInGroup:Group1"`
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data1", "requireOneInGroup"))
	require.True(t, errorBag.HasFailedKeyAndRule("Data2", "requireOneInGroup"))
	require.Equal(t, 2, errorBag.CountErrors())
}

func Test_it_fails_require_one_in_group_if_all_fields_are_null(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data1": null, "Data2": null}`)
	type testData struct {
		Data1 any `validation:"requireOneInGroup:Group1"`
		Data2 any `validation:"requireOneInGroup:Group1"`
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data1", "requireOneInGroup"))
	require.True(t, errorBag.HasFailedKeyAndRule("Data2", "requireOneInGroup"))
	require.Equal(t, 2, errorBag.CountErrors())
}

func Test_it_accepts_require_one_in_group_if_exactly_one_field_is_present(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data1": true}`)
	type testData struct {
		Data1 any `validation:"requireOneInGroup:Group1"`
		Data2 any `validation:"requireOneInGroup:Group1"`
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}

func Test_it_accepts_require_one_in_group_if_exactly_one_field_is_present_and_another_is_null(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data1": true, "Data2": null}`)
	type testData struct {
		Data1 any `validation:"requireOneInGroup:Group1"`
		Data2 any `validation:"requireOneInGroup:Group1"`
	}

	// Act
	var data testData
	err := JsonValidator.New().Validate(jsonString, &data)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, errorBag.CountErrors())
}
