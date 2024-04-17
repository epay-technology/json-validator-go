package Tests

import (
	"errors"
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_use_aliases_for_minLen_to_lenMin(t *testing.T) {
	// Arrange
	JsonValidator.RegisterAlias("MyAlias", "minLen")

	var errorBag *JsonValidator.ErrorBag
	jsonString := []byte(`{}`)
	type testData struct {
		Data []any `validation:"MyAlias:2"`
	}

	// Act
	_, err := JsonValidator.Validate[testData](jsonString)
	_ = errors.As(err, &errorBag)

	// Assert
	require.Error(t, err)
	require.True(t, errorBag.HasFailedKeyAndRule("Data", "MyAlias"))
}
