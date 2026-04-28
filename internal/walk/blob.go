package walk

import (
	"fmt"
	"reflect"

	"github.com/BreakPointSoftware/annon/internal/redactcore"
)

var omitValue = &struct{}{}

func (w *Walker) BlobFromValue(input any, format string) (any, error) {
	return w.blob.BlobFromValue(input, format)
}

func (w *Walker) BlobFromNeutral(input any) (any, error) {
	return w.blob.BlobFromNeutral(input)
}

type BlobBuilder struct {
	cfg     Config
	decider *Decider
	cache   *TypeCache
}

func NewBlobBuilder(cfg Config, decider *Decider, cache *TypeCache) *BlobBuilder {
	return &BlobBuilder{cfg: cfg, decider: decider, cache: cache}
}

func (b *BlobBuilder) BlobFromValue(input any, format string) (any, error) {
	if input == nil { return nil, nil }
	return b.blobFromReflect(reflect.ValueOf(input), format, "", "", true)
}

func (b *BlobBuilder) BlobFromNeutral(input any) (any, error) { return b.blobFromNeutralValue(input, "", "") }

func (b *BlobBuilder) blobFromReflect(v reflect.Value, format, fieldName, tag string, allowDecision bool) (any, error) {
	if !v.IsValid() { return nil, nil }
	if allowDecision {
		dec, err := b.decider.Decide(fieldName, tag, valueInterface(v)); if err != nil { return nil, err }
		if dec.skip { return b.plainFromReflect(v, format) }
		if dec.remove { return omitValue, nil }
		if dec.strategyName != "" {
			strategyValue, err := b.applyBlobAction(v, dec.strategyName); if err != nil { return nil, err }
			return valueInterface(strategyValue), nil
		}
	}
	return b.buildReflectValue(v, format, fieldName, tag)
}

func (b *BlobBuilder) buildReflectValue(v reflect.Value, format, fieldName, tag string) (any, error) {
	switch v.Kind() {
	case reflect.Interface:
		return b.buildInterface(v, format, fieldName, tag)
	case reflect.Pointer:
		return b.buildPointer(v, format, fieldName, tag)
	case reflect.Struct:
		return b.buildStructMap(v, format)
	case reflect.Map:
		return b.buildMap(v, format)
	case reflect.Slice, reflect.Array:
		return b.buildSlice(v, format)
	default:
		return valueInterface(v), nil
	}
}

func (b *BlobBuilder) buildInterface(v reflect.Value, format, fieldName, tag string) (any, error) {
	if v.IsNil() {
		return nil, nil
	}
	return b.blobFromReflect(v.Elem(), format, fieldName, tag, true)
}

func (b *BlobBuilder) buildPointer(v reflect.Value, format, fieldName, tag string) (any, error) {
	if v.IsNil() {
		return nil, nil
	}
	return b.blobFromReflect(v.Elem(), format, fieldName, tag, true)
}

func (b *BlobBuilder) buildStructMap(v reflect.Value, format string) (any, error) {
	out := map[string]any{}
	for _, meta := range b.cache.StructFields(v.Type()) {
		name := meta.OutputName(format)
		if name == "" {
			continue
		}
		value, err := b.blobFromReflect(v.FieldByIndex(meta.Index), format, meta.DetectionName(format), meta.AnonymiseTag, true)
		if err != nil {
			return nil, err
		}
		if value == omitValue {
			continue
		}
		out[name] = value
	}
	return out, nil
}

func (b *BlobBuilder) buildMap(v reflect.Value, format string) (any, error) {
	if v.IsNil() {
		return nil, nil
	}
	out := map[string]any{}
	iter := v.MapRange()
	for iter.Next() {
		name := fmt.Sprint(iter.Key().Interface())
		value, err := b.blobFromReflect(iter.Value(), format, name, "", true)
		if err != nil {
			return nil, err
		}
		if value == omitValue {
			continue
		}
		out[name] = value
	}
	return out, nil
}

func (b *BlobBuilder) buildSlice(v reflect.Value, format string) (any, error) {
	out := make([]any, 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		value, err := b.blobFromReflect(v.Index(i), format, "", "", true)
		if err != nil {
			return nil, err
		}
		out = append(out, value)
	}
	return out, nil
}

func (b *BlobBuilder) plainFromReflect(v reflect.Value, format string) (any, error) { return b.blobFromReflect(v, format, "", "", false) }

func (b *BlobBuilder) blobFromNeutralValue(input any, fieldName, tag string) (any, error) {
	dec, err := b.decider.Decide(fieldName, tag, input); if err != nil { return nil, err }
	if dec.skip { return cloneNeutral(input), nil }
	if dec.remove { return omitValue, nil }
	if dec.strategyName != "" { return redactcore.Apply(dec.strategyName, input, b.cfg.Preservation) }
	return b.buildNeutralValue(input)
}

func (b *BlobBuilder) buildNeutralValue(input any) (any, error) {
	switch value := input.(type) {
	case map[string]any:
		return b.buildNeutralMap(value)
	case []any:
		return b.buildNeutralSlice(value)
	default:
		return input, nil
	}
}

func (b *BlobBuilder) buildNeutralMap(value map[string]any) (any, error) {
	out := make(map[string]any, len(value))
	for key, item := range value {
		resolved, err := b.blobFromNeutralValue(item, key, "")
		if err != nil {
			return nil, err
		}
		if resolved == omitValue {
			continue
		}
		out[key] = resolved
	}
	return out, nil
}

func (b *BlobBuilder) buildNeutralSlice(value []any) (any, error) {
	out := make([]any, 0, len(value))
	for _, item := range value {
		resolved, err := b.blobFromNeutralValue(item, "", "")
		if err != nil {
			return nil, err
		}
		out = append(out, resolved)
	}
	return out, nil
}

func (b *BlobBuilder) applyBlobAction(v reflect.Value, strategyName string) (reflect.Value, error) {
	return applyAction(v, strategyName, b.cfg)
}

func cloneNeutral(input any) any {
	switch value := input.(type) {
	case map[string]any:
		out := make(map[string]any, len(value)); for k, v := range value { out[k] = cloneNeutral(v) }; return out
	case []any:
		out := make([]any, len(value)); for i, v := range value { out[i] = cloneNeutral(v) }; return out
	default:
		return value
	}
}
