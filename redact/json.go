package redact

import "github.com/BreakPointSoftware/annon/internal/engine"

func JSON(input any) []byte {
	return engine.New(engine.DefaultConfig()).JSON(input)
}

func YAML(input any) []byte {
	return engine.New(engine.DefaultConfig()).YAML(input)
}

func JSONBytes(input []byte) []byte {
	return engine.New(engine.DefaultConfig()).JSONBytes(input)
}

func YAMLBytes(input []byte) []byte {
	return engine.New(engine.DefaultConfig()).YAMLBytes(input)
}
