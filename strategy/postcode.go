package strategy

import "strings"

type PostcodeStrategy struct{}

func (PostcodeStrategy) Name() string { return "postcode" }

func (PostcodeStrategy) Anonymise(value any, ctx Context) (any, error) {
	s, ok := value.(string)
	if !ok {
		return value, nil
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return s, nil
	}
	compact := strings.ToUpper(strings.ReplaceAll(s, " ", ""))
	if len(compact) < 3 {
		return s, nil
	}
	if !ctx.Preservation.Postcode.KeepOutward {
		return repeatRune(ctx.Preservation.RedactChar, len(compact)), nil
	}
	outward, inward := splitPostcode(compact)
	if outward == "" || inward == "" {
		return s, nil
	}
	return outward + " " + repeatRune(ctx.Preservation.RedactChar, len(inward)), nil
}

func splitPostcode(compact string) (string, string) {
	if len(compact) <= 3 {
		return "", ""
	}
	return compact[:len(compact)-3], compact[len(compact)-3:]
}
