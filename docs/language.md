# Portugol-Go Language Reference

This implementation targets the VisuAlg 3.x dialect used by Apoio Informática.
Other Portugol dialects are intentionally out of scope for v1.

The finite automated example sweep excludes complete runs of the
[two nonterminating bundled examples](bundled-external-stop.md). Their original
sources and partial reference observations are retained. Their language features
remain required, and a timeout is never a successful completion result.

The [September 9 reference observations](reference-observations-2026-09-09.md)
record further compatibility gaps in syntax, literals, declarations, and calls.
Their implementation remains pending unless a corpus entry is explicitly
verified. The reference accepts vectors larger than 500 elements; the draft
500-slot compatibility restriction has therefore been withdrawn. The recordings
do not establish the upper storage limit in every declaration context.

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
not emulated. Each interpreter has its own source. The inclusive interval
covering every signed 64-bit integer is rejected with `R007` because its draw
size cannot be represented by the random-source interface.

## Execution diagnostics and safeguards

`run`, `check`, `fmt`, and `repl` use exit status 0 for success, 1 for source,
execution, or operational failures, and 2 for invalid command usage. Diagnostics
go to stderr with the source filename, line, column, and a stable code. Runtime
failure preserves preceding stdout. `check` and `fmt` never execute the program.
REPL submissions share one buffered input stream with `leia`; a failed
submission does not prevent a later submission, but the session exits 1 if any
submission failed. A blank line still submits a program in the current REPL.

| Runtime code | Category |
| --- | --- |
| R001 | Type or coercion; missing semantic information |
| R002 | Arithmetic |
| R003 | Storage/indexing; text size or expression depth |
| R004 | Input |
| R005 | Call/return; active call depth |
| R006 | Loop state; execution budget |
| R007 | Built-in function failure |
| R008 | Host or output I/O |

`run --max-steps N file.alg` and `repl --max-steps N` allow a nonnegative work
budget. Zero, the default, is unlimited. A step is charged before each statement,
expression/designator evaluation, and loop iteration, even with an empty body.
Calls share the run's counter; each REPL submission starts a fresh counter.
The next operation stops before its side effects when the budget is exhausted.
This work budget cannot interrupt a blocking reader or host operation.
For example, `portugol run --max-steps 20 examples/execution_budget.alg`
stops the sample loop early; omitting the flag lets it finish.

Regardless of that budget, execution permits at most 256 active language calls,
256 expression-evaluation levels within one call frame, and 16 MiB in one text
value, input token/line, or formatted item. Concatenation, case conversion, and
width/precision expansion check sizes before creating the result. Vector layouts
are checked for address-space overflow; reference-specific storage quotas remain
pending. These are project safeguards, not measured VisuAlg limits. A safeguard
hit in an accepted reference example remains a conformance failure.

Source files and REPL submissions are capped at 4 MiB of original bytes,
including a BOM or CRLF bytes. Decoding checks this before allocating the UTF-8
result, including Windows-1252 expansion. Syntax nesting and AST traversal are
limited to 256 levels. A root statement or declared type starts at level one;
its nested statements, expressions, or element types add levels. Redundant
parentheses also count during parsing. Long flat expression chains are checked
after parsing because their resulting AST can still be deep. Analysis and
printing check the same AST bound before traversal or output.

Source, syntax, and AST limits report positioned `E900` diagnostics. Exact
boundary inputs remain accepted. An oversized REPL submission ends the session
with status 1; recovery from that limit remains part of the later REPL work.

Formatting omits unnecessary operator parentheses while preserving precedence
and association. The CLI checks the complete formatted source against the same
byte and syntax limits before writing stdout. If formatting expands an accepted
input beyond a limit, `fmt` reports `E900`, exits 1, and writes no source output.

Internally, `sema.Analyze` supplies immutable resolved types, vector layouts, and
declaration/use bindings. `interp.New(Options).Run(program, info)` requires that
successful result for the same unchanged AST and returns positioned diagnostics.
Nil input/output mean empty input and discarded output; the working directory
resolves at construction. Hosts and randomness can be injected. The default
headless host uses real time and silent, nonblocking UI operations. Display
options and language host commands await reference recordings.

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
is the smaller generated-input test profile. The source, traversal, and
execution safeguards above also apply to ordinary use.

See [development checks](development.md) and the
[quality baseline](quality-baseline.md) for commands and measured coverage.
