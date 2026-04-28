package strategy

import "strings"

type RedactStrategy struct{}

func (RedactStrategy) Name() string { return "redact" }

func (RedactStrategy) Anonymise(value any, ctx Context) (any, error) {
	if _, ok := value.(string); !ok {
		return value, nil
	}
	if strings.TrimSpace(ctx.Preservation.RedactionText) == "" {
		return "[REDACTED]", nil
	}
	return ctx.Preservation.RedactionText, nil
}

func clampKeep(length, keep int) int {
	if keep < 0 {
		return 0
	}
	if keep > length {
		return length
	}
	return keep
}

func keepPrefix(input string, keep int, redactChar rune) string {
	keep = clampKeep(len(input), keep)
	if keep == len(input) {
		return input
	}
	return input[:keep] + repeatRune(redactChar, len(input)-keep)
}

func repeatRune(r rune, count int) string {
	if count <= 0 {
		return ""
	}
	return strings.Repeat(string(r), count)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
