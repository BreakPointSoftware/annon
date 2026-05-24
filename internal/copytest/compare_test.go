package copytest

import (
	"reflect"
	"testing"

	"github.com/BreakPointSoftware/annon/internal/copytest/baseline"
	"github.com/BreakPointSoftware/annon/internal/copytest/hybrid"
	"github.com/BreakPointSoftware/annon/internal/copytest/hybridopt"
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
	hybridOptimisedResult, err := hybridopt.Copy(input)
	require.NoError(t, err)

	// Both approaches should preserve pure value-only structs exactly.
	assert.Equal(t, input, baselineCopy)
	assert.Equal(t, input, hybridResult.Copy)
	assert.Equal(t, input, hybridOptimisedResult.Copy)

	// Both hybrid variants may still record that a field name looked sensitive.
	assert.True(t, hasReason(hybridResult.Flags, hybrid.SensitiveFieldName))
	assert.True(t, hasReasonOpt(hybridOptimisedResult.Flags, hybridopt.SensitiveFieldName))
}

func TestUnexportedValueFieldPreservation(t *testing.T) {
	input := testdata.NewWithPrivateValue()

	baselineCopy, err := baseline.Copy(input)
	require.NoError(t, err)

	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)
	hybridOptimisedResult, err := hybridopt.Copy(input)
	require.NoError(t, err)

	// Baseline rebuilds exported fields only, so the private hash value is lost.
	assert.NotEqual(t, input.Hash(), baselineCopy.Hash())

	// Both hybrid variants copy structs by value first, so the private hash is preserved.
	assert.Equal(t, input.Hash(), hybridResult.Copy.Hash())
	assert.Equal(t, input.Hash(), hybridOptimisedResult.Copy.Hash())
}

func TestExportedReferenceFieldsDetach(t *testing.T) {
	input := testdata.DemoExportedRefs()

	baselineCopy, err := baseline.Copy(input)
	require.NoError(t, err)

	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)
	hybridOptimisedResult, err := hybridopt.Copy(input)
	require.NoError(t, err)

	// Both approaches should detach exported references from the original object graph.
	assert.NotSame(t, input.Child, baselineCopy.Child)
	assert.NotSame(t, input.Child, hybridResult.Copy.Child)
	assert.NotSame(t, input.Child, hybridOptimisedResult.Copy.Child)
	assert.NotSame(t, reflect.ValueOf(input.Names).Pointer(), reflect.ValueOf(baselineCopy.Names).Pointer())
	assert.NotSame(t, reflect.ValueOf(input.Names).Pointer(), reflect.ValueOf(hybridResult.Copy.Names).Pointer())
	assert.NotSame(t, reflect.ValueOf(input.Names).Pointer(), reflect.ValueOf(hybridOptimisedResult.Copy.Names).Pointer())
	assert.NotSame(t, reflect.ValueOf(input.Meta).Pointer(), reflect.ValueOf(baselineCopy.Meta).Pointer())
	assert.NotSame(t, reflect.ValueOf(input.Meta).Pointer(), reflect.ValueOf(hybridResult.Copy.Meta).Pointer())
	assert.NotSame(t, reflect.ValueOf(input.Meta).Pointer(), reflect.ValueOf(hybridOptimisedResult.Copy.Meta).Pointer())

	// Detachment should not alter the stored values carried by those references.
	assert.Equal(t, input.Child.Name, baselineCopy.Child.Name)
	assert.Equal(t, input.Child.Name, hybridResult.Copy.Child.Name)
	assert.Equal(t, input.Child.Name, hybridOptimisedResult.Copy.Child.Name)
	assert.Equal(t, input.Names, baselineCopy.Names)
	assert.Equal(t, input.Names, hybridResult.Copy.Names)
	assert.Equal(t, input.Names, hybridOptimisedResult.Copy.Names)
	assert.Equal(t, input.Meta, baselineCopy.Meta)
	assert.Equal(t, input.Meta, hybridResult.Copy.Meta)
	assert.Equal(t, input.Meta, hybridOptimisedResult.Copy.Meta)
}

