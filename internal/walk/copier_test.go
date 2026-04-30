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

type nestedCustomer struct {
	Customer *typedCustomer `json:"customer"`
	Items    []typedCustomer `json:"items"`
	Values   [2]string      `json:"values"`
	Notes    any            `json:"notes"`
	private  string
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

func TestWalkHandlesNestedPointersInterfacesAndArrays(t *testing.T) {
	cache := NewTypeCache()
	decider := decision.New(testWalkConfig())
	walker := New(testWalkConfig(), decider, cache)

	input := nestedCustomer{
		Customer: &typedCustomer{Email: "greg@example.com", Secret: "secret"},
		Items:    []typedCustomer{{Email: "one@example.com"}, {Email: "two@example.com"}},
		Values:   [2]string{"TN9 1XA", "AB12 CDE"},
		Notes:    map[string]any{"email": "greg@example.com"},
		private:  "keep-local-only",
	}

	resultAny, err := walker.Copy(input)
	if err != nil {
		t.Fatal(err)
	}

	result := resultAny.(nestedCustomer)
	if result.Customer == nil || result.Customer.Email != "g***@example.com" || result.Customer.Secret != "" {
		t.Fatalf("unexpected nested pointer result: %+v", result.Customer)
	}

	if result.Items[0].Email != "o**@example.com" || result.Items[1].Email != "t**@example.com" {
		t.Fatalf("unexpected slice result: %+v", result.Items)
	}

	if result.Values[0] != "TN9 1XA" || result.Values[1] != "AB12 CDE" {
		t.Fatalf("unexpected array preservation result: %+v", result.Values)
	}

	notesMap, ok := result.Notes.(map[string]any)
	if !ok || notesMap["email"] != "g***@example.com" {
		t.Fatalf("unexpected interface/map result: %#v", result.Notes)
	}

	if input.Customer.Email != "greg@example.com" || input.Items[0].Email != "one@example.com" {
		t.Fatalf("input mutated: %+v", input)
	}
}

func TestWalkHandlesNilPointerAndNilSlice(t *testing.T) {
	cache := NewTypeCache()
	decider := decision.New(testWalkConfig())
	walker := New(testWalkConfig(), decider, cache)

	input := nestedCustomer{}
	resultAny, err := walker.Copy(input)
	if err != nil {
		t.Fatal(err)
	}

	result := resultAny.(nestedCustomer)
	if result.Customer != nil || result.Items != nil || result.Notes != nil {
		t.Fatalf("expected nil fields to remain nil: %+v", result)
	}
}
