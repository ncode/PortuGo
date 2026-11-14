## Context

See `proposal.md` for motivation and scope. The audit baseline is commit `cb78635`; ordinary/race tests, vet, pinned linters, formatting checks, lexer/parser fuzzing, and the conformance replay pass. Replay reports 1,716 verified entries and two reviewed non-applicable entries. Existing allocation evidence and correctness-checked workloads are in `docs/performance.md` and `performance_test.go`.

`ast.Fprint` currently emits global declarations before subprograms even when the supported source places `var` after them. Its forward-only comment stream then attaches comments to unrelated constructs. `runtime.Type.Slots` bounds each aggregate's scalar leaves, while declaration initialization independently allocates every variable. Full semantic queries clone mutable type layouts, scope setup allocates growing maps/catalog copies, and builtin calls allocate argument slices.

## Goals / Non-Goals

**Goals:** Fix the reproduced comment association, bound live language storage before allocation, and retain only simple allocation reductions supported by matched measurements. Keep semantic information safe to share across interpreter instances and preserve all call, output, diagnostic, and recovery contracts.

**Non-Goals:** A whole-process heap limit, source/semantic-memory accounting, configurable memory controls, pooling, new concurrency, a VM, language extensions, changes to oracle evidence, or speculative architecture work. State snapshots requested by Go callers remain independent copies under their current API; they are outside execution-time accounting.

## Decisions

### 1. Preserve the supported global declaration order

Use the retained section and subprogram positions to emit a trailing global `var` section after subprograms when that was its source position. Keep constants/types and other declaration rules intact. This lets the existing ordered comment machinery retain leading, body, inline, separator, and closing comments without a new ownership map. Preserve the existing canonical order for programs that already declare globals first.

Regress both orders, procedures/functions, multiple subprograms, empty `var`, multiline comment fragments, and all formatter modes. Compare output behavior and formatter idempotence. Reordering comments independently would require reconstructing ownership already represented by source positions and is unnecessary.

### 2. Bound concrete live storage, using the existing slot safeguard as a separate check

Use a fixed live budget of `2 * runtime.MaxVectorSlots` units: 2,097,152 cells/containers. One ordinary vector at the existing scalar-slot maximum still fits; multiple large live aggregates and deeply nested container layouts receive a controlled failure. This is a portable count, not a promise about RSS or the Go heap.

Compute units from the concrete resolved layout: a scalar cell costs one; a record costs its containing cell, its record container, and its field cells; a vector costs its containing cell, its vector container, and the recursive cost of each element multiplied by checked dimension widths. Reject arithmetic overflow or a cost above the remaining budget without expanding elements. Preserve `Slots`, its scalar-leaf policy, its existing codes, and constant-dependent bound resolution.

Keep reservation ownership in the interpreter's existing environment/call lifetimes. Reserve before variable/parameter/result allocation and before aggregate copies; release when their owning frame or temporary operation ends. Captured references charge the original storage once. Do not release storage still retained by a caller, global, pending value, or result. Audit every interpreter call to `Zero`, `NewVector`, `NewRecord`, `Clone`, and `AssignFrom`; do not introduce an unchecked alternate allocation path.

Allow a small checked accounting helper in the runtime layer, but keep live counters per interpreter and avoid making standalone runtime helpers depend on interpreter/global state. Use the smallest checked allocation/copy seam needed. A work-step charge is not a replacement: a single allocation can already request excessive storage. A cumulative allocation counter is also insufficient because repeated completed calls must be able to reuse released storage.

Test arithmetic/accounting at small limits and integration with rejecting layouts, recursion, normal/error/halt unwinding, copies, retained field references and builtin arguments, and fresh runs. Track aggregate argument roots separately from language operands so later argument calls cannot cause premature release or alter operand semantics. Avoid allocating gigabytes merely to test a boundary. Preserve destination values if copy reservation fails.

### 3. Reduce semantic query copying without weakening immutability

Keep `Info.TypeOf` and `Info.Binding` independent-copy behavior for callers that request full mutable layouts. Add only the narrow immutable facts needed by interpreter hot paths, or an equally small ownership-safe access path, so routine existence/type-kind/binding-ID checks do not deep-clone vector ranges and element types repeatedly. Full declaration/layout resolution still receives isolated metadata where mutation is needed.

Retain mutation-isolation tests and shared-info multi-interpreter tests. Reject an optimization that returns aliased mutable slices or adds per-expression caches whose setup cost outweighs the measured gain.

### 4. Optimize known setup and call allocations locally

Capture a fresh matched baseline before production changes. Size symbol/fact maps using cheap known declaration/catalog counts where profitable; do not add a separate expensive AST traversal just to calculate exact capacity. Provide read-only builtin catalog traversal while retaining the existing copy-returning `Catalog` API and collision/per-analysis isolation checks.

Evaluate small bounded argument storage for ordinary builtin arities with a safe dynamic path for any longer expression list reachable under current error/absence rules. Keep buffers local to the invocation so nested calls cannot overwrite arguments. Verify escape analysis and the existing builtin inventory, argument-order, absence, retained-operand, and halt regressions. Do not introduce pooling or change evaluation/validation order to make a benchmark faster.

The invocation-local three-value array was rejected after implementation measurements: it escapes through the existing call interface, leaves string-execution allocations unchanged at 1,018 per operation, and increases bytes from 111,180 to 162,382 per operation. Retain the existing argument slice rather than add reusable buffer ownership or pooling to pursue this target. The metadata-query and semantic-setup changes passed the allocation gate.

### 5. Gate each optimization on matched evidence

Use the existing six-workload/six-stage benchmark harness. Record before/after on the same toolchain and host with `GOMAXPROCS=1`, `-cpu=1`, `-benchmem`, at least six repeated samples, and comparable durations. Focus first on small/declarations analysis, vector execution, and string builtin execution; verify the full workload set for regressions. Compare with benchstat and retain an optimization only when allocation improvement is reproducible without a material unexplained regression. Make no latency claim unsupported by the timing uncertainty.

Keep raw results, profile binaries, local build paths, and operational metadata outside the repository. Publish only privacy-screened samples or concise measurements with their exact method. The baseline audit's existing percentages identify candidates, not promised speedups.

## Risks / Trade-offs

- **Missed ownership or copy path** → Audit allocation callers and test shared references, copy failure, recursive unwinding, retained values, and sequential storage reuse before marking the safeguard complete.
- **A new guard rejects a verified program** → Full replay remains mandatory; do not change corpus states or expectations to hide the failure. Review the policy explicitly if the recorded corpus establishes a conflict.
- **Optimization exposes mutable layouts** → Preserve existing copy APIs and adversarial mutation-isolation/race tests.
- **Comment fix changes canonical goldens** → Update only fixtures whose supported trailing declaration order requires the change; verify the original/formatted program pair.
- **Noisy timing creates a false claim** → Use repeated matched runs, report allocation counts separately, and document uncertainty.

## Migration Plan

No dependency or source-language migration is needed. Implement the formatter fix, storage safeguard, and three optimization targets in separately reviewable task groups. Update the language reference/changelog with the observable formatter and guard changes, and performance documentation only after measurement. Preserve source/input/evidence hashes and historical conformance task mappings. Add project-test traceability for the new tooling safeguard through the repository's established mechanism without inventing reference observations. Archive the change separately after integration.
