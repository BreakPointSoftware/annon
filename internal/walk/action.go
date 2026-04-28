package walk

import (
	"reflect"

	"github.com/BreakPointSoftware/annon/internal/redactcore"
)

func applyAction(v reflect.Value, strategyName string, cfg Config) (reflect.Value, error) {
	result, err := redactcore.Apply(strategyName, valueInterface(v), cfg.Preservation)
	if err != nil {
		return reflect.Value{}, err
	}
	if result == nil {
		return reflect.Zero(v.Type()), nil
	}
	resultValue := reflect.ValueOf(result)
	if resultValue.Type().AssignableTo(v.Type()) {
		return resultValue, nil
	}
	if resultValue.Type().ConvertibleTo(v.Type()) {
		return resultValue.Convert(v.Type()), nil
	}
	if v.Kind() == reflect.Interface {
		out := reflect.New(v.Type()).Elem()
		out.Set(resultValue)
		return out, nil
	}
	return v, nil
}
