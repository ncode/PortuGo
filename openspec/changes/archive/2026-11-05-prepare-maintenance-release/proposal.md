## Why

An untracked diagnostic experiment imports the former module path and prevents the working checkout from running its tests. Current documentation still reports completed conformance work as pending, and the merged runtime fixes have not reached the published downloads.

## What Changes

- Preserve `tmp_diag_test.go` privately outside the checkout and remove it from Go test discovery, without committing the experiment or adding a replacement test that only prints diagnostics.
- Reconcile current status claims in `docs/language.md`, `docs/development.md`, and the corpus/example progress documents with the manifest and archived task list. Report zero pending accepted programs and retain the four reviewed rejected-program mismatches and two nonterminating exclusions.
- Separate the changelog into an empty `Unreleased` section, a dated `0.0.2` section for changes since `v0.0.1`, and a dated `0.0.1` section for already shipped entries. Prepare focused patch release notes.
- Update `docs/releases.md` with the maintenance-release procedure and recovery steps, distinguishing it from the historical conformance checklist and its archive-before-tag sequence.
- Qualify a clean, merged candidate and publish the next patch release through the existing release workflow. The planned version is `v0.0.2`; reject an unexpected existing version and allow an interrupted publication to resume after verifying its recorded intent and target commit.
- Verify all six platform archives, `SHA256SUMS`, the embedded version, and the release workflow result before reporting publication complete.

## Capabilities

### New Capabilities

None. This maintenance change uses the existing test, validation, and distribution tooling.

### Modified Capabilities

None. Language, CLI, and release behavior remain as specified. `.openspec.yaml` declares `skip_specs: true` because the work corrects documentation and performs an existing release process.

## Impact

The implementation affects local scratch-file placement, current documentation including `docs/releases.md`, `CHANGELOG.md`, and release metadata/assets. It reuses `scripts/dist.sh`, `make dist`, and `.github/workflows/release.yml`; no dependencies or runtime changes are planned. Historical quality reports retain their original dates and results, and release evidence must identify the exact candidate commit.

## Non-goals

- Resolving the four rejected-program compatibility differences or expanding conformance coverage.
- Rewriting reference evidence, historical task records, or existing release tags.
- Adding packaging infrastructure, package-manager publication, language features, or a stronger compatibility claim.
