package redact

import core "github.com/BreakPointSoftware/annon/internal/redactcore"

type Config = core.Config
type EmailConfig = core.EmailConfig
type PhoneConfig = core.PhoneConfig
type NameConfig = core.NameConfig
type PostcodeConfig = core.PostcodeConfig
type VehicleRegistrationConfig = core.VehicleRegistrationConfig

type Option func(*Config) error

type Redactor struct { config Config }
