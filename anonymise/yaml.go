package anonymise

import "github.com/BreakPointSoftware/annon/encoder"

func (a *Anonymiser) YAML(input any) ([]byte, error) {
	neutral, err := a.walker.BlobFromValue(input, "yaml")
	if err != nil {
		return nil, err
	}
	return encoder.EncodeYAML(neutral)
}

func (a *Anonymiser) FromYAML(input []byte) ([]byte, error) {
	decoded, err := encoder.DecodeYAML(input)
	if err != nil {
		return nil, err
	}
	neutral, err := a.walker.BlobFromNeutral(decoded)
	if err != nil {
		return nil, err
	}
	return encoder.EncodeYAML(neutral)
}
