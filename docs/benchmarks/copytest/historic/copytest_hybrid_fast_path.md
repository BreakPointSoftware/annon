# copytest benchmark: hybridopt fast path

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

## Change under test
- keep the cached repair-plan optimisation in `hybridopt`
- add a real fast path for repair-free structs
- still emit lightweight field-name flags on that path
- expand the benchmark scope with larger copy-volume batch scenarios

## Raw output
```text
goos: linux
goarch: amd64
pkg: github.com/BreakPointSoftware/annon/internal/copytest
cpu: Intel(R) Core(TM) i7-6700K CPU @ 4.00GHz
BenchmarkBaselineSmallStruct-8                   933229      1262 ns/op     224 B/op       4 allocs/op
BenchmarkHybridSmallStruct-8                     126297      9532 ns/op    7264 B/op      91 allocs/op
BenchmarkHybridOptimisedSmallStruct-8           3922754       308.8 ns/op   320 B/op       5 allocs/op
BenchmarkBaselineLargeStruct-8                   89738     13386 ns/op    2240 B/op      11 allocs/op
BenchmarkHybridLargeStruct-8                    114158     10732 ns/op    8848 B/op      97 allocs/op
BenchmarkHybridOptimisedLargeStruct-8          1897347       621.3 ns/op  1856 B/op       5 allocs/op
BenchmarkBaselineManyPointers-8                 408866      2725 ns/op     448 B/op      19 allocs/op
BenchmarkHybridManyPointers-8                    84360     14589 ns/op   10656 B/op     139 allocs/op
BenchmarkHybridOptimisedManyPointers-8          343297      3485 ns/op    3448 B/op      34 allocs/op
BenchmarkBaselineNestedSlices-8                 438732      2803 ns/op    1352 B/op      29 allocs/op
BenchmarkHybridNestedSlices-8                    88940     12760 ns/op   10140 B/op     135 allocs/op
BenchmarkHybridOptimisedNestedSlices-8          260338      4419 ns/op    3201 B/op      50 allocs/op
BenchmarkBaselineNestedMaps-8                  1058523      1227 ns/op     536 B/op      13 allocs/op
BenchmarkHybridNestedMaps-8                      95144     12509 ns/op    9267 B/op     117 allocs/op
BenchmarkHybridOptimisedNestedMaps-8            401407      2559 ns/op    2296 B/op      30 allocs/op
BenchmarkBaselineDeepTree-8                     525690      1958 ns/op     584 B/op      16 allocs/op
BenchmarkHybridDeepTree-8                        74463     17562 ns/op   10924 B/op     135 allocs/op
BenchmarkHybridOptimisedDeepTree-8              290524      4326 ns/op    3833 B/op      39 allocs/op
BenchmarkBaselineCyclicGraph-8                 4303690       296.9 ns/op    24 B/op       1 allocs/op
BenchmarkHybridCyclicGraph-8                    120289     11067 ns/op    7952 B/op      91 allocs/op
BenchmarkHybridOptimisedCyclicGraph-8          1430044       816.4 ns/op  1024 B/op       7 allocs/op
BenchmarkBaselineMixedDomainObject-8            130636      9063 ns/op    3528 B/op      62 allocs/op
BenchmarkHybridMixedDomainObject-8              126261     11353 ns/op    8032 B/op      90 allocs/op
BenchmarkHybridOptimisedMixedDomainObject-8    2576352       459.4 ns/op  1104 B/op       4 allocs/op
BenchmarkBaselineSmallStructBatch100-8            7760    129890 ns/op   22400 B/op      400 allocs/op
BenchmarkHybridSmallStructBatch100-8              1329    908460 ns/op  726402 B/op     9100 allocs/op
BenchmarkHybridOptimisedSmallStructBatch100-8    34662     31543 ns/op   32000 B/op      500 allocs/op
BenchmarkBaselineLargeStructBatch100-8             880   1263941 ns/op  224000 B/op     1100 allocs/op
BenchmarkHybridLargeStructBatch100-8              1137    995524 ns/op  884803 B/op     9700 allocs/op
BenchmarkHybridOptimisedLargeStructBatch100-8    18596     79934 ns/op  185600 B/op      500 allocs/op
BenchmarkBaselineValueOnlyBatch1000-8              985   1188142 ns/op  224000 B/op     4000 allocs/op
BenchmarkHybridValueOnlyBatch1000-8                135   8809469 ns/op 7264038 B/op    91000 allocs/op
BenchmarkHybridOptimisedValueOnlyBatch1000-8      3726    311260 ns/op  320000 B/op     5000 allocs/op
BenchmarkBaselineMixedDomainBatch100-8            1207    963763 ns/op  352803 B/op     6200 allocs/op
BenchmarkHybridMixedDomainBatch100-8              1130    970697 ns/op  803204 B/op     9000 allocs/op
BenchmarkHybridOptimisedMixedDomainBatch100-8    27190     46165 ns/op  110400 B/op      400 allocs/op
```

## Notes
- The fast path gives `hybridopt` a dramatic advantage on value-heavy shapes and at larger copy volumes.
- The larger-copy batch scenarios make the improvement much more obvious than the single-copy runs.
- `hybridopt` is now substantially faster than the semantic `hybrid` across all measured scenarios.
- The biggest wins are on:
  - small value-only structs
  - large value-only structs
  - mixed domain objects
  - repeated batch copies of the same value-heavy shapes
- `hybridopt` still remains slower than the baseline on pointer-heavy, nested container, tree, and cyclic workloads, but the gap is now narrow enough to make the semantic trade-off much more plausible.
- This phase strongly supports the original optimisation idea: copy first, then only walk/repair fields that actually require reference detachment or runtime-state handling.
