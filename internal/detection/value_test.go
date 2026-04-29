package detection

import "testing"

func TestDetectorDetectValue(t *testing.T) {
	detector := NewDetector(nil, true)
	if match := detector.DetectValue("greg@example.com"); match.Strategy != Email {
		t.Fatalf("unexpected match: %+v", match)
	}
}
