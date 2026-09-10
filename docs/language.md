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

The [bundled example sweep](bundled-examples-progress.md) verifies 30 original
programs against reference output, including formatting and execution. Seven supplied files have
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

An optional `;` may end a `var` line or a scalar/vector declaration, including
local declarations and an empty `var` block. It must be on that physical line
and followed only by whitespace or a comment. Repeated semicolons, another
declaration after `;`, or a semicolon on its own line receive `P001` on that line.
Adjacent declarations without a separator remain accepted. Formatting emits
one declaration per line and omits the optional semicolon; see the
[declaration example](../examples/declaration_semicolons.alg).

## Types

A `const` section precedes `var`, globally or inside a subprogram. Each
declaration has the form `name = expression`, optionally ending with `;`.
The following `var` section is required, even when it is empty. Constants use
case-insensitive names and infer their scalar type from the initializer.
Initializers can use scalar literals, ordinary arithmetic, earlier constants,
and accepted built-ins, including `pi`, `abs`, and `randi`. Integer arithmetic
retains its signed 32-bit wrapping behavior.

Global constants initialize before the main body. Local constants initialize
once per call, after parameter setup and before local variables; their values
can depend on the current parameters and global variables. Local declarations
may shadow global constants. Repeated reads retain the initialized value.
Formatting preserves initializer expressions, declaration order, and the
required empty `var` section. See the [constant example](../examples/constants.alg).

Unknown or forward constant dependencies, including cycles, receive `P001`.
Duplicate names or collisions with a variable in the same scope receive `E003`.
Assignment to a constant receives `E002`. Aggregate-valued initializers and
no-value initializer results remain unqualified. The CLI currently requires a
scalar initializer type as a project guard. Existing expression-depth, step,
and arithmetic guards apply during initialization; an initializer failure does
not enter the subprogram body or copy reference parameters back.

- `inteiro`: signed integer
- `real`: floating point number
- `caractere` or `caracter`: string
- `logico`: boolean
- `vetor[a..b] de T`: one or two declared dimensions with checked indexing

Variables are initialized to the zero value of their type.

Recorded keyword aliases are `função` (`funcao`), `fimfunção` (`fimfuncao`),
`então` (`entao`), `senão`
(`senao`), `faça` (`faca`), `até` (`ate`), and `não` (`nao`). Recognition ignores
case and retains token spelling and positions. Formatting emits unaccented
keywords and the canonical type name `caractere`, including parameters and
function results. See the [keyword example](../examples/keyword_aliases.alg).
The function terminator `fimfunção` is also reserved as a variable name;
using it in a declaration receives `P001` on that line.

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

Digit-only decimal literals from `0` through `2147483647` have integer type.
Larger finite values have real type, including values beyond the signed 64-bit
range. A decimal point or exponent also makes a literal real. Leading zeros
do not change the type. The minus sign is a separate unary operator, so
`-2147483648` is real, while `-2147483647 - 1` is an integer expression.
The formatter retains a decimal point for integral real literals, preserving
their type when the result is parsed again. A number outside the finite
64-bit floating-point range receives positioned `P001` as a project guard.
Integer addition, subtraction, multiplication, and negation wrap to signed
32-bit results, including inside assignments and function returns. For example,
`2147483647 + 1` gives `-2147483648`, and negating that minimum gives the same
minimum. Mixing a real operand into addition, subtraction, or multiplication
uses real arithmetic. Division of the signed minimum by `-1` receives positioned
`R002` as a project guard. The recorded real-to-integer assignment diagnostic
behavior remains pending conformance work.

## Vectors

Recorded vector declarations accept one or two dimensions. Literal bounds are
unsigned integers with the upper bound at least as large as the lower bound;
`0..0`, `0..2`, and `2..4` are accepted. Signed forms (including `-0` and `+1`),
fractional literals, parentheses, arithmetic bound expressions, reversed ranges,
and a third dimension receive `P001` on the declaration line.

A bound may instead name an earlier integer constant, ignoring case. Such a
constant can contain arithmetic or be negative; `lo = -2` permits `lo..2`,
although a signed literal directly in the brackets is rejected. A bound's
syntax remains one literal or name: `1..n+1` is rejected. Real, text, and logical
constants, variable names, and direct parameter names receive `P001`.

Constant-based bounds resolve after constant initialization and before vector
allocation. A local constant may capture a parameter, so different calls can
allocate different layouts without changing shared semantic information.
Each allocated vector retains its own concrete bounds. Reversed resolved bounds
receive `P001` at the declaration; output from earlier successful calls remains
visible. See the [constant-bound example](../examples/constant_bounds.alg).

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
[vector example](../examples/vector_bounds.alg). Inline `vetor[...] de ...`
types in procedure/function parameters and function results receive `P001`,
as recorded. Named aggregate parameter forms remain under validation.

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

