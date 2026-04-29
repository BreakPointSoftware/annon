package decision

import (
	"github.com/BreakPointSoftware/annon/internal/detection"
	"github.com/BreakPointSoftware/annon/internal/redactcore"
)

type Config struct {
	UseTags           bool
	UseFieldDetection bool
	UseValueDetection bool
	Detector          *detection.Detector
	Preservation      redactcore.Config
}

type Result struct {
	Skip         bool
	Remove       bool
	StrategyName string
}
