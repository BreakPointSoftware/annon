package preservation

import "testing"

func TestConfigStoresValues(t *testing.T) {
	cfg := Config{
		RedactionText: "[hidden]",
		RedactChar:    'x',
		Email:         EmailConfig{KeepLocalPrefix: 2, KeepDomain: true},
		Phone:         PhoneConfig{KeepLast: 3},
		Name:          NameConfig{KeepPrefix: 2},
		Postcode:      PostcodeConfig{KeepOutward: false},
		VehicleRegistration: VehicleRegistrationConfig{
			KeepPrefix: 2,
		},
	}

	if cfg.RedactionText != "[hidden]" || cfg.RedactChar != 'x' {
		t.Fatalf("unexpected root config: %+v", cfg)
	}
	if cfg.Email.KeepLocalPrefix != 2 || !cfg.Email.KeepDomain {
		t.Fatalf("unexpected email config: %+v", cfg.Email)
	}
}
