package redact

import "github.com/BreakPointSoftware/annon/internal/engine"

func String(input string) string {
	return engine.New(engine.DefaultConfig()).String(input)
}
