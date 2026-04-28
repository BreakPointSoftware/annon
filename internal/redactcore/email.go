package redactcore

import "strings"

func Email(value string, cfg Config) string {
	parts := strings.SplitN(value, "@", 2)
	if len(parts) != 2 {
		return value
	}
	localPart, domain := parts[0], parts[1]
	keep := clampKeep(len(localPart), cfg.Email.KeepLocalPrefix)
	masked := keepPrefix(localPart, keep, cfg.RedactChar)
	if cfg.Email.KeepDomain {
		return masked + "@" + domain
	}
	return masked
}
