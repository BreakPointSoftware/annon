package walk

import "testing"

func TestBlobFromNeutral(t *testing.T) {
	w := New(testWalkerConfig(), nil)
	result, err := w.BlobFromNeutral(map[string]any{"email": "greg@example.com"})
	if err != nil { t.Fatal(err) }
	got := result.(map[string]any)
	if got["email"] != "g***@example.com" { t.Fatalf("unexpected result: %#v", got) }
}
