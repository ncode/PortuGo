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

1. Merge the release changes and select a clean, tested commit. Complete the
   repository's quality checks and the applicable
   [conformance release checklist](../openspec/changes/archive/2026-09-15-complete-visualg-3-0-7-conformance/release.md).
2. Create an annotated version tag on that commit and push it. Do not replace an
   existing release tag.
3. Create and publish a GitHub Release for that existing tag. For example, after
   pushing `v0.1.0`, run `gh release create v0.1.0 --verify-tag --generate-notes`,
   or use GitHub's Releases page. Prereleases use the same build process.
4. Wait for **Release artifacts** to finish and confirm the six archives and
   `SHA256SUMS` appear in the release's Assets. The release page exists while
   its downloads are still being built.

Only the attachment job has repository write permission, using the built-in
`GITHUB_TOKEN`. A rerun replaces assets with the same filenames, matching the
facts workflow. No separate publishing secret is required for GitHub downloads.
Homebrew publication would additionally need a PortuGo formula and tap token.
