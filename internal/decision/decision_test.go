package decision

import (
	"testing"

	"github.com/BreakPointSoftware/annon/internal/detection"
	"github.com/BreakPointSoftware/annon/internal/redactcore"
)

func testDecisionConfig() Config {
	return Config{
		UseTags:           true,
		UseFieldDetection: true,
		UseValueDetection: false,
		Detector:          detection.NewDetector(detection.DefaultRules(), false),
		Preservation:      redactcore.DefaultConfig(),
	}
}

func TestDeciderPrecedence(t *testing.T) {
	decider := New(testDecisionConfig())

	decisionResult, err := decider.Decide("email", "false", "greg@example.com")
	if err != nil {
		t.Fatal(err)
	}

	if !decisionResult.Skip {
		t.Fatalf("expected skip decision, got %+v", decisionResult)
	}

	decisionResult, err = decider.Decide("note", "email", "plain")
	if err != nil {
		t.Fatal(err)
	}

	if decisionResult.StrategyName != "email" {
		t.Fatalf("expected explicit email strategy, got %+v", decisionResult)
	}
}

func TestDeciderRespectsDetectionFlags(t *testing.T) {
	decisionConfig := testDecisionConfig()
	decisionConfig.UseFieldDetection = false
	decisionConfig.UseValueDetection = true
	decisionConfig.Detector = detection.NewDetector(nil, true)

	decider := New(decisionConfig)
	decisionResult, err := decider.Decide("email", "", "greg@example.com")
	if err != nil {
		t.Fatal(err)
	}

	if decisionResult.StrategyName != "email" {
		t.Fatalf("expected value detection strategy, got %+v", decisionResult)
	}
}

func TestDeciderRejectsUnknownTagStrategy(t *testing.T) {
	decider := New(testDecisionConfig())
	if _, err := decider.Decide("note", "notReal", "value"); err == nil {
		t.Fatal("expected invalid tag error")
	}
}
