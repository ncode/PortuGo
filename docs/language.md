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

## I/O

`leia` consumes one whitespace-delimited token per destination. Real input may
use either `,` or `.` as the decimal separator.

`escreva` writes without a newline and `escreval` writes with a newline.
Real output uses a comma decimal separator by default.

Format specifiers are supported as `expr:width` and `expr:width:decimals`.

## Built-ins

Numeric built-ins: `abs`, `raizq`, `exp`, `log`, `logn`, `pi`, `sen`, `cos`,
`tan`, `int`, `frac`, and `aleatorio`.

String built-ins: `copia`, `maiusc`, `minusc`, `asc`, `carac`, `compr`, and
`pos`. String positions are 1-indexed.

`aleatorio()` returns a real in `[0, 1)`. `aleatorio(n)` returns an integer in
`[0, n)`. `aleatorio(a, b)` returns an integer in the inclusive range `[a, b]`.
The generator uses Go's standard pseudo-random source; VisuAlg's exact RNG is
not emulated.

## Out Of Scope

File I/O and GUI primitives are not implemented in v1.

## Validation status

The current implementation is a VisuAlg-like subset. Its behavior has not yet
been qualified against recorded VisuAlg 3.0.7 observations. The active OpenSpec
change tracks that work; passing implementation tests alone is not a claim of
reference conformance.

Runtime fixture output is compared byte for byte, including decimal separators,
whitespace, and newlines. Git preserves committed fixture bytes on every
platform, including Windows. Lexer/parser fuzz tests use a 64 KiB generated-source
profile and adversarial cases have failing subprocess watchdogs. This profile
is a test bound, not an enforced language source or nesting limit. Production
resource guards are planned in groups 3 and 4 of the conformance change.

See [development checks](development.md) and the
[quality baseline](quality-baseline.md) for commands and measured coverage.
