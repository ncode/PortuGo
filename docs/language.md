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

The [bundled example sweep](bundled-examples-progress.md) verifies six original
programs against reference output, including formatting and execution. Other
accepted examples still expose missing features, and three supplied files have
recorded reference errors. Bundled-file presence alone does not establish that
its syntax is accepted by this release.

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

The quoted algorithm name must follow `algoritmo` on the same physical line.
Text after the closing quote on that line is ignored, including another quoted
string or an unmatched quote. The executable body starts on a later line.
A missing or misplaced header produces one `P001` diagnostic; parsing stops
when the program header or the required main `inicio` is invalid. Formatting
emits only the algorithm name on the header line and discards ignored header text.
LF and CRLF inputs retain equivalent token kinds, line/column positions, and
syntax. An unrecognized statement start is diagnosed once; parsing resumes on
the next physical line to retain independent later errors.

The `var` block may be omitted or left empty. Top-level `procedimento` and
`funcao` declarations must appear before `inicio`. Recorded words and statements
after `fimalgoritmo` are ignored during execution. Formatting currently discards
that suffix; preserving its text remains pending.

## Types

- `inteiro`: signed integer
- `real`: floating point number
- `caractere` or `caracter`: string
- `logico`: boolean
- `vetor[a..b] de T`: one or two declared dimensions with checked indexing

Variables are initialized to the zero value of their type.

Recorded keyword aliases are `função` (`funcao`), `então` (`entao`), `senão`
(`senao`), `faça` (`faca`), `até` (`ate`), and `não` (`nao`). Recognition ignores
case and retains token spelling and positions. Formatting emits unaccented
keywords and the canonical type name `caractere`, including parameters and
function results. See the [keyword example](../examples/keyword_aliases.alg).

These are specific accepted spellings. The recorded type spelling `lógico` is
rejected with `P001` on its declaration line. Further vocabulary and physical-line
grammar work remains pending; see the [keyword observations](keyword-forms-progress.md).

Recorded identifiers accept leading and internal underscores. Identifier
matching ignores case: a declaration named `SoMa` can be assigned through `soma`
and read through `SOMA`. Broader identifier character rules remain under
validation.

Recorded real literals accept exponent notation such as `1e2` and a trailing
decimal point such as `5.`. A leading decimal point, as in `.5`, is rejected with
`L001` on its source line. A quoted string reaching a newline without its closing
quote is also rejected with `L001`.

## Vectors

Recorded vector declarations accept one or two dimensions. Literal bounds are
unsigned integers with the upper bound at least as large as the lower bound;
`0..0`, `0..2`, and `2..4` are accepted. Signed forms (including `-0` and `+1`),
fractional literals, parentheses, arithmetic bound expressions, reversed ranges,
and a third dimension receive `P001` on the declaration line. Named constants
in bounds remain pending implementation and are not excluded by these literal
syntax rules.

Each index uses its dimension's declared bounds. A two-dimensional vector
access with one index selects the second dimension's lower bound: for
`vetor[2..3,4..5]`, `v[2]` and `v[2,4]` select the same element. Explicit second
indices do not change this default. An extra index receives `E001` before
execution. An out-of-bounds index receives `R003` at the indexing expression;
preceding program output is retained. Integer elements begin at zero.

Whole-vector assignment, including self-assignment and assignment to a scalar,
is rejected with `E001`. Element assignment remains supported. Recorded vectors
with 500, 501, 5000, and 5001 elements execute successfully; these observations
do not establish a universal storage maximum. See the
[vector example](../examples/vector_bounds.alg). Aggregate parameter behavior
and allocation accounting remain under validation.

Storage lookup checks layout metadata and the exact backing length before
computing an offset. Invalid ranges, overflowing dimension products, and
inconsistent backing storage produce a controlled failure. Interpreter reads,
writes, input, and reference arguments report this as `R003` at the indexing
expression without changing elements or consuming input. These corruption
checks are project safeguards; malformed Go storage objects have no source
language counterpart.

## String Literals

Strings use double quotes. A backslash is an ordinary character: `"a\nb"`
contains four characters, and `"tail\"` ends with a literal backslash. There
are no Go-style escape sequences. An embedded quote cannot be written as `\"`
or doubled quotes. The recorded backslash forms are covered by byte-exact
reference replay; formatting preserves their contents and program names.
See [the literal-string example](../examples/literal_strings.alg).

## Comments

