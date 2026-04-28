package preservation

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
