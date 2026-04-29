package redact

import "github.com/BreakPointSoftware/annon/internal/engine"

func Data(input any) any {
	return engine.New(engine.DefaultConfig()).Data(input)
}
