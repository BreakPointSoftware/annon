package detection

type ValueDetector interface {
	DetectValue(value any) Match
}

type PatternValueDetector struct{}

func (PatternValueDetector) DetectValue(value any) Match {
	s, ok := value.(string)
	if !ok {
		return NoMatchResult()
	}
	if IsEmail(s) {
		return Match{Strategy: Email, Confidence: ValuePatternMatch, MatchedBy: "value:email"}
	}
	if IsUKPhoneNumber(s) {
		return Match{Strategy: Phone, Confidence: ValuePatternMatch, MatchedBy: "value:phone"}
	}
	if IsUKPostcode(s) {
		return Match{Strategy: Postcode, Confidence: ValuePatternMatch, MatchedBy: "value:postcode"}
	}
	if IsVehicleRegistration(s) {
		return Match{Strategy: VehicleRegistration, Confidence: ValuePatternMatch, MatchedBy: "value:vehicleRegistration"}
	}
	return NoMatchResult()
}
