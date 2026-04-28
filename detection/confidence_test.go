package detection

import "testing"

func TestConfidenceOrdering(t *testing.T) {
	if !(NoMatch < ValuePatternMatch && ValuePatternMatch < ContainsMatch && ContainsMatch < FallbackMatch && FallbackMatch < StrongMatch && StrongMatch < ExplicitMatch) {
		t.Fatal("confidence ordering does not match expected precedence")
	}
}

func TestMatchFound(t *testing.T) {
	if NoMatchResult().Found() {
		t.Fatal("no match should not be found")
	}
	if !(Match{Strategy: Email, Confidence: StrongMatch}).Found() {
		t.Fatal("strong match should be found")
	}
}
