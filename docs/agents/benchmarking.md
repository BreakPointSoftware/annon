# Benchmarking Process

This file defines the benchmark recording process for the repository.

## Purpose

Benchmark outputs are used as a source of truth for optimisation work.

The process deliberately separates:

- raw benchmark output
- interpreted analysis

This makes the results:

- human reviewable
- AI reviewable
- easier to validate when optimisation decisions are being made

If the raw benchmark output and the analysis disagree, the raw benchmark output wins.

## Folder structure

Benchmarks are grouped by subsystem.

```text
docs/benchmarks/<subsystem>/
├── README.md
├── historic/
│   └── *.md
├── 1/
│   ├── raw_results.md
│   └── results_anylsis.md
├── 2/
│   ├── raw_results.md
│   └── results_anylsis.md
└── ...
```

## Meaning of each location

### `historic/`
Contains legacy benchmark markdown files created before the numbered-run process existed.

These files are preserved for history, but new benchmark runs must not be added there.

### Numbered run folders
Each new benchmark run gets the next integer folder:

- `1/`
- `2/`
- `3/`

Do not use timestamps or descriptive folder names.

## Required files per run

Each numbered run folder must contain exactly:

- `raw_results.md`
- `results_anylsis.md`

The spelling `results_anylsis.md` is intentional and should be kept consistent.

## `raw_results.md` requirements

Every raw results file must include:

1. subsystem under test
2. run number
3. commit SHA
4. branch name
5. Go version
6. OS/arch
7. CPU
8. date/time
9. exact benchmark command
10. raw benchmark output in a fenced code block

### Template

```md
# Raw Benchmark Results

## Context
- System under test: copytest
- Run number: 1
- Commit: <sha>
- Branch: <branch>
- Go version: <version>
- OS/Arch: <os/arch>
- CPU: <cpu>
- Date: <timestamp>

## Command
```bash
go test -bench=. ./internal/copytest/...
```

## Raw output
```text
<exact benchmark output>
```
```

## `results_anylsis.md` requirements

Every analysis file must include:

1. what changed since the previous run
2. the previous run number explicitly
3. a short summary
4. grouped benchmark tables
5. delta formatting
6. commentary
7. conclusion

### Required sections

```md
# Benchmark Analysis

## Scope
What changed in this run.

## Comparison Base
- Previous run: <n>
- Current run: <n+1>

## Summary
Short summary of the net result.

## Detailed Tables
### Value-heavy structs
| Benchmark | Previous ns/op | Current ns/op | Delta | Verdict |
|-----------|----------------|---------------|-------|---------|

### Pointer-heavy structs
...

### Nested containers
...

### Cyclic graphs
...

### Mixed domain objects
...

## Allocation Tables
Repeat the grouping for:
- B/op
- allocs/op

## Commentary
Explain likely causes of wins and losses.

## Conclusion
State whether the change is a net gain, net loss, or mixed.
```

## Grouping rules

Use the same benchmark groups in every analysis where applicable:

1. small value-heavy
2. large value-heavy
3. pointer-heavy
4. nested slices/maps
5. deep tree / recursive
6. cyclic graph
7. mixed domain object
8. batch / scale scenarios

## Delta formatting rules

### Time
- lower is faster
- higher is slower

Examples:
- `-38.2% faster`
- `+12.4% slower`

### Memory
Examples:
- `-1024 B/op`
- `+4096 B/op`

### Allocations
Examples:
- `-17 allocs/op`
- `+3 allocs/op`

### Verdict values
Use one of:
- `better`
- `worse`
- `flat`
- `mixed`

## Comparison rule

Each run compares only to the immediately previous run.

Examples:
- run `2` compares against run `1`
- run `3` compares against run `2`
- run `4` compares against run `3`

Do not compare against all historical runs by default.

## Review rule

- raw benchmark output is the source of truth
- analysis must be traceable back to the raw output
- if the interpretation is unclear, preserve that uncertainty in the analysis rather than over-claiming

## Command discipline

Before recording a benchmark run:

1. make the intended change
2. run relevant tests
3. ensure tests are green
4. run the benchmark command
5. create the next numbered run folder
6. save raw output first
7. then write the analysis against the previous run

## Example benchmark command

```bash
go test -bench=. ./internal/copytest/...
```

Subsytems may define their own commands in their local `README.md` under `docs/benchmarks/<subsystem>/`.
