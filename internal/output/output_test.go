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

func TestOutputFromNeutralHandlesNestedCollections(t *testing.T) {
	config := decision.Config{UseTags: true, UseFieldDetection: true, UseValueDetection: false, Detector: detection.NewDetector(detection.DefaultRules(), false), Preservation: redactcore.DefaultConfig()}
	builder := New(config, decision.New(config), walk.NewTypeCache())
	result, err := builder.OutputFromNeutral(map[string]any{
		"items": []any{
			map[string]any{"email": "greg@example.com"},
			map[string]any{"reg": "AB12 CDE"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	outputMap := result.(map[string]any)
	items := outputMap["items"].([]any)
	if items[0].(map[string]any)["email"] != "g***@example.com" {
		t.Fatalf("unexpected nested email output: %#v", items[0])
	}
	if items[1].(map[string]any)["reg"] != "AB12 ***" {
		t.Fatalf("unexpected nested vehicle output: %#v", items[1])
	}
}
