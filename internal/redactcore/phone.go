package redactcore

import "strings"

func Phone(value string, cfg Config) string {
	compact := strings.NewReplacer(" ", "", "-", "", "(", "", ")", "").Replace(value)
	keep := clampKeep(len(compact), cfg.Phone.KeepLast)
	if keep == len(compact) {
		return compact
	}
	return repeatRune(cfg.RedactChar, len(compact)-keep) + compact[len(compact)-keep:]
}
