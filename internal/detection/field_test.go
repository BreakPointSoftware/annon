package detection

import "testing"

func TestFieldDetectionHierarchy(t *testing.T) {
	d := NewCompiledDetector(DefaultRules(), nil, false)
	if match := d.DetectField("customerName"); match.Strategy != Name || match.Confidence != ContainsMatch {
		t.Fatalf("unexpected match: %+v", match)
	}
	for _, field := range []string{"username", "fileName", "hostName", "domainName"} {
		if match := d.DetectField(field); match.Found() {
			t.Fatalf("expected no match for %q, got %+v", field, match)
		}
	}
}
