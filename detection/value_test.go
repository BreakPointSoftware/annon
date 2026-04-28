package detection

import "testing"

func TestPatternValueDetector(t *testing.T) {
	d := PatternValueDetector{}
	if match := d.DetectValue(123); match.Found() {
		t.Fatalf("non-string should not match: %+v", match)
	}
	if match := d.DetectValue("greg@example.com"); match.Strategy != Email || match.Confidence != ValuePatternMatch {
		t.Fatalf("unexpected email match: %+v", match)
	}
}