func TestBaselineRejectsCyclesWhileHybridPreservesThem(t *testing.T) {
	input := testdata.DemoCycle()

	// Baseline is intentionally simple and should reject recursive pointer cycles.
	_, err := baseline.Copy(input)
	require.ErrorIs(t, err, baseline.ErrCycleUnsupported)

	// Both hybrid variants should terminate safely and produce a copied cycle.
	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)
	require.NotNil(t, hybridResult.Copy.Sibling)
	hybridOptimisedResult, err := hybridopt.Copy(input)
	require.NoError(t, err)
	require.NotNil(t, hybridOptimisedResult.Copy.Sibling)

	// The copied cycle must be detached from the original node.
	assert.NotSame(t, input, hybridResult.Copy)
	assert.NotSame(t, input, hybridResult.Copy.Sibling)
	assert.NotSame(t, input, hybridOptimisedResult.Copy)
	assert.NotSame(t, input, hybridOptimisedResult.Copy.Sibling)

	// The copied graph must preserve the original self-reference shape.
	assert.Same(t, hybridResult.Copy, hybridResult.Copy.Sibling)
	assert.Same(t, hybridOptimisedResult.Copy, hybridOptimisedResult.Copy.Sibling)

	// The copied node should still carry the same stored value as the original node.
	assert.Equal(t, input.Name, hybridResult.Copy.Name)
	assert.Equal(t, input.Name, hybridResult.Copy.Sibling.Name)
	assert.Equal(t, input.Name, hybridOptimisedResult.Copy.Name)
	assert.Equal(t, input.Name, hybridOptimisedResult.Copy.Sibling.Name)

	// Both hybrid variants should record that the recursive sibling reference was reused.
	assert.True(t, hasFlag(hybridResult.Flags, func(flag hybrid.FieldFlag) bool {
		return flag.Path == "Sibling" &&
			flag.Type == reflect.TypeOf(input.Sibling) &&
			flag.Kind == reflect.Pointer &&
			flag.Reason == hybrid.RecursiveReferenceReused &&
			flag.Action == hybrid.ActionReused
	}))
	assert.True(t, hasFlagOpt(hybridOptimisedResult.Flags, func(flag hybridopt.FieldFlag) bool {
		return flag.Path == "Sibling" &&
			flag.Type == reflect.TypeOf(input.Sibling) &&
			flag.Kind == reflect.Pointer &&
			flag.Reason == hybridopt.RecursiveReferenceReused &&
			flag.Action == hybridopt.ActionReused
	}))
}

func TestSiblingCycle(t *testing.T) {
	input := testdata.DemoSiblingCycle()

	// Baseline is intentionally simple and should reject recursive sibling cycles.
	_, err := baseline.Copy(input)
	require.ErrorIs(t, err, baseline.ErrCycleUnsupported)

	// Both hybrid variants should terminate safely and rebuild the sibling cycle.
	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)
	require.NotNil(t, hybridResult.Copy.Sibling)
	hybridOptimisedResult, err := hybridopt.Copy(input)
	require.NoError(t, err)
	require.NotNil(t, hybridOptimisedResult.Copy.Sibling)

	// The copied left and right nodes must be detached from the original pair.
	assert.NotSame(t, input, hybridResult.Copy)
	assert.NotSame(t, input.Sibling, hybridResult.Copy.Sibling)
	assert.NotSame(t, input, hybridOptimisedResult.Copy)
	assert.NotSame(t, input.Sibling, hybridOptimisedResult.Copy.Sibling)

	// The copied graph should preserve the sibling cycle internally.
	assert.Same(t, hybridResult.Copy, hybridResult.Copy.Sibling.Sibling)
	assert.Same(t, hybridOptimisedResult.Copy, hybridOptimisedResult.Copy.Sibling.Sibling)

	// Stored values on both copied nodes should match the original pair.
	assert.Equal(t, input.Name, hybridResult.Copy.Name)
	assert.Equal(t, input.Sibling.Name, hybridResult.Copy.Sibling.Name)
	assert.Equal(t, input.Name, hybridOptimisedResult.Copy.Name)
	assert.Equal(t, input.Sibling.Name, hybridOptimisedResult.Copy.Sibling.Name)

	// Both hybrid variants should report that recursive references were reused while closing the cycle.
	assert.True(t, hasReason(hybridResult.Flags, hybrid.RecursiveReferenceReused))
	assert.True(t, hasReasonOpt(hybridOptimisedResult.Flags, hybridopt.RecursiveReferenceReused))
}

