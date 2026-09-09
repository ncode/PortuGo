# Portugol-Go Language Reference

This implementation targets the VisuAlg 3.x dialect used by Apoio Informática.
Other Portugol dialects are intentionally out of scope for v1.

## Program Structure

A complete program has the shape:

```portugol
algoritmo "nome"
var
  x: inteiro
inicio
  escreval(x)
fimalgoritmo
```

The `var` block may be omitted. Top-level `procedimento` and `funcao`
declarations must appear before `inicio`.

## Types

- `inteiro`: signed integer
- `real`: floating point number
- `caractere`: string
- `logico`: boolean
- `vetor[a..b] de T`: one or more declared bounds with checked indexing

Variables are initialized to the zero value of their type.

## Expressions

Arithmetic operators are `+`, `-`, `*`, `/`, `\`, `%`, `MOD`, and `^`.
The `/` operator always returns `real`; `\` truncates integer division toward
zero. Integer operands are promoted to real where needed.

Power is left-associative: `2^3^2` means `(2^3)^2` and produces 64.
Unary minus binds more tightly than power: `-2^2` produces 4. Use `-(2^2)`
for -4. These rules are pinned by recorded VisuAlg 3.0.7 probes.
Even exact `/` results remain real: assigning `4/2` to an integer is rejected.

Relational operators are `=`, `<>`, `<`, `>`, `<=`, and `>=`.
Logical operators are `e`, `ou`, `xou`, and unary `nao`.

The operators `e` and `ou` intentionally do not short-circuit; both operands
are evaluated to match VisuAlg behavior.

## Statements

Supported statements:

- Assignment with `<-`
- `leia`, `escreva`, and `escreval`
- `se ... entao ... senao ... fimse`
- `escolha ... caso ... outrocaso ... fimescolha`
- `enquanto ... faca ... fimenquanto`
- `repita ... ate`
- `para ... de ... ate ... passo ... faca ... fimpara`
- `interrompa`
- procedure calls
- `retorne` inside functions

`interrompa` outside a loop is a semantic error.

`para` evaluates its bounds and step once. The default step is 1; zero is
rejected. Assignments to the loop variable are visible within the body but do
not change the iteration sequence. An empty loop leaves the variable at its
initial bound. On normal completion of a nonempty loop, its exposed value is
the smaller of the next iteration value and the terminal bound. On `interrompa`,
it is the smaller of the current body value and the terminal bound. This also
applies to descending loops: `5 ate 1 passo -2` visits 5, 3, 1 and leaves -1;
`1 ate 6 passo 2` visits 1, 3, 5 and leaves 6. These unusual exit rules follow
the recorded program output; the reference GUI memory grid can disagree.

## I/O

`leia` consumes one whitespace-delimited token per destination. Real input may
use either `,` or `.` as the decimal separator.

`escreva` writes without a newline and `escreval` writes with a newline.
The portable CLI uses a deterministic profile matching the recorded `en-US`
reference environment: `.` as decimal separator and LF for newlines. It does
not change formatting with the host locale. Other VisuAlg locales remain
unverified; input continues to accept either decimal separator.

Format specifiers are supported as `expr:width` and `expr:width:decimals`.
Without a positive width, each number has one leading space, reals omit
unnecessary fractional zeros, and a decimal count is ignored. With a positive
width, numbers are right-aligned, the decimal count defaults to zero, and
decimal ties round away from zero. Integer values can also request decimals;
numeric fields expand if needed. Strings are left-aligned and truncated to a
positive width. Logical output is ` VERDADEIRO` or ` FALSO`, including its
leading space; field widths on logical values are rejected.

## Built-ins

Numeric built-ins: `abs`, `raizq`, `exp`, `log`, `logn`, `pi`, `sen`, `cos`,
`tan`, `int`, `frac`, and `aleatorio`.

`exp(base, exponent)` takes two numeric arguments and returns real-valued
exponentiation. The one-argument natural-exponential form is not supported.

String built-ins: `copia`, `maiusc`, `minusc`, `asc`, `carac`, `compr`, and
`pos`. String positions are 1-indexed.

`aleatorio()` returns a real in `[0, 1)`. `aleatorio(n)` returns an integer in
`[0, n)`. `aleatorio(a, b)` returns an integer in the inclusive range `[a, b]`.
The generator uses Go's standard pseudo-random source; VisuAlg's exact RNG is
not emulated.

## Out Of Scope

File I/O and GUI primitives are not implemented in v1.

## Validation status

The current implementation is a VisuAlg-like subset. Fourteen selected VisuAlg
3.0.7 observations are replayed by `TestRecordedWindowsProbes`; their sources,
raw panel text, and validation results are linked from the [Windows report](windows-lab-validation-2026-09-07.md).
Accepted program output is compared after removing only the two application
notices and converting CRLF to the CLI's LF. Invalid programs receive positioned
static diagnostics and are never executed: unlike the reference GUI, the CLI
does not print preceding statements before reporting a statically known error.
`run` exits 1 for lexical, parsing, semantic, or runtime failures and 0 on success.
The active OpenSpec change tracks the remaining work; these selected observations
are not a claim of full reference conformance.

The [reference corpus progress](reference-corpus-progress.md) tracks the broader
recording work. New observations confirm particular constant, named-type and
record forms in the reference application; those forms remain pending in this
implementation. Corpus validation requires expected acceptance/rejection and
generated files to agree with the reference, and an explicit history base to
detect unreviewed regressions in coverage. Removing a probe, inventory entry,
or linked task requires a reviewed retirement, even for pending behavior.
Retirement records remain in later manifests; retired IDs cannot be reused.
Recording verifies input-only file hashes before accepting a capture.
Requirement traceability scans the whole change's specification tree, including
files omitted from the manifest's declared source list.
The [project tooling evidence](conformance-project-evidence.md) separately traces
provenance, recording, normalization, and validation tests. Those mappings do not
substitute for language recordings or establish complete reference conformance.

Runtime fixture output is compared byte for byte, including decimal separators,
whitespace, and newlines. Git preserves committed fixture bytes on every
platform, including Windows. Lexer/parser fuzz tests use a 64 KiB generated-source
profile and adversarial cases have failing subprocess watchdogs. This profile
is a test bound, not an enforced language source or nesting limit. Production
resource guards are planned in groups 3 and 4 of the conformance change.

See [development checks](development.md) and the
[quality baseline](quality-baseline.md) for commands and measured coverage.
