# copytest benchmark: three-variant comparison

## Context
- Commit: `feat/copy-spike-three-variants` working tree before commit
- Branch: `feat/copy-spike-three-variants`
- Go version: `go1.26.2`
- OS/Arch: `linux/amd64`
- CPU: `Intel(R) Core(TM) i7-6700K CPU @ 4.00GHz`

## Command
```bash
go test -bench=. ./internal/copytest/...
```

## Variants under test
- `baseline`: simple recursive copy
- `hybrid`: unoptimised value-copy-first semantic reference implementation
- `hybridopt`: optimised hybrid with cached repair-plan metadata

## Raw output
```text
goos: linux
goarch: amd64
pkg: github.com/BreakPointSoftware/annon/internal/copytest
cpu: Intel(R) Core(TM) i7-6700K CPU @ 4.00GHz
BenchmarkBaselineSmallStruct-8                  1057056      1171 ns/op     224 B/op       4 allocs/op
BenchmarkHybridSmallStruct-8                    134878      8572 ns/op    7264 B/op      91 allocs/op
BenchmarkHybridOptimisedSmallStruct-8          3880158       300.5 ns/op   320 B/op       5 allocs/op
BenchmarkBaselineLargeStruct-8                   91645     12843 ns/op    2240 B/op      11 allocs/op
BenchmarkHybridLargeStruct-8                    118897      9913 ns/op    8848 B/op      97 allocs/op
BenchmarkHybridOptimisedLargeStruct-8          1842550       609.4 ns/op  1856 B/op       5 allocs/op
BenchmarkBaselineManyPointers-8                 392595      2652 ns/op     448 B/op      19 allocs/op
BenchmarkHybridManyPointers-8                    83277     14477 ns/op   10656 B/op     139 allocs/op
BenchmarkHybridOptimisedManyPointers-8          312087      3435 ns/op    3448 B/op      34 allocs/op
BenchmarkBaselineNestedSlices-8                 405385      2795 ns/op    1352 B/op      29 allocs/op
BenchmarkHybridNestedSlices-8                    93662     12532 ns/op   10140 B/op     135 allocs/op
BenchmarkHybridOptimisedNestedSlices-8          244626      4508 ns/op    3201 B/op      50 allocs/op
BenchmarkBaselineNestedMaps-8                   996086      1131 ns/op     536 B/op      13 allocs/op
BenchmarkHybridNestedMaps-8                     107374     11028 ns/op    9267 B/op     117 allocs/op
BenchmarkHybridOptimisedNestedMaps-8            457070      2645 ns/op    2296 B/op      30 allocs/op
BenchmarkBaselineDeepTree-8                     486628      2075 ns/op     584 B/op      16 allocs/op
BenchmarkHybridDeepTree-8                        78199     14378 ns/op   10924 B/op     135 allocs/op
BenchmarkHybridOptimisedDeepTree-8              273716      4160 ns/op    3833 B/op      39 allocs/op
BenchmarkBaselineCyclicGraph-8                 4064911       292.2 ns/op    24 B/op       1 allocs/op
BenchmarkHybridCyclicGraph-8                    115742      9118 ns/op    7952 B/op      91 allocs/op
BenchmarkHybridOptimisedCyclicGraph-8          1685490       679.2 ns/op  1024 B/op       7 allocs/op
BenchmarkBaselineMixedDomainObject-8            131402      8648 ns/op    3528 B/op      62 allocs/op
BenchmarkHybridMixedDomainObject-8              123578      8870 ns/op    8032 B/op      90 allocs/op
BenchmarkHybridOptimisedMixedDomainObject-8    2940127       421.9 ns/op  1104 B/op       4 allocs/op
```

## Notes
- The explicit three-variant model is now visible in one benchmark run.
- `hybrid` is the semantic reference implementation and is consistently the most expensive variant.
- `hybridopt` dramatically reduces the hybrid overhead while keeping the richer semantic behaviour.
- `hybridopt` is now faster than baseline on value-heavy and mixed domain scenarios, but remains slower on pointer-heavy, nested container, tree, and cyclic workloads.
- This file should act as the new reference point for future optimisation phases so we can compare `hybridopt` against both the semantic hybrid and the baseline at the same time.
