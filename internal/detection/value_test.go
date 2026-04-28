package detection

import "testing"

func TestPatternValueDetector(t *testing.T) {
	d := PatternValueDetector{}
	if match := d.DetectValue("greg@example.com"); match.Strategy != Email {
		t.Fatalf("unexpected match: %+v", match)
	}
}
