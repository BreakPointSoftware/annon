package preservation

import "testing"

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.RedactionText != "[REDACTED]" {
		t.Fatalf("unexpected redaction text: %q", cfg.RedactionText)
	}
	if cfg.RedactChar != '*' {
		t.Fatalf("unexpected redact char: %q", cfg.RedactChar)
	}
	if cfg.Email.KeepLocalPrefix != 1 || !cfg.Email.KeepDomain {
		t.Fatalf("unexpected email defaults: %+v", cfg.Email)
	}
	if cfg.Phone.KeepLast != 4 {
		t.Fatalf("unexpected phone defaults: %+v", cfg.Phone)
	}
	if cfg.Name.KeepPrefix != 1 {
		t.Fatalf("unexpected name defaults: %+v", cfg.Name)
	}
	if !cfg.Postcode.KeepOutward {
		t.Fatalf("unexpected postcode defaults: %+v", cfg.Postcode)
	}
	if cfg.VehicleRegistration.KeepPrefix != 4 {
		t.Fatalf("unexpected vehicle defaults: %+v", cfg.VehicleRegistration)
	}
}
