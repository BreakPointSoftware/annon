package detection

type Detector interface {
	Detect(fieldName string, value any) Match
}

type FieldDetector interface {
	DetectField(fieldName string) Match
}

type CompiledDetector struct {
	strong               map[string]Match
	fallback             map[string]Match
	contains             []compiledContainsRule
	valueDetector        ValueDetector
	enableValueDetection bool
}

func NewCompiledDetector(rules []Rule, valueDetector ValueDetector, enableValueDetection bool) *CompiledDetector {
	d := &CompiledDetector{
		strong:               map[string]Match{},
		fallback:             map[string]Match{},
		valueDetector:        valueDetector,
		enableValueDetection: enableValueDetection,
	}
	for _, rule := range rules {
		d.addRule(rule)
	}
	return d
}

func (d *CompiledDetector) Detect(fieldName string, value any) Match {
	if match := d.DetectField(fieldName); match.Found() {
		return match
	}
	if d.enableValueDetection && d.valueDetector != nil {
		if match := d.valueDetector.DetectValue(value); match.Found() {
			return match
		}
	}
	return NoMatchResult()
}

func (d *CompiledDetector) DetectValue(value any) Match {
	if d.valueDetector == nil {
		return NoMatchResult()
	}
	return d.valueDetector.DetectValue(value)
}
