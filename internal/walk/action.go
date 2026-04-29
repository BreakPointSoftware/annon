package walk

import (
	"reflect"

	"github.com/BreakPointSoftware/annon/internal/decision"
	"github.com/BreakPointSoftware/annon/internal/redactcore"
)

func applyAction(inputValue reflect.Value, strategyName string, copyConfig decision.Config) (reflect.Value, error) {
	result, err := redactcore.Apply(strategyName, valueInterface(inputValue), copyConfig.Preservation)
	if err != nil {
		return reflect.Value{}, err
	}

	if result == nil {
		return reflect.Zero(inputValue.Type()), nil
	}

	resultValue := reflect.ValueOf(result)
	if resultValue.Type().AssignableTo(inputValue.Type()) {
		return resultValue, nil
	}

	if resultValue.Type().ConvertibleTo(inputValue.Type()) {
		return resultValue.Convert(inputValue.Type()), nil
	}

	if inputValue.Kind() == reflect.Interface {
		outputValue := reflect.New(inputValue.Type()).Elem()
		outputValue.Set(resultValue)
		return outputValue, nil
	}

	return inputValue, nil
}

func valueInterface(inputValue reflect.Value) any {
	if !inputValue.IsValid() {
		return nil
	}

	return inputValue.Interface()
}
