package strategy

import (
	"testing"

	"github.com/your-org/annon/preservation"
)

func TestNameStrategy(t *testing.T) {
	ctx := Context{Preservation: preservation.Default()}
	got, _ := (NameStrategy{strategyName: "name"}).Anonymise("Greg Bryant", ctx)
	if got != "G*** B*****" {
		t.Fatalf("unexpected name: %v", got)
	}
	ctx.Preservation.Name.KeepPrefix = 2
	got, _ = (NameStrategy{strategyName: "firstName"}).Anonymise("Greg", ctx)
	if got != "Gr**" {
		t.Fatalf("unexpected first name: %v", got)
	}
}
