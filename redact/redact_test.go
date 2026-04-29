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
