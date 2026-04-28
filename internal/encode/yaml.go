package encode

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func DecodeYAML(input []byte) (any, error) {
	var out any
	if err := yaml.Unmarshal(input, &out); err != nil {
		return nil, err
	}
	return normaliseYAML(out), nil
}

func EncodeYAML(input any) ([]byte, error) {
	return yaml.Marshal(input)
}

func normaliseYAML(input any) any {
	switch value := input.(type) {
	case map[string]any:
		out := make(map[string]any, len(value))
		for k, v := range value {
			out[k] = normaliseYAML(v)
		}
		return out
	case map[any]any:
		out := make(map[string]any, len(value))
		for k, v := range value {
			out[asString(k)] = normaliseYAML(v)
		}
		return out
	case []any:
		out := make([]any, len(value))
		for i, v := range value {
			out[i] = normaliseYAML(v)
		}
		return out
	default:
		return input
	}
}

func asString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}
