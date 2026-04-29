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

type Walker = Copier

func New(config decision.Config, decider *decision.Decider, cache *TypeCache) *Walker {
	return &Walker{config: config, decider: decider, cache: cache}
}

func (w *Walker) Copy(input any) (any, error) {
	if input == nil {
		return nil, nil
	}

	inputValue := reflect.ValueOf(input)

	// Build a fresh typed value so redaction never mutates caller-owned data.
	copiedValue, err := w.copyValue(inputValue, "", "", true)
	if err != nil {
		return nil, err
	}

	return copiedValue.Interface(), nil
}

func (w *Walker) copyValue(inputValue reflect.Value, fieldName string, tag string, allowDecision bool) (reflect.Value, error) {
	if !inputValue.IsValid() {
		return inputValue, nil
	}

	if allowDecision {
		// Apply tag and detector decisions before traversing child values.
		fieldDecision, err := w.decider.Decide(fieldName, tag, reflectx.Interface(inputValue))
		if err != nil {
			return reflect.Value{}, err
		}

		if fieldDecision.Skip {
			return w.cloneValue(inputValue)
		}

		if fieldDecision.Remove {
			return reflect.Zero(inputValue.Type()), nil
		}

		if fieldDecision.StrategyName != "" {
			return w.applyTypedAction(inputValue, fieldDecision.StrategyName)
		}
	}

	return w.copyChildren(inputValue, fieldName, tag)
}

func (w *Walker) copyChildren(inputValue reflect.Value, fieldName string, tag string) (reflect.Value, error) {
	switch inputValue.Kind() {
	case reflect.Interface:
		return w.copyInterface(inputValue, fieldName, tag)
	case reflect.Pointer:
		return w.copyPointer(inputValue)
	case reflect.Struct:
		return w.copyStruct(inputValue)
	case reflect.Map:
		return w.copyMap(inputValue)
	case reflect.Slice:
		return w.copySlice(inputValue)
	case reflect.Array:
		return w.copyArray(inputValue)
	default:
		return inputValue, nil
	}
}

func (w *Walker) copyInterface(inputValue reflect.Value, fieldName string, tag string) (reflect.Value, error) {
	if inputValue.IsNil() {
		return reflect.Zero(inputValue.Type()), nil
	}

	copiedValue, err := w.copyValue(inputValue.Elem(), fieldName, tag, true)
	if err != nil {
		return reflect.Value{}, err
	}

	outputValue := reflect.New(inputValue.Type()).Elem()
	outputValue.Set(copiedValue)
	return outputValue, nil
}

func (w *Walker) copyPointer(inputValue reflect.Value) (reflect.Value, error) {
	if inputValue.IsNil() {
		return reflect.Zero(inputValue.Type()), nil
	}

	copiedValue, err := w.copyValue(inputValue.Elem(), "", "", true)
	if err != nil {
		return reflect.Value{}, err
	}

	outputPointer := reflect.New(inputValue.Type().Elem())
	outputPointer.Elem().Set(copiedValue)
	return outputPointer, nil
}

func (w *Walker) copyStruct(structValue reflect.Value) (reflect.Value, error) {
	copiedStruct := reflect.New(structValue.Type()).Elem()

	// Struct fields are copied one by one so tags and detection can be applied per field.
	for _, fieldMeta := range w.cache.StructFields(structValue.Type()) {
		copiedValue, err := w.copyValue(structValue.FieldByIndex(fieldMeta.Index), fieldMeta.DetectionName("typed"), fieldMeta.AnonymiseTag, true)
		if err != nil {
			return reflect.Value{}, err
		}

		copiedStruct.FieldByIndex(fieldMeta.Index).Set(copiedValue)
	}

	return copiedStruct, nil
}

func (w *Walker) copyMap(mapValue reflect.Value) (reflect.Value, error) {
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

		copiedValue, err := w.copyValue(mapIterator.Value(), fieldName, "", true)
		if err != nil {
			return reflect.Value{}, err
		}

		copiedMap.SetMapIndex(mapKey, copiedValue)
	}

	return copiedMap, nil
}

func (w *Walker) copySlice(sliceValue reflect.Value) (reflect.Value, error) {
	if sliceValue.IsNil() {
		return reflect.Zero(sliceValue.Type()), nil
	}

	copiedSlice := reflect.MakeSlice(sliceValue.Type(), sliceValue.Len(), sliceValue.Len())

	for index := 0; index < sliceValue.Len(); index++ {
		copiedValue, err := w.copyValue(sliceValue.Index(index), "", "", true)
		if err != nil {
			return reflect.Value{}, err
		}

		copiedSlice.Index(index).Set(copiedValue)
	}

	return copiedSlice, nil
}

func (w *Walker) copyArray(arrayValue reflect.Value) (reflect.Value, error) {
	copiedArray := reflect.New(arrayValue.Type()).Elem()

	for index := 0; index < arrayValue.Len(); index++ {
		copiedValue, err := w.copyValue(arrayValue.Index(index), "", "", true)
		if err != nil {
			return reflect.Value{}, err
		}

		copiedArray.Index(index).Set(copiedValue)
	}

	return copiedArray, nil
}

func (w *Walker) cloneValue(inputValue reflect.Value) (reflect.Value, error) {
	return w.copyValue(inputValue, "", "", false)
}

func (w *Walker) applyTypedAction(inputValue reflect.Value, strategyName string) (reflect.Value, error) {
	return applyAction(inputValue, strategyName, w.config)
}
