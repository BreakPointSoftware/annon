package hybridopt

import (
	"testing"

	"github.com/BreakPointSoftware/annon/internal/copytest/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHybridFlagsDisabledProducesNoFlags(t *testing.T) {
	input := testdata.DemoSharedChild()
	result, err := CopyWithOptions(input, Options{CollectFlags: false})
	require.NoError(t, err)

	assert.Empty(t, result.Flags)
}

func TestHybridFlagsDisabledPreservesCopySemantics(t *testing.T) {
	input := testdata.DemoExportedRefs()
	flagsOn, err := Copy(input)
	require.NoError(t, err)
	flagsOff, err := CopyWithOptions(input, Options{CollectFlags: false})
	require.NoError(t, err)

	assert.Equal(t, flagsOn.Copy.Child.Name, flagsOff.Copy.Child.Name)
	assert.Equal(t, flagsOn.Copy.Names, flagsOff.Copy.Names)
	assert.Equal(t, flagsOn.Copy.Meta, flagsOff.Copy.Meta)

	assert.NotSame(t, input.Child, flagsOff.Copy.Child)
	assert.NotSame(t, input.Child, flagsOn.Copy.Child)

	assert.Empty(t, flagsOff.Flags)
	assert.NotEmpty(t, flagsOn.Flags)
}
