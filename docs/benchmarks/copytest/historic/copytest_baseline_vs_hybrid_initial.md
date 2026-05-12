# copytest benchmark: baseline vs hybrid initial

## Context
- Commit: `14d046c`
- Branch: `feat/copy-spike-test-clarity`
- Go version: `go1.26.2`
- OS/Arch: `linux/amd64`
- CPU: `Intel(R) Core(TM) i7-6700K CPU @ 4.00GHz`

## Command
```bash
go test -bench=. ./internal/copytest/...
```

## Change under test
- baseline recursive reflection copy
- initial hybrid value-copy-first copy with graph repair and flag collection

## Raw output
```text
goos: linux
goarch: amd64
pkg: github.com/BreakPointSoftware/annon/internal/copytest
cpu: Intel(R) Core(TM) i7-6700K CPU @ 4.00GHz
BenchmarkBaselineSmallStruct-8          1077772      1164 ns/op     224 B/op       4 allocs/op
BenchmarkHybridSmallStruct-8             131674      8540 ns/op    7264 B/op      91 allocs/op
BenchmarkBaselineLargeStruct-8            82869     13481 ns/op    2240 B/op      11 allocs/op
BenchmarkHybridLargeStruct-8             101374     10866 ns/op    8848 B/op      97 allocs/op
BenchmarkBaselineManyPointers-8          444321      2830 ns/op     448 B/op      19 allocs/op
BenchmarkHybridManyPointers-8             74182     15643 ns/op   10656 B/op     139 allocs/op
BenchmarkBaselineNestedSlices-8          329096      3163 ns/op    1352 B/op      29 allocs/op
BenchmarkHybridNestedSlices-8             83473     14107 ns/op   10140 B/op     135 allocs/op
BenchmarkBaselineNestedMaps-8            836865      1206 ns/op     536 B/op      13 allocs/op
BenchmarkHybridNestedMaps-8               94587     12828 ns/op    9267 B/op     117 allocs/op
BenchmarkBaselineDeepTree-8              575776      2314 ns/op     584 B/op      16 allocs/op
BenchmarkHybridDeepTree-8                 74166     15520 ns/op   10924 B/op     135 allocs/op
BenchmarkBaselineCyclicGraph-8          3887242       316.9 ns/op    24 B/op       1 allocs/op
BenchmarkHybridCyclicGraph-8             123954      9241 ns/op    7952 B/op      91 allocs/op
BenchmarkBaselineMixedDomainObject-8     125991      9315 ns/op    3528 B/op      62 allocs/op
BenchmarkHybridMixedDomainObject-8       127683      9387 ns/op    8032 B/op      90 allocs/op
```

## Notes
- The hybrid implementation preserves significantly more semantic information than the baseline.
- The baseline is substantially cheaper on allocations in almost every benchmark.
- The hybrid is competitive on the large value-heavy and mixed domain scenarios, but still materially more allocation-heavy.
- This file acts as the pre-optimisation starting point for later metadata caching and fast-path work.
