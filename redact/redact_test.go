package redact

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDirectFunctions(t *testing.T) {
	testCases := []struct {
		name     string
		run      func() string
		expected string
	}{
		{
			name: "email uses default masking",
			run: func() string {
				return Email("greg@example.com")
			},
			expected: "g***@example.com",
		},
		{
			name: "postcode uses default masking",
			run: func() string {
				return Postcode("TN9 1XA")
			},
			expected: "TN9 ***",
		},
		{
			name: "vehicle registration uses default masking",
			run: func() string {
				return VehicleRegistration("AB12 CDE")
			},
			expected: "AB12 ***",
		},
		{
			name: "redact uses default text",
			run: func() string {
				return Redact("secret")
			},
			expected: "[REDACTED]",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, testCase.run())
		})
	}
}

func TestConfiguredRedactor(t *testing.T) {
	redactor, err := New(WithRedactChar('x'))
	require.NoError(t, err)

	assert.Equal(t, "gxxx@example.com", redactor.Email("greg@example.com"))
}

func TestData(t *testing.T) {
	testCases := []struct {
		name     string
		input    any
		assertFn func(t *testing.T, result any)
	}{
		{
			name: "struct input returns redacted copy",
			input: struct {
				Email string `json:"email"`
				Note  string `json:"note" anonymise:"false"`
			}{Email: "greg@example.com", Note: "greg@example.com"},
			assertFn: func(t *testing.T, result any) {
				redactedResult := result.(struct {
					Email string `json:"email"`
					Note  string `json:"note" anonymise:"false"`
				})
				assert.Equal(t, "g***@example.com", redactedResult.Email)
				assert.Equal(t, "greg@example.com", redactedResult.Note)
			},
		},
		{
			name:  "string input uses generic string redaction",
			input: "greg@example.com",
			assertFn: func(t *testing.T, result any) {
				assert.Equal(t, "g***@example.com", result)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := Data(testCase.input)
			testCase.assertFn(t, result)
		})
	}
}

func TestJSONAndYAMLHelpers(t *testing.T) {
	testCases := []struct {
		name     string
		run      func() []byte
		contains string
	}{
		{
			name: "json helper produces redacted json",
			run: func() []byte {
				return JSON(map[string]any{"email": "greg@example.com"})
			},
			contains: `"email":"g***@example.com"`,
		},
		{
			name: "yaml helper produces redacted yaml",
			run: func() []byte {
				return YAML(map[string]any{"email": "greg@example.com"})
			},
			contains: "email: g***@example.com",
		},
		{
			name: "invalid json bytes return fallback",
			run: func() []byte {
				return JSONBytes([]byte(`{"email":`))
			},
			contains: `"redaction_error":true`,
		},
		{
			name: "invalid yaml bytes return fallback",
			run: func() []byte {
				return YAMLBytes([]byte("email: ["))
			},
			contains: "redaction_error: true",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Contains(t, string(testCase.run()), testCase.contains)
		})
	}
}

func TestString(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "email string uses email strategy", input: "greg@example.com", expected: "g***@example.com"},
		{name: "phone string uses phone strategy", input: "07700 900123", expected: "*******0123"},
		{name: "postcode string uses postcode strategy", input: "TN9 1XA", expected: "TN9 ***"},
		{name: "unknown string uses generic redact", input: "secret", expected: "[REDACTED]"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, String(testCase.input))
		})
	}
}
