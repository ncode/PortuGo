# Quality baseline: conformance group 1

Recorded 2026-09-06 on macOS arm64 with Go 1.25.0, starting from merged plan
commit `90365a2`. This report describes the group 1 changes, before the oracle
corpus or language changes exist.

| Check | Observed result |
| --- | --- |
| Policy regression before implementation | Failed: CI workflow, lint configuration, and tool documentation were absent |
| Golden and watchdog regressions before implementation | Failed for byte differences, missing/updated fixtures, process isolation, crashes, and timeouts |
| Focused helpers and policy | Pass |
| Build and non-mutating gofmt verification | Pass |
| `go vet ./...` | Pass |
| Staticcheck 2025.1.1 (`v0.6.1`) | Pass |
| Pinned golangci-lint `v2.4.0` | Pass, zero issues (also passed installed `v2.13.2`) |
| Ordinary tests, uncached | Pass |
| Race tests, uncached | Pass |
| Lexer fuzzing, 30 seconds, two workers | Pass, 2,304,243 executions |
| Parser fuzzing, 30 seconds, two workers | Pass, 1,618,667 executions |
| Strict OpenSpec validation | Pass: one change, zero failures |
| Module hygiene and whitespace check | Pass; no language-module dependencies added |
| Windows/Linux execution and Go 1.22 minimum | Configured in CI; not executed on this macOS host |
| Oracle validation | Explicitly pending group 2; no conformance result claimed |

Commands are in [development checks](development.md). Some local commands
needed access to the normal Go build cache outside the workspace sandbox.
Fuzz execution counts include the existing local Go fuzz cache and are not
promised to reproduce exactly on a fresh runner. No crash or hang was found
in these campaigns. The tests deliberately verify that a watchdog expiration
or child panic fails rather than becoming a successful diagnostic.

## Coverage inventory

There are six runtime `.alg` fixtures with byte-exact `.out` files and one
optional `.in` file. They exercise arithmetic, console I/O, eager logical
evaluation, loop-variable mutation, subprograms, and vectors. Two example
programs exist but are not yet swept automatically. These are implementation
regressions; none currently has official reference provenance.

Package-local statement coverage from `go test -count=1 -cover ./...` is:

| Package | Coverage |
| --- | ---: |
| lexer | 73.3% |
| parser | 48.4% |
| sema | 31.9% |
| source | 18.6% |
| golden | 91.3% |
| testprocess | 81.8% |
| scripts/checkfmt | 74.1% |

CLI, AST, diagnostics, interpreter, REPL, runtime, stdlib, and token have no
package-local tests and report 0% in that command. The root runtime fixtures
do exercise several of them, but ordinary per-package coverage does not
attribute those calls to the dependency packages. This report does not equate
those zeros with unreachable code or claim broad conformance coverage.

Outstanding work includes authoritative evidence, catalog completeness,
parser/printer goldens and round trips, positioned runtime diagnostics, CLI
and REPL transcripts, production resource limits, file/host behavior, and
bundled-example acceptance. The remaining OpenSpec groups own those gaps.