Arithmetic operators are `+`, `-`, `*`, `/`, `\`, `DIV`, `%`, `MOD`, and `^`.
The case-insensitive word `DIV` is an alias for `\`, with the same precedence,
operand types, and evaluation order. Formatting emits `\`. The word is reserved:
using `div` as a variable name receives `P001` on its declaration line.
The `/` operator returns `real` for numeric operands. Other scalar pairs return
the right operand unchanged: `"7" / 2` gives integer `2`, and `7 / "2"` gives
text `"2"`. With two integer operands, `\` truncates
division toward zero. For other scalar pairs, the recorded reference returns
the right operand unchanged, including its type: `8.5 \ 3` gives integer `3`,
`7 \ 2.5` gives real `2.5`, and `"7" \ 2` gives integer `2`. Both operands are
evaluated once, left to right, even when the result is just the right operand.

For numeric operands, `%` and `MOD` return an integer. A positive integer divisor
produces the remainder after the left operand is truncated and narrowed to
signed 32-bit, using the same conversion as `int`. Zero, negative, and real
divisors produce `-1`, including integral real divisors such as `2.0`. A numeric
or logical pair containing a logical operand returns the right operand unchanged.
Thus `(-7.5) MOD 2` gives `-1`, `7 MOD 0` gives `-1`, and `7 MOD verdadeiro`
gives `verdadeiro`. Text operands for remainder remain rejected with `E001` as a
project guard; their reference application failures are not language diagnostics.

Integer division by zero and real-to-integer remainder conversions outside the
supported finite conversion range receive positioned `R002` project guards.
The pass-through rule still applies to `7.5 \ 0`, which produces integer zero.

Power is left-associative: `2^3^2` means `(2^3)^2` and produces 64.
Unary signs bind more tightly than power: `-2^2` produces 4. Use `-(2^2)`
for -4. Unary `+` preserves a numeric operand and its type. Numeric power returns
real values, including `0^0 = 1`, `2^(-3) = 0.125`, and `10^(-400) = 0`.
Negative bases with fractional exponents, zero with negative exponents, and
overflowing powers report positioned `R002`. Other nonfinite arithmetic results
also receive positioned `R002` as a project guard.

Power with a text or logical operand, and unary minus applied to text or a
logical value, produce no value. In output statements this uses the same
recorded no-value handling as `numpcarac`: the current statement emits nothing.
Additional nonnumeric arithmetic combinations remain pending conformance work.
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

`escolha` evaluates its selector once. Numeric selectors are truncated toward
zero before matching; labels and range bounds retain their numeric values.
Thus `escolha 2.5` matches `caso 2` and `caso 2 ate 2.4`, but not `caso 2.5`.
Recorded selectors outside the signed 32-bit range do not match numeric labels.
Numeric ranges include both endpoints, using `caso lower ate upper` or the
accented `até` spelling. Bounds can be expressions, variables, or function
calls. They run from left to right even when a numeric lower bound already
excludes the selector. Descending ranges do not match.

Comma-separated labels may mix single values and ranges. The first match
executes its body and stops evaluating labels; duplicate or overlapping labels
are accepted. There is no fall-through. `outrocaso` runs only if nothing matches.
For the recorded ASCII text cases, single labels are uppercased and compared
with the unchanged selector: `"B"` matches `caso "b"`, while `"b"` does not.
Logical single labels use equality. Recorded text and logical ranges do not
match. Incompatible label types produce a positioned `E001`.

The range separator must follow its lower expression on the same physical
line. A missing upper expression or the unsupported `1..5` spelling produces
`P001`. A lower bound that produces no value does not match and skips its upper
bound; an evaluated upper bound without a value produces `P001`. A reference
case that skips malformed syntax in an unselected body, and a no-value selector
combined with a function bound, remain pending control-flow conformance cases.
See the [choice-range example](../examples/choice_ranges.alg).

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
Without a positive width, each number has one leading space and a decimal
count is ignored. Reals use 15 significant digits, omit unnecessary fractional
zeros, and discard the sign of zero. Scientific notation uses uppercase `E`
without a plus sign or exponent-leading zeros; for example, `0.00001` prints
as ` 1E-5` and `1000000000000000.0` as ` 1E15`. With a positive
width, numbers are right-aligned, the decimal count defaults to zero, and
decimal ties round away from zero. Integer values can also request decimals;
numeric fields expand if needed. Positive widths are capped at 255 characters
before padding or size checks, including widths larger than the project byte
budget. Strings are left-aligned and truncated to that width. Very large
decimal counts remain pending precision-conformance work; the width cap does
not bound every formatted numeric result. Logical output is ` VERDADEIRO` or
` FALSO`, including its leading space; logical field widths are rejected.

## Display commands

`limpatela` requests one screen clear. It preserves the program's captured text
output. `mudacor(color, target)` selects a foreground (`"frente"`) or background
(`"fundos"`, plural) color from character expressions, evaluated left to right.
Color and target names are case-insensitive. An unknown target, including the
singular `"fundo"`, makes no display change.

The accepted colors are `preto` (`#000000`), `azul` (`#0000FF`), `verde`
(`#008000`), `vermelho` (`#FF0000`), `roxo` (`#800080`), `amarelo` (`#FFFF00`),
and `branco` (`#FFFFFF`). Unknown colors leave the display unchanged. Leading
or trailing spaces are significant in both arguments.

