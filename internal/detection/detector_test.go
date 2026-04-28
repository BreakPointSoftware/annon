package detection

import "testing"

func TestCompiledDetectorDetect(t *testing.T) {
	d := NewCompiledDetector(DefaultRules(), PatternValueDetector{}, true)
	if match := d.Detect("email", "not-an-email"); match.Strategy != Email {
		t.Fatalf("unexpected field match: %+v", match)
	}
}
