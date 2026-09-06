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
bash scripts/check-oracle.sh
go test ./internal/lexer -run=TestFuzz -fuzz=FuzzLexer -fuzztime=30s -timeout=2m -parallel=2
go test ./internal/parser -run=TestFuzz -fuzz=FuzzParser -fuzztime=30s -timeout=2m -parallel=2
```

The formatter check reports files without changing them. Ordinary tests run
all committed fuzz seeds. Fuzz campaigns use a 64 KiB source profile; larger
generated inputs are outside that profile. Adversarial lexer/parser cases run
in child test processes with five-second watchdogs. Crashes and watchdog
expiration fail the test. CI uploads newly discovered failing fuzz inputs;
reduce and commit them as regression seeds after fixing the cause.

Runtime fixtures compare exact bytes with `internal/golden`. Missing `.out`
files fail. To intentionally regenerate those fixtures, run
`go test . -run '^TestRunFixtures$' -update`, then review the diff. Comparison
errors report a repository-relative path, the first differing byte, lengths,
and quoted context without normalizing whitespace or encoding.

The oracle gate is a visible placeholder until group 2 installs the corpus and
validator. It then runs evidence validation, which allows recorded probes with
pending implementations. Groups 3–16 will additionally require incremental
replay of completed groups; final acceptance rejects every pending behavior.
Missing reference evidence is never proof of compatibility.

See the [baseline quality report](quality-baseline.md) for measured coverage,
checks actually executed, and outstanding conformance work.
