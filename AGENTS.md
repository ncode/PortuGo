# AGENTS.md — portugol-go

A Go implementation of **Portugol** (VisuAlg dialect): lexer, parser, semantic analyzer, and tree-walking interpreter, with a CLI and REPL.

This file is the contract between the codebase and any coding agent (Claude Code, Cursor, etc.) operating on it. Read it fully before making changes.

---

## 1. Dialect & scope

**Target dialect:** VisuAlg 3.x (Apoio Informática). Most other "Portugol" dialects (Portugol Studio / UNIVALI, Portugol Webstudio, "academic" textbook Portugol) are intentionally **out of scope** unless explicitly added later. They differ enough in vector syntax, subprogram syntax, and standard library that supporting them as one language is a tax, not a feature.

**In scope (v1):**
- Source files: `.alg`, UTF-8 (also accept Windows-1252 — VisuAlg's native encoding — and transcode on read)
- Full language: types, expressions, control flow, procedures/functions, vectors (one or two dimensions), pass-by-reference (`var` parameters)
- Standard library: numeric, string, conversion, random
- CLI: `run`, `check`, `fmt`, `repl`
- Diagnostics with stable error codes and source positions

**Out of scope (v1):**
- File I/O (`arqabertura`, etc.) — historical, niche, design later
- GUI primitives — VisuAlg has none; some dialects do, ignore them
- Native code emission, bytecode VM — tree-walker is enough until profiling says otherwise
- LSP, debugger, formatter beyond `fmt` — defer

---

## 2. Repository layout

```
cmd/portugol/          CLI entry point (subcommands: run, check, fmt, repl)
internal/token/        TokenKind, keyword table, position type
internal/lexer/        Scanner — produces []Token
internal/ast/          AST node definitions, visitor interface
internal/parser/       Recursive-descent parser — []Token → *ast.Program
internal/sema/         Symbol table, scope resolution, type checker
internal/runtime/      Value types, Environment, call frames
internal/interp/       Tree-walking evaluator
internal/stdlib/       Built-in functions (abs, raizq, copia, ...)
internal/diag/         Diagnostic struct, error codes, source position rendering
internal/repl/         Interactive shell
testdata/              .alg fixtures + .out / .err golden files
examples/              Sample programs (also used as integration tests)
docs/language.md       Authoritative language reference for this implementation
```

`internal/` is used aggressively. Public API stays empty until the language is stable.

---

## 3. Build, test, lint

```sh
go build ./...
go test ./...
go test -race -count=1 ./...
go vet ./...
staticcheck ./...                  # honnef.co/go/tools/cmd/staticcheck
golangci-lint run                  # config in .golangci.yml
```

CI must run: `go vet`, `staticcheck`, `go test -race`, `go test -run=TestFuzz -fuzz=. -fuzztime=30s` on lexer + parser.

No third-party dependencies in core language packages (`token`, `lexer`, `ast`, `parser`, `sema`, `runtime`, `interp`, `stdlib`, `diag`). CLI may depend on `cobra` or stdlib `flag` — prefer `flag` unless we have multi-level subcommands.

---

## 4. Architecture

Pipeline, single direction, no cycles:

```
source bytes
  → lexer        (token stream, lazy)
  → parser       (AST, errors collected, not panicked)
  → sema         (symbols + types, errors collected)
  → interp       (Environment + Value)
```

Each stage:
- **Never** panics on bad input; returns `[]diag.Diagnostic` instead.
- May panic only on internal-invariant violations (use `panic(fmt.Errorf("internal: ..."))`).
- Has a single public entry function and an opaque internal state struct.

The lexer is **case-preserving** in token text but Sema canonicalizes identifiers and keywords to lowercase. Portugol is case-insensitive.

---

## 5. Coding conventions

- `gofmt` enforced. `go vet` clean. `staticcheck` clean.
- Errors are values. The only `panic` allowed is for impossible internal states.
- Table-driven tests. No `testify` unless we hit a real ergonomic wall — prefer stdlib `testing` + small helpers.
- Source positions are mandatory on every AST node and every diagnostic. Use a single `token.Pos` (offset) + a `token.File` for line/col rendering — same shape as `go/token`.
- One concept per file. Long files (>500 lines) are a smell.
- Exported identifiers in `internal/...` still get doc comments — agents read them.
- No global mutable state. The interpreter is an instantiable struct.
- Concurrency: the interpreter is single-goroutine. Don't add goroutines unless the language requires it (it doesn't).

OSS hygiene: keep PR diffs minimal and focused. One feature or fix per PR. Don't reformat unrelated code.

---

## 6. Language reference (VisuAlg 3.x)

This section is a quick reference for agents. The authoritative spec lives in `docs/language.md`. If the two disagree, `docs/language.md` wins and this file gets a PR.

### 6.1 Program structure

```
algoritmo "nome do algoritmo"
var
  x, y: inteiro
  nome: caractere
  v: vetor[1..10] de real
  m: vetor[1..3, 1..3] de inteiro
inicio
  // statements
fimalgoritmo
```

