package annon

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type customer struct {
	ID     string `json:"id" yaml:"id"`
	Email  string `json:"email" yaml:"email"`
	Secret string `json:"secret" yaml:"secret" anonymise:"remove"`
	Note   string `json:"note" yaml:"note" anonymise:"false"`
}

func TestCopy(t *testing.T) {
	testCases := []struct {
		name     string
		input    customer
		assertFn func(t *testing.T, result customer)
	}{
		{
			name:  "copy applies field detection and remove tag",
			input: customer{ID: "123", Email: "greg@example.com", Secret: "secret", Note: "greg@example.com"},
			assertFn: func(t *testing.T, result customer) {
				assert.Equal(t, "g***@example.com", result.Email)
				assert.Empty(t, result.Secret)
				assert.Equal(t, "greg@example.com", result.Note)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			anonymiser, err := New()
			require.NoError(t, err)

			originalInput := testCase.input
			copyResult, err := anonymiser.Copy(testCase.input)
			require.NoError(t, err)

			copiedCustomer := copyResult.(customer)
			testCase.assertFn(t, copiedCustomer)
			assert.Equal(t, originalInput, testCase.input)
		})
	}
}

func TestJSONAndYAML(t *testing.T) {
	anonymiser, err := New()
	require.NoError(t, err)

	inputCustomer := customer{ID: "123", Email: "greg@example.com", Secret: "secret", Note: "greg@example.com"}

	testCases := []struct {
		name        string
		run         func() ([]byte, error)
		contains    []string
		notContains []string
	}{
		{
			name: "json output omits removed field",
			run: func() ([]byte, error) {
				return anonymiser.JSON(inputCustomer)
			},
			contains:    []string{`"email":"g***@example.com"`, `"note":"greg@example.com"`},
			notContains: []string{`"secret"`},
		},
		{
			name: "yaml output omits removed field",
			run: func() ([]byte, error) {
				return anonymiser.YAML(inputCustomer)
			},
			contains:    []string{"email: g***@example.com", "note: greg@example.com"},
			notContains: []string{"secret:"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resultBytes, err := testCase.run()
			require.NoError(t, err)

			for _, expectedSubstring := range testCase.contains {
				assert.Contains(t, string(resultBytes), expectedSubstring)
			}

			for _, unexpectedSubstring := range testCase.notContains {
				assert.NotContains(t, string(resultBytes), unexpectedSubstring)
			}
		})
	}
}

func TestRawInputModes(t *testing.T) {
	testCases := []struct {
		name          string
		run           func() ([]byte, error)
		expectedError error
		contains      []string
	}{
		{
			name: "from json anonymises nested content",
			run: func() ([]byte, error) {
				return FromJSON([]byte(`{"email":"greg@example.com","vehicle":{"reg":"AB12 CDE"}}`))
			},
			contains: []string{`"email":"g***@example.com"`, `"reg":"AB12 ***"`},
		},
		{
			name: "from yaml anonymises nested content",
			run: func() ([]byte, error) {
				return FromYAML([]byte("email: greg@example.com\nvehicle:\n  reg: AB12 CDE\n"))
			},
			contains: []string{"email: g***@example.com", "reg: AB12 ***"},
		},
		{
			name: "from json returns sentinel error for invalid input",
			run: func() ([]byte, error) {
				return FromJSON([]byte(`{"email":`))
			},
			expectedError: ErrInvalidJSON,
		},
		{
			name: "from yaml returns sentinel error for invalid input",
			run: func() ([]byte, error) {
				return FromYAML([]byte("email: ["))
			},
			expectedError: ErrInvalidYAML,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resultBytes, err := testCase.run()

			if testCase.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, testCase.expectedError)
				return
			}

			require.NoError(t, err)

			for _, expectedSubstring := range testCase.contains {
				assert.Contains(t, string(resultBytes), expectedSubstring)
			}
		})
	}
}

func TestFieldRulesAndValueDetection(t *testing.T) {
	testCases := []struct {
		name     string
		run      func() ([]byte, error)
		contains string
	}{
		{
			name: "custom field rule uses redact strategy",
			run: func() ([]byte, error) {
				type alias struct{ Alias string `json:"customerAlias"` }
				return JSON(alias{Alias: "secret"}, WithFieldRules(StrongRule(RedactStrategy, "customerAlias")))
			},
			contains: `"customerAlias":"[REDACTED]"`,
		},
		{
			name: "value detection runs when field detection disabled",
			run: func() ([]byte, error) {
				return FromJSON([]byte(`{"note":"greg@example.com"}`), WithFieldDetection(false), WithValueDetection(true))
			},
			contains: `"note":"g***@example.com"`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resultBytes, err := testCase.run()
			require.NoError(t, err)
			assert.Contains(t, string(resultBytes), testCase.contains)
		})
	}
}

func TestRawInputBytesAreNotMutated(t *testing.T) {
	jsonInput := []byte(`{"email":"greg@example.com"}`)
	jsonOriginal := append([]byte(nil), jsonInput...)
	_, err := FromJSON(jsonInput)
	require.NoError(t, err)
	assert.True(t, bytes.Equal(jsonInput, jsonOriginal))

	yamlInput := []byte("email: greg@example.com\n")
	yamlOriginal := append([]byte(nil), yamlInput...)
	_, err = FromYAML(yamlInput)
	require.NoError(t, err)
	assert.True(t, bytes.Equal(yamlInput, yamlOriginal))
}
