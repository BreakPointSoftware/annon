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

	// Both approaches should preserve pure value-only structs exactly.
	assert.Equal(t, input, baselineCopy)
	assert.Equal(t, input, hybridResult.Copy)

	// The hybrid implementation may still record that a field name looked sensitive.
	assert.True(t, hasReason(hybridResult.Flags, hybrid.SensitiveFieldName))
}

func TestUnexportedValueFieldPreservation(t *testing.T) {
	input := testdata.NewWithPrivateValue()

	baselineCopy, err := baseline.Copy(input)
	require.NoError(t, err)

	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)

	// Baseline rebuilds exported fields only, so the private hash value is lost.
	assert.NotEqual(t, input.Hash(), baselineCopy.Hash())

	// Hybrid copies structs by value first, so the private hash is preserved.
	assert.Equal(t, input.Hash(), hybridResult.Copy.Hash())
}

func TestExportedReferenceFieldsDetach(t *testing.T) {
	input := testdata.DemoExportedRefs()

	baselineCopy, err := baseline.Copy(input)
	require.NoError(t, err)

	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)

	// Both approaches should detach exported references from the original object graph.
	assert.NotSame(t, input.Child, baselineCopy.Child)
	assert.NotSame(t, input.Child, hybridResult.Copy.Child)
	assert.NotSame(t, reflect.ValueOf(input.Names).Pointer(), reflect.ValueOf(baselineCopy.Names).Pointer())
	assert.NotSame(t, reflect.ValueOf(input.Names).Pointer(), reflect.ValueOf(hybridResult.Copy.Names).Pointer())
	assert.NotSame(t, reflect.ValueOf(input.Meta).Pointer(), reflect.ValueOf(baselineCopy.Meta).Pointer())
	assert.NotSame(t, reflect.ValueOf(input.Meta).Pointer(), reflect.ValueOf(hybridResult.Copy.Meta).Pointer())

	// Detachment should not alter the stored values carried by those references.
	assert.Equal(t, input.Child.Name, baselineCopy.Child.Name)
	assert.Equal(t, input.Child.Name, hybridResult.Copy.Child.Name)
	assert.Equal(t, input.Names, baselineCopy.Names)
	assert.Equal(t, input.Names, hybridResult.Copy.Names)
	assert.Equal(t, input.Meta, baselineCopy.Meta)
	assert.Equal(t, input.Meta, hybridResult.Copy.Meta)
}

func TestBaselineRejectsCyclesWhileHybridPreservesThem(t *testing.T) {
	input := testdata.DemoCycle()

	// Baseline is intentionally simple and should reject recursive pointer cycles.
	_, err := baseline.Copy(input)
	require.ErrorIs(t, err, baseline.ErrCycleUnsupported)

	// Hybrid should terminate safely and produce a copied cycle.
	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)
	require.NotNil(t, hybridResult.Copy.Sibling)

	// The copied cycle must be detached from the original node.
	assert.NotSame(t, input, hybridResult.Copy)
	assert.NotSame(t, input, hybridResult.Copy.Sibling)

	// The copied graph must preserve the original self-reference shape.
	assert.Same(t, hybridResult.Copy, hybridResult.Copy.Sibling)

	// The copied node should still carry the same stored value as the original node.
	assert.Equal(t, input.Name, hybridResult.Copy.Name)
	assert.Equal(t, input.Name, hybridResult.Copy.Sibling.Name)

	// Hybrid should record that the recursive sibling reference was reused.
	assert.True(t, hasFlag(hybridResult.Flags, func(flag hybrid.FieldFlag) bool {
		return flag.Path == "Sibling" &&
			flag.Type == reflect.TypeOf(input.Sibling) &&
			flag.Kind == reflect.Pointer &&
			flag.Reason == hybrid.RecursiveReferenceReused &&
			flag.Action == hybrid.ActionReused
	}))
}

func TestSiblingCycle(t *testing.T) {
	input := testdata.DemoSiblingCycle()

	// Baseline is intentionally simple and should reject recursive sibling cycles.
	_, err := baseline.Copy(input)
	require.ErrorIs(t, err, baseline.ErrCycleUnsupported)

	// Hybrid should terminate safely and rebuild the sibling cycle.
	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)
	require.NotNil(t, hybridResult.Copy.Sibling)

	// The copied left and right nodes must be detached from the original pair.
	assert.NotSame(t, input, hybridResult.Copy)
	assert.NotSame(t, input.Sibling, hybridResult.Copy.Sibling)

	// The copied graph should preserve the sibling cycle internally.
	assert.Same(t, hybridResult.Copy, hybridResult.Copy.Sibling.Sibling)

	// Stored values on both copied nodes should match the original pair.
	assert.Equal(t, input.Name, hybridResult.Copy.Name)
	assert.Equal(t, input.Sibling.Name, hybridResult.Copy.Sibling.Name)

	// Hybrid should report that recursive references were reused while closing the cycle.
	assert.True(t, hasReason(hybridResult.Flags, hybrid.RecursiveReferenceReused))
}

func TestSharedChildPointerPreserved(t *testing.T) {
	input := testdata.DemoSharedChild()

	baselineCopy, err := baseline.Copy(input)
	require.NoError(t, err)

	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)

	// Baseline does not preserve shared graph semantics and will duplicate the child.
	assert.NotSame(t, baselineCopy.Left, baselineCopy.Right)

	// Hybrid should detach the child from the original graph.
	assert.NotSame(t, input.Left, hybridResult.Copy.Left)
	assert.NotSame(t, input.Right, hybridResult.Copy.Right)

	// The copied graph should preserve the original sharing relationship internally.
	assert.Same(t, hybridResult.Copy.Left, hybridResult.Copy.Right)

	// The copied shared child should still describe the same logical child.
	assert.Equal(t, input.Left.Name, hybridResult.Copy.Left.Name)
	assert.Equal(t, input.Right.Name, hybridResult.Copy.Right.Name)
}

