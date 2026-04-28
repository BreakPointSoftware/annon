package strategy

import (
	"testing"

	"github.com/BreakPointSoftware/annon/preservation"
)

func TestEmailStrategy(t *testing.T) {
	ctx := Context{Preservation: preservation.Default()}
	got, _ := (EmailStrategy{}).Anonymise("greg.bryant@example.com", ctx)
	if got != "g**********@example.com" {
		t.Fatalf("unexpected masked email: %v", got)
	}
	ctx.Preservation.Email.KeepLocalPrefix = 2
	ctx.Preservation.RedactChar = 'x'
	got, _ = (EmailStrategy{}).Anonymise("greg@example.com", ctx)
	if got != "grxx@example.com" {
		t.Fatalf("unexpected custom masked email: %v", got)
	}
	ctx.Preservation.Email.KeepDomain = false
	got, _ = (EmailStrategy{}).Anonymise("greg@example.com", ctx)
	if got != "grxx" {
		t.Fatalf("unexpected domain-hidden email: %v", got)
	}
}
