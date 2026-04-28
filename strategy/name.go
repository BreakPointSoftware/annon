package strategy

import "strings"

type NameStrategy struct {
	strategyName string
}

func (s NameStrategy) Name() string { return s.strategyName }

func (s NameStrategy) Anonymise(value any, ctx Context) (any, error) {
	name, ok := value.(string)
	if !ok {
		return value, nil
	}
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return name, nil
	}
	for i, part := range parts {
		parts[i] = keepPrefix(part, clampKeep(len(part), ctx.Preservation.Name.KeepPrefix), ctx.Preservation.RedactChar)
	}
	return strings.Join(parts, " "), nil
}
