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
Arithmetic SHALL use the reference operand compatibility, promotion, result type, rounding, truncation, division, modulo, exponentiation, unary sign, and overflow behavior. Statically invalid combinations SHALL be semantic diagnostics; runtime-only failures SHALL return positioned arithmetic diagnostics and SHALL never surface a Go panic, infinity, NaN, or wraparound unless the oracle explicitly produces the corresponding observable value.

#### Scenario: Mix integer and real operands
- **WHEN** a numeric operator receives an integer and a real in a combination accepted by the reference
- **THEN** the value and result type match the oracle-confirmed coercion rule

#### Scenario: Divide by zero
- **WHEN** a division or modulo denominator evaluates to zero
- **THEN** execution returns the reference-compatible positioned arithmetic diagnostic without panicking

#### Scenario: Overflow an integer operation
- **WHEN** an integer operation exceeds the reference integer domain
- **THEN** the result or positioned diagnostic matches the oracle evidence and loop or allocation control cannot wrap silently

### Requirement: Logical evaluation
Logical operators SHALL require the oracle-confirmed operand types and SHALL reproduce the reference evaluation order and short-circuit or eager behavior, including operand side effects and failures. Unary logical negation and exclusive-or SHALL return the reference boolean values.

#### Scenario: Observe right-operand evaluation
- **WHEN** the right operand of `e` or `ou` has an observable side effect or runtime failure
- **THEN** whether and when that operand executes matches the committed VisuAlg 3.0.7 probe

### Requirement: Boolean, numeric, and string comparison
Relational operators SHALL accept exactly the type pairs accepted by VisuAlg 3.0.7 and SHALL reproduce its numeric coercion, boolean ordering or equality, string comparison, case sensitivity, and locale behavior. Unsupported comparisons SHALL be diagnosed at the operator.

#### Scenario: Compare character values
- **WHEN** two `caractere` expressions are compared with a reference-supported relational operator
- **THEN** the result matches the oracle for accents, case, prefix length, and differing characters

#### Scenario: Reject an unsupported comparison
- **WHEN** operands or an operator combination are not comparable in the reference
- **THEN** analysis emits a positioned type diagnostic and execution does not evaluate the invalid comparison

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

#### Scenario: Match no case label
- **WHEN** no ordinary case value or range matches and `outrocaso` is present
- **THEN** exactly the `outrocaso` body executes

#### Scenario: Reject overlapping or invalid labels
- **WHEN** case labels are duplicated, overlap, are non-constant, or use incompatible types in a way rejected by the reference
- **THEN** the implementation emits the traced positioned diagnostic

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
