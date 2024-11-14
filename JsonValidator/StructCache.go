package JsonValidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
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
	IsMap         bool
}

type StructCache struct {
	Cache     map[reflect.Type]*FieldCache
	rootLock  *sync.Mutex
	typeLocks map[reflect.Type]*sync.Mutex
}

func newStructCache() *StructCache {
	return &StructCache{Cache: map[reflect.Type]*FieldCache{}, rootLock: new(sync.Mutex)}
}

func (fieldCache *FieldCache) GetChildByName(name string) *FieldCache {
	for _, child := range fieldCache.Children {
		if child.StructKey == name {
			return child
		}
	}

	return nil
}

func (structCache *StructCache) Analyze(rulebook *Rulebook, targetType reflect.Type) (*FieldCache, error) {
	// Unwrap pointer types, since we only focus on the underlying struct type
	targetType = structCache.typeIndirect(targetType)

	// If the type has already been analyzed, then fetch from the cache
	// This is the code path for 99.9999% of requests.
	if cache, present := structCache.Cache[targetType]; present {
		return cache, nil
	}

	// Otherwise, it is the first time we see this struct, and therefor has to perform the actual analysis
	lock := structCache.acquireTypeLock(targetType)
	defer lock.Unlock()

	// There might have been another concurrent analyze call to the struct cache for the same time,
	// while we were waiting for the type lock. In this case we can skip the additional analysis and simply use the cache
	if cache, present := structCache.Cache[targetType]; present {
		return cache, nil
	}

	// We only allow analyzing structs as root data types
	if targetType.Kind() != reflect.Struct {
		return nil, errors.New(fmt.Sprintf("the struct cache can only Analyze struct types %s given", targetType.Kind().String()))
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
		IsMap:    false,
	}

	structCache.traverseType(root, rulebook)
	structCache.Cache[targetType] = root

	return root, nil
}

func (structCache *StructCache) acquireTypeLock(targetType reflect.Type) *sync.Mutex {
	// We first need the root lock, so we can get or create the required type lock without conflicts
	structCache.rootLock.Lock()
	defer structCache.rootLock.Unlock()

	if typeLock, present := structCache.typeLocks[targetType]; present {
		typeLock.Lock()
		return typeLock
	}

	typeLock := new(sync.Mutex)
	typeLock.Lock()

	return typeLock
}

func (structCache *StructCache) typeIndirect(targetType reflect.Type) reflect.Type {
	if targetType.Kind() == reflect.Pointer {
		return targetType.Elem()
	}

	return targetType
}

func (structCache *StructCache) traverseType(parent *FieldCache, rulebook *Rulebook) {
	if parent.IsStruct {
		structCache.traverseStruct(parent, rulebook)
	} else if parent.IsSlice {
		structCache.traverseSlice(parent, rulebook)
	} else if parent.IsMap {
		structCache.traverseMap(parent, rulebook)
	}
}

func (structCache *StructCache) traverseStruct(parent *FieldCache, rulebook *Rulebook) {
	numFields := parent.Reflection.NumField()

	for i := 0; i < numFields; i++ {
		structField := parent.Reflection.Field(i)
		structType := structCache.typeIndirect(structField.Type)

		field := &FieldCache{
			Parent:        parent,
			Children:      []*FieldCache{},
			Reflection:    structType,
			JsonKey:       structCache.getJsonTagForStructField(structField).JsonKey,
			StructKey:     structField.Name,
			ValidationTag: structCache.getValidationTag(structField, rulebook),
			IsStruct:      structCache.typeIsStruct(structType),
			IsSlice:       structCache.typeIsSlice(structType),
			IsMap:         structCache.typeIsMap(structType),
		}

		structCache.traverseType(field, rulebook)

		if structField.Anonymous { // Embedded fields are a part of the parent type itself
			parent.Children = append(parent.Children, field.Children...)
		} else {
			parent.Children = append(parent.Children, field)
		}
	}
}

func (structCache *StructCache) traverseMap(parent *FieldCache, rulebook *Rulebook) {
	sliceElem := structCache.typeIndirect(parent.Reflection)
	mapSubType := structCache.typeIndirect(sliceElem.Elem())

	field := &FieldCache{
		Parent:        parent,
		Children:      []*FieldCache{},
		Reflection:    mapSubType,
		JsonKey:       "{index}",
		StructKey:     "{index}",
		ValidationTag: newValidationTag(rulebook, ""),
		IsStruct:      structCache.typeIsStruct(mapSubType),
		IsSlice:       structCache.typeIsSlice(mapSubType),
		IsMap:         structCache.typeIsMap(mapSubType),
	}

	parent.Children = append(parent.Children, field)
	structCache.traverseType(field, rulebook)
}

func (structCache *StructCache) typeIsStruct(reflectType reflect.Type) bool {
	return reflectType.Kind() == reflect.Struct
}

func (structCache *StructCache) typeIsSlice(reflectType reflect.Type) bool {
	kind := reflectType.Kind()

	return kind == reflect.Slice || kind == reflect.Array
}

func (structCache *StructCache) typeIsMap(reflectType reflect.Type) bool {
	kind := reflectType.Kind()

	return kind == reflect.Map
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

func (structCache *StructCache) traverseSlice(parent *FieldCache, rulebook *Rulebook) {
	sliceElem := structCache.typeIndirect(parent.Reflection)
	sliceSubtype := structCache.typeIndirect(sliceElem.Elem())

	field := &FieldCache{
		Parent:        parent,
		Children:      []*FieldCache{},
		Reflection:    sliceSubtype,
		JsonKey:       "{index}",
		StructKey:     "{index}",
		ValidationTag: newValidationTag(rulebook, ""),
		IsStruct:      structCache.typeIsStruct(sliceSubtype),
		IsSlice:       structCache.typeIsSlice(sliceSubtype),
		IsMap:         structCache.typeIsMap(sliceSubtype),
	}

	parent.Children = append(parent.Children, field)
	structCache.traverseType(field, rulebook)
}
