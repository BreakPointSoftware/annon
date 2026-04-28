package strategy

import (
	"testing"

	"github.com/your-org/annon/preservation"
)

func TestPhoneStrategy(t *testing.T) {
	ctx := Context{Preservation: preservation.Default()}
	got, _ := (PhoneStrategy{}).Anonymise("07700 900123", ctx)
	if got != "*******0123" {
		t.Fatalf("unexpected phone: %v", got)
	}
	got, _ = (PhoneStrategy{}).Anonymise("+44 7700 900123", ctx)
	if got != "*********0123" {
		t.Fatalf("unexpected international phone: %v", got)
	}
}
