package detection

import "testing"

func TestRuleHelpers(t *testing.T) {
	if got := StrongRule(Email, "email"); got.Type != Strong {
		t.Fatalf("unexpected rule: %+v", got)
	}
}
