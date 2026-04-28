package strategy

import "testing"

func TestDefaultStrategies(t *testing.T) {
	seen := map[string]bool{}
	for _, s := range DefaultStrategies() {
		if seen[s.Name()] {
			t.Fatalf("duplicate strategy: %s", s.Name())
		}
		seen[s.Name()] = true
	}
	for _, name := range []string{"email", "phone", "postcode", "name", "firstName", "surname", "vehicleRegistration", "redact"} {
		if !seen[name] {
			t.Fatalf("missing strategy %q", name)
		}
	}
}
