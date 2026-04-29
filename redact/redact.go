package redact

import "github.com/BreakPointSoftware/annon/internal/redactcore"

func defaultConfig() Config { return redactcore.DefaultConfig() }

func Email(value string, opts ...Option) string { r, _ := New(opts...); return r.Email(value) }
func Phone(value string, opts ...Option) string { r, _ := New(opts...); return r.Phone(value) }
func Postcode(value string, opts ...Option) string { r, _ := New(opts...); return r.Postcode(value) }
func VehicleRegistration(value string, opts ...Option) string { r, _ := New(opts...); return r.VehicleRegistration(value) }
func Name(value string, opts ...Option) string { r, _ := New(opts...); return r.Name(value) }
func FirstName(value string, opts ...Option) string { r, _ := New(opts...); return r.FirstName(value) }
func Surname(value string, opts ...Option) string { r, _ := New(opts...); return r.Surname(value) }
func Redact(value string, opts ...Option) string { r, _ := New(opts...); return r.Redact(value) }

func (r *Redactor) Email(value string) string { return redactcore.Email(value, r.config) }
func (r *Redactor) Phone(value string) string { return redactcore.Phone(value, r.config) }
func (r *Redactor) Postcode(value string) string { return redactcore.Postcode(value, r.config) }
func (r *Redactor) VehicleRegistration(value string) string { return redactcore.VehicleRegistration(value, r.config) }
func (r *Redactor) Name(value string) string { return redactcore.Name(value, r.config) }
func (r *Redactor) FirstName(value string) string { return redactcore.FirstName(value, r.config) }
func (r *Redactor) Surname(value string) string { return redactcore.Surname(value, r.config) }
func (r *Redactor) Redact(value string) string { return redactcore.Redact(value, r.config) }
