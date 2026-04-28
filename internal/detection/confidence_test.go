package detection

import "testing"

func TestConfidenceOrdering(t *testing.T) {
	if !(NoMatch < ValuePatternMatch && ValuePatternMatch < ContainsMatch && ContainsMatch < FallbackMatch && FallbackMatch < StrongMatch && StrongMatch < ExplicitMatch) {
		t.Fatal("confidence ordering does not match expected precedence")
	}
}
