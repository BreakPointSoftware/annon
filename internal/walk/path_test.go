package walk

import "testing"

func TestPath(t *testing.T) {
	var p Path
	p = p.Append("customer").Append("email")
	if p.String() != "customer.email" { t.Fatalf("unexpected path: %q", p.String()) }
}
