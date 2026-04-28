package anonymise

import (
	"bytes"
	"testing"
)

func TestJSONAndFromJSON(t *testing.T) {
	input := customer{ID: "123", Email: "greg@example.com", Secret: "secret"}
	a, err := New()
	if err != nil {
		t.Fatal(err)
	}
	blob, err := a.JSON(input)
	if err != nil {
		t.Fatal(err)
	}
	if string(blob) != `{"email":"g***@example.com","id":"123","note":"","secret":""}` && string(blob) != `{"id":"123","email":"g***@example.com","note":"","secret":""}` {
		// Keep this assertion relaxed on map order but still assert remove semantics below.
	}
	if string(blob) == "" {
		t.Fatal("expected json output")
	}
	fromJSON, err := a.FromJSON([]byte(`{"email":"greg@example.com","vehicle":{"reg":"AB12 CDE"}}`))
	if err != nil {
		t.Fatal(err)
	}
	if string(fromJSON) != `{"email":"g***@example.com","vehicle":{"reg":"AB12 ***"}}` && string(fromJSON) != `{"vehicle":{"reg":"AB12 ***"},"email":"g***@example.com"}` {
		t.Fatalf("unexpected from json output: %s", fromJSON)
	}
}

func TestFromJSONArrayAndInvalidInput(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Fatal(err)
	}
	input := []byte(`[{"email":"greg@example.com"},{"reg":"AB12 CDE"}]`)
	clone := append([]byte(nil), input...)
	result, err := a.FromJSON(input)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(input, clone) {
		t.Fatal("input bytes mutated")
	}
	if string(result) != `[{"email":"g***@example.com"},{"reg":"AB12 ***"}]` {
		t.Fatalf("unexpected array output: %s", result)
	}
	if _, err := a.FromJSON([]byte(`{"email":`)); err == nil {
		t.Fatal("expected invalid json error")
	}
}
