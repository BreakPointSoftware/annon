package walk

import (
	"fmt"
	"reflect"

	"github.com/BreakPointSoftware/annon/internal/redactcore"
)

var omittedValue = &struct{}{}

func (w *Walker) OutputFromValue(input any, format string) (any, error) {
	return w.output.OutputFromValue(input, format)
}

func (w *Walker) OutputFromNeutral(input any) (any, error) {
	return w.output.OutputFromNeutral(input)
}

type OutputBuilder struct {
	config  Config
	decider *Decider
	cache   *TypeCache
}

func NewOutputBuilder(config Config, decider *Decider, cache *TypeCache) *OutputBuilder {
	return &OutputBuilder{config: config, decider: decider, cache: cache}
}

func (b *OutputBuilder) OutputFromValue(input any, format string) (any, error) {
	if input == nil {
		return nil, nil
	}

	return b.buildOutputFromReflect(reflect.ValueOf(input), format, "", "", true)
}

func (b *OutputBuilder) OutputFromNeutral(input any) (any, error) {
	return b.buildOutputFromNeutralValue(input, "", "")
}

func (b *OutputBuilder) buildOutputFromReflect(inputValue reflect.Value, format, fieldName, tag string, allowDecision bool) (any, error) {
	if !inputValue.IsValid() {
		return nil, nil
	}

	if allowDecision {
		fieldDecision, err := b.decider.Decide(fieldName, tag, valueInterface(inputValue))
		if err != nil {
			return nil, err
		}

		if fieldDecision.skip {
			return b.buildPlainOutputFromReflect(inputValue, format)
		}

		if fieldDecision.remove {
			return omittedValue, nil
		}

		if fieldDecision.strategyName != "" {
			transformedValue, err := b.applyOutputAction(inputValue, fieldDecision.strategyName)
			if err != nil {
				return nil, err
			}

			return valueInterface(transformedValue), nil
		}
	}

	return b.buildReflectOutput(inputValue, format, fieldName, tag)
}

func (b *OutputBuilder) buildReflectOutput(inputValue reflect.Value, format, fieldName, tag string) (any, error) {
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
		return valueInterface(inputValue), nil
	}
}

func (b *OutputBuilder) buildOutputInterface(inputValue reflect.Value, format, fieldName, tag string) (any, error) {
	if inputValue.IsNil() {
		return nil, nil
	}

	return b.buildOutputFromReflect(inputValue.Elem(), format, fieldName, tag, true)
}

func (b *OutputBuilder) buildOutputPointer(inputValue reflect.Value, format, fieldName, tag string) (any, error) {
	if inputValue.IsNil() {
		return nil, nil
	}

	return b.buildOutputFromReflect(inputValue.Elem(), format, fieldName, tag, true)
}

func (b *OutputBuilder) buildOutputStructMap(structValue reflect.Value, format string) (any, error) {
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

func (b *OutputBuilder) buildOutputMap(mapValue reflect.Value, format string) (any, error) {
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

func (b *OutputBuilder) buildOutputSlice(sliceValue reflect.Value, format string) (any, error) {
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

func (b *OutputBuilder) buildPlainOutputFromReflect(inputValue reflect.Value, format string) (any, error) {
	return b.buildOutputFromReflect(inputValue, format, "", "", false)
}

func (b *OutputBuilder) buildOutputFromNeutralValue(input any, fieldName, tag string) (any, error) {
	fieldDecision, err := b.decider.Decide(fieldName, tag, input)
	if err != nil {
		return nil, err
	}

	if fieldDecision.skip {
		return cloneNeutral(input), nil
	}

	if fieldDecision.remove {
		return omittedValue, nil
	}

	if fieldDecision.strategyName != "" {
		return redactcore.Apply(fieldDecision.strategyName, input, b.config.Preservation)
	}

	return b.buildNeutralOutput(input)
}

func (b *OutputBuilder) buildNeutralOutput(input any) (any, error) {
	switch neutralValue := input.(type) {
	case map[string]any:
		return b.buildNeutralOutputMap(neutralValue)
	case []any:
		return b.buildNeutralOutputSlice(neutralValue)
	default:
		return input, nil
	}
}

func (b *OutputBuilder) buildNeutralOutputMap(neutralMap map[string]any) (any, error) {
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

func (b *OutputBuilder) buildNeutralOutputSlice(neutralSlice []any) (any, error) {
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

func (b *OutputBuilder) applyOutputAction(inputValue reflect.Value, strategyName string) (reflect.Value, error) {
	return applyAction(inputValue, strategyName, b.config)
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
