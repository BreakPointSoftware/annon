package annon

import (
	copyinternal "github.com/BreakPointSoftware/annon/internal/copy"
	decisionpkg "github.com/BreakPointSoftware/annon/internal/decision"
	"github.com/BreakPointSoftware/annon/internal/detection"
	outputpkg "github.com/BreakPointSoftware/annon/internal/output"
)

type Anonymiser struct {
	config        config
	copier        *copyinternal.Copier
	outputBuilder *outputpkg.Builder
	cache         *copyinternal.TypeCache
}

func New(opts ...Option) (*Anonymiser, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			return nil, err
		}
	}
	var detector *detection.Detector
	if cfg.UseFieldDetection || cfg.UseValueDetection {
		rules := []detection.Rule(nil)
		if cfg.UseFieldDetection {
			rules = append(rules, detection.DefaultRules()...)
			rules = append(rules, cfg.FieldRules...)
		}
		detector = detection.NewDetector(rules, cfg.UseValueDetection)
	}
	decisionConfig := decisionpkg.Config{UseTags: cfg.UseTags, UseFieldDetection: cfg.UseFieldDetection, UseValueDetection: cfg.UseValueDetection, Detector: detector, Preservation: cfg.Preservation}
	cache := copyinternal.NewTypeCache()
	decider := decisionpkg.New(decisionConfig)

	return &Anonymiser{
		config:        cfg,
		copier:        copyinternal.New(decisionConfig, decider, cache),
		outputBuilder: outputpkg.New(decisionConfig, decider, cache),
		cache:         cache,
	}, nil
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
