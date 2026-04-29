package walk

import (
	"reflect"
	"sync"

	"github.com/BreakPointSoftware/annon/internal/detection"
	"github.com/BreakPointSoftware/annon/internal/support/tags"
)

type FieldMeta struct {
	Index        []int
	Name         string
	JSONName     string
	YAMLName     string
	AnonymiseTag string
	IgnoreJSON   bool
	IgnoreYAML   bool
}

func (f FieldMeta) DetectionName(format string) string {
	switch format {
	case "json":
		if f.JSONName != "" && !f.IgnoreJSON {
			return f.JSONName
		}
	case "yaml":
		if f.YAMLName != "" && !f.IgnoreYAML {
			return f.YAMLName
		}
	}

	if f.JSONName != "" && !f.IgnoreJSON {
		return f.JSONName
	}

	if f.YAMLName != "" && !f.IgnoreYAML {
		return f.YAMLName
	}

	return f.Name
}

func (f FieldMeta) OutputName(format string) string {
	if format == "yaml" {
		if f.IgnoreYAML {
			return ""
		}

		if f.YAMLName != "" {
			return f.YAMLName
		}
	}

	if f.IgnoreJSON {
		return ""
	}

	if f.JSONName != "" {
		return f.JSONName
	}

	return f.Name
}

type TypeCache struct{ fields sync.Map }

func NewTypeCache() *TypeCache { return &TypeCache{} }

func (c *TypeCache) StructFields(t reflect.Type) []FieldMeta {
	if cached, ok := c.fields.Load(t); ok {
		return cached.([]FieldMeta)
	}

	fields := compileStructFields(t)
	c.fields.Store(t, fields)
	return fields
}

func compileStructFields(t reflect.Type) []FieldMeta {
	fields := make([]FieldMeta, 0, t.NumField())

	for index := 0; index < t.NumField(); index++ {
		field := t.Field(index)
		if field.PkgPath != "" {
			continue
		}

		jsonName, ignoreJSON := tags.ParseSerialiseTag(field.Tag.Get("json"))
		yamlName, ignoreYAML := tags.ParseSerialiseTag(field.Tag.Get("yaml"))
		_ = detection.Normalise(field.Name)

		fields = append(fields, FieldMeta{
			Index:        field.Index,
			Name:         field.Name,
			JSONName:     jsonName,
			YAMLName:     yamlName,
			AnonymiseTag: field.Tag.Get("anonymise"),
			IgnoreJSON:   ignoreJSON,
			IgnoreYAML:   ignoreYAML,
		})
	}

	return fields
}
