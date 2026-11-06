## Context

See `proposal.md` for scope and motivation. The manifest currently contains 1,679 probes: 1,673 verified, four pending reviewed rejections, and two reviewed non-applicable implementations. Its accepted programs have no pending implementation. All 216 tasks in the archived conformance change are complete.

The working checkout additionally contains the untracked `tmp_diag_test.go`, which imports the former module path and only prints observations. Committed source passes ordinary tests without that file. The latest published release is `v0.0.1`; later commits include runtime fixes and the bilingual README.

The existing release workflow builds six archives and `SHA256SUMS` when a GitHub Release is published. Conformance acceptance requires a completely clean checkout, including ignored files, and private quality evidence bound to its exact commit. These constraints make the order of documentation, merge, validation, and publication significant.

## Goals / Non-Goals

**Goals:** Use the existing tooling, preserve local investigative work, and make current status and published downloads agree with a qualified commit.

**Non-Goals:** Change interpreter behavior, weaken acceptance gates, add automation for a one-time cleanup, or revise historical evidence to look current.

## Decisions

### Preserve the scratch experiment outside test discovery

During implementation, move the untracked file to private storage outside the checkout and verify the preserved bytes before removing the checkout copy. Do not commit its contents, its private destination, or a broad ignore rule. Updating its imports would retain a print-only experiment in the test suite; deleting it outright would lose the investigation. No permanent regression test is needed because no product logic changes.

### Derive current documentation from existing evidence

Recalculate counts and dispositions from the manifest and bundled catalog when applying the change. Reconcile current-status passages in `docs/language.md`, `docs/reference-corpus-progress.md`, `docs/bundled-examples-progress.md`, and `docs/development.md`, including obsolete active-change references and claims that acceptance still fails.

State the distinction between required accepted programs, the four reviewed rejection differences, and the two nonterminating exclusions. Retain limitations outside the recorded corpus. Preserve anchors used by traceability links. Dated reports keep their historical results; link to them as historical evidence instead of presenting them as fresh qualification. No status generator or manifest-state changes are needed.

### Publish a patch containing only changes since the prior release

Use `v0.0.2` as the planned patch version. Before the first publication attempt, check both remote tag and release availability and stop version selection if either is unexpectedly occupied. Initialize a private release handoff record with the version, qualified candidate SHA, and reviewed release notes before creating the tag. This record distinguishes a resumable publication from an unrelated version collision.

Audit `v0.0.1..candidate` to select the changelog and release-note entries. The release candidate's changelog must have an empty `Unreleased` section, a dated `0.0.2` section containing only newly shipped changes, and a dated `0.0.1` section containing the already shipped entries. Verify the prior release date from its publication record and preserve the older entries' text and order. No already shipped entry may remain under `Unreleased` or be attributed to `0.0.2`. Exclude machine-local scratch handling from public release notes.

Use `make dist VERSION=v0.0.2` for rehearsal, with distribution output outside the checkout. Keep the existing packaging script and workflow unchanged unless verification discovers a concrete blocker. Six archives plus their checksum manifest are the required published assets; verify checksums for every downloaded archive and verify `--version` and the hello program on the available native target and in the existing Linux workflow smoke check.

### Qualify and tag the same merged commit

Prepare the focused documentation/release PR with the required user attribution and authenticated account, then complete the required change review and CI checks. After merge, select a clean candidate containing the runtime fixes and documentation updates. Assemble fresh results for that exact commit using `docs/development.md`, including all platform results, both 30-second fuzz campaigns, and the required quality-report check IDs.

Run evidence validation, implementation acceptance, and release acceptance with a preserved explicit history base and private quality report. Use a separate clean candidate checkout, keeping task-progress edits, reports, build caches, and distribution output elsewhere. A prior report or a PR merge-test commit cannot qualify a different merged commit.

Create and push an annotated tag only after those checks pass, then publish the GitHub Release for that existing tag. A tag alone does not attach downloads. Confirm the publication workflow ran against the tagged commit and finished successfully before verifying assets.

### Resume publication from verified remote state

Keep tag creation and release publication as separate tasks. Record each verified result privately as it occurs, and on restart inspect remote state even when a previous command reported failure or timed out. An existing version is resumable only when the prior handoff record identifies this publication, its remote annotated tag resolves to the recorded qualified commit, and any release names that tag and matches the reviewed notes. If those checks fail, stop without moving the tag or overwriting the release.

| Verified state on restart | Next action |
| --- | --- |
| Neither remote tag nor release exists | Create or verify the local annotated tag, then push it to the recorded commit. |
| The matching remote tag exists and the release is absent | Reuse the tag and create the release; do not recreate or move the tag. |
| The matching release is a draft | Publish that draft after verifying its tag and notes. |
| The matching release is published | Resume workflow and asset verification; do not create a duplicate release. Rerun failed publication jobs for the same tagged source when appropriate. |

Review the tag-pushed/release-create-failed, release-published/client-timeout, attachment-failed, and unexpected-version-collision cases against this table before handoff. Check remote state before retrying any mutation. Existing publication is not proof of successful asset verification.

### Keep the maintained release guide consistent

Update `docs/releases.md` with the patch procedure, the recovery table, and the requirement to qualify the exact merged commit before tagging. Identify the archived `v0.1.0` checklist as the historical conformance-change procedure; its archive-before-tag sequence does not require re-archiving that completed change or creating `v0.1.0` for this patch. Retain the historical checklist unchanged. Keep the maintenance instructions in the guide itself so later archival of this proposal cannot leave a broken process link.

Publication operations are grouped separately from implementation and review tasks. Keep this maintenance change active until its publication tasks are actually complete; its archival can follow publication. The conformance manifest continues to reference the already completed historical task list. Do not repoint it or precheck publication tasks to make a gate pass. Bookkeeping added after publication must not move the release tag.

## Risks / Trade-offs

- Scratch-file loss → preserve and verify its bytes privately before removing the checkout copy.
- Stale counts or overstated compatibility → recompute manifest/catalog totals and retain reviewed mismatches, exclusions, and unverified scope explicitly.
- Documentation edits breaking evidence links → preserve referenced anchors and run strict specification and traceability validation.
- Validation against the wrong commit → record the candidate SHA, require matching CI/report results, and requalify if the candidate changes.
- Release page appearing before assets finish → report success only after the workflow and all asset checks pass.
- Publication failure → preserve the handoff record and failing evidence, inspect remote state, and resume from the verified recovery-table state. Never move a published tag to changed source; source changes require a separately qualified version.
- Disclosure in notes or evidence → inspect all outgoing commits and content, retain raw evidence privately, and publish only sanitized results with hashes of the published bytes.

## Migration Plan

1. Preserve the scratch file privately, refresh current documentation including the maintenance release guide, and separate the changelog's released history from the focused patch notes.
2. Review, verify, and merge the change; qualify a clean merged candidate using the existing checks and private evidence format.
3. Initialize the private publication record, publish the annotated patch tag and GitHub Release as separate steps, and verify the workflow, downloaded assets, checksums, and smoke results. Resume interrupted work from verified remote state.
4. Record the tag, target commit, workflow and release links, and sanitized verification outcome. Archive this maintenance change after all work is complete.

Before publication, documentation changes can be reverted and the scratch file restored from its private copy. After publication, preserve tag identity; repair packaging for the same source or issue a new qualified patch as appropriate.
