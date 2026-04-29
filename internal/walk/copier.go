package walk

import (
	"reflect"

	"github.com/BreakPointSoftware/annon/internal/decision"
	"github.com/BreakPointSoftware/annon/internal/support/reflectx"
)

type Copier struct {
	config  decision.Config
	decider *decision.Decider
	cache   *TypeCache
}

func New(config decision.Config, decider *decision.Decider, cache *TypeCache) *Copier {
	return &Copier{config: config, decider: decider, cache: cache}
}

func (c *Copier) Copy(input any) (any, error) {
	if input == nil {
		return nil, nil
	}

	inputValue := reflect.ValueOf(input)
	copiedValue, err := c.copyValue(inputValue, "", "", true)
	if err != nil {
		return nil, err
	}

	return copiedValue.Interface(), nil
}

func (c *Copier) copyValue(inputValue reflect.Value, fieldName string, tag string, allowDecision bool) (reflect.Value, error) {
	if !inputValue.IsValid() {
		return inputValue, nil
	}

	if allowDecision {
		fieldDecision, err := c.decider.Decide(fieldName, tag, reflectx.Interface(inputValue))
		if err != nil {
			return reflect.Value{}, err
		}

		if fieldDecision.Skip {
			return c.cloneValue(inputValue)
		}

		if fieldDecision.Remove {
			return reflect.Zero(inputValue.Type()), nil
		}

		if fieldDecision.StrategyName != "" {
			return c.applyTypedAction(inputValue, fieldDecision.StrategyName)
		}
	}

	return c.copyChildren(inputValue, fieldName, tag)
}

func (c *Copier) copyChildren(inputValue reflect.Value, fieldName string, tag string) (reflect.Value, error) {
	switch inputValue.Kind() {
	case reflect.Interface:
		return c.copyInterface(inputValue, fieldName, tag)
	case reflect.Pointer:
		return c.copyPointer(inputValue)
	case reflect.Struct:
		return c.copyStruct(inputValue)
	case reflect.Map:
		return c.copyMap(inputValue)
	case reflect.Slice:
		return c.copySlice(inputValue)
	case reflect.Array:
		return c.copyArray(inputValue)
	default:
		return inputValue, nil
	}
}

func (c *Copier) copyInterface(inputValue reflect.Value, fieldName string, tag string) (reflect.Value, error) {
	if inputValue.IsNil() {
		return reflect.Zero(inputValue.Type()), nil
	}

	copiedValue, err := c.copyValue(inputValue.Elem(), fieldName, tag, true)
	if err != nil {
		return reflect.Value{}, err
	}

	outputValue := reflect.New(inputValue.Type()).Elem()
	outputValue.Set(copiedValue)
	return outputValue, nil
}

func (c *Copier) copyPointer(inputValue reflect.Value) (reflect.Value, error) {
	if inputValue.IsNil() {
		return reflect.Zero(inputValue.Type()), nil
	}

	copiedValue, err := c.copyValue(inputValue.Elem(), "", "", true)
	if err != nil {
		return reflect.Value{}, err
	}

	outputPointer := reflect.New(inputValue.Type().Elem())
	outputPointer.Elem().Set(copiedValue)
	return outputPointer, nil
}

func (c *Copier) copyStruct(structValue reflect.Value) (reflect.Value, error) {
	copiedStruct := reflect.New(structValue.Type()).Elem()

	for _, fieldMeta := range c.cache.StructFields(structValue.Type()) {
		copiedValue, err := c.copyValue(structValue.FieldByIndex(fieldMeta.Index), fieldMeta.DetectionName("typed"), fieldMeta.AnonymiseTag, true)
		if err != nil {
			return reflect.Value{}, err
		}

		copiedStruct.FieldByIndex(fieldMeta.Index).Set(copiedValue)
	}

	return copiedStruct, nil
}

func (c *Copier) copyMap(mapValue reflect.Value) (reflect.Value, error) {
	if mapValue.IsNil() {
		return reflect.Zero(mapValue.Type()), nil
	}

	copiedMap := reflect.MakeMapWithSize(mapValue.Type(), mapValue.Len())
	mapIterator := mapValue.MapRange()

	for mapIterator.Next() {
		mapKey := mapIterator.Key()
		fieldName := ""
		if mapKey.Kind() == reflect.String {
			fieldName = mapKey.String()
		}

		copiedValue, err := c.copyValue(mapIterator.Value(), fieldName, "", true)
		if err != nil {
			return reflect.Value{}, err
		}

		copiedMap.SetMapIndex(mapKey, copiedValue)
	}

	return copiedMap, nil
}

func (c *Copier) copySlice(sliceValue reflect.Value) (reflect.Value, error) {
	if sliceValue.IsNil() {
		return reflect.Zero(sliceValue.Type()), nil
	}

	copiedSlice := reflect.MakeSlice(sliceValue.Type(), sliceValue.Len(), sliceValue.Len())
	for index := 0; index < sliceValue.Len(); index++ {
		copiedValue, err := c.copyValue(sliceValue.Index(index), "", "", true)
		if err != nil {
			return reflect.Value{}, err
		}

		copiedSlice.Index(index).Set(copiedValue)
	}

	return copiedSlice, nil
}

func (c *Copier) copyArray(arrayValue reflect.Value) (reflect.Value, error) {
	copiedArray := reflect.New(arrayValue.Type()).Elem()
	for index := 0; index < arrayValue.Len(); index++ {
		copiedValue, err := c.copyValue(arrayValue.Index(index), "", "", true)
		if err != nil {
			return reflect.Value{}, err
		}

		copiedArray.Index(index).Set(copiedValue)
	}

	return copiedArray, nil
}

func (c *Copier) cloneValue(inputValue reflect.Value) (reflect.Value, error) {
	return c.copyValue(inputValue, "", "", false)
}

func (c *Copier) applyTypedAction(inputValue reflect.Value, strategyName string) (reflect.Value, error) {
	return applyAction(inputValue, strategyName, c.config)
}
