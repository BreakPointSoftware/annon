package detection

import "github.com/BreakPointSoftware/annon/internal/support/normalise"

func Normalise(input string) string {
	return normalise.FieldName(input)
}
