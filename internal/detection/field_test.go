package detection

import "testing"

func TestFieldDetectionHierarchy(t *testing.T) {
	detector := NewDetector(DefaultRules(), false)
	if match := detector.DetectField("customerName"); match.Strategy != Name || match.Confidence != ContainsMatch {
		t.Fatalf("unexpected match: %+v", match)
	}
	for _, field := range []string{"username", "fileName", "hostName", "domainName"} {
		if match := detector.DetectField(field); match.Found() {
			t.Fatalf("expected no match for %q, got %+v", field, match)
		}
	}
}
