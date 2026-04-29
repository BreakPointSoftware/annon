package redact

import "github.com/BreakPointSoftware/annon/internal/redactcore"

func New(opts ...Option) (*Redactor, error) {
	cfg := redactcore.DefaultConfig()
	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			return nil, err
		}
	}
	return &Redactor{config: cfg}, nil
}

func WithConfig(config Config) Option {
	return func(cfg *Config) error { *cfg = config; return nil }
}

func WithRedactChar(char rune) Option {
	return func(cfg *Config) error { cfg.RedactChar = char; return nil }
}

func WithRedactionText(text string) Option {
	return func(cfg *Config) error { cfg.RedactionText = text; return nil }
}
