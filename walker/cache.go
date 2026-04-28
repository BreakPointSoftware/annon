package walker

import (
	"reflect"
	"strings"
	"sync"

	"github.com/BreakPointSoftware/annon/detection"
)

type FieldMeta struct {
	Index              []int
	Name               string
	JSONName           string
	YAMLName           string
	AnonymiseTag       string
	NormalisedName     string
	NormalisedJSONName string
	NormalisedYAMLName string
	IgnoreJSON         bool
	IgnoreYAML         bool
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

type TypeCache struct {
	fields sync.Map
}

func NewTypeCache() *TypeCache {
	return &TypeCache{}
}

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
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}
		jsonName, ignoreJSON := parseSerialiseTag(field.Tag.Get("json"))
		yamlName, ignoreYAML := parseSerialiseTag(field.Tag.Get("yaml"))
		meta := FieldMeta{
			Index:              field.Index,
			Name:               field.Name,
			JSONName:           jsonName,
			YAMLName:           yamlName,
			AnonymiseTag:       field.Tag.Get("anonymise"),
			NormalisedName:     detection.Normalise(field.Name),
			NormalisedJSONName: detection.Normalise(jsonName),
			NormalisedYAMLName: detection.Normalise(yamlName),
			IgnoreJSON:         ignoreJSON,
			IgnoreYAML:         ignoreYAML,
		}
		fields = append(fields, meta)
	}
	return fields
}

func parseSerialiseTag(tag string) (string, bool) {
	if tag == "-" {
		return "", true
	}
	if tag == "" {
		return "", false
	}
	parts := strings.Split(tag, ",")
	if parts[0] == "-" {
		return "", true
	}
	if parts[0] == "" {
		return "", false
	}
	return parts[0], false
}
