package encoder

import "testing"

func TestJSONCodec(t *testing.T) {
	decoded, err := DecodeJSON([]byte(`{"email":"greg@example.com"}`))
	if err != nil {
		t.Fatal(err)
	}
	blob, err := EncodeJSON(decoded)
	if err != nil {
		t.Fatal(err)
	}
	if string(blob) != `{"email":"greg@example.com"}` {
		t.Fatalf("unexpected json: %s", blob)
	}
}