func TestSharedChildPointerPreserved(t *testing.T) {
	input := testdata.DemoSharedChild()

	baselineCopy, err := baseline.Copy(input)
	require.NoError(t, err)

	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)
	hybridOptimisedResult, err := hybridopt.Copy(input)
	require.NoError(t, err)

	// Baseline does not preserve shared graph semantics and will duplicate the child.
	assert.NotSame(t, baselineCopy.Left, baselineCopy.Right)

	// Both hybrid variants should detach the child from the original graph.
	assert.NotSame(t, input.Left, hybridResult.Copy.Left)
	assert.NotSame(t, input.Right, hybridResult.Copy.Right)
	assert.NotSame(t, input.Left, hybridOptimisedResult.Copy.Left)
	assert.NotSame(t, input.Right, hybridOptimisedResult.Copy.Right)

	// The copied graph should preserve the original sharing relationship internally.
	assert.Same(t, hybridResult.Copy.Left, hybridResult.Copy.Right)
	assert.Same(t, hybridOptimisedResult.Copy.Left, hybridOptimisedResult.Copy.Right)

	// The copied shared child should still describe the same logical child.
	assert.Equal(t, input.Left.Name, hybridResult.Copy.Left.Name)
	assert.Equal(t, input.Right.Name, hybridResult.Copy.Right.Name)
	assert.Equal(t, input.Left.Name, hybridOptimisedResult.Copy.Left.Name)
	assert.Equal(t, input.Right.Name, hybridOptimisedResult.Copy.Right.Name)
}

func TestPrivateReferenceSharedInHybrid(t *testing.T) {
	input := testdata.NewWithPrivateRef()

	baselineCopy, err := baseline.Copy(input)
	require.NoError(t, err)

	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)
	hybridOptimisedResult, err := hybridopt.Copy(input)
	require.NoError(t, err)

	// Baseline only copies exported fields and therefore loses the private reference.
	assert.Nil(t, baselineCopy.Child())

	// Without unsafe, both hybrid variants cannot repair private references and currently leave them shared.
	assert.Same(t, input.Child(), hybridResult.Copy.Child())
	assert.Same(t, input.Child(), hybridOptimisedResult.Copy.Child())

	// The shared private reference should still carry the same stored value.
	assert.Equal(t, input.Child().Name, hybridResult.Copy.Child().Name)
	assert.Equal(t, input.Child().Name, hybridOptimisedResult.Copy.Child().Name)

	// Flags should make that shared private-reference trade-off explicit.
	assert.True(t, hasReason(hybridResult.Flags, hybrid.UnexportedReferenceShared))
	assert.True(t, hasReasonOpt(hybridOptimisedResult.Flags, hybridopt.UnexportedReferenceShared))
}

func TestRuntimeStatePolicy(t *testing.T) {
	input := testdata.DemoRuntimeState()

	baselineCopy, err := baseline.Copy(input)
	require.NoError(t, err)

	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)
	hybridOptimisedResult, err := hybridopt.Copy(input)
	require.NoError(t, err)

	// Runtime and concurrency state should not be reproduced as ordinary data.
	assert.Nil(t, baselineCopy.Func)
	assert.Nil(t, hybridResult.Copy.Func)
	assert.Nil(t, hybridOptimisedResult.Copy.Func)
	assert.Nil(t, hybridResult.Copy.Chan)
	assert.Nil(t, hybridResult.Copy.Ctx)
	assert.Nil(t, hybridResult.Copy.Client)
	assert.Nil(t, hybridOptimisedResult.Copy.Chan)
	assert.Nil(t, hybridOptimisedResult.Copy.Ctx)
	assert.Nil(t, hybridOptimisedResult.Copy.Client)

	// Ordinary data fields should still preserve their stored values.
	assert.Equal(t, input.Name, baselineCopy.Name)
	assert.Equal(t, input.Name, hybridResult.Copy.Name)
	assert.Equal(t, input.Name, hybridOptimisedResult.Copy.Name)

	// The hybrid flag streams should explain why runtime-like fields were zeroed/shared.
	assert.True(t, hasReason(hybridResult.Flags, hybrid.RuntimeStateZeroed))
	assert.True(t, hasReasonOpt(hybridOptimisedResult.Flags, hybridopt.RuntimeStateZeroed))
}

