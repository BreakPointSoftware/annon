package anonymise

import (
	internaldetection "github.com/BreakPointSoftware/annon/internal/detection"
	"github.com/BreakPointSoftware/annon/internal/redactcore"
)

type PreservationConfig = redactcore.Config
type EmailConfig = redactcore.EmailConfig
type PhoneConfig = redactcore.PhoneConfig
type NameConfig = redactcore.NameConfig
type PostcodeConfig = redactcore.PostcodeConfig
type VehicleRegistrationConfig = redactcore.VehicleRegistrationConfig

type FieldStrategy = internaldetection.Strategy
type MatchType = internaldetection.MatchType
type FieldRule = internaldetection.Rule

const (
	EmailStrategy               FieldStrategy = internaldetection.Email
	PhoneStrategy               FieldStrategy = internaldetection.Phone
	PostcodeStrategy            FieldStrategy = internaldetection.Postcode
	NameStrategy                FieldStrategy = internaldetection.Name
	FirstNameStrategy           FieldStrategy = internaldetection.FirstName
	SurnameStrategy             FieldStrategy = internaldetection.Surname
	VehicleRegistrationStrategy FieldStrategy = internaldetection.VehicleRegistration
	RedactStrategy              FieldStrategy = internaldetection.Redact

	StrongMatchType   MatchType = internaldetection.Strong
	FallbackMatchType MatchType = internaldetection.Fallback
	ContainsMatchType MatchType = internaldetection.Contains
)

func StrongRule(strategy FieldStrategy, fields ...string) FieldRule {
	return internaldetection.StrongRule(internaldetection.Strategy(strategy), fields...)
}

func FallbackRule(strategy FieldStrategy, fields ...string) FieldRule {
	return internaldetection.FallbackRule(internaldetection.Strategy(strategy), fields...)
}

func ContainsRule(strategy FieldStrategy, fields []string, exclude []string) FieldRule {
	return internaldetection.ContainsRule(internaldetection.Strategy(strategy), fields, exclude)
}

type config struct {
	UseTags           bool
	UseFieldDetection bool
	UseValueDetection bool
	FieldRules        []internaldetection.Rule
	Preservation      redactcore.Config
}

func defaultConfig() config {
	return config{UseTags: true, UseFieldDetection: true, UseValueDetection: false, Preservation: redactcore.DefaultConfig()}
}