`algoritmo`, `var`, `inicio`, `fimalgoritmo` are all required for a complete program. `var` block may be empty (omit it entirely, or `var` with no declarations).

### 6.2 Types

| Portugol     | Go representation |
|--------------|-------------------|
| `inteiro`    | `int64`           |
| `real`       | `float64`         |
| `caractere` / `caracter` | `string` |
| `logico`     | `bool`            |
| `vetor[a..b] de T` | slice with index offset; bounds checked |

Booleans literals: `verdadeiro`, `falso`. String literals use `"..."`. No char type.

Recorded keyword aliases include `função`, `então`, `senão`, `faça`, `até`, and
`não`, in any letter case. Canonical output uses unaccented keywords and
`caractere`. Do not remove accents indiscriminately: `lógico` is not an accepted
type spelling.

### 6.3 Operators

- Arithmetic: `+ - *`, `/` (numeric pairs return real; other scalar pairs return the right operand), `\` (integer pairs truncate; other scalar pairs return the right operand), `%` or `MOD` (recorded remainder rules in `docs/language.md`), `^` (numeric power, with recorded no-value and domain rules)
- Relational: `=`, `<>`, `<`, `>`, `<=`, `>=`
- Logical: `e`, `ou`, `nao`, `xou`
- String concat: `+` (when both operands are `caractere`)
- Assignment: `<-`

Precedence (high → low): unary `+ -`, left-associative `^`, `nao`, `* / \ % MOD`, `+ -`, relational, `e`, `xou`, `ou`. Recorded VisuAlg 3.0.7 probes pin `2^3^2 = 64` and `-2^2 = 4`; parentheses override those rules.

### 6.4 Control flow

```
se <cond> entao ... senao ... fimse

escolha <expr>
  caso <v1>, <v2>: ...
  caso <v3>: ...
  outrocaso: ...
fimescolha

enquanto <cond> faca ... fimenquanto

repita ... ate <cond>

para <i> de <a> ate <b> [passo <p>] faca ... fimpara

interrompa     // break out of innermost loop
```

`para` semantics: `i` is `inteiro`, `passo` defaults to 1, supports negative step. Loop variable is mutable inside the body but reassigning it does not affect iteration count. Recorded exit-state rules, including descending, empty, and interrupted loops, are defined and tested in `docs/language.md`.

### 6.5 I/O

- `leia(x, y, ...)` — read one input line per variable, preserving character input and applying the recorded scalar conversion rules
- `escreva(...)` — write without newline
- `escreval(...)` — write with newline
- Format specifiers: `x:n` for width, `x:n:m` for real width and decimals

The CLI uses the recorded VisuAlg `en-US` output profile deterministically: decimal `.`, numeric/logical leading spaces, uppercase logical output, and LF. Input accepts comma or dot. Width and rounding rules are documented in `docs/language.md`; other reference locales remain unverified.

### 6.6 Subprograms

```
procedimento P(a: inteiro; var b: real)
var
  // local declarations
inicio
  ...
fimprocedimento

funcao F(x: inteiro): real
var
  ...
inicio
  ...
  retorne <expr>
