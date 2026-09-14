# Clean-checkout release rehearsal — 2026-09-14

The rehearsal ran from a fresh checkout with no working-tree changes. It
validated the build, command entry points, fixtures, manifest replay, and
specification checks before the final handoff.

| Area | Result |
| --- | --- |
| `go build ./...` | PASS |
| `go test ./... -count=1` | PASS |
| `go test -race -count=1 ./...` | PASS |
| `go vet ./...` and `go run ./scripts/checkfmt` | PASS |
| `run`, `check`, `fmt --check`, `fmt -w`, and `repl` smoke commands | PASS |
| Incremental corpus replay | PASS; 1,673 verified, 4 reviewed rejected-source records pending, 2 reviewed not-applicable, zero verified regressions |
| Strict OpenSpec validation | PASS |

The supported-platform CI matrix remains green in the quality report, and
implementation acceptance passes on the same candidate. Release acceptance
is the final candidate gate before the post-merge archive and release-tag
operations; the four reviewed rejected-source records remain visible without
representing pending accepted behavior.
