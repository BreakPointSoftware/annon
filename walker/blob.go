package walker

import (
	"fmt"
	"reflect"

	strategypkg "github.com/BreakPointSoftware/annon/strategy"
)

var omitValue = &struct{}{}

func (w *Walker) BlobFromValue(input any, format string) (any, error) {
	if input == nil {
		return nil, nil
	}
	return w.blobFromReflect(reflect.ValueOf(input), format, "", "", true)
}

func (w *Walker) BlobFromNeutral(input any) (any, error) {
	return w.blobFromNeutralValue(input, "", "")
}

func (w *Walker) blobFromReflect(v reflect.Value, format, fieldName, tag string, allowDecision bool) (any, error) {
	if !v.IsValid() {
		return nil, nil
	}
	if allowDecision {
		dec, err := w.decide(fieldName, tag, valueInterface(v))
		if err != nil {
			return nil, err
		}
		if dec.skip {
			return w.plainFromReflect(v, format)
		}
		if dec.remove {
			return omitValue, nil
		}
		if dec.strategyName != "" {
			strategyValue, err := w.applyStrategy(v, dec.strategyName)
			if err != nil {
				return nil, err
			}
			return valueInterface(strategyValue), nil
		}
	}

	switch v.Kind() {
	case reflect.Interface:
		if v.IsNil() {
			return nil, nil
		}
		return w.blobFromReflect(v.Elem(), format, fieldName, tag, true)
	case reflect.Pointer:
		if v.IsNil() {
			return nil, nil
		}
		return w.blobFromReflect(v.Elem(), format, fieldName, tag, true)
	case reflect.Struct:
		out := map[string]any{}
		for _, meta := range w.cache.StructFields(v.Type()) {
			name := meta.OutputName(format)
			if name == "" {
				continue
			}
			value, err := w.blobFromReflect(v.FieldByIndex(meta.Index), format, meta.DetectionName(format), meta.AnonymiseTag, true)
			if err != nil {
				return nil, err
			}
			if value == omitValue {
				continue
			}
			out[name] = value
		}
		return out, nil
	case reflect.Map:
		if v.IsNil() {
			return nil, nil
		}
		out := map[string]any{}
		iter := v.MapRange()
		for iter.Next() {
			key := iter.Key()
			name := fmt.Sprint(key.Interface())
			value, err := w.blobFromReflect(iter.Value(), format, name, "", true)
			if err != nil {
				return nil, err
			}
			if value == omitValue {
				continue
			}
			out[name] = value
		}
		return out, nil
	case reflect.Slice, reflect.Array:
		out := make([]any, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			value, err := w.blobFromReflect(v.Index(i), format, "", "", true)
			if err != nil {
				return nil, err
			}
			out = append(out, value)
		}
		return out, nil
	default:
		return valueInterface(v), nil
	}
}

func (w *Walker) plainFromReflect(v reflect.Value, format string) (any, error) {
	return w.blobFromReflect(v, format, "", "", false)
}

func (w *Walker) blobFromNeutralValue(input any, fieldName, tag string) (any, error) {
	dec, err := w.decide(fieldName, tag, input)
	if err != nil {
		return nil, err
	}
	if dec.skip {
		return cloneNeutral(input), nil
	}
	if dec.remove {
		return omitValue, nil
	}
	if dec.strategyName != "" {
		value, err := w.applyNeutralStrategy(input, dec.strategyName)
		if err != nil {
			return nil, err
		}
		return value, nil
	}
	switch value := input.(type) {
	case map[string]any:
		out := make(map[string]any, len(value))
		for key, item := range value {
			resolved, err := w.blobFromNeutralValue(item, key, "")
			if err != nil {
				return nil, err
			}
			if resolved == omitValue {
				continue
			}
			out[key] = resolved
		}
		return out, nil
	case []any:
		out := make([]any, 0, len(value))
		for _, item := range value {
			resolved, err := w.blobFromNeutralValue(item, "", "")
			if err != nil {
				return nil, err
			}
			out = append(out, resolved)
		}
		return out, nil
	default:
		return input, nil
	}
}

func (w *Walker) applyNeutralStrategy(value any, strategyName string) (any, error) {
	strategyImpl, ok := w.cfg.Strategies[strategyName]
	if !ok {
		return nil, fmt.Errorf("unknown strategy %q", strategyName)
	}
	return strategyImpl.Anonymise(value, strategypkg.Context{Preservation: w.cfg.Preservation})
}

func (w *Walker) validateTag(parsed parsedTag) error {
	if parsed.empty || parsed.skip || parsed.auto || parsed.remove {
		return nil
	}
	if _, ok := w.cfg.Strategies[parsed.strategyName]; !ok {
		return fmt.Errorf("unknown anonymise tag strategy %q", parsed.strategyName)
	}
	return nil
}

func cloneNeutral(input any) any {
	switch value := input.(type) {
	case map[string]any:
		out := make(map[string]any, len(value))
		for k, v := range value {
			out[k] = cloneNeutral(v)
		}
		return out
	case []any:
		out := make([]any, len(value))
		for i, v := range value {
			out[i] = cloneNeutral(v)
		}
		return out
	default:
		return value
	}
}
