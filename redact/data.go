package redact

import (
	"reflect"

	annon "github.com/BreakPointSoftware/annon"
)

func Data(input any) any {
	defer func() {
		if recover() != nil {
		}
	}()

	if stringInput, isString := input.(string); isString {
		return String(stringInput)
	}

	redactedValue, err := annon.New()
	if err != nil {
		return fallbackValue(input)
	}

	result, err := redactedValue.Copy(input)
	if err != nil {
		return fallbackValue(input)
	}

	return result
}

func fallbackValue(input any) any {
	if input == nil {
		return nil
	}

	inputValue := reflect.ValueOf(input)
	if !inputValue.IsValid() {
		return nil
	}

	switch inputValue.Kind() {
	case reflect.String:
		return Redact(inputValue.String())
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		return input
	case reflect.Pointer, reflect.Interface, reflect.Map, reflect.Slice:
		return reflect.Zero(inputValue.Type()).Interface()
	case reflect.Array, reflect.Struct:
		return reflect.New(inputValue.Type()).Elem().Interface()
	default:
		return "[REDACTED]"
	}
}
