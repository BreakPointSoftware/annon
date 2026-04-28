package detection

import "testing"

func TestNormaliseExamples(t *testing.T) {
	tests := map[string]string{
		"FirstName":        "firstname",
		"first_name":       "firstname",
		"first-name":       "firstname",
		"first name":       "firstname",
		"FIRST.NAME":       "firstname",
		"vehicle_reg":      "vehiclereg",
		"addressLine1":     "addressline1",
		" address line 1 ": "addressline1",
	}
	for input, want := range tests {
		if got := Normalise(input); got != want {
			t.Fatalf("Normalise(%q) = %q, want %q", input, got, want)
		}
	}
}