fimfuncao
```

- Default pass-by-value. `var` parameter = pass-by-reference.
- `retorne` sets the function result and continues execution; `fimfuncao` ends the call. Paths without `retorne` are accepted. Recorded result initialization and reuse rules are defined in `docs/language.md`.
- Recursion is allowed.
- Forward declarations are not part of the language; declarations must precede use, but the parser collects all top-level declarations first so order within a file does not matter.

### 6.7 Built-in functions (initial set)

Numeric: `abs`, `arccos`, `arcsen`, `arctan`, `cos`, `cotan`, `exp`, `grauprad`, `int`, `log`, `logn`, `pi`, `quad`, `radpgrau`, `raizq`, `sen`, `tan`, `randi`; legacy `aleatorio` corrections remain pending. `frac` is rejected as undeclared.
`pi` is written without parentheses. The new inverse-trigonometric and angle-conversion calls, `cotan`, and `quad` follow the recorded optional-argument and no-value rules in `docs/language.md`.
`exp(base, exponent)` takes two numeric arguments and returns real-valued power. `log` is base ten and `logn` is the one-argument natural logarithm. Empty calls, no-value propagation, argument order, and domain failures follow `docs/language.md`.
String: `copia(s, p, n)`, `maiusc`, `minusc`, `asc`, `carac`, `compr`, `pos`
Conversion: implicit `inteiro` → `real`; explicit elsewhere via built-ins (`int`, etc.)
`numpcarac(x)` converts a number to text with 15 significant digits; `numpcarac()` returns `"0"`. Its recorded no-value behavior for nonnumeric input is documented in `docs/language.md`.

VisuAlg strings are **1-indexed** in `copia` and `pos`. Don't make it 0-indexed "because Go". This is the language we're implementing.

---

## 7. Pitfalls and behavioral decisions

These cost time when wrong. Each must have a regression test.

1. **Declared vector indexing.** Bounds may start at zero or a positive integer; store an offset, do not assume 0 or 1. Vectors have at most two dimensions. An omitted second index selects that dimension's lower bound. Whole-vector assignment is rejected.
2. **Integer vs real division.** `/` produces `real` for numeric pairs and otherwise returns the right scalar operand. `\` truncates two integers toward zero; other scalar pairs return the right operand and its type after evaluating both operands. Mixing integer and real operands in `+ - *` promotes to real. `%` and `MOD` follow the recorded divisor and conversion rules in `docs/language.md`.
3. **Short-circuit `e` / `ou`.** VisuAlg historically does **not** short-circuit. Decide once, document, test both branches always evaluate. This is a common source of student bugs and we should not silently change it.
4. **Case-insensitivity.** `Soma`, `soma`, `SOMA` all refer to the same identifier. Canonicalize at the symbol-table boundary. Keywords likewise.
5. **Encoding.** Real VisuAlg files are Windows-1252. Detect BOM / UTF-8 validity; otherwise assume CP1252 and transcode. Never read as raw bytes into a Go string and hope.
6. **Number formatting on output.** Use the recorded deterministic profile in `docs/language.md`; compare fixture bytes exactly.
7. **Reading multiple values with `leia`.** Each variable consumes one complete input line. Preserve spaces and empty character lines, share unread input across calls, and emit the recorded typed input echo.
8. **`escolha` fall-through.** Does **not** fall through. Each `caso` is independent.
9. **`interrompa` outside a loop** is a sema error, not a runtime error.
10. **Uninitialized variables.** VisuAlg gives them zero values per type. Match this; do not error on read-before-write.

---

## 8. Testing strategy

- **Lexer:** golden token streams per fixture. Tokens printed one per line: `KIND text @line:col`.
- **Parser:** golden AST via a deterministic compact printer in `internal/ast/print.go`. No JSON. No struct tag noise.
- **Sema:** every diagnostic has a stable error code (`E001` type mismatch, `E002` undeclared identifier, ...). Tests assert on the code, not the message.
- **Interpreter:** end-to-end. `testdata/run/<name>.alg` paired with `<name>.in` (stdin) and `<name>.out` (expected stdout). A single `TestRun` walks the directory.
- **Fuzz:** `FuzzLexer` and `FuzzParser` must run in CI. Goal: no panics on arbitrary input, ever.
- **Property test:** `parse(print(ast)) == ast` for the canonical printer (round-trip). Add this once the printer exists.
- **Benchmarks:** `Benchmark*` colocated with code. We don't optimize until we have a benchmark to point at.

Never land a bug fix without a regression fixture.

---

## 9. Workflow

- One feature per PR. Diffs minimal.
- Every user-visible change updates `docs/language.md` and `CHANGELOG.md` in the same PR.
- New language feature → at least one fixture in `testdata/` and one program in `examples/`.
- Tests, vet, staticcheck, race detector all green before requesting review.
- Don't reformat or rename across the tree opportunistically. Focused diffs only.

---

## 10. Open questions

These are real decisions, not rhetorical. Resolve before implementing the affected area.

1. **Dialect target.** VisuAlg only, or also Portugol Studio (UNIVALI)? They differ on vector syntax (`vetor[10]` vs `vetor[1..10]`), subprogram syntax, and stdlib. Pick one for v1.
2. **Short-circuit `e` / `ou`.** Spec-faithful (no SC) or pragmatic (SC)? Affects observable behavior of programs with side effects in conditions.
3. **Decimal separator on I/O — resolved for the current profile.** Output uses the recorded `en-US` decimal dot on every host; input accepts comma or dot. Other reference locales remain unverified.
4. **`aleatorio` semantics.** Match VisuAlg's RNG exactly (would need to reverse-engineer it) or use Go's `math/rand/v2` with documented seeding?
5. **File I/O.** v1 = no, v2 = maybe. Confirm.
6. **CLI framework.** stdlib `flag` or `cobra`? Default: `flag`.

---

## 11. References

- VisuAlg manual / reference (Apoio Informática). Confirm canonical URL before relying on it; the project moved hosts.
- Forbellone & Eberspächer, *Lógica de Programação* — the textbook most VisuAlg dialect choices trace back to.
- Nystrom, *Crafting Interpreters* — design reference for the tree-walker.
- `go/token`, `go/ast`, `go/parser` in the Go standard library — model for our package shapes.

---

## 12. Conventions for agents working on this repo

- Read `docs/language.md` before changing language behavior. If it doesn't cover the case, update it as part of the same PR.
- When adding a built-in, update: `internal/stdlib/`, `internal/sema/` (signature), `docs/language.md`, `testdata/run/`.
- When adding an AST node, update: `internal/ast/`, the printer, the parser, sema, interp, and at least one fixture.
- Don't introduce dependencies in core packages. Ask first.
- Don't paper over an ambiguity in the spec by picking a behavior silently. Add a test that pins the choice and a note in `docs/language.md`.
