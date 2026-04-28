package detection

import "testing"

func TestPatternHelpers(t *testing.T) {
	if !IsEmail("greg.bryant+test@example.com") || IsEmail("not-an-email") {
		t.Fatal("email detection mismatch")
	}
	if !IsUKPhoneNumber("07700 900123") || !IsUKPhoneNumber("+44 7700 900123") || !IsUKPhoneNumber("0044 7700 900123") || IsUKPhoneNumber("12345") {
		t.Fatal("phone detection mismatch")
	}
	if !IsUKPostcode("TN9 1XA") || !IsUKPostcode("tn91xa") || IsUKPostcode("ABC123") {
		t.Fatal("postcode detection mismatch")
	}
	if !IsVehicleRegistration("AB12 CDE") || !IsVehicleRegistration("A123 ABC") || IsVehicleRegistration("HELLO") {
		t.Fatal("vehicle registration detection mismatch")
	}
}
