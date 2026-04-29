package walk

import (
	"fmt"

	"github.com/BreakPointSoftware/annon/internal/detection"
	"github.com/BreakPointSoftware/annon/internal/redactcore"
)

type decision struct {
	skip         bool
	remove       bool
	strategyName string
}

type Decider struct {
	cfg Config
}

func NewDecider(cfg Config) *Decider {
	return &Decider{cfg: cfg}
}

func (d *Decider) Decide(fieldName, tag string, value any) (decision, error) {
	if d.cfg.UseTags {
		parsed := parseTag(tag)
		if err := d.validateTag(parsed); err != nil {
			return decision{}, err
		}
		if parsed.skip {
			return decision{skip: true}, nil
		}
		if parsed.remove {
			return decision{remove: true, strategyName: parsed.strategyName}, nil
		}
		if parsed.strategyName != "" && !parsed.auto {
			return decision{strategyName: parsed.strategyName}, nil
		}
		if parsed.auto {
			match := d.detect(fieldName, value)
			if match.Found() {
				return decision{strategyName: string(match.Strategy)}, nil
			}
			return decision{}, nil
		}
	}

	match := d.detect(fieldName, value)
	if match.Found() {
		if match.Strategy == detection.Remove {
			return decision{remove: true, strategyName: string(match.Strategy)}, nil
		}
		return decision{strategyName: string(match.Strategy)}, nil
	}
	return decision{}, nil
}

func (d *Decider) detect(fieldName string, value any) detection.Match {
	if d.cfg.Detector == nil {
		return detection.NoMatchResult()
	}

	if d.cfg.UseFieldDetection && d.cfg.UseValueDetection {
		return d.cfg.Detector.Detect(fieldName, value)
	}

	if d.cfg.UseFieldDetection {
		return d.cfg.Detector.DetectField(fieldName)
	}

	if d.cfg.UseValueDetection {
		return d.cfg.Detector.DetectValue(value)
	}

	return detection.NoMatchResult()
}

func (d *Decider) validateTag(parsed parsedTag) error {
	if parsed.empty || parsed.skip || parsed.auto || parsed.remove {
		return nil
	}
	if !redactcore.SupportedStrategy(parsed.strategyName) {
		return fmt.Errorf("unknown anonymise tag strategy %q", parsed.strategyName)
	}
	return nil
}
