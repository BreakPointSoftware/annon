package anonymise

import (
	"fmt"

	"github.com/your-org/annon/detection"
	"github.com/your-org/annon/walker"
)

type Anonymiser struct {
	config Config
	walker *walker.Walker
	cache  *walker.TypeCache
}

func New(opts ...Option) (*Anonymiser, error) {
	cfg, err := resolveConfig(opts...)
	if err != nil {
		return nil, err
	}
	cache := walker.NewTypeCache()
	return &Anonymiser{
		config: cfg,
		walker: walker.New(walker.Config{
			UseTags:           cfg.UseTags,
			UseFieldDetection: cfg.UseFieldDetection,
			UseValueDetection: cfg.UseValueDetection,
			Detector:          cfg.Detector,
			Strategies:        cfg.Strategies,
			Preservation:      cfg.Preservation,
		}, cache),
		cache: cache,
	}, nil
}

func Copy[T any](input T, opts ...Option) (T, error) {
	a, err := New(opts...)
	if err != nil {
		var zero T
		return zero, err
	}
	result, err := a.Copy(input)
	if err != nil {
		var zero T
		return zero, err
	}
	return result.(T), nil
}

func JSON(input any, opts ...Option) ([]byte, error) {
	a, err := New(opts...)
	if err != nil {
		return nil, err
	}
	return a.JSON(input)
}

func YAML(input any, opts ...Option) ([]byte, error) {
	a, err := New(opts...)
	if err != nil {
		return nil, err
	}
	return a.YAML(input)
}

func FromJSON(input []byte, opts ...Option) ([]byte, error) {
	a, err := New(opts...)
	if err != nil {
		return nil, err
	}
	return a.FromJSON(input)
}

func FromYAML(input []byte, opts ...Option) ([]byte, error) {
	a, err := New(opts...)
	if err != nil {
		return nil, err
	}
	return a.FromYAML(input)
}

func resolveConfig(opts ...Option) (Config, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			return Config{}, err
		}
	}
	if cfg.Detector != nil && len(cfg.FieldRules) > 0 {
		return Config{}, fmt.Errorf("cannot combine custom detector with additional field rules")
	}
	if cfg.Detector == nil && (cfg.UseFieldDetection || cfg.UseValueDetection) {
		rules := []detection.Rule(nil)
		if cfg.UseFieldDetection {
			rules = append(rules, detection.DefaultRules()...)
			rules = append(rules, cfg.FieldRules...)
		}
		cfg.Detector = detection.NewCompiledDetector(rules, detection.PatternValueDetector{}, cfg.UseValueDetection)
	}
	return cfg, nil
}
