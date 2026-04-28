package annon

import (
	"github.com/BreakPointSoftware/annon/internal/detection"
	"github.com/BreakPointSoftware/annon/internal/walk"
)

type Anonymiser struct {
	config config
	walker *walk.Walker
	cache  *walk.TypeCache
}

func New(opts ...Option) (*Anonymiser, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			return nil, err
		}
	}
	var detector *detection.CompiledDetector
	if cfg.UseFieldDetection || cfg.UseValueDetection {
		rules := []detection.Rule(nil)
		if cfg.UseFieldDetection {
			rules = append(rules, detection.DefaultRules()...)
			rules = append(rules, cfg.FieldRules...)
		}
		detector = detection.NewCompiledDetector(rules, detection.PatternValueDetector{}, cfg.UseValueDetection)
	}
	cache := walk.NewTypeCache()
	return &Anonymiser{config: cfg, walker: walk.New(walk.Config{UseTags: cfg.UseTags, UseFieldDetection: cfg.UseFieldDetection, UseValueDetection: cfg.UseValueDetection, Detector: detector, Preservation: cfg.Preservation}, cache), cache: cache}, nil
}

func Copy[T any](input T, opts ...Option) (T, error) {
	a, err := New(opts...)
	if err != nil { var zero T; return zero, err }
	result, err := a.Copy(input)
	if err != nil { var zero T; return zero, err }
	return result.(T), nil
}

func JSON(input any, opts ...Option) ([]byte, error) {
	a, err := New(opts...)
	if err != nil { return nil, err }
	return a.JSON(input)
}

func YAML(input any, opts ...Option) ([]byte, error) {
	a, err := New(opts...)
	if err != nil { return nil, err }
	return a.YAML(input)
}

func FromJSON(input []byte, opts ...Option) ([]byte, error) {
	a, err := New(opts...)
	if err != nil { return nil, err }
	return a.FromJSON(input)
}

func FromYAML(input []byte, opts ...Option) ([]byte, error) {
	a, err := New(opts...)
	if err != nil { return nil, err }
	return a.FromYAML(input)
}
