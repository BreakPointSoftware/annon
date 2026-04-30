package baseline

import (
	"errors"
	"reflect"
)

var ErrCycleUnsupported = errors.New("baseline copy does not support cycles")

type Walker struct {
	activePointers map[visitKey]bool
}

type visitKey struct {
	typ reflect.Type
	ptr uintptr
}

func New() *Walker {
	return &Walker{activePointers: map[visitKey]bool{}}
}

func Copy[T any](input T) (T, error) {
	walker := New()
	inputValue := reflect.ValueOf(input)
	copiedValue, err := walker.copyValue(inputValue)
	if err != nil {
		var zero T
		return zero, err
	}
	if !copiedValue.IsValid() {
		var zero T
		return zero, nil
	}
	return copiedValue.Interface().(T), nil
}

func (w *Walker) copyValue(inputValue reflect.Value) (reflect.Value, error) {
	if !inputValue.IsValid() {
		return inputValue, nil
	}

	switch inputValue.Kind() {
	case reflect.Interface:
		if inputValue.IsNil() {
			return reflect.Zero(inputValue.Type()), nil
		}
		copiedValue, err := w.copyValue(inputValue.Elem())
		if err != nil {
			return reflect.Value{}, err
		}
		outputValue := reflect.New(inputValue.Type()).Elem()
		outputValue.Set(copiedValue)
		return outputValue, nil
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
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return reflect.Zero(inputValue.Type()), nil
	default:
		return inputValue, nil
	}
}

func (w *Walker) copyPointer(inputValue reflect.Value) (reflect.Value, error) {
	if inputValue.IsNil() {
		return reflect.Zero(inputValue.Type()), nil
	}
	key := visitKey{typ: inputValue.Type(), ptr: inputValue.Pointer()}
	if w.activePointers[key] {
		return reflect.Value{}, ErrCycleUnsupported
	}
	w.activePointers[key] = true
	defer delete(w.activePointers, key)

	copiedValue, err := w.copyValue(inputValue.Elem())
	if err != nil {
		return reflect.Value{}, err
	}
	outputPointer := reflect.New(inputValue.Type().Elem())
	outputPointer.Elem().Set(copiedValue)
	return outputPointer, nil
}

func (w *Walker) copyStruct(structValue reflect.Value) (reflect.Value, error) {
	outputStruct := reflect.New(structValue.Type()).Elem()
	for fieldIndex := 0; fieldIndex < structValue.NumField(); fieldIndex++ {
		structField := structValue.Type().Field(fieldIndex)
		if structField.PkgPath != "" {
			continue
		}
		copiedValue, err := w.copyValue(structValue.Field(fieldIndex))
		if err != nil {
			return reflect.Value{}, err
		}
		outputStruct.Field(fieldIndex).Set(copiedValue)
	}
	return outputStruct, nil
}

func (w *Walker) copyMap(mapValue reflect.Value) (reflect.Value, error) {
	if mapValue.IsNil() {
		return reflect.Zero(mapValue.Type()), nil
	}
	outputMap := reflect.MakeMapWithSize(mapValue.Type(), mapValue.Len())
	iter := mapValue.MapRange()
	for iter.Next() {
		copiedValue, err := w.copyValue(iter.Value())
		if err != nil {
			return reflect.Value{}, err
		}
		outputMap.SetMapIndex(iter.Key(), copiedValue)
	}
	return outputMap, nil
}

func (w *Walker) copySlice(sliceValue reflect.Value) (reflect.Value, error) {
	if sliceValue.IsNil() {
		return reflect.Zero(sliceValue.Type()), nil
	}
	outputSlice := reflect.MakeSlice(sliceValue.Type(), sliceValue.Len(), sliceValue.Len())
	for index := 0; index < sliceValue.Len(); index++ {
		copiedValue, err := w.copyValue(sliceValue.Index(index))
		if err != nil {
			return reflect.Value{}, err
		}
		outputSlice.Index(index).Set(copiedValue)
	}
	return outputSlice, nil
}

func (w *Walker) copyArray(arrayValue reflect.Value) (reflect.Value, error) {
	outputArray := reflect.New(arrayValue.Type()).Elem()
	for index := 0; index < arrayValue.Len(); index++ {
		copiedValue, err := w.copyValue(arrayValue.Index(index))
		if err != nil {
			return reflect.Value{}, err
		}
		outputArray.Index(index).Set(copiedValue)
	}
	return outputArray, nil
}
