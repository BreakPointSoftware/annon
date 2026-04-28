package detection

import "testing"

func TestPatternHelpers(t *testing.T) {
	if !IsEmail("greg@example.com") || IsEmail("not-an-email") {
		t.Fatal("email mismatch")
	}
}
