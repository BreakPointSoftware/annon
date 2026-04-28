package annon

import "github.com/BreakPointSoftware/annon/internal/encode"

func (a *Anonymiser) YAML(input any) ([]byte, error) {
	neutral, err := a.walker.BlobFromValue(input, "yaml")
	if err != nil { return nil, err }
	return encode.EncodeYAML(neutral)
}

func (a *Anonymiser) FromYAML(input []byte) ([]byte, error) {
	decoded, err := encode.DecodeYAML(input)
	if err != nil { return nil, err }
	neutral, err := a.walker.BlobFromNeutral(decoded)
	if err != nil { return nil, err }
	return encode.EncodeYAML(neutral)
}
