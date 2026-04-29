package engine

import (
	"github.com/BreakPointSoftware/annon/internal/detection"
	"github.com/BreakPointSoftware/annon/internal/redactcore"
)

func redactString(input string, config Config) string {
	switch {
	case detection.IsEmail(input):
		return redactcore.Email(input, config.Preservation)
	case detection.IsUKPhoneNumber(input):
		return redactcore.Phone(input, config.Preservation)
	case detection.IsUKPostcode(input):
		return redactcore.Postcode(input, config.Preservation)
	case detection.IsVehicleRegistration(input):
		return redactcore.VehicleRegistration(input, config.Preservation)
	default:
		return redactcore.Redact(input, config.Preservation)
	}
}
