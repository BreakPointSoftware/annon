package encode

import "testing"

func TestYAMLCodec(t *testing.T) {
	decoded, err := DecodeYAML([]byte("email: greg@example.com\n"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := EncodeYAML(decoded); err != nil {
		t.Fatal(err)
	}
}
