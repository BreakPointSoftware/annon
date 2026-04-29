package walk

import (
	"reflect"

	"github.com/BreakPointSoftware/annon/internal/decision"
	"github.com/BreakPointSoftware/annon/internal/redactcore"
	"github.com/BreakPointSoftware/annon/internal/support/reflectx"
)

func applyAction(inputValue reflect.Value, strategyName string, copyConfig decision.Config) (reflect.Value, error) {
	result, err := redactcore.Apply(strategyName, reflectx.Interface(inputValue), copyConfig.Preservation)
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
