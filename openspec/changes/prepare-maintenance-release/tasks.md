## 1. Restore the working checkout

- [x] 1.1 Preserve `tmp_diag_test.go` in private storage outside the checkout, verify identical bytes, and remove its checkout copy without staging the experiment or its private destination.
- [x] 1.2 Run `go test -count=1 ./...` from the working checkout and confirm that old-module imports no longer prevent test discovery.

## 2. Reconcile documentation and patch notes

- [x] 2.1 Recompute manifest/catalog counts and dispositions, verify the archived task completion, and identify current-status claims contradicted by that evidence.
- [x] 2.2 Refresh `docs/language.md`, `docs/reference-corpus-progress.md`, `docs/bundled-examples-progress.md`, and `docs/development.md` as needed: zero pending accepted implementations, four reviewed rejection mismatches, two nonterminating exclusions, correct archive references, and retained limits outside the recorded corpus. Preserve traceability anchors and dated historical results.
- [x] 2.3 Audit changes since `v0.0.1`; create an empty `Unreleased` section, a dated `0.0.2` section for newly shipped changes including the runtime fixes, and a dated `0.0.1` section for already shipped entries. Verify the prior publication date, preserve older entry text/order, and prepare focused patch release notes.
- [x] 2.4 Update `docs/releases.md` with the maintenance patch procedure and publication recovery table. Distinguish the historical conformance checklist from this patch's sequence, retain that checklist unchanged, and state that this maintenance change can be archived after publication.

## 3. Verify and merge the maintenance change

- [ ] 3.1 Run the required build, formatting, vet, pinned lint, ordinary-test, race, module-hygiene, and lexer/parser fuzz checks from `docs/development.md`; use platform CI results for the selected Go version and fix confirmed failures within scope.
- [x] 3.2 Run strict OpenSpec validation and the existing oracle/traceability checks with an explicit history base; confirm the documentation edits preserve evidence links and accepted-program coverage. Check the changelog sections and walk the design's four publication recovery cases for consistency with the maintained release guide.
- [x] 3.3 Rehearse `make dist VERSION=v0.0.2` with output outside the checkout; verify all six archives and `SHA256SUMS`, then smoke-test the native executable's version and hello output. Reuse the existing packaging implementation.
- [x] 3.4 Run the required change review on the privacy-screened final diff, address confirmed findings, and report files excluded from review.
- [ ] 3.5 Inspect every outgoing commit and publication text for private content, verify the required commit identities and authenticated PR account, and create the focused PR. Merge after required review and CI checks pass.

## 4. Qualify the merged release candidate

- [ ] 4.1 Select a separate clean checkout of the merged candidate containing the runtime fixes and documentation updates. Record its SHA and explicit conformance history base privately; keep progress edits and generated output outside that checkout.
- [ ] 4.2 Assemble fresh, hash-verified quality evidence for that exact candidate: `build`, `gofmt`, `vet`, `staticcheck`, `golangci-lint`, `tests`, `race`, `windows`, `macos`, `linux`, `fuzz-lexer`, `fuzz-parser`, and `openspec`. Verify CI head SHAs and retain raw logs privately.
- [ ] 4.3 Run evidence validation, implementation acceptance, and release acceptance on the clean candidate with the private quality report and explicit base. Require all applicable gates to pass and retain the four reviewed rejections without treating them as pending accepted behavior.

## 5. Publish and verify the patch release

- [ ] 5.1 Before the first attempt, confirm that remote tag and release `v0.0.2` are unused and initialize a private handoff record with the version, qualified candidate SHA, and reviewed notes. On restart, inspect remote state and use the recovery table; stop for an unrelated version collision or mismatched tag, target commit, or release notes.
- [ ] 5.2 Create and push an annotated `v0.0.2` tag on the qualified commit, or verify and reuse this publication's matching remote annotated tag. Confirm its target SHA and record the result before proceeding; never move or overwrite the tag.
- [ ] 5.3 Publish the GitHub Release for that tag with the reviewed patch notes. On restart, create a missing release, publish a matching draft, or reuse the matching published release after verification; record its URL and confirmed state before continuing.
- [ ] 5.4 Wait for the existing Release artifacts workflow to succeed against the tagged commit, resuming verification or rerunning failed publication jobs for that source as appropriate, then verify the six expected platform archives and `SHA256SUMS` are attached.
- [ ] 5.5 Download all published archives and verify their checksums; confirm native `--version` reports `v0.0.2`, the hello program has exact expected output, and the workflow's Linux smoke check passed.
- [ ] 5.6 Complete the handoff record with the immutable tag, target SHA, release/workflow links, asset names, and sanitized verification outcome. Mark publication tasks complete only after verification; keep later bookkeeping separate from the tagged commit and leave this change ready for archival.
