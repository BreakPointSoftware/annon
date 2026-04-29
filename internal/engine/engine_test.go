package engine

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestEngineJSONBytesFallback(t *testing.T) {
	engine := New(DefaultConfig())
	assert.Equal(t, string(jsonFallbackBytes), string(engine.JSONBytes([]byte(`{"email":`))))
}

func TestEngineYAMLBytesFallback(t *testing.T) {
	engine := New(DefaultConfig())
	assert.Equal(t, string(yamlFallbackBytes), string(engine.YAMLBytes([]byte("email: ["))))
}

func TestEngineDataRedactsStringInput(t *testing.T) {
	engine := New(DefaultConfig())
	assert.Equal(t, "g***@example.com", engine.Data("greg@example.com"))
}

func TestEngineDataHandlesUnsupportedValuesSafely(t *testing.T) {
	engine := New(DefaultConfig())

	testCases := []struct {
		name     string
		input    any
		expected any
	}{
		{name: "func falls back to redacted string", input: func() {}, expected: "[REDACTED]"},
		{name: "channel falls back to redacted string", input: make(chan int), expected: "[REDACTED]"},
		{name: "nil stays nil", input: nil, expected: nil},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, engine.Data(testCase.input))
		})
	}
}

func TestEngineDataPreservesSafePrimitiveValues(t *testing.T) {
	engine := New(DefaultConfig())
	assert.Equal(t, 123, engine.Data(123))
	assert.Equal(t, true, engine.Data(true))
}

func TestEngineJSONAndYAMLAlwaysReturnValidPayloads(t *testing.T) {
	engine := New(DefaultConfig())

	jsonBytes := engine.JSON(func() {})
	var jsonValue any
	assert.NoError(t, json.Unmarshal(jsonBytes, &jsonValue))

	yamlBytes := engine.YAML(func() {})
	var yamlValue any
	assert.NoError(t, yaml.Unmarshal(yamlBytes, &yamlValue))
}

func TestRecoverFallbackHelpers(t *testing.T) {
	var recoveredValue any = "value"
	func() {
		defer recoverToValue("secret", &recoveredValue)
		panic("boom")
	}()
	assert.Equal(t, "[REDACTED]", recoveredValue)

	var recoveredJSON []byte
	func() {
		defer recoverToJSONFallback(&recoveredJSON)
		panic("boom")
	}()
	assert.Equal(t, string(jsonFallbackBytes), string(recoveredJSON))

	var recoveredYAML []byte
	func() {
		defer recoverToYAMLFallback(&recoveredYAML)
		panic("boom")
	}()
	assert.Equal(t, string(yamlFallbackBytes), string(recoveredYAML))
}