`limpatela` ignores the remaining tokens on its physical line, including
parentheses and apparent arguments. `mudacor` requires an opening parenthesis
and two character expressions separated by a comma on the same physical line;
it ignores the rest of that line after the second expression, including extra
arguments or a missing closing parenthesis. Ignored expressions are not evaluated.
Missing values and malformed forms produce the recorded `P001` or `E001`.

In an expression, either display keyword produces no value and has no host
effect; apparent call arguments are ignored. This discards the containing write
as with other no-value expressions. Canonical formatting removes the ignored
syntax. The remaining mixed-type expression and diagnostic-timing differences
are listed in the [display recording report](display-commands-progress.md).

Both commands use the injected host. The default headless host makes their UI
effects silent, nonblocking no-ops without emitting terminal escape bytes.
Injected failures return `R008` at the command, retain preceding output, stop
later statements, and retain the underlying error without rendering its details.
The [display example](../examples/display.alg) demonstrates portable output.

## Built-ins

Numeric built-ins include `abs`, `arccos`, `arcsen`, `arctan`, `cos`, `cotan`,
`exp`, `grauprad`, `int`, `log`, `logn`, `pi`, `quad`, `radpgrau`, `raizq`,
`sen`, `tan`, `randi`, and the current legacy `aleatorio` forms. The reference
rejects `frac` as an undeclared name (`E002`).

`arccos`, `arcsen`, `arctan`, `cotan`, `grauprad`, and `radpgrau` return real
values. Trigonometric results and arguments use radians; `grauprad` converts
degrees to radians and `radpgrau` converts radians to degrees. Their empty
parentheses supply zero. `quad(x)` squares its argument, retaining integer or
real type and signed 32-bit wrapping for integers. `quad()` and text or logical
arguments to `quad` produce no value.

Angle conversion keeps finite results for large inputs: `grauprad(1e308)`
prints `1.74532925199433E306` and `radpgrau(1e306)` prints
`5.72957795130823E307`. A nonfinite `radpgrau` result produces no value.
Overflow in real `quad` receives positioned `R007` as a project guard for
the reference's unpositioned application fault.

These seven functions accept at most one numeric argument. Extra arguments
receive `P001`, except that a statically absent argument to the six real functions
propagates no value before later arguments are considered. The six real
functions reject text and logical arguments with `P001`. Out-of-domain `arccos`
and `arcsen`, and
zero-argument or zero-valued `cotan`, produce no value at runtime; output handles
that result as described for `numpcarac` below. This runtime absence is allowed
even when the expression's static result type is real.

`pi` is written without parentheses and uses the recorded real precision:
default output is `3.14159265358979`. `pi()` receives `P001`; the numeric
functions require parentheses.

`exp(base, exponent)` takes two numeric arguments and returns real-valued
exponentiation. A lone numeric argument or extra numeric arguments receive
`P001`. `exp()` produces no value. Arguments are evaluated left to right.
A text, logical, or absent base ends the call without evaluating later arguments.
Text, logical, and generic no-value exponents also stop immediately.

Absence caused by a numeric domain has a distinct rule: unary numeric calls
preserve its origin, whereas `exp` and `numpcarac` turn it into generic absence.
When the exponent has numeric-domain absence, `exp` evaluates at most one
extra argument, then produces no value. Thus `exp(2, arccos(2), f(), g())`
calls only `f`; `exp(2, exp(arccos(2)), f())` skips `f`.
An error in that evaluated extra argument is preserved at its source position.

