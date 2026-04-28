package encoder

import "testing"

func TestYAMLCodec(t *testing.T) {
	decoded, err := DecodeYAML([]byte("email: greg@example.com\nitems:\n  - reg: AB12 CDE\n"))
	if err != nil {
		t.Fatal(err)
	}
	m := decoded.(map[string]any)
	if m["email"] != "greg@example.com" {
		t.Fatalf("unexpected decoded yaml: %#v", decoded)
	}
	if _, err := EncodeYAML(decoded); err != nil {
		t.Fatal(err)
	}
}

func TestDecodeYAMLPreservesNonStringKeysSafely(t *testing.T) {
	decoded, err := DecodeYAML([]byte("1: one\ntrue: yes\n"))
	if err != nil {
		t.Fatal(err)
	}
	m := decoded.(map[string]any)
	if m["1"] != "one" || m["true"] != "yes" {
		t.Fatalf("unexpected yaml key normalisation: %#v", decoded)
	}
}
