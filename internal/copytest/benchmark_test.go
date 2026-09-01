package copytest

import (
	"testing"

	"github.com/BreakPointSoftware/annon/internal/copytest/baseline"
	"github.com/BreakPointSoftware/annon/internal/copytest/hybrid"
	"github.com/BreakPointSoftware/annon/internal/copytest/testdata"
)

func BenchmarkBaselineSmallStruct(b *testing.B) {
	input := testdata.DemoValueOnly()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = baseline.Copy(input)
	}
}

func BenchmarkHybridSmallStruct(b *testing.B) {
	input := testdata.DemoValueOnly()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybrid.Copy(input)
	}
}

func BenchmarkBaselineLargeStruct(b *testing.B) {
	input := testdata.DemoLargeValue()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = baseline.Copy(input)
	}
}

func BenchmarkHybridLargeStruct(b *testing.B) {
	input := testdata.DemoLargeValue()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybrid.Copy(input)
	}
}

func BenchmarkBaselineManyPointers(b *testing.B) {
	input := testdata.DemoManyPointers()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = baseline.Copy(input)
	}
}

func BenchmarkHybridManyPointers(b *testing.B) {
	input := testdata.DemoManyPointers()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybrid.Copy(input)
	}
}

func BenchmarkBaselineNestedSlices(b *testing.B) {
	input := testdata.DemoNestedCollections()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = baseline.Copy(input)
	}
}

func BenchmarkHybridNestedSlices(b *testing.B) {
	input := testdata.DemoNestedCollections()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybrid.Copy(input)
	}
}

func BenchmarkBaselineNestedMaps(b *testing.B) {
	input := testdata.DemoExportedRefs()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = baseline.Copy(input)
	}
}

func BenchmarkHybridNestedMaps(b *testing.B) {
	input := testdata.DemoExportedRefs()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybrid.Copy(input)
	}
}

func BenchmarkBaselineDeepTree(b *testing.B) {
	input := testdata.DemoTree()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = baseline.Copy(input)
	}
}

func BenchmarkHybridDeepTree(b *testing.B) {
	input := testdata.DemoTree()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybrid.Copy(input)
	}
}

func BenchmarkBaselineCyclicGraph(b *testing.B) {
	input := testdata.DemoCycle()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = baseline.Copy(input)
	}
}

func BenchmarkHybridCyclicGraph(b *testing.B) {
	input := testdata.DemoCycle()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybrid.Copy(input)
	}
}

func BenchmarkBaselineMixedDomainObject(b *testing.B) {
	input := testdata.DemoDomainObject()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = baseline.Copy(input)
	}
}

func BenchmarkHybridMixedDomainObject(b *testing.B) {
	input := testdata.DemoDomainObject()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybrid.Copy(input)
	}
}
