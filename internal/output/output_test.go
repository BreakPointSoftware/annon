package output

import (
	"testing"

	copyinternal "github.com/BreakPointSoftware/annon/internal/copy"
	"github.com/BreakPointSoftware/annon/internal/decision"
	"github.com/BreakPointSoftware/annon/internal/detection"
	"github.com/BreakPointSoftware/annon/internal/redactcore"
)

func TestOutputFromNeutral(t *testing.T) {
	config := decision.Config{UseTags: true, UseFieldDetection: true, UseValueDetection: false, Detector: detection.NewDetector(detection.DefaultRules(), false), Preservation: redactcore.DefaultConfig()}
	builder := New(config, decision.New(config), copyinternal.NewTypeCache())
	result, err := builder.OutputFromNeutral(map[string]any{"email": "greg@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	outputMap := result.(map[string]any)
	if outputMap["email"] != "g***@example.com" {
		t.Fatalf("unexpected result: %#v", outputMap)
	}
}
