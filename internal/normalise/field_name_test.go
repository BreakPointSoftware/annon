package normalise

import "testing"

func TestFieldName(t *testing.T) {
	if got := FieldName("FIRST.NAME"); got != "firstname" {
		t.Fatalf("unexpected normalised field name: %q", got)
	}
}