`//` consumes the rest of its physical line. The reference also treats it as
a comment inside quotes, so a string containing `//` is rejected as unterminated.
Single slashes, `/*`, `*/`, and braces otherwise remain literal inside strings.

Outside strings, `{` and `}` consume the rest of the line. At the first
non-whitespace position, `/` and `*` also consume the line, including `/*` and
`*/`. These forms do not open multiline blocks and need no closing delimiter.
Statements on subsequent lines still execute. After a complete statement,
`/` and `*` retain their operator meaning; `/*` is not an inline comment there.
See [the comment example](../examples/comment_lines.alg).

Thirty-one recorded cases cover these forms, quoted delimiters, and formatting
without changing execution. Comments are currently discarded when formatting.
The reference's handling of comments within incomplete expressions and the
remaining syntax-recovery behavior are still pending.

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

- Assignment with `<-` or `:=`; formatting uses the canonical `<-` spelling
- `leia`, `escreva`, and `escreval`
- `se ... entao ... senao ... fimse`
- `escolha ... [faca] ... caso ... outrocaso ... fimescolha`
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

## Subprogram Calls

Parameterless procedure and function declarations may omit the empty `()`.
Procedures can be invoked as `P` or `P()`. Functions can be used as `F` or
`F()` in expressions. Formatting emits parentheses on declarations and
procedure calls, and retains the accepted bare spelling of function values.
See [the parameterless-call example](../examples/parameterless_calls.alg).

Each function occurrence executes a call: a recorded counter function used
twice produces successive values. Function names take priority in expressions,
even over a local variable or parameter with the same spelling. Assignments
still target the declared variable. Global variables may also share a function
or procedure name. Functions and procedures share a separate namespace: duplicate
callable names are rejected with `E003` at the second declaration, before bodies
are analyzed. See the [name-priority example](../examples/call_name_priority.alg).

Procedure names take priority in bare and parenthesized call statements, even
over a local variable or parameter. An assignment starting with a procedure
name, including an indexed assignment, is rejected with `E004` at the procedure
declaration. Reading a colliding variable as a value still reads that variable.

Variables cannot be called, procedures cannot produce values, and user functions
cannot be invoked as statements. Bare functions that require arguments receive
`E004`. Function results are not writable reference-argument storage; the CLI
rejects them statically. Reference probes for temporary reference arguments
caused application faults and remain inconclusive, so this guard is not claimed
as a matching reference rejection.

Arguments are evaluated once from left to right; a reference argument's vector
element is selected before later arguments execute. Calls use the callee's
lexical globals, parameters, and locals, independent of the caller's locals.
Recorded mutually recursive functions resolve later declarations.

Numeric value and `var` parameters accept integers or reals. Passing a real to
an integer parameter truncates toward zero. A `var` parameter receives its own
converted copy: writes through it become visible to the caller on successful
return. Direct writes to a global remain separate during the call. Copies return
in parameter order, so the last parameter targeting the same location wins.
Nested and recursive calls retain independent parameter copies. See the
[argument example](../examples/call_arguments.alg).

Copy-back also carries the parameter's numeric type. For example, passing a real
variable to an integer `var` parameter leaves an integer value on return;
assigning `2.5` to it then fails with `R001`, preserving earlier output. Passing
an integer variable to a real `var` parameter can leave a fractional value.
Out-of-range or non-finite integer arguments fail before the body executes;
failed call setup does not copy partially prepared parameters back.

Recorded procedure argument-count and type errors point to the declaration
line; function argument-count errors point to the call. Empty-argument edge
cases and integer operators after numeric reference conversion remain under
validation.

## Function Results

`retorne <expression>` updates the active function's result and execution
continues. It does not leave an `se`, `escolha`, `enquanto`, `repita`, or `para`
body. Later statements run and a later `retorne` replaces the earlier result.
`fimfuncao` ends the call and supplies the result; successful completion then
copies `var` parameters back, including changes made after `retorne`. See the
[function-result example](../examples/function_results.alg).
The recorded choice body also accepts a `caso` label without a trailing colon;
formatting emits the colon consistently.

Functions may reach their end without executing `retorne`. A result that has
never been assigned starts at the scalar zero value: `0`, empty text, or
`falso`. Active recursive calls have independent results, parameters, and locals.
Completed calls at the same call depth retain the previous result when the
result type matches, even across different function names and frame sizes.
Procedure calls do not replace that retained result. Starting a new program
clears it. Argument expressions run before the callee selects its retained
result, so a function call inside an argument can replace that result.
Cross-type fallthrough uses a fresh typed zero as a project guard;
the reference exposes internal storage in these cases, so this behavior is not
claimed as reference-equivalent and those captures are not published.

