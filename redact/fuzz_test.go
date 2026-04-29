package redact

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func FuzzData(f *testing.F) {
	f.Add("greg@example.com")
	f.Add("07700 900123")
	f.Add("TN9 1XA")
	f.Add("plain text")

	f.Fuzz(func(t *testing.T, input string) {
		_ = Data(input)
	})
}

func FuzzJSONBytes(f *testing.F) {
	f.Add([]byte(`{"email":"greg@example.com"}`))
	f.Add([]byte(`{"email":`))
	f.Add([]byte(`[]`))

	f.Fuzz(func(t *testing.T, input []byte) {
		outputBytes := JSONBytes(input)
		var decoded any
		if err := json.Unmarshal(outputBytes, &decoded); err != nil {
			t.Fatalf("JSONBytes returned invalid json: %v\n%s", err, outputBytes)
		}
	})
}

func FuzzYAMLBytes(f *testing.F) {
	f.Add([]byte("email: greg@example.com\n"))
	f.Add([]byte("email: ["))
	f.Add([]byte("[]\n"))

	f.Fuzz(func(t *testing.T, input []byte) {
		outputBytes := YAMLBytes(input)
		var decoded any
		if err := yaml.Unmarshal(outputBytes, &decoded); err != nil {
			t.Fatalf("YAMLBytes returned invalid yaml: %v\n%s", err, outputBytes)
		}
	})
}
