package redactcore

import "testing"

func TestHelpers(t *testing.T) {
	cfg := DefaultConfig()
	if got := Email("greg@example.com", cfg); got != "g***@example.com" {
		t.Fatalf("unexpected email: %q", got)
	}
	if got := Phone("07700 900123", cfg); got != "*******0123" {
		t.Fatalf("unexpected phone: %q", got)
	}
	if got := Postcode("tn91xa", cfg); got != "TN9 ***" {
		t.Fatalf("unexpected postcode: %q", got)
	}
	if got := VehicleRegistration("AB12 CDE", cfg); got != "AB12 ***" {
		t.Fatalf("unexpected vehicle registration: %q", got)
	}
	if got := Name("Greg Bryant", cfg); got != "G*** B*****" {
		t.Fatalf("unexpected name: %q", got)
	}
	if got := Redact("secret", cfg); got != "[REDACTED]" {
		t.Fatalf("unexpected redaction: %q", got)
	}
}
