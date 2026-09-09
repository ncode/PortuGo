# Project conformance evidence

These mappings cover project tooling and defensive guards. Running
the reference application cannot validate a JSON manifest, a Git history check,
or the project's recording and replay code. Each reference non-applicability
decision below therefore names existing project tests. The quality suite runs
those tests; corpus validation checks that their trace links still exist.

No language feature or bundled example is exempted by these decisions. Missing
reference recordings, inventory entries, and implementation work still block the
corresponding conformance gates.

## Provenance validation

`project.provenance` covers required provenance fields, artifact hashes, safe
paths, and exclusion of reference binaries and archives. These are properties of
the corpus validator. `TestManifestValidation`, `TestManifestRejectsSymlink`, and
`TestLoadRejectsUnknownFieldsAndTrailingJSON` exercise their failure cases.
Reference version and acquisition metadata remain mandatory for recorded probes.

## Recording and normalization

`project.recording` covers private staging, source/input integrity, capture
metadata, required GUI attachments, and generated-file bytes. Its tests are
`TestPrepareRecording`, `TestCaptureRecording`, `TestCaptureGeneratedBytes`, and
`TestCaptureInitialFiles`. `project.normalization` uses `TestNormalize` to verify
that only the declared envelope and line endings change, including rejection of
invalid envelopes and preservation of significant output bytes.

These tests verify the recording tools. They do not establish the reference
result of an unrecorded program. Recorded language probes still require the
actual source, input, capture time, and reviewed observation artifacts.

## Traceability and phase validation

`project.traceability` covers missing and stale links, independent discovery of
specifications, unique IDs, and preservation of probe, inventory, task, and
retirement history. Its tests are `TestManifestValidation`,
`TestManifestPreservesVerifiedHistory`, `TestManifestPreservesRetiredHistory`,
`TestManifestRetirementIDs`, and `TestCommandHistoryValidation`.

`project.validation-phases` covers recorded/pending evidence, completed owner
groups, acceptance rejection of pending implementation, reference disposition,
and the distinction between pending mismatches and verified regressions. Its
tests are `TestManifestValidation`, `TestManifestReferenceDisposition`,
`TestReplay`, and `TestReplaySummary`.

These mappings verify enforcement by the tooling. They do not declare the
inventory complete, satisfy missing checklist/audit inputs, or authorize a
conformance release. Release acceptance and the example sweep remain pending.

## CLI contracts

`project.cli-contracts` links `TestCommandContracts` to the CLI command and exit
status requirements. The subprocess table checks `run`, `check`, `fmt`, and
`repl`, stdout/stderr separation, statuses 0/1/2, invalid arguments, source and
formatting limits, retained output before failure, and execution budgets.
These command-line interfaces do not exist in the reference editor.

## Positioned diagnostics

`project.positioned-diagnostics` covers diagnostic severity, stable codes,
deterministic rendering, inspectable causes, runtime failure positions, and
rejection of missing or stale semantic information. `TestDiagnosticContract`
checks rendering and prevents underlying host error text from leaking into the
stable message. `TestExecutionDiagnostics` exercises arithmetic, indexing,
input, loop, built-in, and output failures; `TestRejectMissingOrStaleInfo` checks
the semantic handoff. The parser recovery tests below check ordered syntax
diagnostics. Reference message and language-behavior observations remain
separate; these mappings cover the project's diagnostic API.

## Execution limits

`project.execution-limits` covers the explicit work budget, active-call limit,
evaluation-depth limit, text and format allocation limits, and reuse after a
failed run. Its tests are `TestStepBudgetAndReuse`, `TestCallAndValueLimits`,
`TestDefensiveEvaluationDepth`, `TestFormattedItemBoundary`, and
`TestBuiltinTextLimits`. These safeguards are project policies, not inferred
reference limits. Recorded programs still fail acceptance if a configured
limit prevents their required outcome.

## Front-end recovery and limits

`project.frontend-recovery-limits` covers retained statements between independent
syntax errors, vector-bound recovery, source and syntax-depth limits, and
controlled processing of the bounded adversarial test profile. It links
`TestStatementLineRecovery`, `TestVectorBoundRecovery`, `TestStructuralLimits`,
the lexer and parser `TestFuzzAdversarial` tests, and REPL `TestSubmissionLimit`.
These tests qualify the project's collected diagnostics and resource policy;
they do not claim that reference execution discovers errors at the same phase.

The tooling requirement links now include these existing tests and the release
phase validator. This closes eight missing trace links without exempting any
language behavior or example. Forty recorded implementation differences and
55 missing mappings remain: 19 requirements and 36 bundled examples. In-place
formatting, comment preservation, source-byte mapping, and the remaining
language and host behavior continue to require their own implementation and
evidence. A linked requirement does not by itself mean its tasks are complete.

## Defensive vector storage

`project.vector-storage` covers invalid or inconsistent Go storage objects,
including missing metadata, truncated or oversized backing slices, reversed
ranges, and overflowing widths or dimension products. Such objects cannot be
injected through a reference source program. `TestVectorCellRejectsInvalidStorage`
checks controlled rejection without mutation, and `TestVectorCellOffsets` checks
every element in small layouts and offsets near both integer extremes.

`TestCorruptedVectorStorage` injects truncated backing storage during execution
and checks reads, assignments, input, and reference arguments. Each operation
returns `R003` at its indexing expression, preserves preceding output and
elements, and leaves input unread. Recorded ordinary bounds failures remain
separately required and linked; this project guard does not exempt language
recordings, allocation accounting, or maximum-size validation.
