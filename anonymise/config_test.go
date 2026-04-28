package anonymise

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()
	if !cfg.UseTags || !cfg.UseFieldDetection || cfg.UseValueDetection {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
	if cfg.Preservation.RedactionText != "[REDACTED]" {
		t.Fatalf("unexpected preservation defaults: %+v", cfg.Preservation)
	}
}
