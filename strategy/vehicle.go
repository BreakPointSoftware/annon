package strategy

import "strings"

type VehicleStrategy struct{}

func (VehicleStrategy) Name() string { return "vehicleRegistration" }

func (VehicleStrategy) Anonymise(value any, ctx Context) (any, error) {
	s, ok := value.(string)
	if !ok {
		return value, nil
	}
	compact := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(s), " ", ""))
	keep := clampKeep(len(compact), ctx.Preservation.VehicleRegistration.KeepPrefix)
	if compact == "" {
		return s, nil
	}
	masked := keepPrefix(compact, keep, ctx.Preservation.RedactChar)
	if len(compact) > 4 && keep >= 4 {
		return masked[:4] + " " + masked[4:], nil
	}
	return masked, nil
}
