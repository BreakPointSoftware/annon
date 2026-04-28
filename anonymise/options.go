package anonymise

import (
	"fmt"

	"github.com/BreakPointSoftware/annon/detection"
	"github.com/BreakPointSoftware/annon/preservation"
	strategypkg "github.com/BreakPointSoftware/annon/strategy"
)

type Option func(*Config) error

func WithTags(enabled bool) Option {
	return func(cfg *Config) error {
		cfg.UseTags = enabled
		return nil
	}
}

func WithFieldDetection(enabled bool) Option {
	return func(cfg *Config) error {
		cfg.UseFieldDetection = enabled
		return nil
	}
}

func WithValueDetection(enabled bool) Option {
	return func(cfg *Config) error {
		cfg.UseValueDetection = enabled
		return nil
	}
}

func WithDetector(detectorImpl detection.Detector) Option {
	return func(cfg *Config) error {
		cfg.Detector = detectorImpl
		return nil
	}
}

func WithStrategies(strategies ...strategypkg.Strategy) Option {
	return func(cfg *Config) error {
		if cfg.Strategies == nil {
			cfg.Strategies = map[string]strategypkg.Strategy{}
		}
		for _, strategyImpl := range strategies {
			cfg.Strategies[strategyImpl.Name()] = strategyImpl
		}
		return nil
	}
}

func WithDefaultStrategies() Option {
	return func(cfg *Config) error {
		if cfg.Strategies == nil {
			cfg.Strategies = map[string]strategypkg.Strategy{}
		}
		for _, strategyImpl := range strategypkg.DefaultStrategies() {
			cfg.Strategies[strategyImpl.Name()] = strategyImpl
		}
		return nil
	}
}

func WithDefaultFieldDetection() Option {
	return func(cfg *Config) error {
		cfg.UseFieldDetection = true
		cfg.Detector = nil
		return nil
	}
}

func WithFieldRules(rules ...detection.Rule) Option {
	return func(cfg *Config) error {
		if cfg.Detector != nil {
			return fmt.Errorf("cannot combine WithFieldRules with WithDetector")
		}
		cfg.FieldRules = append(cfg.FieldRules, rules...)
		cfg.UseFieldDetection = true
		return nil
	}
}

func WithPreservation(config preservation.Config) Option {
	return func(cfg *Config) error {
		cfg.Preservation = config
		return nil
	}
}

func WithRedactChar(char rune) Option {
	return func(cfg *Config) error {
		cfg.Preservation.RedactChar = char
		return nil
	}
}

func WithRedactionText(text string) Option {
	return func(cfg *Config) error {
		cfg.Preservation.RedactionText = text
		return nil
	}
}
