# Release builds

PortuGo follows the [facts release workflow](https://github.com/ncode/facts/blob/main/.github/workflows/release.yaml):
publishing a GitHub Release builds and attaches executable archives and
`SHA256SUMS`. Pushing a tag alone does not publish downloads.

## Build locally

Run from the repository root on Linux or macOS with Go 1.27, Make, Bash,
`tar`, `zip`, and either `sha256sum` or `shasum` installed:

```sh
make dist
```

The version is the exact tag at HEAD, or `dev-<short commit>` for an untagged
checkout. You can set the version explicitly for a packaging rehearsal:

```sh
make dist VERSION=dev-local
```

Output goes to `dist/`. Each archive contains a versioned folder and its
`portugo` executable (`portugo.exe` on Windows). Binaries embed the version
printed by `portugo --version`; a plain `go build` reports `dev`.

| Targets | Format |
| --- | --- |
| Linux amd64 and arm64 | `portugo-VERSION-linux-ARCH.tar.gz` |
| macOS amd64 and arm64 | `portugo-VERSION-darwin-ARCH.tar.gz` |
| Windows amd64 and arm64 | `portugo-VERSION-windows-ARCH.zip` |

Builds disable cgo and trim local source paths. `SHA256SUMS` lists only the
archives produced by that invocation. Staging happens in a temporary directory;
a compiler failure leaves any preceding distribution intact.

`DIST_DIR` changes the output directory. `DIST_TARGETS` selects a space-separated
subset of the six supported targets, for example:

```sh
make dist VERSION=dev-local DIST_TARGETS='linux/amd64 darwin/arm64'
```

## Verify before publishing

The **Release artifacts** workflow runs on pull requests and can also be run
manually from GitHub Actions. It runs the tests, builds all six archives, checks
their SHA-256 hashes, and runs the extracted Linux executable to verify the
embedded version and the hello example. Download the `portugo-dist` Actions
artifact to inspect the result. These runs do not attach assets to a release.

For a local distribution, verify all generated archives from its directory:

```sh
cd dist
sha256sum --check SHA256SUMS
```

On macOS, use `shasum -a 256 --check SHA256SUMS` instead.

## Publish a release

Use this procedure for a patch release from the already archived conformance
baseline. The [historical conformance checklist](../openspec/changes/archive/2026-09-15-complete-visualg-3-0-7-conformance/release.md)
records the original change's planned `v0.1.0` archive-before-tag sequence. It
does not require re-archiving that change or choosing `v0.1.0` for a patch.
The maintenance change stays active through publication verification and may
be archived afterward. Keep the historical checklist and manifest task links
unchanged.

1. Complete the maintenance implementation and review tasks, merge the changes,
   and select a separate clean checkout of the exact candidate commit. Follow
   the [development checks](development.md), including fresh platform and quality
   evidence plus evidence, implementation-acceptance, and release validation.
   Record an explicit history base; the previous release commit is suitable for
   checking the complete patch. Keep evidence, task-progress edits, and build
   output outside this checkout.
2. Before the first publication attempt, verify that the intended tag and release
   are unused. Initialize a private handoff record with the version, qualified
   candidate commit, and reviewed release notes. An unrelated existing version
   is a collision; choose no replacement tag automatically and never overwrite it.
3. Create and push an annotated version tag on the qualified commit. Verify its
   remote target and record the result before creating the release.
4. Create and publish the GitHub Release for that existing tag using the reviewed
   notes. For example, after pushing `v0.0.2`, run
   `gh release create v0.0.2 --verify-tag --notes-file /path/to/release-notes.md`,
   or use GitHub's Releases page. Record the release URL and confirmed state.
5. Wait for **Release artifacts** to succeed against the tagged commit. Confirm
   all six archives and `SHA256SUMS` appear in Assets, download them, and verify
   every checksum. Run the downloaded native executable to check its version
   and hello output; verify the workflow's Linux smoke result as well.
6. Record the tag, target commit, workflow/release links, asset names, and
   verification outcome. Keep raw evidence private and publish only a sanitized
   summary. Complete the publication checklist before archiving the maintenance
   change. Later bookkeeping must not move the tag.

The release page exists while its downloads are still being built. Publication
is complete only after workflow and asset verification, not when the page first
appears.

## Resume an interrupted publication

Inspect remote state before retrying any mutation, including after a command
reports a timeout. Reuse an existing version only when the prior handoff record
identifies this publication, the remote annotated tag resolves to its qualified
commit, and any release names that tag and matches the reviewed notes. Stop on
missing provenance or a mismatch; do not move tags or overwrite an unrelated
release.

| Verified state | Next action |
| --- | --- |
| Neither remote tag nor release exists | Create or verify the local annotated tag, then push it to the recorded commit. |
| Matching remote tag exists; release is absent | Reuse the tag and create the release. |
| Matching release is a draft | Verify its tag and notes, then publish the draft. |
| Matching release is published | Resume workflow and asset verification without creating a duplicate release. Rerun failed publication jobs for the same tagged source when appropriate. |

A successful tag push followed by failed release creation resumes at release
creation. If release publication succeeds but the client times out, discover
and reuse the published release. An attachment failure resumes the existing
workflow. A version collision without the matching handoff record stops the
operation. None of these cases bypasses final asset verification.

Only the attachment job has repository write permission, using the built-in
`GITHUB_TOKEN`. A rerun replaces assets with the same filenames, matching the
existing workflow. No separate publishing secret is required for GitHub downloads.
Source changes require a newly qualified version; do not move a published tag.
Homebrew publication would additionally need a PortuGo formula and tap token.
