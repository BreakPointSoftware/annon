package copytest

import (
	"reflect"
	"testing"

	"github.com/BreakPointSoftware/annon/internal/copytest/baseline"
	"github.com/BreakPointSoftware/annon/internal/copytest/hybrid"
	"github.com/BreakPointSoftware/annon/internal/copytest/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValueOnlyStructs(t *testing.T) {
	input := testdata.DemoValueOnly()
	baselineCopy, err := baseline.Copy(input)
	require.NoError(t, err)
	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)

	assert.Equal(t, input, baselineCopy)
	assert.Equal(t, input, hybridResult.Copy)
	assert.True(t, hasReason(hybridResult.Flags, hybrid.SensitiveFieldName))
}

func TestUnexportedValueFieldPreservation(t *testing.T) {
	input := testdata.NewWithPrivateValue()
	baselineCopy, err := baseline.Copy(input)
	require.NoError(t, err)
	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)

	assert.NotEqual(t, input.Hash(), baselineCopy.Hash())
	assert.Equal(t, input.Hash(), hybridResult.Copy.Hash())
}

func TestExportedReferenceFieldsDetach(t *testing.T) {
	input := testdata.DemoExportedRefs()
	baselineCopy, err := baseline.Copy(input)
	require.NoError(t, err)
	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)

	assert.NotSame(t, input.Child, baselineCopy.Child)
	assert.NotSame(t, input.Child, hybridResult.Copy.Child)
	assert.NotSame(t, reflect.ValueOf(input.Names).Pointer(), reflect.ValueOf(baselineCopy.Names).Pointer())
	assert.NotSame(t, reflect.ValueOf(input.Names).Pointer(), reflect.ValueOf(hybridResult.Copy.Names).Pointer())
	assert.NotSame(t, reflect.ValueOf(input.Meta).Pointer(), reflect.ValueOf(baselineCopy.Meta).Pointer())
	assert.NotSame(t, reflect.ValueOf(input.Meta).Pointer(), reflect.ValueOf(hybridResult.Copy.Meta).Pointer())
}

func TestBaselineRejectsCyclesWhileHybridPreservesThem(t *testing.T) {
	input := testdata.DemoCycle()
	_, err := baseline.Copy(input)
	require.ErrorIs(t, err, baseline.ErrCycleUnsupported)

	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)
	require.NotNil(t, hybridResult.Copy.Sibling)
	assert.Same(t, hybridResult.Copy, hybridResult.Copy.Sibling)
	assert.Contains(t, hybridResult.Flags, hybrid.FieldFlag{Path: "Sibling", Type: reflect.TypeOf(input.Sibling), Kind: reflect.Pointer, Reason: hybrid.RecursiveReferenceReused, Action: hybrid.ActionReused})
}

func TestSiblingCycle(t *testing.T) {
	input := testdata.DemoSiblingCycle()
	_, err := baseline.Copy(input)
	require.ErrorIs(t, err, baseline.ErrCycleUnsupported)

	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)
	require.NotNil(t, hybridResult.Copy.Sibling)
	assert.Same(t, hybridResult.Copy, hybridResult.Copy.Sibling.Sibling)
}

func TestSharedChildPointerPreserved(t *testing.T) {
	input := testdata.DemoSharedChild()
	baselineCopy, err := baseline.Copy(input)
	require.NoError(t, err)
	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)

	assert.NotSame(t, baselineCopy.Left, baselineCopy.Right)
	assert.Same(t, hybridResult.Copy.Left, hybridResult.Copy.Right)
}

func TestPrivateReferenceSharedInHybrid(t *testing.T) {
	input := testdata.NewWithPrivateRef()
	baselineCopy, err := baseline.Copy(input)
	require.NoError(t, err)
	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)

	assert.Nil(t, baselineCopy.Child())
	assert.Same(t, input.Child(), hybridResult.Copy.Child())
	assert.True(t, hasReason(hybridResult.Flags, hybrid.UnexportedReferenceShared))
}

func TestRuntimeStatePolicy(t *testing.T) {
	input := testdata.DemoRuntimeState()
	baselineCopy, err := baseline.Copy(input)
	require.NoError(t, err)
	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)

	assert.Nil(t, baselineCopy.Func)
	assert.Nil(t, hybridResult.Copy.Func)
	assert.Nil(t, hybridResult.Copy.Chan)
	assert.Nil(t, hybridResult.Copy.Ctx)
	assert.Nil(t, hybridResult.Copy.Client)
	assert.True(t, hasReason(hybridResult.Flags, hybrid.RuntimeStateZeroed))
}

func TestInterfaceFields(t *testing.T) {
	input := testdata.DemoInterfaceFields()
	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)

	childCopy := hybridResult.Copy.Pointer.(*testdata.Child)
	assert.NotSame(t, input.Pointer, childCopy)
	ptrSlice := hybridResult.Copy.Slice.([]*testdata.Child)
	assert.NotSame(t, input.Slice.([]*testdata.Child)[0], ptrSlice[0])
	ptrMap := hybridResult.Copy.Map.(map[string]*testdata.Child)
	assert.NotSame(t, input.Map.(map[string]*testdata.Child)["primary"], ptrMap["primary"])
}

func TestNestedCollections(t *testing.T) {
	input := testdata.DemoNestedCollections()
	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)
	assert.Same(t, hybridResult.Copy.Groups["one"][0], hybridResult.Copy.Groups["two"][0])
	assert.Same(t, hybridResult.Copy.Items[0]["first"], hybridResult.Copy.Items[1]["second"])
}

func TestNilReferences(t *testing.T) {
	input := testdata.DemoNilRefs()
	baselineCopy, err := baseline.Copy(input)
	require.NoError(t, err)
	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)
	assert.Nil(t, baselineCopy.Ptr)
	assert.Nil(t, hybridResult.Copy.Ptr)
	assert.Nil(t, baselineCopy.Slice)
	assert.Nil(t, hybridResult.Copy.Slice)
	assert.Nil(t, baselineCopy.Map)
	assert.Nil(t, hybridResult.Copy.Map)
	assert.Nil(t, baselineCopy.Any)
	assert.Nil(t, hybridResult.Copy.Any)
}

func TestSourceObjectsNotMutated(t *testing.T) {
	input := testdata.DemoExportedRefs()
	_, err := baseline.Copy(input)
	require.NoError(t, err)
	_, err = hybrid.Copy(input)
	require.NoError(t, err)
	assert.Equal(t, "child", input.Child.Name)
	assert.Equal(t, []string{"one", "two"}, input.Names)
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, input.Meta)
}

func hasReason(flags []hybrid.FieldFlag, reason hybrid.FlagReason) bool {
	for _, flag := range flags {
		if flag.Reason == reason {
			return true
		}
	}
	return false
}
