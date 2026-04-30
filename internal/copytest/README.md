# Copy Traversal A/B Spike

This internal spike compares two object copy/redaction traversal strategies.

## Why hybrid copy exists

The existing simple traversal walks values recursively and rebuilds structures one field at a time.
That is straightforward, but it loses unexported value fields and does not preserve shared graph
shape or cycles.

The hybrid approach exists to test whether a value-first struct copy gives us a better baseline for:

- preserving unexported value fields
- preserving shared graph structure
- handling cycles safely
- separating traversal/copy concerns from later redaction decisions

## Why structs are copied by value first

Struct assignment in Go preserves exported and unexported value fields automatically.
That means strings, ints, bools, arrays, and private fixed-size values can come across without
manually rebuilding the whole struct from zero.

## Why exported reference fields are repaired

After a shallow struct copy, exported reference fields would still point back to the original graph.
The hybrid walker repairs exported pointers, maps, slices, and interfaces so those references detach
from the source and preserve the copied graph shape.

## Why unexported reference fields are not deeply copied

Without `unsafe`, unexported fields cannot be set directly through reflection from another package.
The spike therefore preserves unexported value fields through the initial struct copy, but leaves
unexported reference fields shared and records flags so the trade-off is visible.

## Why runtime and resource types are treated specially

Types such as locks, channels, contexts, files, database handles, HTTP clients, and connection-like
objects represent runtime state rather than serialisable data.
Copying them blindly would be misleading for logging/export scenarios, so the hybrid walker zeros them
where it can and records explicit flags when it cannot.

## What the flag model is for

The hybrid copy records field-level flags so later stages can make decisions without forcing the full
redaction engine into this spike. The flags capture:

- where notable fields were found
- whether references were deep-copied or reused
- where unexported references remained shared
- where runtime state was zeroed
- where unsupported kinds were encountered

## What the tests are proving

The comparison tests are intended to make the trade-offs explicit.

- baseline:
  - preserves simple exported data
  - rejects recursive cycles rather than trying to preserve them
  - does not preserve shared graph structure
  - may lose unexported value semantics because it rebuilds structs field by field
- hybrid:
  - preserves unexported value fields through value-first struct copying
  - detaches copied references from the original object graph
  - preserves shared references and cycles inside the copied graph
  - zeroes or shares runtime-state fields according to explicit policy
  - emits flags that explain notable copy decisions

The tests therefore check three things together wherever relevant:

1. the copied graph is detached from the original
2. the copied graph preserves the intended internal shape
3. the copied graph still carries the same meaningful stored values unless policy says otherwise

## What is intentionally not solved yet

- no production integration with the public `redact` API
- no final redaction decision/apply pipeline built on top of the flags
- no `unsafe`-based deep copy of unexported references
- no claim that the hybrid strategy is production-ready

This spike exists to compare behaviour, correctness, and benchmark characteristics before committing
to a final internal design.
