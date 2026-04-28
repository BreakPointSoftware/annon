package strategy

import (
	"testing"

	"github.com/BreakPointSoftware/annon/preservation"
)

func TestPostcodeStrategy(t *testing.T) {
	ctx := Context{Preservation: preservation.Default()}
	got, _ := (PostcodeStrategy{}).Anonymise("TN9 1XA", ctx)
	if got != "TN9 ***" {
		t.Fatalf("unexpected postcode: %v", got)
	}
	got, _ = (PostcodeStrategy{}).Anonymise("tn91xa", ctx)
	if got != "TN9 ***" {
		t.Fatalf("unexpected compact postcode: %v", got)
	}
	ctx.Preservation.Postcode.KeepOutward = false
	got, _ = (PostcodeStrategy{}).Anonymise("TN9 1XA", ctx)
	if got != "******" {
		t.Fatalf("unexpected fully redacted postcode: %v", got)
	}
}
