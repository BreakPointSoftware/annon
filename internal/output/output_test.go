package output

import (
	"testing"

	"github.com/BreakPointSoftware/annon/internal/decision"
	"github.com/BreakPointSoftware/annon/internal/detection"
	"github.com/BreakPointSoftware/annon/internal/redactcore"
	"github.com/BreakPointSoftware/annon/internal/walk"
)

func TestOutputFromNeutral(t *testing.T) {
	config := decision.Config{UseTags: true, UseFieldDetection: true, UseValueDetection: false, Detector: detection.NewDetector(detection.DefaultRules(), false), Preservation: redactcore.DefaultConfig()}
	builder := New(config, decision.New(config), walk.NewTypeCache())
	result, err := builder.OutputFromNeutral(map[string]any{"email": "greg@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	outputMap := result.(map[string]any)
	if outputMap["email"] != "g***@example.com" {
		t.Fatalf("unexpected result: %#v", outputMap)
	}
}
