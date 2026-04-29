package annon

import (
	"fmt"

	"github.com/BreakPointSoftware/annon/internal/encode"
)

func (a *Anonymiser) YAML(input any) ([]byte, error) {
	neutral, err := a.outputBuilder.OutputFromValue(input, "yaml")
	if err != nil {
		return nil, err
	}
	return encode.EncodeYAML(neutral)
}

func (a *Anonymiser) FromYAML(input []byte) ([]byte, error) {
	decoded, err := encode.DecodeYAML(input)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidYAML, err)
	}
	neutral, err := a.outputBuilder.OutputFromNeutral(decoded)
	if err != nil {
		return nil, err
	}
	return encode.EncodeYAML(neutral)
}
