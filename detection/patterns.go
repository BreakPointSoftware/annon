package detection

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	emailRegex                = regexp.MustCompile(`^[A-Za-z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[A-Za-z0-9-]+(?:\.[A-Za-z0-9-]+)+$`)
	ukPostcodeRegex           = regexp.MustCompile(`^(GIR0AA|[A-Z]{1,2}[0-9][0-9A-Z]?[0-9][A-Z]{2})$`)
	ukPhoneRegex              = regexp.MustCompile(`^(?:0|\+44|0044)7\d{9}$`)
	ukVehicleRegistrationRegex = regexp.MustCompile(`^(?:[A-Z]{2}[0-9]{2}[A-Z]{3}|[A-Z][0-9]{1,3}[A-Z]{3})$`)
)

func IsEmail(value string) bool {
	if len(value) < 3 || !strings.Contains(value, "@") {
		return false
	}
	return emailRegex.MatchString(value)
}

func IsUKPhoneNumber(value string) bool {
	if len(value) < 10 {
		return false
	}
	var b strings.Builder
	b.Grow(len(value))
	for i, r := range value {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
			continue
		}
		if r == '+' && i == 0 {
			b.WriteRune(r)
			continue
		}
		if r == ' ' || r == '-' || r == '(' || r == ')' {
			continue
		}
		return false
	}
	normalised := b.String()
	if len(normalised) < 11 {
		return false
	}
	return ukPhoneRegex.MatchString(normalised)
}

func IsUKPostcode(value string) bool {
	if len(value) < 5 {
		return false
	}
	normalised := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(value), " ", ""))
	return ukPostcodeRegex.MatchString(normalised)
}

func IsVehicleRegistration(value string) bool {
	if len(value) < 6 {
		return false
	}
	normalised := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(value), " ", ""))
	return ukVehicleRegistrationRegex.MatchString(normalised)
}
