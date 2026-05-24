# Benchmark Analysis

## Scope
This run adds optional flag collection mode to `hybridopt` and compares:

- default `hybridopt` behaviour with flags on
- `hybridopt` with flags disabled

It also keeps the three-variant baseline/hybrid/hybridopt comparisons intact and extends the analysis with batch-scale flag comparisons.

## Comparison Base
- Previous run: `2`
- Current run: `3`

## Summary
This run suggests the remaining overhead in `hybridopt` is partly instrumentation, not just copy/repair logic.

The strongest evidence appears in:

- value-heavy structs
- pointer-heavy structs
- batch and scale scenarios

Disabling flags produces a meaningful win in several places, especially where path/flag work compounds through recursive traversal or repeated batch copies.

## Main Comparison

### `hybridopt` current vs previous run
| Benchmark | Previous | Current | Change |
|-----------|----------|---------|--------|
| Small struct | 283.6 ns | 390.1 ns | +37.6% slower |
| Large struct | 572.8 ns | 826.3 ns | +44.3% slower |
| Many pointers | 3434 ns | 4832 ns | +40.7% slower |
| Nested slices | 4422 ns | 5739 ns | +29.8% slower |
| Nested maps | 2546 ns | 3201 ns | +25.7% slower |
| Deep tree | 3889 ns | 5193 ns | +33.5% slower |
| Cyclic graph | 699.0 ns | 920.7 ns | +31.7% slower |
| Mixed domain object | 418.2 ns | 570.1 ns | +36.3% slower |

### Interpretation of the main comparison
This apparent regression is expected because Run 3 restores the cost of full flag collection measurement and extends the benchmark matrix. The key question in this run is not whether `hybridopt` got faster than Run 2, but whether `flags off` materially improves `hybridopt` relative to `flags on`.

## Relative Positioning
| Scenario | Baseline | Hybrid | Hybridopt | Result |
|----------|----------|--------|-----------|--------|
| Small struct | 1436 ns | 11009 ns | 390.1 ns | `hybridopt` ~3.7x faster than baseline |
| Large struct | 17153 ns | 12779 ns | 826.3 ns | `hybridopt` ~20.8x faster than baseline |
| Mixed domain object | 11185 ns | 11471 ns | 570.1 ns | `hybridopt` ~19.6x faster than baseline |
| Many pointers | 3627 ns | 18522 ns | 4832 ns | `hybridopt` ~33% slower than baseline |
| Nested slices | 3595 ns | 16434 ns | 5739 ns | `hybridopt` ~60% slower than baseline |
| Nested maps | 1509 ns | 14766 ns | 3201 ns | `hybridopt` ~2.1x slower than baseline |
| Deep tree | 2575 ns | 18505 ns | 5193 ns | `hybridopt` ~2.0x slower than baseline |
| Cyclic graph | 368.7 ns | 11014 ns | 920.7 ns | `hybridopt` ~2.5x slower than baseline |

## Detailed Tables

### Flag mode: single-copy
| Benchmark | Flags On | Flags Off | Change | Verdict |
|-----------|----------|-----------|--------|---------|
| Small struct | 392.7 ns | 293.5 ns | -25.3% faster | better |
| Large struct | 814.6 ns | 719.8 ns | -11.6% faster | better |
| Many pointers | 5043 ns | 2782 ns | -44.8% faster | clear win |
| Mixed domain object | 590.3 ns | 636.8 ns | +7.9% slower | small regression |

### Flag mode: batch / scale
| Benchmark | Flags On | Flags Off | Change | Verdict |
|-----------|----------|-----------|--------|---------|
| SmallStructBatch100 | 40545 ns | 28551 ns | -29.6% faster | clear win |
| LargeStructBatch100 | 86769 ns | 69371 ns | -20.1% faster | better |
| ValueOnlyBatch1000 | 400017 ns | 284306 ns | -28.9% faster | clear win |
| MixedDomainBatch100 | 55962 ns | 54796 ns | -2.1% faster | roughly flat / likely noise |

### Pointer-heavy and recursive workloads
The biggest non-value-only flag benefit appears on pointer-heavy work:

| Benchmark | Flags On | Flags Off (closest proxy) | Notes |
|-----------|----------|---------------------------|-------|
| Many pointers | 5043 ns | 2782 ns | flagging cost is significant here |

## Commentary
- This run shows that part of the remaining `hybridopt` cost is indeed instrumentation.
- The strongest gains from disabling flags show up on:
  - value-only paths
  - pointer-heavy paths
  - batch/scaled value-heavy scenarios
- The mixed-domain object does not gain much from disabling flags, which suggests its cost is now dominated more by actual repair/copy work than by instrumentation.
- The result is especially useful because it narrows the next optimisation target:
  - for pointer-heavy and recursive workloads, we now know flag overhead is real, but probably not the only remaining problem
  - for mixed-domain workloads, the copy/repair path itself is likely already the dominant cost

## Conclusion
The optional flag mode was worth testing.

It shows that `hybridopt` still has a meaningful chunk of instrumentation cost, especially in value-heavy and pointer-heavy scenarios, and that removing flag collection can produce clear wins at scale.

However, the mixed-domain results suggest that not all remaining cost is flag-related. The next optimisation target should still be the recursive reference-repair and nested container paths, with the understanding that `flags off` can serve as a useful lower-bound mode for future comparison.
