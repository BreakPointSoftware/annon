package walker

import (
	"reflect"
	"testing"
)

func TestBlobFromValue(t *testing.T) {
	w := New(testWalkerConfig(), nil)
	input := typedCustomer{ID: "123", Email: "greg@example.com", Secret: "secret", Note: "plain"}
	result, err := w.BlobFromValue(input, "json")
	if err != nil {
		t.Fatal(err)
	}
	got := result.(map[string]any)
	if got["email"] != "g***@example.com" {
		t.Fatalf("unexpected blob email: %#v", got)
	}
	if _, ok := got["Secret"]; ok {
		t.Fatalf("expected removed field to be omitted: %#v", got)
	}
	if !reflect.DeepEqual(input, typedCustomer{ID: "123", Email: "greg@example.com", Secret: "secret", Note: "plain"}) {
		t.Fatal("input mutated")
	}
}

func TestBlobFromNeutral(t *testing.T) {
	w := New(testWalkerConfig(), nil)
	input := map[string]any{"email": "greg@example.com", "items": []any{map[string]any{"reg": "AB12 CDE"}}}
	result, err := w.BlobFromNeutral(input)
	if err != nil {
		t.Fatal(err)
	}
	got := result.(map[string]any)
	if got["email"] != "g***@example.com" {
		t.Fatalf("unexpected blob result: %#v", got)
	}
	items := got["items"].([]any)
	vehicle := items[0].(map[string]any)
	if vehicle["reg"] != "AB12 ***" {
		t.Fatalf("unexpected nested vehicle result: %#v", vehicle)
	}
}

func TestBlobUsesYAMLDetectionName(t *testing.T) {
	type yamlNamed struct {
		Contact string `json:"nickname" yaml:"email"`
	}
	w := New(testWalkerConfig(), nil)
	result, err := w.BlobFromValue(yamlNamed{Contact: "greg@example.com"}, "yaml")
	if err != nil {
		t.Fatal(err)
	}
	got := result.(map[string]any)
	if got["email"] != "g***@example.com" {
		t.Fatalf("expected yaml detection name to drive masking, got %#v", got)
	}
}

func TestBlobFromNeutralArrayOfObjects(t *testing.T) {
	w := New(testWalkerConfig(), nil)
	input := []any{
		map[string]any{"email": "greg@example.com"},
		map[string]any{"reg": "AB12 CDE"},
	}
	result, err := w.BlobFromNeutral(input)
	if err != nil {
		t.Fatal(err)
	}
	items := result.([]any)
	if items[0].(map[string]any)["email"] != "g***@example.com" {
		t.Fatalf("unexpected first item: %#v", items[0])
	}
	if items[1].(map[string]any)["reg"] != "AB12 ***" {
		t.Fatalf("unexpected second item: %#v", items[1])
	}
}
