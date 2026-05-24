# copytest benchmark: hybrid cached repair plan

## Context
- Commit: `14d046c` + local phase 1 metadata cache changes
- Branch: `feat/copy-spike-test-clarity`
- Go version: `go1.26.2`
- OS/Arch: `linux/amd64`
- CPU: `Intel(R) Core(TM) i7-6700K CPU @ 4.00GHz`

## Command
```bash
go test -bench=. ./internal/copytest/...
```

## Change under test
- added cached per-type repair plans for the hybrid strategy
- moved field classification work out of the hot copy path and into one-time metadata compilation

## Raw output
```text
goos: linux
goarch: amd64
pkg: github.com/BreakPointSoftware/annon/internal/copytest
cpu: Intel(R) Core(TM) i7-6700K CPU @ 4.00GHz
BenchmarkBaselineSmallStruct-8          901602      1144 ns/op     224 B/op       4 allocs/op
BenchmarkHybridSmallStruct-8           3810298       337.7 ns/op   320 B/op       5 allocs/op
BenchmarkBaselineLargeStruct-8           78036     13139 ns/op    2240 B/op      11 allocs/op
BenchmarkHybridLargeStruct-8           1953075       615.6 ns/op  1856 B/op       5 allocs/op
BenchmarkBaselineManyPointers-8         407937      2699 ns/op     448 B/op      19 allocs/op
BenchmarkHybridManyPointers-8           295788      3459 ns/op    3448 B/op      34 allocs/op
BenchmarkBaselineNestedSlices-8         446452      3014 ns/op    1352 B/op      29 allocs/op
BenchmarkHybridNestedSlices-8           243122      4574 ns/op    3201 B/op      50 allocs/op
BenchmarkBaselineNestedMaps-8           872635      1163 ns/op     536 B/op      13 allocs/op
BenchmarkHybridNestedMaps-8             489481      2545 ns/op    2296 B/op      30 allocs/op
BenchmarkBaselineDeepTree-8             491055      2040 ns/op     584 B/op      16 allocs/op
BenchmarkHybridDeepTree-8               272416      4050 ns/op    3833 B/op      39 allocs/op
BenchmarkBaselineCyclicGraph-8         4161991       283.5 ns/op    24 B/op       1 allocs/op
BenchmarkHybridCyclicGraph-8           1597614       700.3 ns/op  1024 B/op       7 allocs/op
BenchmarkBaselineMixedDomainObject-8    134280      8885 ns/op    3528 B/op      62 allocs/op
BenchmarkHybridMixedDomainObject-8     2923461       417.9 ns/op  1104 B/op       4 allocs/op
```

## Notes
- This optimisation dramatically reduced hybrid overhead for value-heavy structs and mixed domain objects.
- The biggest gains came from avoiding repeated field classification work in `copyStruct`.
- Hybrid is now faster than baseline on:
  - small value-only structs
  - large value-only structs
  - mixed domain object
- Hybrid remains slower on pointer-heavy, slice-heavy, map-heavy, tree, and cyclic workloads, but the gap is much smaller than in the initial run.
- This phase strongly suggests that metadata caching is worth keeping even if later optimisation steps stop here.
