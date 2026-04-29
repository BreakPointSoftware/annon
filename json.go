package annon

import (
	"fmt"

	"github.com/BreakPointSoftware/annon/internal/encode"
)

func (a *Anonymiser) JSON(input any) ([]byte, error) {
	neutral, err := a.outputBuilder.OutputFromValue(input, "json")
	if err != nil {
		return nil, err
	}
	return encode.EncodeJSON(neutral)
}

func (a *Anonymiser) FromJSON(input []byte) ([]byte, error) {
	decoded, err := encode.DecodeJSON(input)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}
	neutral, err := a.outputBuilder.OutputFromNeutral(decoded)
	if err != nil {
		return nil, err
	}
	return encode.EncodeJSON(neutral)
}
