package anonymise

import "github.com/BreakPointSoftware/annon/encoder"

func (a *Anonymiser) JSON(input any) ([]byte, error) {
	neutral, err := a.walker.BlobFromValue(input, "json")
	if err != nil {
		return nil, err
	}
	return encoder.EncodeJSON(neutral)
}

func (a *Anonymiser) FromJSON(input []byte) ([]byte, error) {
	decoded, err := encoder.DecodeJSON(input)
	if err != nil {
		return nil, err
	}
	neutral, err := a.walker.BlobFromNeutral(decoded)
	if err != nil {
		return nil, err
	}
	return encoder.EncodeJSON(neutral)
}
