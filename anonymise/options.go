package anonymise

import internaldetection "github.com/BreakPointSoftware/annon/internal/detection"

type Option func(*config) error

func WithTags(enabled bool) Option {
	return func(cfg *config) error { cfg.UseTags = enabled; return nil }
}

func WithFieldDetection(enabled bool) Option {
	return func(cfg *config) error { cfg.UseFieldDetection = enabled; return nil }
}

func WithValueDetection(enabled bool) Option {
	return func(cfg *config) error { cfg.UseValueDetection = enabled; return nil }
}

func WithDefaultFieldDetection() Option {
	return func(cfg *config) error { cfg.UseFieldDetection = true; return nil }
}

func WithFieldRules(rules ...FieldRule) Option {
	return func(cfg *config) error {
		for _, rule := range rules {
			cfg.FieldRules = append(cfg.FieldRules, internaldetection.Rule(rule))
		}
		cfg.UseFieldDetection = true
		return nil
	}
}

func WithPreservation(preservation PreservationConfig) Option {
	return func(cfg *config) error { cfg.Preservation = preservation; return nil }
}

func WithRedactChar(char rune) Option {
	return func(cfg *config) error { cfg.Preservation.RedactChar = char; return nil }
}

func WithRedactionText(text string) Option {
	return func(cfg *config) error { cfg.Preservation.RedactionText = text; return nil }
}
