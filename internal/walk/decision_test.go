package walk

import (
	"testing"

	"github.com/BreakPointSoftware/annon/internal/detection"
	"github.com/BreakPointSoftware/annon/internal/redactcore"
)

func testDeciderConfig() Config {
	return Config{
		UseTags:           true,
		UseFieldDetection: true,
		UseValueDetection: false,
		Detector:          detection.NewDetector(detection.DefaultRules(), false),
		Preservation:      redactcore.DefaultConfig(),
	}
}

func TestDeciderPrecedence(t *testing.T) {
	d := NewDecider(testDeciderConfig())
	dec, err := d.Decide("email", "false", "greg@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if !dec.skip {
		t.Fatalf("expected skip decision, got %+v", dec)
	}

	dec, err = d.Decide("note", "email", "plain")
	if err != nil {
		t.Fatal(err)
	}
	if dec.strategyName != "email" {
		t.Fatalf("expected explicit email strategy, got %+v", dec)
	}
}

func TestDeciderRespectsDetectionFlags(t *testing.T) {
	cfg := testDeciderConfig()
	cfg.UseFieldDetection = false
	cfg.UseValueDetection = true
	cfg.Detector = detection.NewDetector(nil, true)
	d := NewDecider(cfg)
	dec, err := d.Decide("email", "", "greg@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if dec.strategyName != "email" {
		t.Fatalf("expected value detection strategy, got %+v", dec)
	}
}

func TestDeciderRejectsUnknownTagStrategy(t *testing.T) {
	d := NewDecider(testDeciderConfig())
	if _, err := d.Decide("note", "notReal", "value"); err == nil {
		t.Fatal("expected invalid tag error")
	}
}
