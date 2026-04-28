package strategy

import (
	"testing"

	"github.com/your-org/annon/preservation"
)

func TestRedactStrategy(t *testing.T) {
	ctx := Context{Preservation: preservation.Default()}
	got, _ := (RedactStrategy{}).Anonymise("secret", ctx)
	if got != "[REDACTED]" {
		t.Fatalf("unexpected redaction: %v", got)
	}
	ctx.Preservation.RedactionText = ""
	got, _ = (RedactStrategy{}).Anonymise("secret", ctx)
	if got != "[REDACTED]" {
		t.Fatalf("unexpected fallback redaction: %v", got)
	}
}
