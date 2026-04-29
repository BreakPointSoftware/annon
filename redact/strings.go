package redact

import "github.com/BreakPointSoftware/annon/internal/detection"

func String(input string) string {
	switch {
	case detection.IsEmail(input):
		return Email(input)
	case detection.IsUKPhoneNumber(input):
		return Phone(input)
	case detection.IsUKPostcode(input):
		return Postcode(input)
	case detection.IsVehicleRegistration(input):
		return VehicleRegistration(input)
	default:
		return Redact(input)
	}
}
