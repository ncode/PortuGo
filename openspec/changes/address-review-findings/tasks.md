## 1. Establish implementation and measurement baselines

- [x] 1.1 Read the change artifacts, `docs/language.md`, `docs/performance.md`, applicable repository/skill instructions, and the affected callers. Confirm the isolated worktree and unchanged audit baseline; preserve unrelated work.
- [x] 1.2 Capture fresh privacy-screened before measurements with the existing benchmark harness on a fixed host/toolchain, `GOMAXPROCS=1`, `-cpu=1`, `-benchmem`, and six repeated samples. Keep raw output outside the checkout and record the commands/environment needed for a matched after run.

## 2. Preserve formatter comment association

- [x] 2.1 Add a failing regression fixture for a procedure/function followed by global `var`, with leading, body, inline, closing, variable, and separator comments; verify the reproduced relocation before production changes. Include controls for globals-first, multiple subprograms, empty sections, and multiline fragments where relevant.
- [x] 2.2 Fix the shared printer's declaration emission order using existing section positions and comment handling. Keep constants/types and the existing globals-first canonical behavior intact; update only affected goldens.
- [x] 2.3 Verify stdout, `--check`, and `-w` modes, unchanged original/formatted execution, second-pass idempotence, and byte preservation on malformed/resource-limited inputs.

## 3. Bound live runtime storage

- [x] 3.1 Audit all interpreter allocation/copy callers and write failing small-limit or pre-allocation regressions for cumulative declarations, nested containers, checked arithmetic, and constant-dependent layouts. Keep existing scalar-slot and exact-boundary controls.
- [x] 3.2 Implement checked concrete-layout cost calculation and the fixed 2,097,152-unit per-run budget, charging cells and containers before allocation and reporting `R003` at the responsible declaration/call/copy. Keep counters interpreter-owned and preserve independent existing guards.
- [x] 3.3 Integrate reservation ownership for globals, parameters, function results, locals, nested aggregates, and execution-time copies. Release completed temporary/frame storage, avoid double charging captured references, retain charges for retained storage, and preserve destination values on failed copy reservation.
- [x] 3.4 Add regressions for recursion sharing the budget, repeated completed calls reusing storage, copied-record independence, captured field references, normal/error/halt unwinding, and fresh interpreter runs/REPL submissions. Verify nested oversize layouts reject without expanding their elements or requiring huge test allocations.

## 4. Remove repeated metadata copying from vector hot paths

- [x] 4.1 Add or retain semantic mutation-isolation and shared-info multi-interpreter regressions; identify the exact interpreter queries that need only immutable scalar facts.
- [x] 4.2 Implement the smallest narrow query/access changes that avoid routine layout clones while preserving independent full-layout results from `TypeOf` and `Binding`. Validate vector bounds, aliases, records, dynamic cells, and diagnostic positions.
- [x] 4.3 Measure vector execution against the matched baseline; retain the change only with a reproducible allocation reduction and no material unexplained regression.

## 5. Reduce semantic setup allocation

- [x] 5.1 Use cheap known capacities for semantic maps where beneficial and add read-only builtin-catalog traversal while preserving the existing copy-returning catalog contract, name-collision rules, and per-analysis isolation.
- [x] 5.2 Run registry/collision/isolation and analysis regressions, then compare repeated small/declarations analysis measurements. Remove complexity that does not produce a reproducible benefit.

## 6. Reduce builtin argument allocation

- [x] 6.1 Evaluate invocation-local bounded storage for ordinary builtin argument counts with a safe path for longer reachable lists; preserve left-to-right effects, conditional arity, nested calls, absence, retained operands, and halts. Reject the candidate if it fails the allocation gate.
- [x] 6.2 Verify builtin inventory and all argument-order/absence/halt regressions, inspect escape analysis, and compare repeated string-execution measurements. Keep only a confirmed allocation improvement. The local array escaped and increased bytes without reducing allocations, so the original slice remains.

## 7. Document, review, and validate the complete change

- [x] 7.1 Update `docs/language.md` and `CHANGELOG.md` for formatter order/comment preservation and the precise live-storage policy, including its separate scope from process memory and work budgets. Add established project-test traceability without changing reference source/input/evidence bytes or historical task mappings.
- [x] 7.2 Complete matched after measurements and benchstat comparisons, including the full benchmark workload/stage set. Update `docs/performance.md` with only privacy-screened measured results, method, and uncertainty; report each retained or rejected optimization explicitly.
- [x] 7.3 Run build, ordinary/race tests, vet, formatting checks, pinned Staticcheck and golangci-lint, 30-second lexer/parser fuzzing, strict OpenSpec validation, and full conformance evidence/trace/replay. Preserve all previously verified results and the two reviewed nonterminating exclusions.
- [x] 7.4 Privacy-screen the exact outgoing implementation diff and complete code review, addressing confirmed findings and reporting excluded files. Report final changed files, checks, benchmark evidence, limitations, and task status. Keep archiving separate from implementation.