Arity that depends on an absent numeric result is checked during execution:
`exp(arccos(2))` produces no value, while `exp(arccos(0))` receives `P001`.
Other invalid numeric domains and nonfinite results receive positioned `R007`.
Recorded exponentiation includes negative bases with integral exponents,
`exp(0, 0) = 1`, and underflow of `exp(10, -400)` to zero.

`log(x)` is base ten and `logn(x)` is the natural logarithm; both take one
numeric argument. Empty parentheses supply zero, which receives `R007`, as
do negative inputs. `raizq`, `sen`, `cos`, and `tan` also supply zero for empty
parentheses. Text and logical inputs to `raizq` produce no value; those inputs
to logarithms or trigonometric functions receive `P001`. A negative square-root
input receives positioned `R007` as a project guard for the reference's
unpositioned application fault.

For integer arguments, `abs` retains signed 32-bit wrapping, so the absolute
value of `-2147483647 - 1` remains `-2147483648`. Real arguments use real
absolute value. `abs()` produces no value. `int()` returns integer zero.
Text, logical, or absent arguments to either function produce no value.
`abs` rejects extra arguments before conversion; `int` stops on a text,
logical, or absent first argument before considering later arguments.
`int(x)` truncates toward zero and narrows the result to
signed 32-bit: `int(2147483648.0)` is `-2147483648`, and
`int(4294967296.0)` is zero. Non-finite arguments and values outside the signed
64-bit intermediate conversion range receive positioned `R007` as project guards.

An absent numeric argument propagates through a containing numeric call,
subject to the `exp` argument rules above. Output handles the result as
described for `numpcarac`.
The full numeric domain matrix and catalog consolidation remain pending.

String built-ins: `copia`, `maiusc`, `minusc`, `asc`, `carac`, `compr`, and
`pos`. Their recorded character repertoire is Windows-1252, decoded to UTF-8
inside the implementation. String positions count characters and are 1-indexed.
UTF-8 source takes priority when bytes are valid in both supported encodings.

String values keep their first 255 characters. The limit applies to literals,
each concatenation, constants, assignments, parameters, function results, and
vector elements. Thus a character appended after position 255 is unavailable
to search, slicing, comparison, or numeric conversion. Accented Windows-1252
characters count once after decoding. Unicode extensions use the same limit
without splitting UTF-8 encodings. Formatting retains the complete source
literal so reparsing preserves its runtime behavior.

Character input consumes and echoes the complete entered line, then stores
its first 255 characters. The line-size safeguard still applies before
conversion. See the [string-limit example](../examples/string_limits.alg).

`copia(s, start, count)` accepts integer or real bounds, truncates real fractions
and narrows them to signed 32-bit values, substitutes zero for no-value bounds,
clamps starts below one to one, and clips the length at the string's end.
Nonpositive lengths and starts past the end return an empty string. `compr(s)`
counts characters, including accented letters and the euro sign. `pos(needle, s)`
returns the first matching position, or zero for missing or empty search text.
Search is case-sensitive. `maiusc` and `minusc` use the recorded accented case
pairs; uppercase preserves `µ` and `ƒ` because their Unicode uppercase forms
are outside Windows-1252. Case conversion and copying preserve empty strings.
Non-finite real bounds and real bounds outside the signed 64-bit conversion
range receive positioned `R007` as a project guard; unpositioned reference
application faults are not treated as language diagnostics.

These text calls require parentheses and their complete argument lists.
Character arguments of the wrong type or with no value receive `E001`.
`copia` accepts numeric or no-value bounds, reporting `E001` for other types.
Missing character arguments receive `E001`; extra arguments and a missing
final `copia` bound receive `P001`. Argument expressions evaluate left to right.

`asc(s)` returns the first character's Windows-1252 byte, not its Unicode
code point: `asc("€")` is 128. An empty string produces no value. As a project
guard, a first character outside Windows-1252 receives positioned `R007`.

`carac(n)` accepts an integer code in `[0, 255]` and uses the complete recorded
character table, including its substitutions for drawing characters.
Codes 0–31, 127, and 255 produce a space; `carac()` also produces a space.
Out-of-domain integers produce no value. A no-value argument is treated as zero.
Real, character, and logical arguments receive `E001`; extra arguments receive
`P001`.
The two code functions are not inverses: `carac(128)` is `"Ç"`, whose `asc`
value is 199. No-value results in output use the discard-and-continue behavior
described below. Ordinary Unicode casing outside the reference repertoire
remains a project extension, subject to the existing text-size limit.

