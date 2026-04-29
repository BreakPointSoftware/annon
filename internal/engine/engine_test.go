package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
