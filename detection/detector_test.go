package detection

import "testing"

func TestCompiledDetectorDetect(t *testing.T) {
	d := NewCompiledDetector(DefaultRules(), PatternValueDetector{}, true)
	if match := d.Detect("email", "not-an-email"); match.Strategy != Email || match.Confidence != StrongMatch {
		t.Fatalf("expected strong field match, got %+v", match)
	}
	if match := d.Detect("note", "greg@example.com"); match.Strategy != Email || match.Confidence != ValuePatternMatch {
		t.Fatalf("expected value match, got %+v", match)
	}
	if match := d.Detect("note", "plain text"); match.Found() {
		t.Fatalf("expected no match, got %+v", match)
	}
}
