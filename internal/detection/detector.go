package detection

type Detector struct {
	strong               map[string]Match
	fallback             map[string]Match
	contains             []compiledContainsRule
	enableValueDetection bool
}

func NewDetector(rules []Rule, enableValueDetection bool) *Detector {
	detector := &Detector{
		strong:               map[string]Match{},
		fallback:             map[string]Match{},
		enableValueDetection: enableValueDetection,
	}

	for _, rule := range rules {
		detector.addRule(rule)
	}

	return detector
}

func (detector *Detector) Detect(fieldName string, value any) Match {
	if match := detector.DetectField(fieldName); match.Found() {
		return match
	}

	if detector.enableValueDetection {
		if match := detector.DetectValue(value); match.Found() {
			return match
		}
	}

	return NoMatchResult()
}

func (detector *Detector) DetectValue(value any) Match {
	stringValue, isString := value.(string)
	if !isString {
		return NoMatchResult()
	}

	if IsEmail(stringValue) {
		return Match{Strategy: Email, Confidence: ValuePatternMatch, MatchedBy: "value:email"}
	}

	if IsUKPhoneNumber(stringValue) {
		return Match{Strategy: Phone, Confidence: ValuePatternMatch, MatchedBy: "value:phone"}
	}

	if IsUKPostcode(stringValue) {
		return Match{Strategy: Postcode, Confidence: ValuePatternMatch, MatchedBy: "value:postcode"}
	}

	if IsVehicleRegistration(stringValue) {
		return Match{Strategy: VehicleRegistration, Confidence: ValuePatternMatch, MatchedBy: "value:vehicleRegistration"}
	}

	return NoMatchResult()
}
