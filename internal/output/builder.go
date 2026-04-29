package output

import (
	"fmt"
	"reflect"

	"github.com/BreakPointSoftware/annon/internal/decision"
	"github.com/BreakPointSoftware/annon/internal/redactcore"
	"github.com/BreakPointSoftware/annon/internal/support/reflectx"
	"github.com/BreakPointSoftware/annon/internal/walk"
)

var omittedValue = &struct{}{}

type Builder struct {
	config  decision.Config
	decider *decision.Decider
	cache   *walk.TypeCache
}

func New(config decision.Config, decider *decision.Decider, cache *walk.TypeCache) *Builder {
	return &Builder{config: config, decider: decider, cache: cache}
}

func (b *Builder) OutputFromValue(input any, format Format) (any, error) {
	if input == nil {
		return nil, nil
	}

	// Build output from typed input so the encoder always sees a neutral shape.
	return b.buildOutputFromReflect(reflect.ValueOf(input), format.String(), "", "", true)
}

func (b *Builder) OutputFromNeutral(input any) (any, error) {
	return b.buildOutputFromNeutralValue(input, "", "")
}

func (b *Builder) buildOutputFromReflect(inputValue reflect.Value, format, fieldName, tag string, allowDecision bool) (any, error) {
	if !inputValue.IsValid() {
		return nil, nil
	}

	if allowDecision {
		// Apply tag and detector decisions before traversing child values.
		fieldDecision, err := b.decider.Decide(fieldName, tag, reflectx.Interface(inputValue))
		if err != nil {
			return nil, err
		}

		if fieldDecision.Skip {
			return b.buildPlainOutputFromReflect(inputValue, format)
		}

		if fieldDecision.Remove {
			return omittedValue, nil
		}

		if fieldDecision.StrategyName != "" {
			transformedValue, err := b.applyOutputAction(inputValue, fieldDecision.StrategyName)
			if err != nil {
				return nil, err
			}
			return reflectx.Interface(transformedValue), nil
		}
	}

	return b.buildReflectOutput(inputValue, format, fieldName, tag)
}

func (b *Builder) buildReflectOutput(inputValue reflect.Value, format, fieldName, tag string) (any, error) {
	switch inputValue.Kind() {
	case reflect.Interface:
		return b.buildOutputInterface(inputValue, format, fieldName, tag)
	case reflect.Pointer:
		return b.buildOutputPointer(inputValue, format, fieldName, tag)
	case reflect.Struct:
		return b.buildOutputStructMap(inputValue, format)
	case reflect.Map:
		return b.buildOutputMap(inputValue, format)
	case reflect.Slice, reflect.Array:
		return b.buildOutputSlice(inputValue, format)
	default:
		return reflectx.Interface(inputValue), nil
	}
}

func (b *Builder) buildOutputInterface(inputValue reflect.Value, format, fieldName, tag string) (any, error) {
	if inputValue.IsNil() {
		return nil, nil
	}
	return b.buildOutputFromReflect(inputValue.Elem(), format, fieldName, tag, true)
}

func (b *Builder) buildOutputPointer(inputValue reflect.Value, format, fieldName, tag string) (any, error) {
	if inputValue.IsNil() {
		return nil, nil
	}
	return b.buildOutputFromReflect(inputValue.Elem(), format, fieldName, tag, true)
}

func (b *Builder) buildOutputStructMap(structValue reflect.Value, format string) (any, error) {
	outputMap := map[string]any{}

	for _, fieldMeta := range b.cache.StructFields(structValue.Type()) {
		outputFieldName := fieldMeta.OutputName(format)
		if outputFieldName == "" {
			continue
		}
		outputValue, err := b.buildOutputFromReflect(structValue.FieldByIndex(fieldMeta.Index), format, fieldMeta.DetectionName(format), fieldMeta.AnonymiseTag, true)
		if err != nil {
			return nil, err
		}
		if outputValue == omittedValue {
			continue
		}

		outputMap[outputFieldName] = outputValue
	}

	return outputMap, nil
}

