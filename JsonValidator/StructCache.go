package JsonValidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type FieldCache struct {
	Parent        *FieldCache
	Children      []*FieldCache
	Reflection    reflect.Type
	JsonKey       string
	StructKey     string
	ValidationTag *ValidationTag
	IsStruct      bool
	IsSlice       bool
}

type StructCache struct {
	Cache map[reflect.Type]*FieldCache
}

func newStructCache() *StructCache {
	return &StructCache{Cache: map[reflect.Type]*FieldCache{}}
}

func (structCache *StructCache) analyze(rulebook *Rulebook, targetType reflect.Type) (*FieldCache, error) {
	// Unwrap pointer types, since we only focus on the underlying struct type
	targetType = structCache.typeIndirect(targetType)

	// If the type has already been analyzed, then fetch from the cache
	if cache, present := structCache.Cache[targetType]; present {
		return cache, nil
	}

	// We only allow analyzing structs as root data types
	if targetType.Kind() != reflect.Struct {
		return nil, errors.New(fmt.Sprintf("the struct cache can only analyze struct types %s given", targetType.Kind().String()))
	}

	root := &FieldCache{
		Parent:     nil,
		Children:   []*FieldCache{},
		Reflection: targetType,
		JsonKey:    "",
		StructKey:  "",
		ValidationTag: &ValidationTag{
			Rules:              []*RuleContext{},
			PresenceRules:      []*RuleContext{},
			ExplicitlyNullable: false,
		},
		IsStruct: true,
		IsSlice:  false,
	}

	result := structCache.traverseType(root, rulebook)
	structCache.Cache[targetType] = result

	return result, nil
}

func (structCache *StructCache) typeIndirect(targetType reflect.Type) reflect.Type {
	if targetType.Kind() == reflect.Pointer {
		return targetType.Elem()
	}

	return targetType
}

func (structCache *StructCache) traverseType(parent *FieldCache, rulebook *Rulebook) *FieldCache {
	if parent.IsStruct {
		structCache.traverseStruct(parent, rulebook)
	} else if parent.IsSlice {
		structCache.traverseSlice(parent, rulebook)
	}

	return parent
}

func (structCache *StructCache) traverseStruct(parent *FieldCache, rulebook *Rulebook) *FieldCache {
	numFields := parent.Reflection.NumField()

	for i := 0; i < numFields; i++ {
		structField := parent.Reflection.Field(i)
		structType := structField.Type

		field := &FieldCache{
			Parent:        parent,
			Children:      []*FieldCache{},
			Reflection:    structType,
			JsonKey:       structCache.getJsonTagForStructField(structField).JsonKey,
			StructKey:     structField.Name,
			ValidationTag: structCache.getValidationTag(structField, rulebook),
			IsStruct:      structCache.typeIsStruct(structType),
			IsSlice:       structCache.typeIsSlice(structType),
		}

		parent.Children = append(parent.Children, field)
		structCache.traverseType(field, rulebook)
	}

	return parent
}

func (structCache *StructCache) typeIsStruct(reflectType reflect.Type) bool {
	return reflectType.Kind() == reflect.Struct
}

func (structCache *StructCache) typeIsSlice(reflectType reflect.Type) bool {
	kind := reflectType.Kind()

	return kind == reflect.Slice || kind == reflect.Array
}

func (structCache *StructCache) getJsonTagForStructField(field reflect.StructField) *JsonTag {
	tagline, ok := field.Tag.Lookup("json")

	if !ok {
		return &JsonTag{JsonKey: field.Name}
	}

	return &JsonTag{JsonKey: strings.Split(tagline, ",")[0]}
}

func (structCache *StructCache) getValidationTag(field reflect.StructField, rulebook *Rulebook) *ValidationTag {
	tagline, ok := field.Tag.Lookup("validation")

	if !ok {
		return newValidationTag(rulebook, "")
	}

	return newValidationTag(rulebook, tagline)
}

func (structCache *StructCache) traverseSlice(parent *FieldCache, rulebook *Rulebook) *FieldCache {
	sliceElem := structCache.typeIndirect(parent.Reflection)
	sliceSubtype := sliceElem.Elem()

	field := &FieldCache{
		Parent:        parent,
		Children:      []*FieldCache{},
		Reflection:    sliceSubtype,
		JsonKey:       "{index}",
		StructKey:     "{index}",
		ValidationTag: newValidationTag(rulebook, ""),
		IsStruct:      structCache.typeIsStruct(sliceSubtype),
		IsSlice:       structCache.typeIsSlice(sliceSubtype),
	}

	parent.Children = append(parent.Children, field)
	structCache.traverseType(field, rulebook)

	return parent
}
