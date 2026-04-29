package detection

import "testing"

func TestDetectorDetect(t *testing.T) {
	detector := NewDetector(DefaultRules(), true)
	if match := detector.Detect("email", "not-an-email"); match.Strategy != Email {
		t.Fatalf("unexpected field match: %+v", match)
	}
}
