package redactcore

import "strings"

func Postcode(value string, cfg Config) string {
	s := strings.TrimSpace(value)
	if s == "" {
		return s
	}
	compact := strings.ToUpper(strings.ReplaceAll(s, " ", ""))
	if len(compact) < 3 {
		return value
	}
	if !cfg.Postcode.KeepOutward {
		return repeatRune(cfg.RedactChar, len(compact))
	}
	outward, inward := splitPostcode(compact)
	if outward == "" || inward == "" {
		return value
	}
	return outward + " " + repeatRune(cfg.RedactChar, len(inward))
}

func splitPostcode(compact string) (string, string) {
	if len(compact) <= 3 {
		return "", ""
	}
	return compact[:len(compact)-3], compact[len(compact)-3:]
}
