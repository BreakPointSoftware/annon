package redactcore

import "strings"

func Redact(value string, cfg Config) string {
	if strings.TrimSpace(cfg.RedactionText) == "" {
		return "[REDACTED]"
	}
	return cfg.RedactionText
}

func Apply(strategyName string, value any, cfg Config) (any, error) {
	s, ok := value.(string)
	if !ok {
		return value, nil
	}
	switch strategyName {
	case "email":
		return Email(s, cfg), nil
	case "phone":
		return Phone(s, cfg), nil
	case "postcode":
		return Postcode(s, cfg), nil
	case "name":
		return Name(s, cfg), nil
	case "firstName":
		return FirstName(s, cfg), nil
	case "surname":
		return Surname(s, cfg), nil
	case "vehicleRegistration":
		return VehicleRegistration(s, cfg), nil
	case "redact":
		return Redact(s, cfg), nil
	default:
		return value, nil
	}
}
