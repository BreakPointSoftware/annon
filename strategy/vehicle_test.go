package strategy

import (
	"testing"

	"github.com/your-org/annon/preservation"
)

func TestVehicleStrategy(t *testing.T) {
	ctx := Context{Preservation: preservation.Default()}
	got, _ := (VehicleStrategy{}).Anonymise("AB12 CDE", ctx)
	if got != "AB12 ***" {
		t.Fatalf("unexpected vehicle registration: %v", got)
	}
	got, _ = (VehicleStrategy{}).Anonymise("ab12cde", ctx)
	if got != "AB12 ***" {
		t.Fatalf("unexpected compact vehicle registration: %v", got)
	}
	ctx.Preservation.VehicleRegistration.KeepPrefix = 2
	got, _ = (VehicleStrategy{}).Anonymise("AB12 CDE", ctx)
	if got != "AB*****" {
		t.Fatalf("unexpected reduced-prefix vehicle registration: %v", got)
	}
}
