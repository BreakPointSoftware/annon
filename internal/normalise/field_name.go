package normalise

import (
	"strings"
	"unicode"
)

func FieldName(input string) string {
	var b strings.Builder
	b.Grow(len(input))
	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}
