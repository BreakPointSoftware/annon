package redactcore

type Config struct {
	RedactionText string
	RedactChar    rune

	Email               EmailConfig
	Phone               PhoneConfig
	Name                NameConfig
	Postcode            PostcodeConfig
	VehicleRegistration VehicleRegistrationConfig
}

type EmailConfig struct {
	KeepLocalPrefix int
	KeepDomain      bool
}

type PhoneConfig struct {
	KeepLast int
}

type NameConfig struct {
	KeepPrefix int
}

type PostcodeConfig struct {
	KeepOutward bool
}

type VehicleRegistrationConfig struct {
	KeepPrefix int
}

func DefaultConfig() Config {
	return Config{
		RedactionText: "[REDACTED]",
		RedactChar:    '*',
		Email: EmailConfig{KeepLocalPrefix: 1, KeepDomain: true},
		Phone: PhoneConfig{KeepLast: 4},
		Name: NameConfig{KeepPrefix: 1},
		Postcode: PostcodeConfig{KeepOutward: true},
		VehicleRegistration: VehicleRegistrationConfig{KeepPrefix: 4},
	}
}

func SupportedStrategy(name string) bool {
	switch name {
	case "email", "phone", "postcode", "name", "firstName", "surname", "vehicleRegistration", "redact":
		return true
	default:
		return false
	}
}
