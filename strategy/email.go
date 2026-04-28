package strategy

import "strings"

type EmailStrategy struct{}

func (EmailStrategy) Name() string { return "email" }

func (EmailStrategy) Anonymise(value any, ctx Context) (any, error) {
	s, ok := value.(string)
	if !ok {
		return value, nil
	}
	parts := strings.SplitN(s, "@", 2)
	if len(parts) != 2 {
		return s, nil
	}
	localPart, domain := parts[0], parts[1]
	keep := clampKeep(len(localPart), ctx.Preservation.Email.KeepLocalPrefix)
	masked := keepPrefix(localPart, keep, ctx.Preservation.RedactChar)
	if ctx.Preservation.Email.KeepDomain {
		return masked + "@" + domain, nil
	}
	return masked, nil
}
