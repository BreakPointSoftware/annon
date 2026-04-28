package redactcore

import "strings"

func VehicleRegistration(value string, cfg Config) string {
	compact := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(value), " ", ""))
	keep := clampKeep(len(compact), cfg.VehicleRegistration.KeepPrefix)
	if compact == "" {
		return value
	}
	masked := keepPrefix(compact, keep, cfg.RedactChar)
	if len(compact) > 4 && keep >= 4 {
		return masked[:4] + " " + masked[4:]
	}
	return masked
}
