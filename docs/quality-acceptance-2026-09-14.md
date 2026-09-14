# Current quality evidence — 2026-09-14

This report records the green GitHub Actions checks for the preceding stacked
documentation candidate (PR #190), together with fresh local checks on the
same source. It is a sanitized summary; raw workflow logs remain outside the
repository.

| Check | Result | Coverage |
| --- | --- | --- |
| Build | PASS | `go build ./...` |
| Formatter verification | PASS | `go run ./scripts/checkfmt` |
| Vet | PASS | `go vet ./...` |
| Staticcheck | PASS | Pinned Staticcheck release in GitHub Actions |
| GolangCI-Lint | PASS | Pinned GolangCI-Lint release in GitHub Actions |
| Ordinary tests | PASS | `go test ./... -count=1` |
| Race tests | PASS | `go test -race -count=1 ./...` |
| Linux | PASS | Go 1.22.12 and 1.25.0 |
| macOS | PASS | Go 1.22.12 and 1.25.0 |
| Windows | PASS | Go 1.22.12 and 1.25.0 |
| Lexer fuzz | PASS | 30-second campaign with two workers |
| Parser fuzz | PASS | 30-second campaign with two workers |
| Recorder, normalizer, replay, and traceability tests | PASS | `go test ./scripts/conformance -count=1` |
| Evidence validation and replay | PASS | 1,679 probes; zero verified regressions |
| Specification validation | PASS | Strict OpenSpec validation |

These checks establish quality and evidence status for the current stack. They
do not establish implementation acceptance: 10 implementation probes remain
pending in the manifest.

## Pending implementation inventory

| Owner | Count | Probe IDs |
| --- | ---: | --- |
| Bundled example qualification (17.3) | 10 | See the manifest for the current pending probe inventory. |

The bundled entries retain their recorded sources and outputs while their
deterministic mismatches remain pending. The two recorded nonassignable `var`
argument guards are verified by project regressions; their reviewed reference
application faults remain separate from language-level diagnostic claims.

Ten bundled transcripts now have verified generated-output contracts. One
checks ten spaced integer lines, nine bounded random values and a fixed zero
tail; the other checks ten alternating bounded integer and uppercase text
lines; the third checks a bounded generated input sequence and its integer
permutation; the fourth and fifth check bounded real input sequences, sorted
permutations, and formatted rows; the sixth and seventh check generated
code/name records with both sort orders; the eighth checks ten numbered integer
rows and a fixed not-found search; the ninth and tenth check a repeated nonzero
integer with fixed prompts and no final LF. All replay the original and formatted
source without comparing random sequences. The remaining generated-output cases
and rejected paths stay pending.
