package redactcore

import "strings"

func Name(value string, cfg Config) string {
	parts := strings.Fields(value)
	if len(parts) == 0 {
		return value
	}
	for i, part := range parts {
		parts[i] = keepPrefix(part, clampKeep(len(part), cfg.Name.KeepPrefix), cfg.RedactChar)
	}
	return strings.Join(parts, " ")
}

func FirstName(value string, cfg Config) string { return Name(value, cfg) }

func Surname(value string, cfg Config) string { return Name(value, cfg) }
