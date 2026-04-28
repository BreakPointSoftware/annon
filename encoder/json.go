package encoder

import (
	"bytes"
	"encoding/json"
)

func DecodeJSON(input []byte) (any, error) {
	decoder := json.NewDecoder(bytes.NewReader(input))
	decoder.UseNumber()
	var out any
	if err := decoder.Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func EncodeJSON(input any) ([]byte, error) {
	return json.Marshal(input)
}
