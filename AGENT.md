# AGENT.md

This file captures repository-specific working guidance for AI/code agents.

## General

- Treat the repository as alpha-stage software.
- Prefer intent-revealing names over generic implementation names.
- Keep the public surface centred on `redact`.
- Internal packages are free to evolve and are not part of the compatibility contract.

## Style

- Prefer self-documenting variable names over short ambiguous names.
- Use vertical spacing to separate logical phases.
- Add short comments above important logic blocks when they improve scanning.
- Do not comment trivial assignments.
- Prefer stronger types over raw strings where that meaningfully improves correctness.

## Refactoring

- Prefer small, reviewable refactors.
- Keep package boundaries explicit.
- Avoid vague dumping-ground packages like `core`.
- Favour internal concepts such as `engine`, `walk`, `decision`, `output`, and `fallback` when they match intent.

## Testing

- Keep `go test ./...` green after each meaningful slice.
- Prefer table-driven tests.
- Use `testify` for assertions.
- Add fuzz tests when working on defensive/no-error public APIs.

## Low-Maintenance Mode

When working in low-maintenance mode, always use the following workflow:

1. Create a feature branch before making changes.
2. Keep the change scoped to a small PR-sized slice.
3. Run the relevant test/build commands.
4. Ensure tests are green before proposing completion.
5. Commit the change with a focused message.
6. Push the branch to `origin`.
7. Raise a GitHub pull request with `gh pr create`.
8. Write a detailed PR body that explains:
   - what changed
   - why it changed
   - behaviour notes
   - testing performed
   - follow-up work

### Low-Maintenance Rules

- Do not work directly on `main`.
- Do not stop after local implementation if the mode expects a PR workflow.
- Do not open a PR before tests are green.
- Prefer one coherent vertical slice per PR.
- Keep README/examples updated when public behaviour changes.

## Example Commands

- Test everything:
  - `go test ./...`
- Build examples:
  - `go build ./examples/...`
- Run benchmarks:
  - `go test -bench=. ./...`

## Current Public Direction

- Public entrypoint: `github.com/BreakPointSoftware/annon/redact`
- Public APIs are intended to be defensive and no-error.
- Examples live under `examples/`.