func (b *Builder) buildOutputMap(mapValue reflect.Value, format string) (any, error) {
	if mapValue.IsNil() {
		return nil, nil
	}
	outputMap := map[string]any{}
	mapIterator := mapValue.MapRange()

	for mapIterator.Next() {
		outputFieldName := fmt.Sprint(mapIterator.Key().Interface())
		outputValue, err := b.buildOutputFromReflect(mapIterator.Value(), format, outputFieldName, "", true)
		if err != nil {
			return nil, err
		}
		if outputValue == omittedValue {
			continue
		}

		outputMap[outputFieldName] = outputValue
	}

	return outputMap, nil
}

func (b *Builder) buildOutputSlice(sliceValue reflect.Value, format string) (any, error) {
	outputValues := make([]any, 0, sliceValue.Len())

	for index := 0; index < sliceValue.Len(); index++ {
		outputValue, err := b.buildOutputFromReflect(sliceValue.Index(index), format, "", "", true)
		if err != nil {
			return nil, err
		}
		outputValues = append(outputValues, outputValue)
	}

	return outputValues, nil
}

func (b *Builder) buildPlainOutputFromReflect(inputValue reflect.Value, format string) (any, error) {
	return b.buildOutputFromReflect(inputValue, format, "", "", false)
}

func (b *Builder) buildOutputFromNeutralValue(input any, fieldName, tag string) (any, error) {
	fieldDecision, err := b.decider.Decide(fieldName, tag, input)
	if err != nil {
		return nil, err
	}

	if fieldDecision.Skip {
		return cloneNeutral(input), nil
	}

	if fieldDecision.Remove {
		return omittedValue, nil
	}

	if fieldDecision.StrategyName != "" {
		return redactcore.Apply(fieldDecision.StrategyName, input, b.config.Preservation)
	}

	return b.buildNeutralOutput(input)
}

func (b *Builder) buildNeutralOutput(input any) (any, error) {
	switch neutralValue := input.(type) {
	case map[string]any:
		return b.buildNeutralOutputMap(neutralValue)
	case []any:
		return b.buildNeutralOutputSlice(neutralValue)
	default:
		return input, nil
	}
}

func (b *Builder) buildNeutralOutputMap(neutralMap map[string]any) (any, error) {
	outputMap := make(map[string]any, len(neutralMap))

	for fieldName, fieldValue := range neutralMap {
		resolvedValue, err := b.buildOutputFromNeutralValue(fieldValue, fieldName, "")
		if err != nil {
			return nil, err
		}
		if resolvedValue == omittedValue {
			continue
		}

		outputMap[fieldName] = resolvedValue
	}

	return outputMap, nil
}

func (b *Builder) buildNeutralOutputSlice(neutralSlice []any) (any, error) {
	outputValues := make([]any, 0, len(neutralSlice))

	for _, elementValue := range neutralSlice {
		resolvedValue, err := b.buildOutputFromNeutralValue(elementValue, "", "")
		if err != nil {
			return nil, err
		}
		outputValues = append(outputValues, resolvedValue)
	}

	return outputValues, nil
}

func (b *Builder) applyOutputAction(inputValue reflect.Value, strategyName string) (reflect.Value, error) {
	return applyAction(inputValue, strategyName, b.config)
}

func applyAction(inputValue reflect.Value, strategyName string, outputConfig decision.Config) (reflect.Value, error) {
	result, err := redactcore.Apply(strategyName, reflectx.Interface(inputValue), outputConfig.Preservation)
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

func cloneNeutral(input any) any {
	switch value := input.(type) {
	case map[string]any:
		outputMap := make(map[string]any, len(value))

		for key, item := range value {
			outputMap[key] = cloneNeutral(item)
		}

		return outputMap
	case []any:
		outputValues := make([]any, len(value))

		for index, item := range value {
			outputValues[index] = cloneNeutral(item)
		}

		return outputValues
	default:
		return value
	}
}
