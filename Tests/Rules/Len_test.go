package Rules

import (
	"errors"
	"github.com/stretchr/testify/require"
	"json-validator-go/JsonValidator"
	"testing"
)

func Test_it_does_not_fail_len_rule_for_right_length_arrays(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data": [1,2]}`)
	type testData struct {
		Data []any `validation:"len:2"`
	}

	// Act
	_, err := JsonValidator.Validate[testData](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.NoError(t, err)
}

func Test_it_fails_len_rule_for_too_short_arrays(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{}`)
	type testData struct {
		Data []any `validation:"len:2"`
	}

	// Act
	_, err := JsonValidator.Validate[testData](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data", "len"))
}

func Test_it_fails_len_rule_for_too_long_arrays(t *testing.T) {
	// Arrange
	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{"Data": [1,2,3]}`)
	type testData struct {
		Data []any `validation:"len:2"`
	}

	// Act
	_, err := JsonValidator.Validate[testData](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data", "len"))
}
