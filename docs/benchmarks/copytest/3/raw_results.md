# Raw Benchmark Results

## Context
- System under test: copytest
- Run number: 3
- Commit: `713b8ba`
- Branch: `feat/copy-spike-three-variants`
- Go version: `go1.26.2`
- OS/Arch: `linux/amd64`
- CPU: `Intel(R) Core(TM) i7-6700K CPU @ 4.00GHz`
- Date: `2026-05-24T21:51:05+01:00`

## Command
```bash
go test -bench=. ./internal/copytest/...
```

## Raw output
```text
goos: linux
goarch: amd64
pkg: github.com/BreakPointSoftware/annon/internal/copytest
cpu: Intel(R) Core(TM) i7-6700K CPU @ 4.00GHz
BenchmarkBaselineSmallStruct-8                   852494      1436 ns/op     224 B/op       4 allocs/op
BenchmarkHybridSmallStruct-8                     102823     11009 ns/op    7264 B/op      91 allocs/op
BenchmarkHybridOptimisedSmallStruct-8           3184226       390.1 ns/op   320 B/op       5 allocs/op
BenchmarkBaselineLargeStruct-8                   61826     17153 ns/op    2240 B/op      11 allocs/op
BenchmarkHybridLargeStruct-8                     85266     12779 ns/op    8848 B/op      97 allocs/op
BenchmarkHybridOptimisedLargeStruct-8          1538436       826.3 ns/op  1856 B/op       5 allocs/op
BenchmarkBaselineManyPointers-8                303898      3627 ns/op     448 B/op      19 allocs/op
BenchmarkHybridManyPointers-8                   57439     18522 ns/op   10656 B/op     139 allocs/op
BenchmarkHybridOptimisedManyPointers-8         263264      4832 ns/op    3448 B/op      34 allocs/op
BenchmarkBaselineNestedSlices-8                358306      3595 ns/op    1352 B/op      29 allocs/op
BenchmarkHybridNestedSlices-8                   68631     16434 ns/op   10140 B/op     135 allocs/op
BenchmarkHybridOptimisedNestedSlices-8         188724      5739 ns/op    3201 B/op      50 allocs/op
BenchmarkBaselineNestedMaps-8                  790747      1509 ns/op     536 B/op      13 allocs/op
BenchmarkHybridNestedMaps-8                     79191     14766 ns/op    9268 B/op     117 allocs/op
BenchmarkHybridOptimisedNestedMaps-8           326344      3201 ns/op    2296 B/op      30 allocs/op
BenchmarkBaselineDeepTree-8                    433684      2575 ns/op     584 B/op      16 allocs/op
BenchmarkHybridDeepTree-8                       62030     18505 ns/op   10924 B/op     135 allocs/op
BenchmarkHybridOptimisedDeepTree-8            208592      5193 ns/op    3833 B/op      39 allocs/op
BenchmarkBaselineCyclicGraph-8                3029588       368.7 ns/op    24 B/op       1 allocs/op
BenchmarkHybridCyclicGraph-8                   112276     11014 ns/op    7952 B/op      91 allocs/op
BenchmarkHybridOptimisedCyclicGraph-8         1301499       920.7 ns/op  1024 B/op       7 allocs/op
BenchmarkBaselineMixedDomainObject-8           103777     11185 ns/op    3528 B/op      62 allocs/op
BenchmarkHybridMixedDomainObject-8             109929     11471 ns/op    8032 B/op      90 allocs/op
BenchmarkHybridOptimisedMixedDomainObject-8   2075560       570.1 ns/op  1104 B/op       4 allocs/op
BenchmarkBaselineSmallStructBatch100-8           8162    147364 ns/op   22400 B/op     400 allocs/op
BenchmarkHybridSmallStructBatch100-8             1053   1168256 ns/op  726405 B/op    9100 allocs/op
BenchmarkHybridOptimisedSmallStructBatch100-8   31916     36630 ns/op   32000 B/op     500 allocs/op
BenchmarkBaselineLargeStructBatch100-8            684   1620106 ns/op  224000 B/op    1100 allocs/op
BenchmarkHybridLargeStructBatch100-8              789   1299852 ns/op  884804 B/op    9700 allocs/op
BenchmarkHybridOptimisedLargeStructBatch100-8   13573     82204 ns/op  185600 B/op     500 allocs/op
BenchmarkBaselineValueOnlyBatch1000-8             703   1509905 ns/op  224000 B/op    4000 allocs/op
BenchmarkHybridValueOnlyBatch1000-8               100  10909407 ns/op 7264032 B/op   91000 allocs/op
BenchmarkHybridOptimisedValueOnlyBatch1000-8     2872    394308 ns/op  320001 B/op    5000 allocs/op
BenchmarkBaselineMixedDomainBatch100-8           1117   1098159 ns/op  352803 B/op    6200 allocs/op
BenchmarkHybridMixedDomainBatch100-8              954   1205015 ns/op  803206 B/op    9000 allocs/op
BenchmarkHybridOptimisedMixedDomainBatch100-8   17960     57241 ns/op  110400 B/op     400 allocs/op
BenchmarkHybridOptimisedFlagsOnSmallStruct-8   3118810       392.7 ns/op   320 B/op       5 allocs/op
BenchmarkHybridOptimisedFlagsOffSmallStruct-8  4125782       293.5 ns/op   240 B/op       4 allocs/op
BenchmarkHybridOptimisedFlagsOnLargeStruct-8   1536249       814.6 ns/op  1856 B/op       5 allocs/op
BenchmarkHybridOptimisedFlagsOffLargeStruct-8  1729322       719.8 ns/op  1776 B/op       4 allocs/op
BenchmarkHybridOptimisedFlagsOnManyPointers-8   215508      5043 ns/op   3448 B/op      34 allocs/op
BenchmarkHybridOptimisedFlagsOffManyPointers-8  394040      2782 ns/op    912 B/op      21 allocs/op
BenchmarkHybridOptimisedFlagsOnMixedDomainObject-8 1986723   590.3 ns/op 1104 B/op       4 allocs/op
BenchmarkHybridOptimisedFlagsOffMixedDomainObject-8 2079709  636.8 ns/op 1104 B/op       4 allocs/op
BenchmarkHybridOptimisedFlagsOnSmallStructBatch100-8 33032 40545 ns/op  32000 B/op     500 allocs/op
BenchmarkHybridOptimisedFlagsOffSmallStructBatch100-8 38049 28551 ns/op 24000 B/op     400 allocs/op
BenchmarkHybridOptimisedFlagsOnLargeStructBatch100-8 15009 86769 ns/op 185600 B/op     500 allocs/op
BenchmarkHybridOptimisedFlagsOffLargeStructBatch100-8 17826 69371 ns/op 177600 B/op     400 allocs/op
BenchmarkHybridOptimisedFlagsOnValueOnlyBatch1000-8 2977 400017 ns/op 320000 B/op 5000 allocs/op
BenchmarkHybridOptimisedFlagsOffValueOnlyBatch1000-8 3544 284306 ns/op 240000 B/op 4000 allocs/op
BenchmarkHybridOptimisedFlagsOnMixedDomainBatch100-8 19392 55962 ns/op 110400 B/op 400 allocs/op
BenchmarkHybridOptimisedFlagsOffMixedDomainBatch100-8 21393 54796 ns/op 110400 B/op 400 allocs/op
```
