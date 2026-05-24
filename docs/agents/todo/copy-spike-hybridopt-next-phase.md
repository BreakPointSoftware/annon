# Copy Spike Optimisation Next Phase

## Status
Paused pending a review of how benchmark markdown files are being handled.

## Current spike state
The copy traversal spike currently has three simultaneous variants:

- `baseline`
- `hybrid`
- `hybridopt`

Benchmark history currently exists under:
- `docs/benchmarks/copytest/historic/copytest_baseline_vs_hybrid_initial.md`
- `docs/benchmarks/copytest/historic/copytest_hybrid_cached_repair_plan.md`
- `docs/benchmarks/copytest/historic/copytest_three_variant_comparison.md`
- `docs/benchmarks/copytest/historic/copytest_hybrid_fast_path.md`

## Next planned optimisation phase

### Goal
Add optional flag collection mode to `hybridopt` so we can compare:
- copy + repair + flags
- copy + repair without flags

### Proposed API
In `internal/copytest/hybridopt`:

```go
type Options struct {
    CollectFlags bool
}
```

```go
func Copy[T any](input T) (HybridCopyResult[T], error)
func CopyWithOptions[T any](input T, options Options) (HybridCopyResult[T], error)
```

Default:
- `Copy(...)` keeps current behaviour
- `CollectFlags: true`

### Implementation tasks
1. Add options type
2. Thread options into `hybridWalker`
3. Guard all flag appends behind `CollectFlags`
4. Avoid flag-only work when flags are off where practical
5. Keep copied output semantics identical regardless of flag mode

### Tests to add
- `TestHybridFlagsDisabledProducesNoFlags`
- `TestHybridFlagsDisabledPreservesCopySemantics`

### Benchmarks to add
- `BenchmarkHybridOptimisedFlagsOnSmallStruct`
- `BenchmarkHybridOptimisedFlagsOffSmallStruct`
- `BenchmarkHybridOptimisedFlagsOnLargeStruct`
- `BenchmarkHybridOptimisedFlagsOffLargeStruct`
- `BenchmarkHybridOptimisedFlagsOnManyPointers`
- `BenchmarkHybridOptimisedFlagsOffManyPointers`
- `BenchmarkHybridOptimisedFlagsOnMixedDomainObject`
- `BenchmarkHybridOptimisedFlagsOffMixedDomainObject`

Batch variants:
- `BenchmarkHybridOptimisedFlagsOnSmallStructBatch100`
- `BenchmarkHybridOptimisedFlagsOffSmallStructBatch100`
- `BenchmarkHybridOptimisedFlagsOnLargeStructBatch100`
- `BenchmarkHybridOptimisedFlagsOffLargeStructBatch100`
- `BenchmarkHybridOptimisedFlagsOnValueOnlyBatch1000`
- `BenchmarkHybridOptimisedFlagsOffValueOnlyBatch1000`
- `BenchmarkHybridOptimisedFlagsOnMixedDomainBatch100`
- `BenchmarkHybridOptimisedFlagsOffMixedDomainBatch100`

### Commands to run
```bash
go test ./internal/copytest/...
go test ./...
go test -bench=. ./internal/copytest/...
```

### Benchmark run folder currently planned
```text
docs/benchmarks/copytest/1/
├── raw_results.md
└── results_anylsis.md
```

### Acceptance criteria
- flag-on and flag-off modes both work
- default behaviour unchanged
- copied values equivalent between both modes
- benchmark output captured
- full suite green

## Benchmark process now chosen

Benchmarks will now be stored under numbered run folders per subsystem.

For `copytest`, the structure is:

```text
docs/benchmarks/copytest/
├── README.md
├── historic/
│   ├── copytest_baseline_vs_hybrid_initial.md
│   ├── copytest_hybrid_cached_repair_plan.md
│   ├── copytest_three_variant_comparison.md
│   └── copytest_hybrid_fast_path.md
├── 1/
│   ├── raw_results.md
│   └── results_anylsis.md
└── ...
```

## Suggested next conversation
Pause optimisation work and first decide:
- benchmark doc structure
- level of detail per file
- whether to keep raw benchmark output inline or split it
