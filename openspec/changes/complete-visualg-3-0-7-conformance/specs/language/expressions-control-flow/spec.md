## Purpose

Defines oracle-compatible expression evaluation, comparisons, selection, loops, interruption, return propagation, statement validation, and arithmetic safety.

## ADDED Requirements

### Requirement: Operator vocabulary and precedence
The language SHALL implement the complete VisuAlg 3.0.7 unary, arithmetic, relational, and logical operator vocabulary with the associativity and precedence established by reduced oracle probes. Textual operators and their oracle-confirmed symbolic aliases SHALL be case-insensitive where applicable, and parentheses SHALL override precedence.

#### Scenario: Evaluate an unparenthesized mixed expression
- **WHEN** an expression combines operator classes whose precedence or associativity could change the result
- **THEN** its syntax tree and value match the committed precedence probe for VisuAlg 3.0.7

#### Scenario: Override precedence explicitly
- **WHEN** parentheses form a valid grouping
- **THEN** the grouped expression is evaluated first and produces the reference result

#### Scenario: Replay recorded power precedence
- **WHEN** the recorded expressions `2^3^2` and `-2^2` are evaluated
- **THEN** they produce real values 64 and 4 respectively, reflecting left-associative power and tighter unary minus

### Requirement: Numeric evaluation and coercion
Integer addition, subtraction, multiplication, and negation SHALL preserve the recorded signed 32-bit wrapping. Addition, subtraction, and multiplication with real operands SHALL retain real arithmetic and real result types.

