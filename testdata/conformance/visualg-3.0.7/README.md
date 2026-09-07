# VisuAlg 3.0.7 reference corpus

This directory is being assembled by OpenSpec group 2. It does not yet contain
the complete evidence inventory and does not establish full conformance.
The earlier selected observations remain in `docs/validation/visualg-2026-09-07`.

## Commands

Run from the repository root with Go 1.22 or newer:

```sh
go test ./scripts/conformance
go run ./scripts/conformance validate --mode evidence --base origin/main
go run ./scripts/conformance validate --mode incremental --base origin/main
go run ./scripts/conformance validate --mode implementation-acceptance --base origin/main
```

`validate` verifies metadata and hashes, builds the CLI from the current checkout,
and replays recorded probes in isolated temporary working directories. Use
`--candidate` to supply an already built executable for CLI-only observations.
Each probe has a 1–30,000 ms deadline and a finite `maxSteps` budget (10,000 when
omitted or zero); source is limited to 64 KiB, each captured output stream to 1 MiB, and
each evidence file to 16 MiB. A deadline or output limit is a failed observation.
Pending results, including mismatches, remain visible in the JSON report.
Verified mismatches fail every phase. Acceptance rejects all pending behavior.
State and host-event observations use the current checkout's deterministic
execution adapter in a subprocess, with the same runtime guards. The adapter
uses a fixed PCG seed, a clock starting at the Unix epoch, and recorded,
nonblocking host operations. State JSON maps canonical global names to scalar
values or vector objects with `bounds` and flattened `values`. Host JSON is an
ordered array of typed operation records; no language host commands are yet
implemented, so current program traces are empty. Each channel is byte-compared
with its hashed expected artifact. Observation files are capped at 1 MiB.
State/host probes require the default current-checkout candidate so an unrelated
external executable cannot be credited with the adapter's observations.

Every validation requires either `--base` or `--previous` for history checks.
`--base` checks an earlier Git manifest for removed probes, inventory entries,
linked tasks, retirement records, and verified-to-pending downgrades;
use the PR's base commit or parent branch for stacked work instead of `origin/main`.
`--previous` accepts a repository-relative manifest file instead. A missing
manifest in the initial base is allowed; invalid refs and other Git failures are
errors. CI fetches history and compares against the PR base or pre-push commit;
manual workflow runs compare against the previous commit. The shell gate accepts
the same flags: `bash scripts/check-oracle.sh --base origin/main`.

## Recording on Windows

The recorder works with the official application in an interactive Windows
session. Install the reference distribution outside the repository and verify
its hashes against the manifest before recording. Keep the acquisition URL,
date, executable/archive hashes, product version, OS and locale in provenance;
machine names, user accounts, connection details and machine paths are private.

For each reduced probe, stage a fresh private directory outside the checkout:

```powershell
go run ./scripts/conformance stage --probe output-format --staging C:\Temp\visualg-output-format
```

Follow the generated `instructions.txt`. Open `source.alg`, check the editor
contents, supply the exact `input.txt` sequence, and run it. Save unchanged
Unicode control text as UTF-8 `raw.txt`; distinguish this observation from native
console bytes. Preserve generated files byte for byte at their declared relative
paths. A GUI-only result needs both `screenshot.png` and `transcription.txt`.
Record acceptance/rejection and the actual capture time; an interrupted or stale
run is not evidence.

The September 9 syntax batch supplied the exact CP1252-decoded source through
the application's editor control, verified the editor readback and source hash,
then checked fresh observation timestamps and the executable hash. This method
records language behavior; it does not establish file-open decoding behavior.
The 20 new GUI rejection screenshots were individually inspected for private
information and preserved without image edits. Their PNG metadata contains only
image format, color, resolution, and pixel data. The accompanying transcriptions
are labeled manual; partial output panels are retained separately and do not
stand in for the rejection diagnostic. Connection details, raw window/control
inventories, bootstrap captures, and inconclusive runs remain private.

```powershell
go run ./scripts/conformance capture --staging C:\Temp\visualg-output-format --accepted true --captured-at 2026-09-07T12:00:00Z --normalizer panel-v1
```

