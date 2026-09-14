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
| Incremental corpus replay | PASS; 1,653 verified, 23 pending, 2 reviewed not-applicable, zero verified regressions |
| Strict OpenSpec validation | PASS |

The supported-platform CI matrix remains green in the quality report. Release
acceptance is still blocked by the pending implementation probes and the
remaining final handoff tasks.
