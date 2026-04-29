package walk

import (
	"testing"

	"github.com/BreakPointSoftware/annon/internal/decision"
	"github.com/BreakPointSoftware/annon/internal/detection"
	"github.com/BreakPointSoftware/annon/internal/redactcore"
)

type typedCustomer struct {
	Email  string `json:"email"`
	Secret string `anonymise:"remove"`
}

func testWalkConfig() decision.Config {
	return decision.Config{UseTags: true, UseFieldDetection: true, UseValueDetection: false, Detector: detection.NewDetector(detection.DefaultRules(), false), Preservation: redactcore.DefaultConfig()}
}

func TestWalkStruct(t *testing.T) {
	cache := NewTypeCache()
	decider := decision.New(testWalkConfig())
	walker := New(testWalkConfig(), decider, cache)
	resultAny, err := walker.Copy(typedCustomer{Email: "greg@example.com", Secret: "secret"})
	if err != nil {
		t.Fatal(err)
	}

	result := resultAny.(typedCustomer)
	if result.Email != "g***@example.com" || result.Secret != "" {
		t.Fatalf("unexpected walk result: %+v", result)
	}
}
