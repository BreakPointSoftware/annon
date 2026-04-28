package walk

import (
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
	typed *TypedCopier
	blob  *BlobBuilder
}

func New(cfg Config, cache *TypeCache) *Walker {
	if cache == nil {
		cache = NewTypeCache()
	}
	decider := NewDecider(cfg)
	return &Walker{
		typed: NewTypedCopier(cfg, decider, cache),
		blob:  NewBlobBuilder(cfg, decider, cache),
	}
}

func (w *Walker) Copy(input any) (any, error) {
	return w.typed.Copy(input)
}

type TypedCopier struct {
	cfg     Config
	decider *Decider
	cache   *TypeCache
}

func NewTypedCopier(cfg Config, decider *Decider, cache *TypeCache) *TypedCopier {
	return &TypedCopier{cfg: cfg, decider: decider, cache: cache}
}

func (c *TypedCopier) Copy(input any) (any, error) {
	if input == nil {
		return nil, nil
	}
	v := reflect.ValueOf(input)
	cloned, err := c.copyValue(v, "", "", true)
	if err != nil {
		return nil, err
	}
	return cloned.Interface(), nil
}

func (c *TypedCopier) copyValue(v reflect.Value, fieldName string, tag string, allowDecision bool) (reflect.Value, error) {
	if !v.IsValid() {
		return v, nil
	}
	if allowDecision {
		dec, err := c.decider.Decide(fieldName, tag, valueInterface(v))
		if err != nil {
			return reflect.Value{}, err
		}
		if dec.skip {
			return c.cloneValue(v)
		}
		if dec.remove {
			return reflect.Zero(v.Type()), nil
		}
		if dec.strategyName != "" {
			return c.applyTypedAction(v, dec.strategyName)
		}
	}
	return c.copyChildren(v, fieldName, tag)
}

func (c *TypedCopier) copyChildren(v reflect.Value, fieldName string, tag string) (reflect.Value, error) {
	switch v.Kind() {
	case reflect.Interface:
		return c.copyInterface(v, fieldName, tag)
	case reflect.Pointer:
		return c.copyPointer(v)
	case reflect.Struct:
		return c.copyStruct(v)
	case reflect.Map:
		return c.copyMap(v)
	case reflect.Slice:
		return c.copySlice(v)
	case reflect.Array:
		return c.copyArray(v)
	default:
		return v, nil
	}
}

func (c *TypedCopier) copyInterface(v reflect.Value, fieldName string, tag string) (reflect.Value, error) {
	if v.IsNil() {
		return reflect.Zero(v.Type()), nil
	}
	copied, err := c.copyValue(v.Elem(), fieldName, tag, true)
	if err != nil {
		return reflect.Value{}, err
	}
	out := reflect.New(v.Type()).Elem()
	out.Set(copied)
	return out, nil
}

func (c *TypedCopier) copyPointer(v reflect.Value) (reflect.Value, error) {
	if v.IsNil() {
		return reflect.Zero(v.Type()), nil
	}
	copied, err := c.copyValue(v.Elem(), "", "", true)
	if err != nil {
		return reflect.Value{}, err
	}
	out := reflect.New(v.Type().Elem())
	out.Elem().Set(copied)
	return out, nil
}

func (c *TypedCopier) copyStruct(v reflect.Value) (reflect.Value, error) {
	out := reflect.New(v.Type()).Elem()
	for _, meta := range c.cache.StructFields(v.Type()) {
		copied, err := c.copyValue(v.FieldByIndex(meta.Index), meta.DetectionName("typed"), meta.AnonymiseTag, true)
		if err != nil {
			return reflect.Value{}, err
		}
		out.FieldByIndex(meta.Index).Set(copied)
	}
	return out, nil
}

func (c *TypedCopier) copyMap(v reflect.Value) (reflect.Value, error) {
	if v.IsNil() {
		return reflect.Zero(v.Type()), nil
	}
	out := reflect.MakeMapWithSize(v.Type(), v.Len())
	iter := v.MapRange()
	for iter.Next() {
		key, value := iter.Key(), iter.Value()
		field := ""
		if key.Kind() == reflect.String {
			field = key.String()
		}
		copied, err := c.copyValue(value, field, "", true)
		if err != nil {
			return reflect.Value{}, err
		}
		out.SetMapIndex(key, copied)
	}
	return out, nil
}

func (c *TypedCopier) copySlice(v reflect.Value) (reflect.Value, error) {
	if v.IsNil() {
		return reflect.Zero(v.Type()), nil
	}
	out := reflect.MakeSlice(v.Type(), v.Len(), v.Len())
	for i := 0; i < v.Len(); i++ {
		copied, err := c.copyValue(v.Index(i), "", "", true)
		if err != nil {
			return reflect.Value{}, err
		}
		out.Index(i).Set(copied)
	}
	return out, nil
}

func (c *TypedCopier) copyArray(v reflect.Value) (reflect.Value, error) {
	out := reflect.New(v.Type()).Elem()
	for i := 0; i < v.Len(); i++ {
		copied, err := c.copyValue(v.Index(i), "", "", true)
		if err != nil {
			return reflect.Value{}, err
		}
		out.Index(i).Set(copied)
	}
	return out, nil
}

func (c *TypedCopier) cloneValue(v reflect.Value) (reflect.Value, error) { return c.copyValue(v, "", "", false) }

func (c *TypedCopier) applyTypedAction(v reflect.Value, strategyName string) (reflect.Value, error) {
	return applyAction(v, strategyName, c.cfg)
}

func valueInterface(v reflect.Value) any { if !v.IsValid() { return nil }; return v.Interface() }