The timestamp above is an example: use the observed time. Add `--gui-only` for
GUI-only results. Staging retains hashes of the source, stdin and auxiliary input
files. Capture verifies source, stdin and input-only files before writing
`normalized.txt` and `evidence.json` without overwriting previous captures.
Files also declared as generated outputs may change; capture hashes their final
bytes and attaches GUI evidence when requested. Restage auxiliary-input probes
with the current recorder before capturing them.

Review all text, images and metadata before copying necessary evidence into the
corpus. Keep unredacted captures private. Document any redaction and recompute
hashes over the published bytes. Update the manifest from the reviewed evidence;
recording never marks an implementation verified. Never commit a reference
executable, installer or distribution archive.

## Manifest version 1

The JSON model is defined by `scripts/conformance/manifest.go`. Unknown fields
and trailing JSON are errors. All artifact paths are portable, repository-relative
paths without parent traversal, drive prefixes, backslashes or symlinks. Artifact
hashes are lowercase SHA-256 of the exact published bytes.

- Top level: `version`, `reference`, `recorderVersion`, `normalizerVersion`,
  `tasksPath`, `inventorySources`, `inventory`, `probes`, optional `retired`.
- Inventory items have stable `id`, `kind`, `link`, and `probes`. Kinds cover
  requirements, checklist items, defects, assumptions, discovered features and
  bundled examples. Links use `file#Requirement: Exact heading` where applicable.
  Every requirement in the `specs/` tree adjacent to `tasksPath`, plus every
  requirement in additional declared source files, must have a requirement-kind
  trace. Omitting a source file from `inventorySources` cannot hide its requirements.
- Probes have stable `id`, `ownerGroup`, existing `tasks`, `source`, `input`,
  optional initial `files`, `timeoutMS`, `evidence`, and `implementation`.
- Recorded evidence includes explicit `accepted`, UTC `capturedAt`, hashed `raw`
  and `normalized` artifacts, `normalizer`, and optional generated/GUI artifacts.
- Evidence states are `unrecorded`, `recorded`, and `not-applicable`.
  Implementation states are independently `pending`, `verified`, and
  `not-applicable`. Reference rejection still requires negative coverage.
- Pending implementation links to its owning tasks; it must not invent future
  test paths. Verified implementation links existing `file_test.go#TestFunction`
  symbols and declares exact stdout, exit status, diagnostic codes/locations,
  and generated bytes. Omitted diagnostic columns mean they were not established.
  Accepted reference runs require exit status 0 and no diagnostics; rejections
  require exit status 1 and mapped diagnostics. Generated paths must match the
  reference inventory exactly once each, with matching bytes.
- Reviewed exceptions and downgrades carry a reason and an existing review link.
  Project-specific safeguards may have reference non-applicability while their
  implementation is pending; verification still requires actual project tests.
  Probe, inventory, and task IDs share a unique namespace. Removing a probe or
  inventory entry, or a task linked by an earlier probe, requires a reviewed
  `retired` disposition regardless of implementation state. A replacement may
  name an active probe, inventory entry, or task. Retirement records must remain
  in later manifests, and retired IDs cannot also be active or reused.

The three earlier rejection probes include their original GUI screenshots and
explicitly labeled manual transcriptions. For those probes, `raw`, `normalized`
and `transcription` refer to the transcription bytes under `text-v1`; the original
screenshot is separately hashed. Their earlier partial output panels remain in
`raw.txt`. Capture timestamps and source hashes are from the original recordings.

Evidence mode requires every required observation or reviewed exception.
Incremental mode also rejects pending behavior owned by completed task groups.
Implementation acceptance requires all behavior verified or explicitly reviewed
as outside scope. Evidence collection cannot certify implementation completion.

## Normalization

`panel-v1` removes exactly the application prefix `Início da execução\r\n`
and suffix `\r\nFim da execução.\r\n`, then converts CRLF to LF. Missing notices
are errors. `text-v1` only converts CRLF to LF, for transcribed error/control
text. `bytes-v1` makes no changes, for generated files. None trims whitespace,
changes casing or decimals, rewrites error positions, or removes matching text
inside the program output.
