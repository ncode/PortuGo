## Purpose

Defines VisuAlg 3.0.7 procedure and function declarations, calls, parameters, scope, recursion, argument order, and return behavior.

## ADDED Requirements

### Requirement: Subprogram declaration forms
The language SHALL accept every oracle-confirmed procedure and function declaration form, including declarations with no parameters, optional empty parentheses where accepted, parameter groups, local declaration sections, and reference-compatible terminators. Forms rejected by the reference SHALL receive positioned diagnostics.

#### Scenario: Declare a parameterless procedure
- **WHEN** a parameterless procedure omits or includes parentheses in a form accepted by VisuAlg 3.0.7
- **THEN** it declares the same callable symbol and its body is analyzed normally

#### Scenario: Reject a non-reference declaration form
- **WHEN** a procedure or function uses punctuation or section placement rejected by the oracle
- **THEN** parsing reports the traced positioned diagnostic and recovers at the next structural boundary

### Requirement: Call syntax and context
Procedures and functions SHALL support the oracle-confirmed parenthesized and bare call forms. A call form SHALL be accepted only in the statement or expression contexts where VisuAlg 3.0.7 accepts that callable kind, and ambiguity with a variable or designator SHALL be resolved according to reference evidence.

Recorded function names SHALL take priority over colliding local variables or
parameters in bare value expressions and parenthesized function calls.

#### Scenario: Invoke a bare procedure call
- **WHEN** a declared procedure is invoked without parentheses in a reference-accepted statement form
- **THEN** it executes once with the same behavior as the corresponding accepted parenthesized form

#### Scenario: Use a procedure as a value
- **WHEN** a procedure call appears where a value is required and the reference rejects it
- **THEN** analysis emits a positioned call diagnostic

#### Scenario: Read a colliding function name
- **WHEN** a function value is read where a local variable or parameter has the same case-insensitive name
- **THEN** the function executes and supplies its result, as in the recorded name-priority probes

### Requirement: Value parameters
Value parameters SHALL accept exactly the argument types and implicit conversions accepted by VisuAlg 3.0.7. Each argument expression SHALL be evaluated once in reference order before the body observes its parameter, and later mutation of a value parameter SHALL follow the scalar or aggregate copy rules.

#### Scenario: Evaluate multiple value arguments
- **WHEN** multiple argument expressions have visible side effects
- **THEN** each is evaluated once in the oracle-confirmed order before the subprogram body begins

### Requirement: Reference parameters
A `var` parameter SHALL follow the oracle-confirmed designator capture,
conversion, and copy-back rules. Recorded scalar parameters receive independent
copies, with real-to-integer conversion truncating toward zero. Successful
return SHALL copy parameter values and their numeric types back to the captured
caller locations in parameter order. Temporaries, constants, and call results
SHALL require explicit acceptance evidence before they become reference storage.

#### Scenario: Capture a vector element before later arguments
- **WHEN** a vector element is passed to a scalar `var` parameter and a later argument changes the index variable
- **THEN** copy-back updates the element selected before that later argument ran

#### Scenario: Convert a numeric reference argument
- **WHEN** a real variable containing 3.5 is passed to an integer `var` parameter
- **THEN** the parameter initially contains 3, the caller remains unchanged during the call, and successful return copies the integer value and type back

#### Scenario: Pass the same location twice
- **WHEN** two `var` parameters capture the same caller location
- **THEN** each parameter retains its own value during execution and the last parameter's copy-back determines the caller's final value

### Requirement: Lexical scope and bindings
Name lookup SHALL follow VisuAlg 3.0.7 lexical scope: parameters and locals SHALL shadow only names the reference allows them to shadow, global declarations SHALL remain visible in subprogram bodies according to declaration rules, and caller locals SHALL never become visible through dynamic scope. Bindings SHALL be case-insensitive and fixed by semantic analysis.

#### Scenario: Read a global from a subprogram
- **WHEN** a subprogram references an unshadowed global declared in a reference-visible position
- **THEN** the reference resolves to that global regardless of which caller invoked the subprogram

#### Scenario: Prevent dynamic-scope capture
- **WHEN** a caller has a local name absent from the callee's lexical environment
- **THEN** the callee does not resolve that caller-local name

### Requirement: Declaration visibility and recursion
Subprograms SHALL be visible before, after, or only following their declarations exactly as established by the VisuAlg 3.0.7 oracle. Direct recursion, mutual recursion, and calls between top-level declarations SHALL be accepted or diagnosed according to those visibility rules without depending on Go declaration order.

#### Scenario: Execute direct recursion
- **WHEN** a valid function calls itself with a terminating input
- **THEN** each call receives an independent frame and the final value matches the reference

#### Scenario: Analyze an unavailable call target
- **WHEN** a call refers to a declaration not visible under the reference ordering rules
- **THEN** analysis returns the positioned undeclared-call diagnostic

### Requirement: Function return behavior
Functions SHALL enforce the oracle-confirmed return syntax, optional or required return expressions, result coercion, and control-flow completeness. Procedures SHALL accept bare or valued returns only where the reference does. A call that reaches its end without a required result SHALL be rejected statically when provable or produce the positioned call/return diagnostic at runtime.

#### Scenario: Return a compatible value
- **WHEN** a function executes a valid return expression compatible with its declared result type
- **THEN** the caller receives the reference-compatible value and statements after the return do not execute

#### Scenario: Reach the end without a required value
- **WHEN** a function path reaches its terminator without the result required by VisuAlg 3.0.7
- **THEN** analysis or execution emits the stable positioned return diagnostic rather than inventing a value

### Requirement: Call-frame safety
Each invocation SHALL have independent parameter and local storage while sharing only globals and explicit aliases allowed by the language. Runtime call setup, recursion, return unwinding, and invalid call state SHALL produce language diagnostics instead of leaking frames or panicking.

#### Scenario: Re-enter a recursive local variable
- **WHEN** recursive calls assign different values to the same local declaration
- **THEN** each active call observes its own local value and returning restores the caller's frame
