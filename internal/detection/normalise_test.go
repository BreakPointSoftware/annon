package detection

import "testing"

func TestNormalise(t *testing.T) {
	if got := Normalise("FirstName"); got != "firstname" {
		t.Fatalf("unexpected normalised value: %q", got)
	}
}
