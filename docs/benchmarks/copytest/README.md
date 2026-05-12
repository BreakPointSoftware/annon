# copytest benchmark history

This directory contains benchmark history for the internal copy traversal spike.

## Subsystem under test

The `copytest` spike currently compares:

- `baseline`
- `hybrid`
- `hybridopt`

## Layout

```text
docs/benchmarks/copytest/
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

## historic

The `historic/` folder contains legacy benchmark markdown files created before the numbered-run process was introduced.

These files are preserved for context only.

## numbered runs

Each new benchmark run gets the next integer folder.

Examples:
- `1/`
- `2/`
- `3/`

Each run folder must contain:

- `raw_results.md`
- `results_anylsis.md`

## comparison rule

Each run compares only to the immediately previous run.

Examples:
- run `2` compares against run `1`
- run `3` compares against run `2`

## next run

The next structured benchmark run for `copytest` should be created under:

```text
docs/benchmarks/copytest/1/
```
