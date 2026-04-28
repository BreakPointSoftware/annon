package detection

import "testing"

func TestFieldDetectionHierarchy(t *testing.T) {
	d := NewCompiledDetector(DefaultRules(), nil, false)

	tests := []struct {
		field      string
		strategy   Strategy
		confidence Confidence
	}{
		{"surname", Surname, StrongMatch},
		{"firstName", FirstName, StrongMatch},
		{"vehicleRegistration", VehicleRegistration, StrongMatch},
		{"reg", VehicleRegistration, FallbackMatch},
		{"customerName", Name, ContainsMatch},
	}

	for _, tt := range tests {
		match := d.DetectField(tt.field)
		if match.Strategy != tt.strategy || match.Confidence != tt.confidence {
			t.Fatalf("DetectField(%q) = %+v", tt.field, match)
		}
	}

	for _, field := range []string{"username", "fileName", "hostName", "domainName"} {
		if match := d.DetectField(field); match.Found() {
			t.Fatalf("expected no name match for %q, got %+v", field, match)
		}
	}
}

func TestCustomStrongOverridesContains(t *testing.T) {
	rules := append(DefaultRules(), StrongRule(Redact, "customerName"))
	d := NewCompiledDetector(rules, nil, false)
	match := d.DetectField("customerName")
	if match.Strategy != Redact || match.Confidence != StrongMatch {
		t.Fatalf("unexpected match: %+v", match)
	}
}

func TestCompiledDetectorRetainsCustomRulesAfterDefaults(t *testing.T) {
	rules := append(DefaultRules(), StrongRule(Redact, "driverAlias"))
	d := NewCompiledDetector(rules, nil, false)
	if got := d.DetectField("driverAlias"); got.Strategy != Redact {
		t.Fatalf("expected custom rule to be compiled, got %+v", got)
	}
}
