# Development checks

The module remains compatible with Go 1.22. CI runs ordinary tests on Go
1.22.12 and 1.25.0 on Linux, macOS, and Windows. Its quality and fuzz jobs use
Go 1.25.0. Install the same pinned tools with that toolchain:

```sh
go install honnef.co/go/tools/cmd/staticcheck@v0.6.1
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.4.0
npm install --global @fission-ai/openspec@1.8.0
```

Tool dependencies stay outside the language module. The lint configuration
explicitly enables errcheck, govet, ineffassign, staticcheck, and unused, with
uncapped diagnostics. This pins the behavior of the
[Go 1.25-compatible golangci-lint release](https://github.com/golangci/golangci-lint/releases/tag/v2.4.0)
and [Staticcheck 2025.1.1](https://staticcheck.dev/changes/2025.1/).

Run the required checks from the repository root:

```sh
go build ./...
go run ./scripts/checkfmt
go vet ./...
staticcheck ./...
golangci-lint run
go test -count=1 ./...
go test -race -count=1 ./...
go mod tidy
git diff --exit-code -- go.mod go.sum
openspec validate --all --strict --no-interactive
bash scripts/check-oracle.sh --base origin/main
go test ./internal/lexer -run=TestFuzz -fuzz=FuzzLexer -fuzztime=30s -timeout=2m -parallel=2
go test ./internal/parser -run=TestFuzz -fuzz=FuzzParser -fuzztime=30s -timeout=2m -parallel=2
```

The formatter check reports files without changing them. Ordinary tests run
all committed fuzz seeds. Fuzz campaigns use a 64 KiB source profile; larger
generated inputs are outside that profile. Adversarial lexer/parser cases run
in child test processes with five-second watchdogs. Explicit 30-second budgets
cover full-size encoded-source decoding and selected formatter boundary checks,
allowing race instrumentation under load. Oversized-source rejection and depth
cases retain five seconds. Crashes and watchdog expiration fail the test, even
if a child printed `PASS` before exiting. CI uploads newly discovered failing
fuzz inputs; reduce and commit them as regression seeds after fixing the cause.

Runtime fixtures compare exact bytes with `internal/golden`. Missing `.out`
files fail. To intentionally regenerate those fixtures, run
`go test . -run '^TestRunFixtures$' -update`, then review the diff. Comparison
errors report a repository-relative path, the first differing byte, lengths,
and quoted context without normalizing whitespace or encoding.

`.gitattributes` disables Git text conversion for every `testdata` directory.
This preserves committed fixture bytes on checkout and staging, including LF,
CRLF, and Windows-1252 data, even when `core.autocrlf=true` on Windows.

The oracle gate runs evidence validation, which allows recorded probes with
pending implementations and currently fails for missing recordings and mappings.
Supply the PR base or parent branch with `--base` (replace `origin/main` for
stacked work); history comparison is required to detect unreviewed coverage
downgrades. CI supplies the PR base or pre-push commit and fetches its history.
Groups 3–16 will additionally require incremental
replay of completed groups; final acceptance rejects every pending behavior.
Missing reference evidence is never proof of compatibility.

Replayed diagnostics must carry a positive source line and column. An omitted
expected column permits any positive column; it never permits an unpositioned
error to count as a matching rejection.

Final acceptance also requires local quality results for the exact candidate
commit. Keep the report and its evidence outside the checkout, then run:

```sh
go run ./scripts/conformance validate --mode implementation-acceptance --base origin/main --quality /path/to/private-results/quality.json
go run ./scripts/conformance validate --mode release --base origin/main --quality /path/to/private-results/quality.json
```

Both modes require a clean checkout, including untracked and ignored files, and
build and replay that checkout themselves; `--candidate` cannot override it.
Results from another commit fail, so regenerate the evidence after a merge,
rebase or archive commit. Evidence and incremental modes retain their existing
requirements and do not require this report.

The version-1 report has this shape (the abbreviated result list is not valid
for acceptance):

```json
{
  "version": 1,
  "commit": "<full candidate commit ID>",
  "results": [
    {
      "check": "build",
      "status": "pass",
      "evidence": {"path": "build.txt", "sha256": "<SHA-256 of evidence bytes>"}
    }
  ]
}
```

Required check IDs are `build`, `gofmt`, `vet`, `staticcheck`, `golangci-lint`,
`tests`, `race`, `windows`, `macos`, `linux`, `fuzz-lexer`, `fuzz-parser`, and
`openspec`. Each must occur exactly once with status `pass` and a nonempty,
hash-matching evidence file. Evidence paths are relative to the report directory;
parent traversal and symlinks are rejected. Unknown fields, unknown checks,
duplicate checks, missing results and failed or skipped results fail acceptance.

These are trusted local result records, not independent attestations. Record
the actual commands, toolchain/platform, exit status and result in each evidence
file, including successful commands with empty output. Platform results must
cover both supported Go versions, and each fuzz campaign must run for 30 seconds.
Preserve raw operational evidence privately; publish only reviewed, sanitized
evidence, with hashes recalculated for its published bytes.

Implementation acceptance reports remaining task IDs on stderr without making
its own future reporting or handoff a prerequisite. Release mode additionally
requires every task checkbox complete; malformed or duplicate task entries fail.
Archive creation and tag creation remain later operations, never prerequisites
for either command. Passing synthetic regression tests does not qualify the
current implementation or complete the release handoff.

See the [baseline quality report](quality-baseline.md) for measured coverage,
checks actually executed, and outstanding conformance work.
