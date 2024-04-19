package Analyzer

import (
	"github.com/epay-technology/json-validator-go/JsonValidator"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_it_can_analyze_simple_structs(t *testing.T) {
	// Arrange
	validator := JsonValidator.New()
	type simpleStruct struct {
		Id   int    `json:"id" validation:"required|integer"`
		Text string `json:"text" validation:"required|string"`
	}

	// Act
	var data simpleStruct
	fieldCache, err := validator.Analyze(&data)

	// Assert
	require.NoError(t, err)
	require.Len(t, fieldCache.Children, 2)

	require.Len(t, fieldCache.Children[0].Children, 0)
	require.Len(t, fieldCache.Children[1].Children, 0)

	require.Nil(t, fieldCache.Parent)

	require.True(t, fieldCache.IsStruct)
	require.False(t, fieldCache.IsSlice)

	require.Same(t, fieldCache, fieldCache.Children[0].Parent)
	require.Same(t, fieldCache, fieldCache.Children[1].Parent)

	require.False(t, fieldCache.Children[0].IsStruct)
	require.False(t, fieldCache.Children[0].IsSlice)

	require.False(t, fieldCache.Children[1].IsStruct)
	require.False(t, fieldCache.Children[1].IsSlice)

	require.Equal(t, "Id", fieldCache.Children[0].StructKey)
	require.Equal(t, "id", fieldCache.Children[0].JsonKey)

	require.Len(t, fieldCache.Children[0].ValidationTag.Rules, 1)
	require.Equal(t, fieldCache.Children[0].ValidationTag.Rules[0].Name, "integer")

	require.Len(t, fieldCache.Children[0].ValidationTag.PresenceRules, 1)
	require.Equal(t, fieldCache.Children[0].ValidationTag.PresenceRules[0].Name, "required")
}

func Test_it_uses_caching(t *testing.T) {
	// Arrange
	validator := JsonValidator.New()
	type simpleStruct struct {
		Id   int    `json:"id" validation:"required|integer"`
		Text string `json:"text" validation:"required|string"`
	}

	// Act
	var data simpleStruct
	fieldCache1, err1 := validator.Analyze(&data)
	fieldCache2, err2 := validator.Analyze(&data)

	// Assert
	require.NoError(t, err1)
	require.NoError(t, err2)

	require.Same(t, fieldCache1, fieldCache2)
}

func Test_it_can_analyze_nested_structs(t *testing.T) {
	// Arrange
	validator := JsonValidator.New()
	type child struct {
		Id int `json:"id" validation:"required|integer"`
	}

	type simpleStruct struct {
		Child child `json:"child" validation:"required|object"`
	}

	// Act
	var data simpleStruct
	fieldCache, err := validator.Analyze(&data)

	// Assert
	require.NoError(t, err)
	require.Len(t, fieldCache.Children, 1)

	require.Len(t, fieldCache.Children[0].Children, 1)

	require.Nil(t, fieldCache.Parent)

	require.True(t, fieldCache.IsStruct)
	require.False(t, fieldCache.IsSlice)

	require.Same(t, fieldCache, fieldCache.Children[0].Parent)

	require.True(t, fieldCache.Children[0].IsStruct)
	require.False(t, fieldCache.Children[0].IsSlice)

	require.False(t, fieldCache.Children[0].Children[0].IsStruct)
	require.False(t, fieldCache.Children[0].Children[0].IsSlice)

	require.Equal(t, "Child", fieldCache.Children[0].StructKey)
	require.Equal(t, "child", fieldCache.Children[0].JsonKey)

	require.Equal(t, "Id", fieldCache.Children[0].Children[0].StructKey)
	require.Equal(t, "id", fieldCache.Children[0].Children[0].JsonKey)

	require.Len(t, fieldCache.Children[0].ValidationTag.Rules, 1)
	require.Equal(t, fieldCache.Children[0].ValidationTag.Rules[0].Name, "object")

	require.Len(t, fieldCache.Children[0].ValidationTag.PresenceRules, 1)
	require.Equal(t, fieldCache.Children[0].ValidationTag.PresenceRules[0].Name, "required")

	require.Len(t, fieldCache.Children[0].Children[0].ValidationTag.Rules, 1)
	require.Equal(t, fieldCache.Children[0].Children[0].ValidationTag.Rules[0].Name, "integer")

	require.Len(t, fieldCache.Children[0].Children[0].ValidationTag.PresenceRules, 1)
	require.Equal(t, fieldCache.Children[0].Children[0].ValidationTag.PresenceRules[0].Name, "required")
}

func Test_it_can_analyze_array_of_structs(t *testing.T) {
	// Arrange
	validator := JsonValidator.New()
	type child struct {
		Id int `json:"id" validation:"required|integer"`
	}

	type simpleStruct struct {
		Children []child `json:"children" validation:"required|array"`
	}

	// Act
	var data simpleStruct
	fieldCache, err := validator.Analyze(&data)

	// Assert
	require.NoError(t, err)
	require.Len(t, fieldCache.Children, 1)

	require.Len(t, fieldCache.Children[0].Children, 1)

	require.Nil(t, fieldCache.Parent)

	require.True(t, fieldCache.IsStruct)
	require.False(t, fieldCache.IsSlice)

	require.Same(t, fieldCache, fieldCache.Children[0].Parent)

	require.False(t, fieldCache.Children[0].IsStruct)
	require.True(t, fieldCache.Children[0].IsSlice)

	require.True(t, fieldCache.Children[0].Children[0].IsStruct)
	require.False(t, fieldCache.Children[0].Children[0].IsSlice)

	require.Equal(t, "Children", fieldCache.Children[0].StructKey)
	require.Equal(t, "children", fieldCache.Children[0].JsonKey)

	require.Equal(t, "Id", fieldCache.Children[0].Children[0].Children[0].StructKey)
	require.Equal(t, "id", fieldCache.Children[0].Children[0].Children[0].JsonKey)

	require.Len(t, fieldCache.Children[0].ValidationTag.Rules, 1)
	require.Equal(t, fieldCache.Children[0].ValidationTag.Rules[0].Name, "array")

	require.Len(t, fieldCache.Children[0].ValidationTag.PresenceRules, 1)
	require.Equal(t, fieldCache.Children[0].ValidationTag.PresenceRules[0].Name, "required")

	require.Len(t, fieldCache.Children[0].Children[0].Children[0].ValidationTag.Rules, 1)
	require.Equal(t, fieldCache.Children[0].Children[0].Children[0].ValidationTag.Rules[0].Name, "integer")

	require.Len(t, fieldCache.Children[0].Children[0].Children[0].ValidationTag.PresenceRules, 1)
	require.Equal(t, fieldCache.Children[0].Children[0].Children[0].ValidationTag.PresenceRules[0].Name, "required")
}

func Test_it_can_analyze_deeply_nested_arrays_of_structs(t *testing.T) {
	// Arrange
	validator := JsonValidator.New()
	type child struct {
		Id int `json:"id" validation:"required|integer"`
	}

	type simpleStruct struct {
		Children [][][]child `json:"children" validation:"required|array"`
	}

	// Act
	var data simpleStruct
	fieldCache, err := validator.Analyze(&data)

	// Assert
	require.NoError(t, err)
	// Root
	require.Len(t, fieldCache.Children, 1)
	require.True(t, fieldCache.IsStruct)
	require.False(t, fieldCache.IsSlice)

	// Children field in root
	require.Len(t, fieldCache.Children[0].Children, 1)
	require.False(t, fieldCache.Children[0].IsStruct)
	require.True(t, fieldCache.Children[0].IsSlice)

	// First sub level array
	require.Len(t, fieldCache.Children[0].Children[0].Children, 1)
	require.False(t, fieldCache.Children[0].Children[0].IsStruct)
	require.True(t, fieldCache.Children[0].Children[0].IsSlice)

	// Second sub level array
	require.Len(t, fieldCache.Children[0].Children[0].Children[0].Children, 1)
	require.False(t, fieldCache.Children[0].Children[0].Children[0].IsStruct)
	require.True(t, fieldCache.Children[0].Children[0].Children[0].IsSlice)

	// The struct type of the inner most level of the nested array
	require.Len(t, fieldCache.Children[0].Children[0].Children[0].Children[0].Children, 1)
	require.True(t, fieldCache.Children[0].Children[0].Children[0].Children[0].IsStruct)
	require.False(t, fieldCache.Children[0].Children[0].Children[0].Children[0].IsSlice)

	// The first field of the inner most struct type
	require.Len(t, fieldCache.Children[0].Children[0].Children[0].Children[0].Children[0].Children, 0)
	require.False(t, fieldCache.Children[0].Children[0].Children[0].Children[0].Children[0].IsStruct)
	require.False(t, fieldCache.Children[0].Children[0].Children[0].Children[0].Children[0].IsSlice)

	innerField := fieldCache.Children[0].Children[0].Children[0].Children[0].Children[0]

	require.Len(t, innerField.ValidationTag.Rules, 1)
	require.Equal(t, innerField.ValidationTag.Rules[0].Name, "integer")

	require.Len(t, innerField.ValidationTag.PresenceRules, 1)
	require.Equal(t, innerField.ValidationTag.PresenceRules[0].Name, "required")
}

func Test_it_can_analyse_nilable_structs(t *testing.T) {
	// Arrange
	validator := JsonValidator.New()
	type child struct {
		Id int `json:"id" validation:"required|integer"`
	}

	type simpleStruct struct {
		Child1 *child `json:"child1" validation:"nullable|object"`
		Child2 *child `json:"child2" validation:"nullable|object"`
	}

	// Act
	var data simpleStruct
	fieldCache, err := validator.Analyze(&data)

	// Assert
	require.NoError(t, err)

	// Root
	require.Len(t, fieldCache.Children, 2)

	require.Len(t, fieldCache.Children[0].Children, 1)
	require.Len(t, fieldCache.Children[1].Children, 1)
}

func Test_it_can_analyse_arrays_of_nilable_structs(t *testing.T) {
	// Arrange
	validator := JsonValidator.New()
	type child struct {
		Id int `json:"id" validation:"required|integer"`
	}

	type simpleStruct struct {
		Child1 []*child `json:"child1" validation:"nullable|object"`
	}

	// Act
	var data simpleStruct
	fieldCache, err := validator.Analyze(&data)

	// Assert
	require.NoError(t, err)

	// Root
	require.Len(t, fieldCache.Children, 1)
	require.Len(t, fieldCache.Children[0].Children, 1)
	require.Len(t, fieldCache.Children[0].Children[0].Children, 1)
	require.Len(t, fieldCache.Children[0].Children[0].Children[0].Children, 0)
}

func Test_it_can_analyse_nilable_arrays(t *testing.T) {
	// Arrange
	validator := JsonValidator.New()
	type child struct {
		Id int `json:"id" validation:"required|integer"`
	}

	type simpleStruct struct {
		Child1 *[]child `json:"child1" validation:"nullable|object"`
	}

	// Act
	var data simpleStruct
	fieldCache, err := validator.Analyze(&data)

	// Assert
	require.NoError(t, err)

	// Root
	require.Len(t, fieldCache.Children, 1)
	require.Len(t, fieldCache.Children[0].Children, 1)
	require.Len(t, fieldCache.Children[0].Children[0].Children, 1)
	require.Len(t, fieldCache.Children[0].Children[0].Children[0].Children, 0)
}

func Test_it_can_analyse_maps(t *testing.T) {
	// Arrange
	validator := JsonValidator.New()
	type child struct {
		Id int `json:"id" validation:"required|integer"`
	}

	type simpleStruct struct {
		Child1 map[string]child `json:"child1" validation:"nullable|object"`
	}

	// Act
	var data simpleStruct
	fieldCache, err := validator.Analyze(&data)

	// Assert
	require.NoError(t, err)

	// Root
	require.Len(t, fieldCache.Children, 1)
	require.Equal(t, "", fieldCache.StructKey)
	require.True(t, fieldCache.IsStruct)
	require.False(t, fieldCache.IsSlice)
	require.False(t, fieldCache.IsMap)

	// Map
	require.Len(t, fieldCache.Children[0].Children, 1)
	require.Equal(t, "Child1", fieldCache.Children[0].StructKey)
	require.False(t, fieldCache.Children[0].IsStruct)
	require.False(t, fieldCache.Children[0].IsSlice)
	require.True(t, fieldCache.Children[0].IsMap)

	// struct
	require.Len(t, fieldCache.Children[0].Children[0].Children, 1)
	require.Equal(t, "{index}", fieldCache.Children[0].Children[0].StructKey)
	require.True(t, fieldCache.Children[0].Children[0].IsStruct)
	require.False(t, fieldCache.Children[0].Children[0].IsSlice)
	require.False(t, fieldCache.Children[0].Children[0].IsMap)

	// struct field
	require.Len(t, fieldCache.Children[0].Children[0].Children[0].Children, 0)
	require.Equal(t, "Id", fieldCache.Children[0].Children[0].Children[0].StructKey)
	require.False(t, fieldCache.Children[0].Children[0].Children[0].IsStruct)
	require.False(t, fieldCache.Children[0].Children[0].Children[0].IsSlice)
	require.False(t, fieldCache.Children[0].Children[0].Children[0].IsMap)
}
