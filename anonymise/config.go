package anonymise

import (
	"github.com/your-org/annon/detection"
	"github.com/your-org/annon/preservation"
	strategypkg "github.com/your-org/annon/strategy"
)

type Config struct {
	UseTags           bool
	UseFieldDetection bool
	UseValueDetection bool
	Detector          detection.Detector
	FieldRules        []detection.Rule
	Strategies        map[string]strategypkg.Strategy
	Preservation      preservation.Config
}

func defaultConfig() Config {
	strategies := map[string]strategypkg.Strategy{}
	for _, strategyImpl := range strategypkg.DefaultStrategies() {
		strategies[strategyImpl.Name()] = strategyImpl
	}
	return Config{
		UseTags:           true,
		UseFieldDetection: true,
		UseValueDetection: false,
		Strategies:        strategies,
		Preservation:      preservation.Default(),
	}
}