func TestInterfaceFields(t *testing.T) {
	input := testdata.DemoInterfaceFields()

	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)
	hybridOptimisedResult, err := hybridopt.Copy(input)
	require.NoError(t, err)

	// Interface-contained references should be detached from the original graph in both hybrid variants.
	childCopy := hybridResult.Copy.Pointer.(*testdata.Child)
	assert.NotSame(t, input.Pointer, childCopy)
	childOptimisedCopy := hybridOptimisedResult.Copy.Pointer.(*testdata.Child)
	assert.NotSame(t, input.Pointer, childOptimisedCopy)
	ptrSlice := hybridResult.Copy.Slice.([]*testdata.Child)
	assert.NotSame(t, input.Slice.([]*testdata.Child)[0], ptrSlice[0])
	ptrOptimisedSlice := hybridOptimisedResult.Copy.Slice.([]*testdata.Child)
	assert.NotSame(t, input.Slice.([]*testdata.Child)[0], ptrOptimisedSlice[0])
	ptrMap := hybridResult.Copy.Map.(map[string]*testdata.Child)
	assert.NotSame(t, input.Map.(map[string]*testdata.Child)["primary"], ptrMap["primary"])
	ptrOptimisedMap := hybridOptimisedResult.Copy.Map.(map[string]*testdata.Child)
	assert.NotSame(t, input.Map.(map[string]*testdata.Child)["primary"], ptrOptimisedMap["primary"])

	// Interface-contained copied values should still preserve the original stored data.
	assert.Equal(t, input.Pointer.(*testdata.Child).Name, childCopy.Name)
	assert.Equal(t, input.Slice.([]*testdata.Child)[0].Name, ptrSlice[0].Name)
	assert.Equal(t, input.Map.(map[string]*testdata.Child)["primary"].Name, ptrMap["primary"].Name)
	assert.Equal(t, input.Pointer.(*testdata.Child).Name, childOptimisedCopy.Name)
	assert.Equal(t, input.Slice.([]*testdata.Child)[0].Name, ptrOptimisedSlice[0].Name)
	assert.Equal(t, input.Map.(map[string]*testdata.Child)["primary"].Name, ptrOptimisedMap["primary"].Name)
}

func TestNestedCollections(t *testing.T) {
	input := testdata.DemoNestedCollections()

	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)
	hybridOptimisedResult, err := hybridopt.Copy(input)
	require.NoError(t, err)

	// Nested references should be detached from the original graph.
	assert.NotSame(t, input.Groups["one"][0], hybridResult.Copy.Groups["one"][0])
	assert.NotSame(t, input.Items[0]["first"], hybridResult.Copy.Items[0]["first"])
	assert.NotSame(t, input.Groups["one"][0], hybridOptimisedResult.Copy.Groups["one"][0])
	assert.NotSame(t, input.Items[0]["first"], hybridOptimisedResult.Copy.Items[0]["first"])

	// Repeated nested references should still be shared within the copied graph.
	assert.Same(t, hybridResult.Copy.Groups["one"][0], hybridResult.Copy.Groups["two"][0])
	assert.Same(t, hybridResult.Copy.Items[0]["first"], hybridResult.Copy.Items[1]["second"])
	assert.Same(t, hybridOptimisedResult.Copy.Groups["one"][0], hybridOptimisedResult.Copy.Groups["two"][0])
	assert.Same(t, hybridOptimisedResult.Copy.Items[0]["first"], hybridOptimisedResult.Copy.Items[1]["second"])

	// The copied nested children should still carry the same stored values.
	assert.Equal(t, input.Groups["one"][0].Name, hybridResult.Copy.Groups["one"][0].Name)
	assert.Equal(t, input.Items[0]["first"].Name, hybridResult.Copy.Items[0]["first"].Name)
	assert.Equal(t, input.Groups["one"][0].Name, hybridOptimisedResult.Copy.Groups["one"][0].Name)
	assert.Equal(t, input.Items[0]["first"].Name, hybridOptimisedResult.Copy.Items[0]["first"].Name)
}

