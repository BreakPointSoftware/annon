package detection

import support "github.com/BreakPointSoftware/annon/internal/support/normalise"

func Normalise(input string) string {
	return support.FieldName(input)
}
