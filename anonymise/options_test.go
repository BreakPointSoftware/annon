package anonymise

import (
	"testing"

	"github.com/BreakPointSoftware/annon/detection"
	"github.com/BreakPointSoftware/annon/preservation"
)

func TestOptions(t *testing.T) {
	cfg := defaultConfig()
	_ = WithTags(false)(&cfg)
	_ = WithFieldDetection(false)(&cfg)
	_ = WithValueDetection(true)(&cfg)
	_ = WithRedactChar('x')(&cfg)
	_ = WithRedactionText("[hidden]")(&cfg)
	_ = WithPreservation(preservation.Default())(&cfg)
	if cfg.UseTags || cfg.UseFieldDetection || !cfg.UseValueDetection {
		t.Fatalf("unexpected config flags: %+v", cfg)
	}
	if cfg.Preservation.RedactChar != '*' {
		t.Fatalf("expected preservation override to apply: %+v", cfg.Preservation)
	}
}

func TestWithFieldRules(t *testing.T) {
	cfg := defaultConfig()
	if err := WithFieldRules(detection.StrongRule(detection.Redact, "customerAlias"))(&cfg); err != nil {
		t.Fatal(err)
	}
	if len(cfg.FieldRules) != 1 {
		t.Fatalf("expected custom field rule to be recorded: %+v", cfg)
	}
}
