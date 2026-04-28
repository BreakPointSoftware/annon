package annon

import "github.com/BreakPointSoftware/annon/internal/encode"

func (a *Anonymiser) JSON(input any) ([]byte, error) {
	neutral, err := a.walker.BlobFromValue(input, "json")
	if err != nil { return nil, err }
	return encode.EncodeJSON(neutral)
}

func (a *Anonymiser) FromJSON(input []byte) ([]byte, error) {
	decoded, err := encode.DecodeJSON(input)
	if err != nil { return nil, err }
	neutral, err := a.walker.BlobFromNeutral(decoded)
	if err != nil { return nil, err }
	return encode.EncodeJSON(neutral)
}
