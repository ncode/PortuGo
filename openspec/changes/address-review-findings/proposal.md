## Why

A repository audit reproduced formatter comment relocation and found that individually bounded aggregates can collectively request unbounded storage. The recorded performance baseline also identifies three allocation costs that can be reduced without changing language behavior.

## What Changes

- Preserve comment attachment when a supported global `var` section follows subprogram declarations; retain equivalent execution and idempotent formatting in every CLI mode.
- Add a fixed per-interpreter aggregate-storage budget across globals, active call frames, and execution-time aggregate copies. Count nested backing cells and containers, reject before the allocation that exceeds the budget, and release reservations with their owning storage.
- Preserve the existing per-aggregate scalar-slot limit, all reference-qualified behavior, positioned diagnostics, and fresh-run/REPL isolation. Document the new budget as a project safeguard, separate from work budgets and host process memory.
- Reduce repeated vector-layout cloning in interpreter metadata queries while retaining independent full-layout results from semantic APIs.
- Reduce semantic symbol-map and builtin-catalog setup allocation, and evaluate bounded storage for ordinary builtin arguments without changing evaluation order, absence propagation, or halt behavior. Retain only measured improvements; the invocation-local argument-array candidate was rejected because it increased allocated bytes without reducing allocations.
- Require regression coverage, matched repeated allocation measurements, full quality checks, conformance replay, and privacy-screened code review before implementation handoff.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `tooling/diagnostics-cli-repl`: Preserve comment association across supported declaration orders and define the runtime aggregate-storage safeguard, its diagnostics, and storage lifetime rules. The three optimizations are implementation changes under existing semantic, call, and output contracts and introduce no new language requirements.

## Impact

Affected seams are canonical printing and comment emission, runtime aggregate-layout accounting, interpreter declaration/call/copy lifetimes, semantic-information queries and scope setup, builtin catalog traversal, and builtin call buffers. Update focused tests/fixtures, `docs/language.md`, `CHANGELOG.md`, and measured performance documentation. Preserve corpus source/input/evidence bytes and historical task mappings. Add no dependency, dialect feature, general heap profiler, user-configurable memory option, release publication, or unrelated refactor.
