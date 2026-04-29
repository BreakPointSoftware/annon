package redact

import annon "github.com/BreakPointSoftware/annon"

var (
	jsonFallback = []byte(`{"redaction_error":true}`)
	yamlFallback = []byte("redaction_error: true\n")
)

func JSON(input any) []byte {
	defer func() {
		if recover() != nil {
		}
	}()

	outputBytes, err := annon.JSON(input)
	if err != nil {
		return append([]byte(nil), jsonFallback...)
	}

	return outputBytes
}

func YAML(input any) []byte {
	defer func() {
		if recover() != nil {
		}
	}()

	outputBytes, err := annon.YAML(input)
	if err != nil {
		return append([]byte(nil), yamlFallback...)
	}

	return outputBytes
}

func JSONBytes(input []byte) []byte {
	defer func() {
		if recover() != nil {
		}
	}()

	outputBytes, err := annon.FromJSON(input)
	if err != nil {
		return append([]byte(nil), jsonFallback...)
	}

	return outputBytes
}

func YAMLBytes(input []byte) []byte {
	defer func() {
		if recover() != nil {
		}
	}()

	outputBytes, err := annon.FromYAML(input)
	if err != nil {
		return append([]byte(nil), yamlFallback...)
	}

	return outputBytes
}
