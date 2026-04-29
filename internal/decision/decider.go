package decision

import (
	"fmt"

	"github.com/BreakPointSoftware/annon/internal/detection"
	"github.com/BreakPointSoftware/annon/internal/redactcore"
)

type Decider struct {
	config Config
}

func New(config Config) *Decider {
	return &Decider{config: config}
}

func (d *Decider) Decide(fieldName, tag string, value any) (Result, error) {
	if d.config.UseTags {
		parsedTag := parseTag(tag)

		if err := d.validateTag(parsedTag); err != nil {
			return Result{}, err
		}

		if parsedTag.skip {
			return Result{Skip: true}, nil
		}

		if parsedTag.remove {
			return Result{Remove: true, StrategyName: parsedTag.strategyName}, nil
		}

		if parsedTag.strategyName != "" && !parsedTag.auto {
			return Result{StrategyName: parsedTag.strategyName}, nil
		}

		if parsedTag.auto {
			match := d.detect(fieldName, value)
			if match.Found() {
				return Result{StrategyName: string(match.Strategy)}, nil
			}

			return Result{}, nil
		}
	}

	match := d.detect(fieldName, value)
	if match.Found() {
		if match.Strategy == detection.Remove {
			return Result{Remove: true, StrategyName: string(match.Strategy)}, nil
		}

		return Result{StrategyName: string(match.Strategy)}, nil
	}

	return Result{}, nil
}

func (d *Decider) detect(fieldName string, value any) detection.Match {
	if d.config.Detector == nil {
		return detection.NoMatchResult()
	}

	if d.config.UseFieldDetection && d.config.UseValueDetection {
		return d.config.Detector.Detect(fieldName, value)
	}

	if d.config.UseFieldDetection {
		return d.config.Detector.DetectField(fieldName)
	}

	if d.config.UseValueDetection {
		return d.config.Detector.DetectValue(value)
	}

	return detection.NoMatchResult()
}

func (d *Decider) validateTag(parsedTag parsedTag) error {
	if parsedTag.empty || parsedTag.skip || parsedTag.auto || parsedTag.remove {
		return nil
	}

	if !redactcore.SupportedStrategy(parsedTag.strategyName) {
		return fmt.Errorf("unknown anonymise tag strategy %q", parsedTag.strategyName)
	}

	return nil
}
