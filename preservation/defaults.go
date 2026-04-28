package preservation

func Default() Config {
	return Config{
		RedactionText: "[REDACTED]",
		RedactChar:    '*',
		Email: EmailConfig{
			KeepLocalPrefix: 1,
			KeepDomain:      true,
		},
		Phone: PhoneConfig{
			KeepLast: 4,
		},
		Name: NameConfig{
			KeepPrefix: 1,
		},
		Postcode: PostcodeConfig{
			KeepOutward: true,
		},
		VehicleRegistration: VehicleRegistrationConfig{
			KeepPrefix: 4,
		},
	}
}
