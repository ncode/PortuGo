# VisuAlg 3.0.7 reference corpus

This directory is being assembled by OpenSpec group 2. It does not yet contain
the complete evidence inventory and does not establish full conformance.
The earlier selected observations remain in `docs/validation/visualg-2026-09-07`.

## Commands

Run from the repository root with Go 1.22 or newer:

```sh
go test ./scripts/conformance
go run ./scripts/conformance validate --mode evidence
go run ./scripts/conformance validate --mode incremental --base origin/main
go run ./scripts/conformance validate --mode implementation-acceptance
```

`validate` verifies metadata and hashes, builds the CLI from the current checkout,
and replays recorded probes in isolated temporary working directories. Use
`--candidate` to supply an already built executable. Each probe has a 1–30,000 ms
deadline; source is limited to 64 KiB, each captured output stream to 1 MiB, and
each evidence file to 16 MiB. A deadline or output limit is a failed observation.
Pending results, including mismatches, remain visible in the JSON report.
Verified mismatches fail every phase. Acceptance rejects all pending behavior.
State and host-event observations need the group 3 adapter; the CLI runner
reports them as unsupported rather than claiming they match.

`--base` checks an earlier Git manifest for removed verified probes or downgrades.
`--previous` accepts a repository-relative manifest file instead. A missing
manifest in the initial base is allowed; invalid refs and other Git failures are
errors. Do not omit this comparison when reviewing changes to an existing corpus.

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

```powershell
go run ./scripts/conformance capture --staging C:\Temp\visualg-output-format --accepted true --captured-at 2026-09-07T12:00:00Z --normalizer panel-v1
```

The timestamp above is an example: use the observed time. Add `--gui-only` for
GUI-only results. The command verifies staged source/input hashes and writes
`normalized.txt` and `evidence.json` without overwriting previous captures.
It hashes every declared generated file and attaches GUI evidence when requested.

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
  Every requirement in the declared source files must be traced.
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
- Reviewed exceptions and downgrades carry a reason and an existing review link.
  Project-specific safeguards may have reference non-applicability while their
  implementation is pending; verification still requires actual project tests.
  Retired IDs retain a reviewed disposition and optional replacement probe.

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
