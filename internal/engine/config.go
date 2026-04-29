package engine

import (
	"github.com/BreakPointSoftware/annon/internal/detection"
	"github.com/BreakPointSoftware/annon/internal/redactcore"
)

type Config struct {
	UseTags           bool
	UseFieldDetection bool
	UseValueDetection bool
	FieldRules        []detection.Rule
	Preservation      redactcore.Config
}

func DefaultConfig() Config {
	return Config{
		UseTags:           true,
		UseFieldDetection: true,
		UseValueDetection: false,
		Preservation:      redactcore.DefaultConfig(),
	}
}
