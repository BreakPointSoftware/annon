package copytest

import (
	"testing"

	"github.com/BreakPointSoftware/annon/internal/copytest/baseline"
	"github.com/BreakPointSoftware/annon/internal/copytest/hybrid"
	"github.com/BreakPointSoftware/annon/internal/copytest/hybridopt"
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

func BenchmarkHybridOptimisedSmallStruct(b *testing.B) {
	input := testdata.DemoValueOnly()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybridopt.Copy(input)
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

func BenchmarkHybridOptimisedLargeStruct(b *testing.B) {
	input := testdata.DemoLargeValue()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybridopt.Copy(input)
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

func BenchmarkHybridOptimisedManyPointers(b *testing.B) {
	input := testdata.DemoManyPointers()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybridopt.Copy(input)
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

func BenchmarkHybridOptimisedNestedSlices(b *testing.B) {
	input := testdata.DemoNestedCollections()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybridopt.Copy(input)
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

func BenchmarkHybridOptimisedNestedMaps(b *testing.B) {
	input := testdata.DemoExportedRefs()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybridopt.Copy(input)
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

func BenchmarkHybridOptimisedDeepTree(b *testing.B) {
	input := testdata.DemoTree()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybridopt.Copy(input)
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

func BenchmarkHybridOptimisedCyclicGraph(b *testing.B) {
	input := testdata.DemoCycle()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybridopt.Copy(input)
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

func BenchmarkHybridOptimisedMixedDomainObject(b *testing.B) {
	input := testdata.DemoDomainObject()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybridopt.Copy(input)
	}
}

func BenchmarkBaselineSmallStructBatch100(b *testing.B) {
	input := testdata.DemoValueOnly()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			_, _ = baseline.Copy(input)
		}
	}
}

func BenchmarkHybridSmallStructBatch100(b *testing.B) {
	input := testdata.DemoValueOnly()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			_, _ = hybrid.Copy(input)
		}
	}
}

func BenchmarkHybridOptimisedSmallStructBatch100(b *testing.B) {
	input := testdata.DemoValueOnly()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			_, _ = hybridopt.Copy(input)
		}
	}
}

func BenchmarkBaselineLargeStructBatch100(b *testing.B) {
	input := testdata.DemoLargeValue()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			_, _ = baseline.Copy(input)
		}
	}
}

func BenchmarkHybridLargeStructBatch100(b *testing.B) {
	input := testdata.DemoLargeValue()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			_, _ = hybrid.Copy(input)
		}
	}
}

func BenchmarkHybridOptimisedLargeStructBatch100(b *testing.B) {
	input := testdata.DemoLargeValue()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			_, _ = hybridopt.Copy(input)
		}
	}
}

func BenchmarkBaselineValueOnlyBatch1000(b *testing.B) {
	input := testdata.DemoValueOnly()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			_, _ = baseline.Copy(input)
		}
	}
}

func BenchmarkHybridValueOnlyBatch1000(b *testing.B) {
	input := testdata.DemoValueOnly()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			_, _ = hybrid.Copy(input)
		}
	}
}

func BenchmarkHybridOptimisedValueOnlyBatch1000(b *testing.B) {
	input := testdata.DemoValueOnly()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			_, _ = hybridopt.Copy(input)
		}
	}
}

func BenchmarkBaselineMixedDomainBatch100(b *testing.B) {
	input := testdata.DemoDomainObject()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			_, _ = baseline.Copy(input)
		}
	}
}

func BenchmarkHybridMixedDomainBatch100(b *testing.B) {
	input := testdata.DemoDomainObject()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			_, _ = hybrid.Copy(input)
		}
	}
}

func BenchmarkHybridOptimisedMixedDomainBatch100(b *testing.B) {
	input := testdata.DemoDomainObject()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			_, _ = hybridopt.Copy(input)
		}
	}
}