Return expressions accept integer-to-real widening, but an integer function
rejects real expressions with `E001`, including `2.0`. Numeric argument narrowing
does not apply to returns. Valued returns in procedures or the main body receive
`E005`. Assigning to a function name without a variable of that name receives
`E002`; it does not assign the result. A return expression must begin on the
same physical line as `retorne`. A bare return in a function receives `E001`
on the return line; a bare return in a procedure or the main body receives
`E005`. The following statement or terminator is retained for analysis.
Diagnostic ordering when a file also contains malformed syntax remains pending.

## I/O

`leia` consumes one complete line per destination. Character input preserves
spaces and accepts an empty line. LF and CRLF terminate input lines; a final
nonempty line also works without a terminator. Multiple destinations,
consecutive reads, and reused interpreters share the same buffered input.
In the REPL, the line immediately after `fimalgoritmo` is available to `leia`;
an extra blank line is input data.

Console input is echoed after conversion: integers use plain decimal text,
reals use ten fractional digits, logical values use `Verdadeiro` or `Falso`,
and characters retain their exact text. Each echo ends with LF and is separate
from `escreva`/`escreval` formatting. Real input accepts either decimal separator.

Numeric input ignores leading ASCII spaces and accepts a sign and decimal
exponent. Malformed text retains the unsigned mantissa before decimal and
exponent scaling: `1.5x`, `-1.5x`, and `1.5 ` each produce `15`. Empty input
or a nonnumeric prefix produces zero. Integer input truncates the numeric
result toward zero and narrows it to a signed 32-bit value, including wrapping:
`2147483648` becomes `-2147483648` and `4294967297` becomes `1`.
Numeric input discards the sign of zero. Invalid text such as `x7` replaces
an existing numeric value with zero.

Logical input is true exactly when its first character is `v` or `V`.
Leading spaces, empty text, English `true`, and numeric `1` therefore produce
false; `v`, `VERD`, and `Verdadeiro ` produce true.

The headless runtime reports positioned `R004` for input exhaustion, a number
outside the finite real range, or integer conversion outside the signed 64-bit
range before narrowing. These are project guards. Failures retain the
destination and preceding output; successful conversions replace the value
before echoing it. An echo writer failure reports `R008` at that destination.
The input line buffer is limited to 16 MiB and reports `R003` at the destination
before growing beyond that limit.

`escreva` writes without a newline and `escreval` writes with a newline.
Either command may stand alone on its physical line: bare `escreva` emits
nothing and bare `escreval` emits one newline. Arguments require parentheses
on that line; unparenthesized values or another command after a bare write
produce `P001` at the write statement.
Each statement evaluates and formats its arguments in order before emitting
its own text. Output from functions called during that evaluation appears first.
`escreval` requests a pending newline, which the next output statement to finish
consumes. This includes a nested `escreva` or a bare output command, so the outer
statement may finish without another newline. Starting a new program clears
this pending newline, including after a failed run.

An argument-evaluation failure discards that statement's buffered items while
preserving output from earlier statements and completed nested writes.
Pending buffered output is limited to 16 MiB across nested statements as a
project resource guard; exceeding it reports `R003` before emitting the
affected statement's items. Completed statements release their buffer capacity.

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
submission failed. A real `fimalgoritmo` token submits the program immediately;
the same text inside a string, comment, or longer identifier does not. Blank
lines inside unfinished source are preserved, including their source positions.
Consecutive programs need no separator, and `leia` consumes following input from
the shared stream before source entry resumes. EOF submits any buffered source
and reports incomplete input; `:sair` cancels the buffer and exits. Submitted
source uses the same UTF-8, UTF-8 BOM, and Windows-1252 decoder as files.
Lexical, syntax, semantic, runtime, and execution-budget failures clear the
submitted buffer and permit another program; each program gets a fresh budget.

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
boundary inputs remain accepted. An oversized REPL submission reports `E900`
once and is discarded. The REPL then waits for a line whose first token is
`algoritmo`, preserving that line and subsequent input for the next program.
Header lookalikes in strings or comments do not restart entry. `:sair` and EOF
still exit during recovery, and the session ultimately returns status 1 even
when a later program succeeds. A transient program output failure likewise
reports `R008` and permits another submission. Unusable prompt/input streams
still end the session with an operational error.

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