func TestNilReferences(t *testing.T) {
	input := testdata.DemoNilRefs()

	baselineCopy, err := baseline.Copy(input)
	require.NoError(t, err)

	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)
	hybridOptimisedResult, err := hybridopt.Copy(input)
	require.NoError(t, err)

	// Absence of data should be preserved too.
	assert.Nil(t, baselineCopy.Ptr)
	assert.Nil(t, hybridResult.Copy.Ptr)
	assert.Nil(t, hybridOptimisedResult.Copy.Ptr)
	assert.Nil(t, baselineCopy.Slice)
	assert.Nil(t, hybridResult.Copy.Slice)
	assert.Nil(t, hybridOptimisedResult.Copy.Slice)
	assert.Nil(t, baselineCopy.Map)
	assert.Nil(t, hybridResult.Copy.Map)
	assert.Nil(t, hybridOptimisedResult.Copy.Map)
	assert.Nil(t, baselineCopy.Any)
	assert.Nil(t, hybridResult.Copy.Any)
	assert.Nil(t, hybridOptimisedResult.Copy.Any)
}

func TestSourceObjectsNotMutated(t *testing.T) {
	input := testdata.DemoExportedRefs()

	// Neither approach should mutate caller-owned input while building copies.
	_, err := baseline.Copy(input)
	require.NoError(t, err)
	_, err = hybrid.Copy(input)
	require.NoError(t, err)
	_, err = hybridopt.Copy(input)
	require.NoError(t, err)
	assert.Equal(t, "child", input.Child.Name)
	assert.Equal(t, []string{"one", "two"}, input.Names)
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, input.Meta)
}

func TestHybridDetachesCopiedGraphFromOriginal(t *testing.T) {
	input := testdata.DemoSharedChild()

	hybridResult, err := hybrid.Copy(input)
	require.NoError(t, err)
	hybridOptimisedResult, err := hybridopt.Copy(input)
	require.NoError(t, err)

	// The copied references should be detached from the original source graph.
	assert.NotSame(t, input.Left, hybridResult.Copy.Left)
	assert.NotSame(t, input.Right, hybridResult.Copy.Right)
	assert.NotSame(t, input.Left, hybridOptimisedResult.Copy.Left)
	assert.NotSame(t, input.Right, hybridOptimisedResult.Copy.Right)

	// The copied graph should still preserve the original sharing relationship.
	assert.Same(t, hybridResult.Copy.Left, hybridResult.Copy.Right)
	assert.Same(t, hybridOptimisedResult.Copy.Left, hybridOptimisedResult.Copy.Right)

	// The detached copied child should still carry the same stored value.
	assert.Equal(t, input.Left.Name, hybridResult.Copy.Left.Name)
	assert.Equal(t, input.Left.Name, hybridOptimisedResult.Copy.Left.Name)
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

func hasReasonOpt(flags []hybridopt.FieldFlag, reason hybridopt.FlagReason) bool {
	for _, flag := range flags {
		if flag.Reason == reason {
			return true
		}
	}
	return false
}

func hasFlagOpt(flags []hybridopt.FieldFlag, match func(hybridopt.FieldFlag) bool) bool {
	for _, flag := range flags {
		if match(flag) {
			return true
		}
	}
	return false
}
