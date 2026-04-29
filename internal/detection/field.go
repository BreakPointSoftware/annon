package detection

import "strings"

type compiledContainsRule struct {
	strategy  Strategy
	fields    []string
	exclude   []string
	matchedBy string
}

func (detector *Detector) addRule(rule Rule) {
	for _, field := range rule.Fields {
		normalised := Normalise(field)
		match := Match{Strategy: rule.Strategy}
		switch rule.Type {
		case Strong:
			match.Confidence = StrongMatch
			match.MatchedBy = "field:strong:" + normalised
			detector.strong[normalised] = match
		case Fallback:
			match.Confidence = FallbackMatch
			match.MatchedBy = "field:fallback:" + normalised
			detector.fallback[normalised] = match
		}
	}
	if rule.Type == Contains {
		compiled := compiledContainsRule{strategy: rule.Strategy, matchedBy: "field:contains"}
		for _, field := range rule.Fields {
			compiled.fields = append(compiled.fields, Normalise(field))
		}
		for _, field := range rule.Exclude {
			compiled.exclude = append(compiled.exclude, Normalise(field))
		}
		detector.contains = append(detector.contains, compiled)
	}
}

func (detector *Detector) DetectField(fieldName string) Match {
	normalised := Normalise(fieldName)
	if match, ok := detector.strong[normalised]; ok {
		return match
	}
	if match, ok := detector.fallback[normalised]; ok {
		return match
	}
	for _, rule := range detector.contains {
		blocked := false
		for _, exclude := range rule.exclude {
			if normalised == exclude || strings.Contains(normalised, exclude) {
				blocked = true
				break
			}
		}
		if blocked {
			continue
		}
		for _, field := range rule.fields {
			if strings.Contains(normalised, field) {
				return Match{Strategy: rule.strategy, Confidence: ContainsMatch, MatchedBy: rule.matchedBy + ":" + field}
			}
		}
	}
	return NoMatchResult()
}
