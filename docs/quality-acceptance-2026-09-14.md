# Candidate quality evidence — 2026-09-14

This report records the green GitHub Actions checks for PR #183, whose exact
head was the pushed candidate for this stack. The report is a sanitized
summary; raw workflow logs remain outside the repository.

| Check | Result | Coverage |
| --- | --- | --- |
| Build | PASS | `go build ./...` in the quality job and every platform matrix job |
| Formatter verification | PASS | `go run ./scripts/checkfmt` |
| Vet | PASS | `go vet ./...` |
| Staticcheck | PASS | Pinned Staticcheck release in the quality job |
| GolangCI-Lint | PASS | Pinned GolangCI-Lint release in the quality job |
| Ordinary tests | PASS | Linux, macOS and Windows with Go 1.22.12 and 1.25.0 |
| Race tests | PASS | `go test -race -count=1 ./...` |
| Linux | PASS | Go 1.22.12 and 1.25.0 |
| macOS | PASS | Go 1.22.12 and 1.25.0 |
| Windows | PASS | Go 1.22.12 and 1.25.0 |
| Lexer fuzz | PASS | 30-second campaign with two workers |
| Parser fuzz | PASS | 30-second campaign with two workers |
| Specification validation | PASS | Strict OpenSpec validation and the oracle inventory check |

These checks establish quality evidence for the exact candidate. They do not
establish implementation acceptance: the manifest still contains 22 pending
implementation probes, and the final acceptance and release gates remain open.
