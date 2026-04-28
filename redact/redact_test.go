package redact

import "testing"

func TestDirectFunctions(t *testing.T) {
	if got := Email("greg@example.com"); got != "g***@example.com" { t.Fatalf("unexpected email: %q", got) }
	if got := Postcode("TN9 1XA"); got != "TN9 ***" { t.Fatalf("unexpected postcode: %q", got) }
}
