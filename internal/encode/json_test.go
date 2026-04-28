package encode

import "testing"

func TestJSONCodec(t *testing.T) {
	decoded, err := DecodeJSON([]byte(`{"email":"greg@example.com"}`))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := EncodeJSON(decoded); err != nil {
		t.Fatal(err)
	}
}
