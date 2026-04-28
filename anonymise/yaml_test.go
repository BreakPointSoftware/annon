package anonymise

import (
	"bytes"
	"strings"
	"testing"
)

func TestYAMLAndFromYAML(t *testing.T) {
	input := customer{ID: "123", Email: "greg@example.com", Secret: "secret"}
	a, err := New()
	if err != nil {
		t.Fatal(err)
	}
	blob, err := a.YAML(input)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(blob), "email: g***@example.com") {
		t.Fatalf("unexpected yaml output: %s", blob)
	}
	fromYAML, err := a.FromYAML([]byte("email: greg@example.com\nvehicle:\n  reg: AB12 CDE\n"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(fromYAML), "email: g***@example.com") || !strings.Contains(string(fromYAML), "reg: AB12 ***") {
		t.Fatalf("unexpected yaml blob output: %s", fromYAML)
	}
}

func TestFromYAMLArrayAndInvalidInput(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Fatal(err)
	}
	input := []byte("- email: greg@example.com\n- reg: AB12 CDE\n")
	clone := append([]byte(nil), input...)
	result, err := a.FromYAML(input)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(input, clone) {
		t.Fatal("input bytes mutated")
	}
	text := string(result)
	if !strings.Contains(text, "email: g***@example.com") || !strings.Contains(text, "reg: AB12 ***") {
		t.Fatalf("unexpected yaml array output: %s", result)
	}
	if _, err := a.FromYAML([]byte("email: [")); err == nil {
		t.Fatal("expected invalid yaml error")
	}
}
