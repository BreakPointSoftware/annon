package redact

import "github.com/BreakPointSoftware/annon/internal/redactcore"

type Config = redactcore.Config
type EmailConfig = redactcore.EmailConfig
type PhoneConfig = redactcore.PhoneConfig
type NameConfig = redactcore.NameConfig
type PostcodeConfig = redactcore.PostcodeConfig
type VehicleRegistrationConfig = redactcore.VehicleRegistrationConfig

type Option func(*Config) error

type Redactor struct { config Config }
