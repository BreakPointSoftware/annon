package walk

import (
	"fmt"
	"reflect"

	"github.com/BreakPointSoftware/annon/internal/detection"
	"github.com/BreakPointSoftware/annon/internal/redactcore"
)

type Config struct {
	UseTags           bool
	UseFieldDetection bool
	UseValueDetection bool
	Detector          detection.Detector
	Preservation      redactcore.Config
}

type Walker struct {
	cfg   Config
	cache *TypeCache
}

type decision struct {
	skip         bool
	remove       bool
	strategyName string
}

func New(cfg Config, cache *TypeCache) *Walker {
	if cache == nil {
		cache = NewTypeCache()
	}
	return &Walker{cfg: cfg, cache: cache}
}

func (w *Walker) Copy(input any) (any, error) {
	if input == nil {
		return nil, nil
	}
	v := reflect.ValueOf(input)
	cloned, err := w.copyValue(v, "", "", true)
	if err != nil {
		return nil, err
	}
	return cloned.Interface(), nil
}

func (w *Walker) copyValue(v reflect.Value, fieldName string, tag string, allowDecision bool) (reflect.Value, error) {
	if !v.IsValid() {
		return v, nil
	}
	if allowDecision {
		dec, err := w.decide(fieldName, tag, valueInterface(v))
		if err != nil {
			return reflect.Value{}, err
		}
		if dec.skip {
			return w.cloneValue(v)
		}
		if dec.remove {
			return reflect.Zero(v.Type()), nil
		}
		if dec.strategyName != "" {
			return w.applyStrategy(v, dec.strategyName)
		}
	}
	switch v.Kind() {
	case reflect.Interface:
		if v.IsNil() {
			return reflect.Zero(v.Type()), nil
		}
		copied, err := w.copyValue(v.Elem(), fieldName, tag, true)
		if err != nil {
			return reflect.Value{}, err
		}
		out := reflect.New(v.Type()).Elem(); out.Set(copied); return out, nil
	case reflect.Pointer:
		if v.IsNil() { return reflect.Zero(v.Type()), nil }
		copied, err := w.copyValue(v.Elem(), "", "", true)
		if err != nil { return reflect.Value{}, err }
		out := reflect.New(v.Type().Elem()); out.Elem().Set(copied); return out, nil
	case reflect.Struct:
		out := reflect.New(v.Type()).Elem()
		for _, meta := range w.cache.StructFields(v.Type()) {
			copied, err := w.copyValue(v.FieldByIndex(meta.Index), meta.DetectionName("typed"), meta.AnonymiseTag, true)
			if err != nil { return reflect.Value{}, err }
			out.FieldByIndex(meta.Index).Set(copied)
		}
		return out, nil
	case reflect.Map:
		if v.IsNil() { return reflect.Zero(v.Type()), nil }
		out := reflect.MakeMapWithSize(v.Type(), v.Len())
		iter := v.MapRange()
		for iter.Next() {
			key, value := iter.Key(), iter.Value(); field := ""
			if key.Kind() == reflect.String { field = key.String() }
			copied, err := w.copyValue(value, field, "", true)
			if err != nil { return reflect.Value{}, err }
			out.SetMapIndex(key, copied)
		}
		return out, nil
	case reflect.Slice:
		if v.IsNil() { return reflect.Zero(v.Type()), nil }
		out := reflect.MakeSlice(v.Type(), v.Len(), v.Len())
		for i := 0; i < v.Len(); i++ {
			copied, err := w.copyValue(v.Index(i), "", "", true)
			if err != nil { return reflect.Value{}, err }
			out.Index(i).Set(copied)
		}
		return out, nil
	case reflect.Array:
		out := reflect.New(v.Type()).Elem()
		for i := 0; i < v.Len(); i++ {
			copied, err := w.copyValue(v.Index(i), "", "", true)
			if err != nil { return reflect.Value{}, err }
			out.Index(i).Set(copied)
		}
		return out, nil
	default:
		return v, nil
	}
}

func (w *Walker) cloneValue(v reflect.Value) (reflect.Value, error) { return w.copyValue(v, "", "", false) }

func (w *Walker) decide(fieldName, tag string, value any) (decision, error) {
	if w.cfg.UseTags {
		parsed := parseTag(tag)
		if err := w.validateTag(parsed); err != nil { return decision{}, err }
		if parsed.skip { return decision{skip: true}, nil }
		if parsed.remove { return decision{remove: true, strategyName: parsed.strategyName}, nil }
		if parsed.strategyName != "" && !parsed.auto { return decision{strategyName: parsed.strategyName}, nil }
		if parsed.auto {
			match := w.detect(fieldName, value)
			if match.Found() { return decision{strategyName: string(match.Strategy)}, nil }
			return decision{}, nil
		}
	}
	match := w.detect(fieldName, value)
	if match.Found() {
		if match.Strategy == detection.Remove { return decision{remove: true, strategyName: string(match.Strategy)}, nil }
		return decision{strategyName: string(match.Strategy)}, nil
	}
	return decision{}, nil
}

func (w *Walker) detect(fieldName string, value any) detection.Match {
	if w.cfg.Detector == nil { return detection.NoMatchResult() }
	if w.cfg.UseFieldDetection && w.cfg.UseValueDetection { return w.cfg.Detector.Detect(fieldName, value) }
	if w.cfg.UseFieldDetection {
		if detector, ok := w.cfg.Detector.(detection.FieldDetector); ok { return detector.DetectField(fieldName) }
		return detection.NoMatchResult()
	}
	if w.cfg.UseValueDetection {
		if detector, ok := w.cfg.Detector.(detection.ValueDetector); ok { return detector.DetectValue(value) }
		return detection.NoMatchResult()
	}
	return detection.NoMatchResult()
}

func (w *Walker) applyStrategy(v reflect.Value, strategyName string) (reflect.Value, error) {
	result, err := redactcore.Apply(strategyName, valueInterface(v), w.cfg.Preservation)
	if err != nil { return reflect.Value{}, err }
	if result == nil { return reflect.Zero(v.Type()), nil }
	resultValue := reflect.ValueOf(result)
	if resultValue.Type().AssignableTo(v.Type()) { return resultValue, nil }
	if resultValue.Type().ConvertibleTo(v.Type()) { return resultValue.Convert(v.Type()), nil }
	if v.Kind() == reflect.Interface { out := reflect.New(v.Type()).Elem(); out.Set(resultValue); return out, nil }
	return v, nil
}

func valueInterface(v reflect.Value) any { if !v.IsValid() { return nil }; return v.Interface() }

func (w *Walker) validateTag(parsed parsedTag) error {
	if parsed.empty || parsed.skip || parsed.auto || parsed.remove { return nil }
	if !redactcore.SupportedStrategy(parsed.strategyName) { return fmt.Errorf("unknown anonymise tag strategy %q", parsed.strategyName) }
	return nil
}
