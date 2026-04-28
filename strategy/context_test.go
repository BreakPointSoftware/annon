package strategy

import (
	"testing"

	"github.com/BreakPointSoftware/annon/preservation"
)

func TestContextCarriesPreservation(t *testing.T) {
	ctx := Context{Preservation: preservation.Default()}
	if ctx.Preservation.RedactionText != "[REDACTED]" {
		t.Fatalf("unexpected context: %+v", ctx)
	}
}