`numpcarac(x)` converts a numeric value to character text with 15 significant
digits, no leading space, and uppercase `E` exponents without a plus sign or
leading exponent zeros. It discards the sign of zero. `numpcarac()` returns
`"0"`; parentheses are required and at most one argument is accepted.
Omitted parentheses, extra arguments, and an unindexed vector argument receive
`P001`.

String and logical arguments produce no value, and nested `numpcarac` calls
preserve that absence. In an output item this discards the entire statement's
buffered text, skips its remaining items and format expressions, and leaves
any pending newline for the next successful output statement. Execution then
continues. A typed return of this absent value receives `E001`; assignment is
also rejected with `E001` as a project guard. Non-finite numeric arguments receive
positioned `R007` as a project guard.

`caracpnum(s)` requires one character argument and selects an integer or real
result at runtime. Analysis keeps this numeric choice unresolved. Numeric
expressions and calls preserve the resulting type; integer-only consumers
check the actual value during execution. A real result used by `carac`,
`randi`, a vector index, or an integer return receives `E001`; real loop or
format bounds receive `P001`. Incompatible assignments receive `R001` after
evaluating the conversion. Integer parameters retain their documented real
argument conversion. A constant conversion can supply an integer vector
bound; a real result receives `P001` when that declaration is reached.

ASCII spaces around the text are ignored; tabs and nonbreaking spaces are
not. Empty text, an absent numeric prefix, leading `.` or `,`, and unsupported
alphabetic spellings such as `NaN` produce integer zero. Decimal points,
commas, or `e`/`E` exponent syntax select real type, including `"1.0"`,
`"1,0"`, and `"1e0"`. An exponent with no digits (`"1e"`, `"1e+"`, or
`"1e-"`) contributes zero. Real underflow produces zero; overflow and
malformed text after a numeric prefix receive `R007` at the conversion call.
Already emitted output is preserved if conversion fails.

Integer spellings from zero through 2147483647 produce integer values;
2147483648 and 2147483649 produce integer zero, and positive integer
spellings from 2147483650 select real type. A leading plus is accepted.
Negative integer spellings through -2147483648 produce zero; more negative
integer spellings select real type. Signed decimal or exponent forms preserve
their sign. Hexadecimal prefixes `$`, `0x`, and `0X` are accepted, including
a leading plus: values through `7FFFFFFF` produce integers, negative and larger
32-bit values produce zero, and overflow beyond 32 bits receives `R007`. An
empty or invalid first hexadecimal digit produces zero. Digit separators such
as underscores are rejected after a numeric prefix.

Omitted parentheses or extra arguments receive `P001`; an empty call,
non-character argument, or no-value argument receives `E001`. Conversion uses
the same bounded string value as other consumers: 400 literal nines become
255 nines before conversion and produce `1E255`. The boundary cases retained
by the conversion slice are now verified; see `docs/string-limits-progress.md`.

`randi(n)` takes one signed 32-bit integer bound and returns an integer. For
positive `n`, its domain is `[0, n)`. Negative bounds use their unsigned
32-bit representation as the exclusive width, then interpret the result as
signed 32-bit; for example, `randi(-7)` produces values in
`[-2147483648, -8]` or `[0, 2147483647]`. Both `randi()` and `randi(0)`
return zero. Arguments are evaluated once before the draw. Parentheses are
required; omitted parentheses, extra arguments, and unindexed vectors receive
`P001`. Real, character, logical, and no-value arguments receive `E001`.

The injected source is per interpreter. Empty and zero-bound calls do not
consume a draw in this implementation; this does not promise reference seed
or draw-count compatibility. Out-of-range injected results and runtime bounds
outside signed 32-bit receive positioned `R007` as project guards. Large
whole-number literals have real type and receive `E001`. `rand` callable forms and
command-form random input are separate pending work.

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
256 expression-evaluation levels within one call frame, and a 16 MiB
allocation safeguard for incoming text, an input token/line, or a formatted
item. Language string values have the separate 255-character reference limit.
Case conversion retains its allocation guard for direct library callers;
format widths are capped before allocation, and precision expansion still
checks its requested size before creating a result. Each vector
aggregate is capped at 1,048,576 scalar slots, including nested elements, with
checked dimension products. Oversized literal layouts receive `E900` during
analysis; constant-dependent layouts receive `R003` at declaration initialization,
before allocation. Reference-specific storage quotas remain pending.
These are project safeguards, not measured VisuAlg limits. A safeguard
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
commands expose typed foreground/background color changes and screen clears;
other language host commands remain pending reference qualification.

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
