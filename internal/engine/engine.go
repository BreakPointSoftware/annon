package engine

import (
	copyinternal "github.com/BreakPointSoftware/annon/internal/copy"
	"github.com/BreakPointSoftware/annon/internal/decision"
	"github.com/BreakPointSoftware/annon/internal/detection"
	"github.com/BreakPointSoftware/annon/internal/encode"
	"github.com/BreakPointSoftware/annon/internal/output"
)

type Engine struct {
	config        Config
	copier        *copyinternal.Copier
	outputBuilder *output.Builder
}

func New(config Config) *Engine {
	var detector *detection.Detector
	if config.UseFieldDetection || config.UseValueDetection {
		rules := []detection.Rule(nil)
		if config.UseFieldDetection {
			rules = append(rules, detection.DefaultRules()...)
			rules = append(rules, config.FieldRules...)
		}
		detector = detection.NewDetector(rules, config.UseValueDetection)
	}

	decisionConfig := decision.Config{
		UseTags:           config.UseTags,
		UseFieldDetection: config.UseFieldDetection,
		UseValueDetection: config.UseValueDetection,
		Detector:          detector,
		Preservation:      config.Preservation,
	}

	cache := copyinternal.NewTypeCache()
	decider := decision.New(decisionConfig)

	return &Engine{
		config:        config,
		copier:        copyinternal.New(decisionConfig, decider, cache),
		outputBuilder: output.New(decisionConfig, decider, cache),
	}
}

func (e *Engine) Data(input any) (result any) {
	defer recoverToValue(input, &result)

	if stringInput, isString := input.(string); isString {
		return redactString(stringInput, e.config)
	}

	redactedValue, err := e.copier.Copy(input)
	if err != nil {
		return fallbackValue(input)
	}

	return redactedValue
}

func (e *Engine) JSON(input any) (result []byte) {
	defer recoverToJSONFallback(&result)

	neutralOutput, err := e.outputBuilder.OutputFromValue(input, "json")
	if err != nil {
		return cloneJSONFallback()
	}

	encodedBytes, err := encode.EncodeJSON(neutralOutput)
	if err != nil {
		return cloneJSONFallback()
	}

	return encodedBytes
}

func (e *Engine) YAML(input any) (result []byte) {
	defer recoverToYAMLFallback(&result)

	neutralOutput, err := e.outputBuilder.OutputFromValue(input, "yaml")
	if err != nil {
		return cloneYAMLFallback()
	}

	encodedBytes, err := encode.EncodeYAML(neutralOutput)
	if err != nil {
		return cloneYAMLFallback()
	}

	return encodedBytes
}

func (e *Engine) JSONBytes(input []byte) (result []byte) {
	defer recoverToJSONFallback(&result)

	decodedValue, err := encode.DecodeJSON(input)
	if err != nil {
		return cloneJSONFallback()
	}

	neutralOutput, err := e.outputBuilder.OutputFromNeutral(decodedValue)
	if err != nil {
		return cloneJSONFallback()
	}

	encodedBytes, err := encode.EncodeJSON(neutralOutput)
	if err != nil {
		return cloneJSONFallback()
	}

	return encodedBytes
}

func (e *Engine) YAMLBytes(input []byte) (result []byte) {
	defer recoverToYAMLFallback(&result)

	decodedValue, err := encode.DecodeYAML(input)
	if err != nil {
		return cloneYAMLFallback()
	}

	neutralOutput, err := e.outputBuilder.OutputFromNeutral(decodedValue)
	if err != nil {
		return cloneYAMLFallback()
	}

	encodedBytes, err := encode.EncodeYAML(neutralOutput)
	if err != nil {
		return cloneYAMLFallback()
	}

	return encodedBytes
}

func (e *Engine) String(input string) (result string) {
	defer func() {
		if recover() != nil {
			result = redactString(input, Config{Preservation: e.config.Preservation})
		}
	}()

	return redactString(input, e.config)
}
