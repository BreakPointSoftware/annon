package annon

import (
	"github.com/BreakPointSoftware/annon/internal/detection"
	"github.com/BreakPointSoftware/annon/internal/redactcore"
)

type PreservationConfig = redactcore.Config
type EmailConfig = redactcore.EmailConfig
type PhoneConfig = redactcore.PhoneConfig
type NameConfig = redactcore.NameConfig
type PostcodeConfig = redactcore.PostcodeConfig
type VehicleRegistrationConfig = redactcore.VehicleRegistrationConfig

type FieldStrategy = detection.Strategy
type MatchType = detection.MatchType
type FieldRule = detection.Rule

const (
	EmailStrategy               FieldStrategy = detection.Email
	PhoneStrategy               FieldStrategy = detection.Phone
	PostcodeStrategy            FieldStrategy = detection.Postcode
	NameStrategy                FieldStrategy = detection.Name
	FirstNameStrategy           FieldStrategy = detection.FirstName
	SurnameStrategy             FieldStrategy = detection.Surname
	VehicleRegistrationStrategy FieldStrategy = detection.VehicleRegistration
	RedactStrategy              FieldStrategy = detection.Redact

	StrongMatchType   MatchType = detection.Strong
	FallbackMatchType MatchType = detection.Fallback
	ContainsMatchType MatchType = detection.Contains
)

func StrongRule(strategy FieldStrategy, fields ...string) FieldRule {
	return detection.StrongRule(detection.Strategy(strategy), fields...)
}

func FallbackRule(strategy FieldStrategy, fields ...string) FieldRule {
	return detection.FallbackRule(detection.Strategy(strategy), fields...)
}

func ContainsRule(strategy FieldStrategy, fields []string, exclude []string) FieldRule {
	return detection.ContainsRule(detection.Strategy(strategy), fields, exclude)
}

type config struct {
	UseTags           bool
	UseFieldDetection bool
	UseValueDetection bool
	FieldRules        []detection.Rule
	Preservation      redactcore.Config
}

func defaultConfig() config {
	return config{UseTags: true, UseFieldDetection: true, UseValueDetection: false, Preservation: redactcore.DefaultConfig()}
}