The `\` operator SHALL truncate division toward zero for two integer operands and SHALL return the right operand with its original type for other scalar pairs. Numeric remainder SHALL truncate and narrow the dividend as `int` does, compute the remainder for a positive integer divisor, and return integer `-1` for zero, negative, or real divisors. Numeric or logical remainder pairs containing a logical operand SHALL return the right operand unchanged. Both operands SHALL be evaluated exactly once, left to right.

The `/` operator SHALL return real division for numeric pairs and the right operand unchanged for other scalar pairs. Numeric power SHALL return real values, including the recorded zero-power and underflow results, and SHALL report invalid power domains at the operator as `R002`. Other nonfinite arithmetic results SHALL receive positioned project `R002` guards. Power involving text or logical scalars, and unary minus of text or logical values, SHALL produce no value. Numeric unary plus SHALL preserve its operand value and type.

Arithmetic SHALL use the reference operand compatibility, promotion, result type, rounding, truncation, division, modulo, exponentiation, unary sign, and overflow behavior. Statically invalid combinations SHALL be semantic diagnostics; runtime-only failures SHALL return positioned arithmetic diagnostics and SHALL never surface a Go panic, infinity, NaN, or wraparound unless the oracle explicitly produces the corresponding observable value.

#### Scenario: Mix integer and real operands
- **WHEN** a numeric operator receives an integer and a real in a combination accepted by the reference
- **THEN** the value and result type match the oracle-confirmed coercion rule

#### Scenario: Divide two integers by zero
- **WHEN** both operands of `\` are integers and the denominator is zero
- **THEN** execution returns a positioned `R002` project diagnostic without reproducing the reference application's unpositioned failure

#### Scenario: Use a zero remainder divisor
- **WHEN** numeric operands use `%` or `MOD` with a zero divisor
- **THEN** the result is integer `-1`

#### Scenario: Preserve the right operand of mixed division
- **WHEN** the recorded expressions `8.5 \ 3`, `7 \ 2.5`, and `7.5 \ 0` are evaluated
- **THEN** they produce integer `3`, real `2.5`, and integer `0`, respectively, after evaluating both operands

#### Scenario: Apply a unary sign to an integer retained by mixed real division
- **WHEN** the recorded mixed `/` expression retains an integer operand
- **THEN** unary plus reports `P001`, while unary minus negates the real value whose low bits hold the unsigned 32-bit integer payload

#### Scenario: Reject invalid power domains
- **WHEN** a negative base has a fractional exponent, zero has a negative exponent, or a numeric power overflows
- **THEN** execution reports `R002` at the power operator and does not emit the incomplete output statement

#### Scenario: Observe a no-value arithmetic result
- **WHEN** the recorded nonnumeric power or unary-minus expressions occur in an output statement
- **THEN** the expression produces no value and that statement emits nothing

#### Scenario: Overflow an integer operation
- **WHEN** an integer operation exceeds the reference integer domain
- **THEN** the result or positioned diagnostic matches the oracle evidence and loop or allocation control cannot wrap silently

### Requirement: Logical evaluation
Logical operators SHALL require the oracle-confirmed operand types and SHALL reproduce the reference evaluation order and short-circuit or eager behavior, including operand side effects and failures. Unary logical negation and exclusive-or SHALL return the reference boolean values.

The recorded mixed logical and comparison reductions SHALL preserve intermediate operands that affect enclosing arithmetic, including across numeric builtin and language-function evaluation. `2 + (3 e 7)` SHALL produce `10`, `20 - (3 e 7)` SHALL produce `-4`, and `2 * (3 e 7)` SHALL produce `21`. Ordinary `2 + (3 + 7)` SHALL remain `12`.

#### Scenario: Retain operands through nested reductions
- **WHEN** the recorded expression `2 + ((3 e 7) + (4 e 9))` is evaluated
- **THEN** it produces integer `20`, and evaluating the formatted program preserves that result

#### Scenario: Reject incompatible retained addition operands
- **WHEN** `2 + (3 e ("x" = 7))` is evaluated
- **THEN** the retained incompatible operand receives positioned `P001`

#### Scenario: Observe right-operand evaluation
- **WHEN** the right operand of `e` or `ou` has an observable side effect or runtime failure
- **THEN** whether and when that operand executes matches the committed VisuAlg 3.0.7 probe

### Requirement: Boolean, numeric, and string comparison
Relational operators SHALL accept exactly the type pairs accepted by VisuAlg 3.0.7 and SHALL reproduce its numeric coercion, boolean ordering or equality, string comparison, case sensitivity, and locale behavior. Unsupported uses SHALL be diagnosed at the recorded operand or operator position and phase.

Numeric pairs, logical pairs, and character pairs SHALL return logical comparison results. Logical ordering SHALL place `falso` before `verdadeiro`. Character ordering SHALL compare Windows-1252 byte values, retaining case, accents, and prefix length. Other scalar pairs SHALL retain the right operand's concrete value for every relational operator. Comparisons SHALL preserve their logical assignment and condition category separately from that value. Ordinary operands SHALL be evaluated exactly once, left to right; no-value handling SHALL match the recorded generic and numeric-domain cases.

#### Scenario: Compare character values
- **WHEN** two `caractere` expressions are compared with a reference-supported relational operator
- **THEN** the result matches the oracle for accents, case, prefix length, and differing characters

#### Scenario: Reject an unsupported comparison
- **WHEN** operands or an operator combination are not comparable in the reference
- **THEN** the recorded positioned diagnostic is emitted, preserving prior output and the recorded operand evaluation order

#### Scenario: Preserve a mixed comparison result
- **WHEN** `1 = "1"` or `"1" >= 1` is evaluated
- **THEN** the result is text `"1"` or integer `1`, respectively, after both operands are evaluated

#### Scenario: Assign a comparison result
- **WHEN** a comparison retaining integer `7` is assigned to a logical variable
- **THEN** output retains `7` and a subsequent bare condition rejects its nonlogical value
- **AND** assigning the comparison to an integer destination reports `R001` during execution

#### Scenario: Retain a value in logical storage
- **WHEN** a logical variable, vector element, record field, or reference parameter receives a comparison retaining a character value
- **THEN** later reads and assignments use that concrete character type, including after parameter copy-back
- **AND** retaining a record gives an empty layout without the original named type or fields

#### Scenario: Consume a comparison as a function result
- **WHEN** a logical function returns `"x" = 7`
- **THEN** its result is `falso`, while an integer function returning the same expression receives positioned `E001`

#### Scenario: Apply a logical operator to a mixed comparison value
- **WHEN** the recorded `e`, `ou`, or `xou` expressions consume a mixed comparison result
- **THEN** both operands are evaluated, `e` retains the right value, and `ou` or `xou` reports positioned `P001`

### Requirement: Deterministic expression order
Subexpressions, designator indices, statement expressions, and format expressions SHALL be evaluated once and in the exact left-to-right or oracle-recorded order. A failure SHALL stop further evaluation only where the reference stops it.

#### Scenario: Observe expression side effects
- **WHEN** sibling expressions call subprograms that each modify visible state
- **THEN** the final state and output match the reference evaluation order and no expression is evaluated twice

### Requirement: Choice matching and ranges
`escolha` SHALL evaluate its selector once and compare it against oracle-confirmed single labels, comma-separated labels, and inclusive ranged labels using the reference type and coercion rules. It SHALL select the first reference-matching arm, SHALL not fall through, and SHALL execute `outrocaso` only when no ordinary arm matches.

#### Scenario: Match a ranged case label
- **WHEN** a selector falls within an accepted inclusive case range
- **THEN** exactly that case body executes and control continues after `fimescolha`

#### Scenario: Truncate a numeric choice selector
- **WHEN** the recorded numeric selector has a fractional part
- **THEN** it is truncated toward zero once before comparisons, while single labels and range bounds retain their numeric values

#### Scenario: Match no case label
- **WHEN** no ordinary case value or range matches and `outrocaso` is present
- **THEN** exactly the `outrocaso` body executes

#### Scenario: Match dynamic and overlapping labels
- **WHEN** numeric label bounds are expressions, variables, or function calls, and accepted ranges overlap or repeat a prior value
- **THEN** bounds are evaluated in source order and the first matching label selects its body without evaluating later labels

#### Scenario: Reject incompatible label types
- **WHEN** a case value or range bound has a type incompatible with the selector in a way rejected by the reference
- **THEN** the implementation emits the traced positioned diagnostic

#### Scenario: Match ordinary text labels
- **WHEN** the recorded ASCII text selector is compared with a single text label
- **THEN** the label is uppercased and compared with the unchanged selector, including labels computed by expressions

### Requirement: Loop semantics
`enquanto`, `repita`, and `para` SHALL implement the complete oracle-confirmed VisuAlg 3.0.7 forms, including `repita ... ate`, any accepted `repita ... fimrepita` infinite form, default and explicit `para` steps, negative steps, bound and step evaluation timing, and loop-variable mutation behavior.

#### Scenario: Run a descending for loop
- **WHEN** a `para` loop uses a valid negative step and descending bounds
- **THEN** it visits exactly the reference sequence of loop-variable values and terminates

#### Scenario: Reject or stop a zero step
- **WHEN** a `para` step evaluates to zero
- **THEN** the reference-compatible diagnostic or zero-iteration behavior occurs without an infinite host process

#### Scenario: Approach an integer boundary
- **WHEN** a `para` iteration reaches the maximum or minimum integer and another increment would overflow
- **THEN** the loop terminates or reports the reference-compatible diagnostic without wrapping

#### Scenario: Replay recorded loop exit state
- **WHEN** the finite integer loops in the 2026-09-07 Windows observations complete normally, execute zero iterations, or reach `interrompa`
- **THEN** normal nonempty completion exposes the smaller of the next iteration value and terminal bound, an empty loop exposes its initial bound, and interruption exposes the smaller of the body value and terminal bound, including the recorded descending cases

### Requirement: Interruption and return propagation
`interrompa` and `retorne` SHALL be valid only in the oracle-confirmed contexts.
`interrompa` SHALL propagate through nested blocks to the correct loop boundary.
`retorne` SHALL update the active function's result without leaving a conditional,
choice, loop, or function body. A misplaced command SHALL be rejected before
execution when statically knowable.

#### Scenario: Interrupt the innermost loop
- **WHEN** `interrompa` executes inside nested loops
- **THEN** only the oracle-designated loop ends and execution resumes at the correct following statement

#### Scenario: Set a result inside nested control flow
- **WHEN** `retorne` executes in a valid nested branch within a subprogram
- **THEN** the active function's result changes and execution continues through the remaining statements and loop iterations until an actual loop or function boundary is reached

### Requirement: I/O statement validation
`leia`, `escreva`, and `escreval` syntax, arity, designator requirements, format fields, and expression types SHALL be validated according to VisuAlg 3.0.7 before execution wherever possible. Invalid read destinations and invalid width or precision forms SHALL not be deferred into unpositioned runtime failures.

#### Scenario: Read into a non-designator
- **WHEN** `leia` receives a literal, call result, or other expression that the reference does not permit as a destination
- **THEN** analysis emits a positioned diagnostic at that argument

#### Scenario: Use an invalid output format
- **WHEN** an output item uses a width or precision count rejected by the reference
- **THEN** syntax or semantic analysis reports the traced positioned diagnostic