func TestPrivateReferenceSharedInHybrid(t *testing.T) {
	input := testdata.NewWithPrivateRef()

	baselineCopy, err := baseline.Copy(input)
	require.NoError(t, err)

	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)

	// Baseline only copies exported fields and therefore loses the private reference.
	assert.Nil(t, baselineCopy.Child())

	// Without unsafe, hybrid cannot repair private references and currently leaves them shared.
	assert.Same(t, input.Child(), hybridResult.Copy.Child())

	// The shared private reference should still carry the same stored value.
	assert.Equal(t, input.Child().Name, hybridResult.Copy.Child().Name)

	// A flag should make that shared private-reference trade-off explicit.
	assert.True(t, hasReason(hybridResult.Flags, hybrid.UnexportedReferenceShared))
}

func TestRuntimeStatePolicy(t *testing.T) {
	input := testdata.DemoRuntimeState()

	baselineCopy, err := baseline.Copy(input)
	require.NoError(t, err)

	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)

	// Runtime and concurrency state should not be reproduced as ordinary data.
	assert.Nil(t, baselineCopy.Func)
	assert.Nil(t, hybridResult.Copy.Func)
	assert.Nil(t, hybridResult.Copy.Chan)
	assert.Nil(t, hybridResult.Copy.Ctx)
	assert.Nil(t, hybridResult.Copy.Client)

	// Ordinary data fields should still preserve their stored values.
	assert.Equal(t, input.Name, baselineCopy.Name)
	assert.Equal(t, input.Name, hybridResult.Copy.Name)

	// The hybrid flag stream should explain why runtime-like fields were zeroed/shared.
	assert.True(t, hasReason(hybridResult.Flags, hybrid.RuntimeStateZeroed))
}

func TestInterfaceFields(t *testing.T) {
	input := testdata.DemoInterfaceFields()

	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)

	// Interface-contained references should be detached from the original graph.
	childCopy := hybridResult.Copy.Pointer.(*testdata.Child)
	assert.NotSame(t, input.Pointer, childCopy)
	ptrSlice := hybridResult.Copy.Slice.([]*testdata.Child)
	assert.NotSame(t, input.Slice.([]*testdata.Child)[0], ptrSlice[0])
	ptrMap := hybridResult.Copy.Map.(map[string]*testdata.Child)
	assert.NotSame(t, input.Map.(map[string]*testdata.Child)["primary"], ptrMap["primary"])

	// Interface-contained copied values should still preserve the original stored data.
	assert.Equal(t, input.Pointer.(*testdata.Child).Name, childCopy.Name)
	assert.Equal(t, input.Slice.([]*testdata.Child)[0].Name, ptrSlice[0].Name)
	assert.Equal(t, input.Map.(map[string]*testdata.Child)["primary"].Name, ptrMap["primary"].Name)
}

func TestNestedCollections(t *testing.T) {
	input := testdata.DemoNestedCollections()

	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)

	// Nested references should be detached from the original graph.
	assert.NotSame(t, input.Groups["one"][0], hybridResult.Copy.Groups["one"][0])
	assert.NotSame(t, input.Items[0]["first"], hybridResult.Copy.Items[0]["first"])

	// Repeated nested references should still be shared within the copied graph.
	assert.Same(t, hybridResult.Copy.Groups["one"][0], hybridResult.Copy.Groups["two"][0])
	assert.Same(t, hybridResult.Copy.Items[0]["first"], hybridResult.Copy.Items[1]["second"])

	// The copied nested children should still carry the same stored values.
	assert.Equal(t, input.Groups["one"][0].Name, hybridResult.Copy.Groups["one"][0].Name)
	assert.Equal(t, input.Items[0]["first"].Name, hybridResult.Copy.Items[0]["first"].Name)
}

func TestNilReferences(t *testing.T) {
	input := testdata.DemoNilRefs()

	baselineCopy, err := baseline.Copy(input)
	require.NoError(t, err)

	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)

	// Absence of data should be preserved too.
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

	// Neither approach should mutate caller-owned input while building copies.
	_, err := baseline.Copy(input)
	require.NoError(t, err)
	_, err = hybrid.Copy(input)
	require.NoError(t, err)
	assert.Equal(t, "child", input.Child.Name)
	assert.Equal(t, []string{"one", "two"}, input.Names)
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, input.Meta)
}

func TestHybridDetachesCopiedGraphFromOriginal(t *testing.T) {
	input := testdata.DemoSharedChild()

	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)

	// The copied references should be detached from the original source graph.
	assert.NotSame(t, input.Left, hybridResult.Copy.Left)
	assert.NotSame(t, input.Right, hybridResult.Copy.Right)

	// The copied graph should still preserve the original sharing relationship.
	assert.Same(t, hybridResult.Copy.Left, hybridResult.Copy.Right)

	// The detached copied child should still carry the same stored value.
	assert.Equal(t, input.Left.Name, hybridResult.Copy.Left.Name)
}

func hasReason(flags []hybrid.FieldFlag, reason hybrid.FlagReason) bool {
	for _, flag := range flags {
		if flag.Reason == reason {
			return true
		}
	}
	return false
}

func hasFlag(flags []hybrid.FieldFlag, match func(hybrid.FieldFlag) bool) bool {
	for _, flag := range flags {
		if match(flag) {
			return true
		}
	}
	return false
}
