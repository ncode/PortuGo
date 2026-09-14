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
| Evidence validation and replay | PASS | 1,678 probes; zero verified regressions |
| Specification validation | PASS | Strict OpenSpec validation |

These checks establish quality and evidence status for the current stack. They
do not establish implementation acceptance: 24 implementation probes remain
pending in the manifest.

## Pending implementation inventory

| Owner | Count | Probe IDs |
| --- | ---: | --- |
| Bundled example qualification (17.3) | 22 | `bundled-215411a6ad9a`, `bundled-24919c4a2e4c`, `bundled-3b4c69f65b60`, `bundled-43742446e585`, `bundled-460092b86448`, `bundled-467736cf246b`, `bundled-49eb7f47065e`, `bundled-54f1e50d1c02`, `bundled-7f44aa03fcc5`, `bundled-86701928f296`, `bundled-8b2d8c702924`, `bundled-991ec2bd1566`, `bundled-a5945bc9ee4d`, `bundled-acf5d46a3a6a`, `bundled-b1c9cc780302`, `bundled-c2c3c37e29ae`, `bundled-c47bf27802a4`, `bundled-c6443beee11c`, `bundled-da71c1ad0344`, `bundled-f251f2cd0c22`, `bundled-f4fa0c514cf7`, `bundled-fa172f3f5f70` |
| Nonassignable `var` arguments (2.2, 5.4) | 2 | `var-expression-argument`, `var-literal-argument` |

The bundled entries retain their recorded sources and outputs while their
deterministic mismatches remain pending. The two `var` entries retain reviewed
reference application faults and are not promoted to language-level
diagnostics.
