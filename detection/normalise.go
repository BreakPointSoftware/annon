package detection

import (
	"strings"
	"unicode"
)

func Normalise(input string) string {
	var b strings.Builder
	b.Grow(len(input))
	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}
