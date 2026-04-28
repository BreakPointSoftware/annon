package detection

import "testing"

func TestRuleHelpers(t *testing.T) {
	strong := StrongRule(Email, "email")
	if strong.Type != Strong || strong.Strategy != Email || len(strong.Fields) != 1 {
		t.Fatalf("unexpected strong rule: %+v", strong)
	}
	contains := ContainsRule(Name, []string{"name"}, []string{"username"})
	if contains.Type != Contains || len(contains.Exclude) != 1 {
		t.Fatalf("unexpected contains rule: %+v", contains)
	}
}