func BenchmarkHybridOptimisedFlagsOnSmallStruct(b *testing.B) {
	input := testdata.DemoValueOnly()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybridopt.CopyWithOptions(input, hybridopt.Options{CollectFlags: true})
	}
}

func BenchmarkHybridOptimisedFlagsOffSmallStruct(b *testing.B) {
	input := testdata.DemoValueOnly()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybridopt.CopyWithOptions(input, hybridopt.Options{CollectFlags: false})
	}
}

func BenchmarkHybridOptimisedFlagsOnLargeStruct(b *testing.B) {
	input := testdata.DemoLargeValue()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybridopt.CopyWithOptions(input, hybridopt.Options{CollectFlags: true})
	}
}

func BenchmarkHybridOptimisedFlagsOffLargeStruct(b *testing.B) {
	input := testdata.DemoLargeValue()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybridopt.CopyWithOptions(input, hybridopt.Options{CollectFlags: false})
	}
}

func BenchmarkHybridOptimisedFlagsOnManyPointers(b *testing.B) {
	input := testdata.DemoManyPointers()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybridopt.CopyWithOptions(input, hybridopt.Options{CollectFlags: true})
	}
}

func BenchmarkHybridOptimisedFlagsOffManyPointers(b *testing.B) {
	input := testdata.DemoManyPointers()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybridopt.CopyWithOptions(input, hybridopt.Options{CollectFlags: false})
	}
}

func BenchmarkHybridOptimisedFlagsOnMixedDomainObject(b *testing.B) {
	input := testdata.DemoDomainObject()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybridopt.CopyWithOptions(input, hybridopt.Options{CollectFlags: true})
	}
}

func BenchmarkHybridOptimisedFlagsOffMixedDomainObject(b *testing.B) {
	input := testdata.DemoDomainObject()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = hybridopt.CopyWithOptions(input, hybridopt.Options{CollectFlags: false})
	}
}

func BenchmarkHybridOptimisedFlagsOnSmallStructBatch100(b *testing.B) {
	input := testdata.DemoValueOnly()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			_, _ = hybridopt.CopyWithOptions(input, hybridopt.Options{CollectFlags: true})
		}
	}
}

func BenchmarkHybridOptimisedFlagsOffSmallStructBatch100(b *testing.B) {
	input := testdata.DemoValueOnly()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			_, _ = hybridopt.CopyWithOptions(input, hybridopt.Options{CollectFlags: false})
		}
	}
}

func BenchmarkHybridOptimisedFlagsOnLargeStructBatch100(b *testing.B) {
	input := testdata.DemoLargeValue()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			_, _ = hybridopt.CopyWithOptions(input, hybridopt.Options{CollectFlags: true})
		}
	}
}

func BenchmarkHybridOptimisedFlagsOffLargeStructBatch100(b *testing.B) {
	input := testdata.DemoLargeValue()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			_, _ = hybridopt.CopyWithOptions(input, hybridopt.Options{CollectFlags: false})
		}
	}
}

func BenchmarkHybridOptimisedFlagsOnValueOnlyBatch1000(b *testing.B) {
	input := testdata.DemoValueOnly()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			_, _ = hybridopt.CopyWithOptions(input, hybridopt.Options{CollectFlags: true})
		}
	}
}

func BenchmarkHybridOptimisedFlagsOffValueOnlyBatch1000(b *testing.B) {
	input := testdata.DemoValueOnly()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			_, _ = hybridopt.CopyWithOptions(input, hybridopt.Options{CollectFlags: false})
		}
	}
}

func BenchmarkHybridOptimisedFlagsOnMixedDomainBatch100(b *testing.B) {
	input := testdata.DemoDomainObject()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			_, _ = hybridopt.CopyWithOptions(input, hybridopt.Options{CollectFlags: true})
		}
	}
}

func BenchmarkHybridOptimisedFlagsOffMixedDomainBatch100(b *testing.B) {
	input := testdata.DemoDomainObject()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			_, _ = hybridopt.CopyWithOptions(input, hybridopt.Options{CollectFlags: false})
		}
	}
}
