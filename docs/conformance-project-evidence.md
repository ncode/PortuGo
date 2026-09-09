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
