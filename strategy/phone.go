package strategy

import "strings"

type PhoneStrategy struct{}

func (PhoneStrategy) Name() string { return "phone" }

func (PhoneStrategy) Anonymise(value any, ctx Context) (any, error) {
	s, ok := value.(string)
	if !ok {
		return value, nil
	}
	compact := strings.NewReplacer(" ", "", "-", "", "(", "", ")", "").Replace(s)
	keep := clampKeep(len(compact), ctx.Preservation.Phone.KeepLast)
	if keep == len(compact) {
		return compact, nil
	}
	return repeatRune(ctx.Preservation.RedactChar, len(compact)-keep) + compact[len(compact)-keep:], nil
}
