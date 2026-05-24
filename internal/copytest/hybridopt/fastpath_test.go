package hybridopt

import (
	"testing"

	"github.com/BreakPointSoftware/annon/internal/copytest/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFastPathValueOnlyStruct(t *testing.T) {
	input := testdata.DemoValueOnly()
	result, err := Copy(input)
	require.NoError(t, err)

	assert.Equal(t, input, result.Copy)
	assert.True(t, hasReason(result.Flags, SensitiveFieldName))
}

func TestFastPathPreservesPrivateValueFields(t *testing.T) {
	input := testdata.NewWithPrivateValue()
	result, err := Copy(input)
	require.NoError(t, err)

	assert.Equal(t, input.Hash(), result.Copy.Hash())
}

func hasReason(flags []FieldFlag, reason FlagReason) bool {
	for _, flag := range flags {
		if flag.Reason == reason {
			return true
		}
	}
	return false
}
