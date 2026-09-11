## Purpose

Defines one complete VisuAlg 3.0.7 built-in catalog with consistent semantic signatures, runtime behavior, conversions, failures, and random domains.

## ADDED Requirements

### Requirement: Authoritative built-in catalog
The project SHALL maintain one authoritative catalog of every VisuAlg 3.0.7 built-in name, alias, callable form, arity, parameter mode and type, result type, domain restriction, and runtime operation. Semantic call validation and runtime dispatch SHALL agree with that catalog. Completeness SHALL also be compared against the independently recorded reference inventory, and validation SHALL fail if a required function is absent from both implementation layers or either layer contradicts an entry.

#### Scenario: Validate the catalog against both stages
- **WHEN** the built-in catalog conformance test enumerates all descriptors
- **THEN** every descriptor has matching semantic acceptance and runtime behavior and no extra semantic-only or runtime-only built-in exists

#### Scenario: Detect a function omitted from both layers
- **WHEN** an accepted reference function has no implementation catalog entry, semantic signature, or runtime evaluator
- **THEN** the independent inventory comparison fails and names that missing function

#### Scenario: Resolve a built-in alias
- **WHEN** a program calls an oracle-confirmed alias using any accepted letter case
- **THEN** it receives the same signature and behavior as the catalog's canonical built-in

### Requirement: Numeric built-ins
The recorded `arccos`, `arcsen`, `arctan`, `cotan`, `grauprad`, and `radpgrau` functions SHALL accept zero or one numeric argument, use zero when omitted, return real values, and preserve the recorded trigonometric and angle-conversion units. Text and logical arguments and extra numeric arguments SHALL receive `P001`; an already absent argument SHALL propagate no value before later arguments are considered. Out-of-domain inverse cosine/sine and cotangent of zero SHALL produce no value. Runtime no-value results SHALL be permitted even for statically real expressions.

`quad` SHALL preserve numeric argument type, square integers with signed 32-bit wrapping, and produce no value for absent, text, or logical arguments. Real squaring overflow SHALL receive the positioned `R007` project guard for the reference's unpositioned fault. Angle conversion SHALL preserve finite large results without premature intermediate overflow; a nonfinite `radpgrau` result SHALL produce numeric-domain absence. The constant `pi` SHALL accept its bare form and reject parentheses with `P001`, preserving the recorded default real output.

For integer input, `abs` SHALL preserve the recorded signed 32-bit result, including the signed minimum. `int` SHALL truncate numeric input toward zero and narrow the result to signed 32-bit; unsupported non-finite or intermediate conversion values SHALL receive the positioned built-in failure guard.

`log` SHALL use base ten and `logn` SHALL use the natural logarithm with one numeric argument. Empty logarithm calls and nonpositive inputs SHALL receive positioned `R007`. Empty `raizq`, `sen`, `cos`, and `tan` calls SHALL supply zero. Empty `abs` and `exp` calls SHALL produce no value, while `int()` SHALL produce integer zero. The rejected legacy name `frac` SHALL receive `E002`.

`exp` SHALL evaluate numeric arguments left to right and return real exponentiation. A text, logical, or absent base SHALL produce no value and skip later arguments. A text, logical, or generically absent exponent SHALL also stop immediately. An exponent with numeric-domain absence SHALL evaluate at most one extra expression before returning generic absence, preserving any failure in that expression. Unary numeric calls SHALL preserve numeric-domain absence; `exp` and `numpcarac` SHALL clear its numeric origin. Invalid numeric arity SHALL receive `P001`, with domain-dependent arity deferred to execution; invalid numeric domains and nonfinite results SHALL receive positioned `R007`. `int` SHALL stop before later arguments on text, logical, or absent input. `abs`, `quad`, and `numpcarac` SHALL reject extra arguments before converting text input.

The catalog SHALL implement the full oracle-confirmed numeric set with the reference arity, accepted numeric types, result types, angle units, rounding or truncation, constants, domains, and exceptional behavior. The candidate inventory SHALL explicitly include `abs`, `arccos`, `arcsen`, `arctan`, `cos`, `cotan`, `exp`, `grauprad`, `int`, `log`, `logn`, `pi`, `quad`, `radpgrau`, `raizq`, `sen`, `tan`, and the existing `frac`, plus any additional discovered functions and aliases. Every candidate SHALL receive a recorded disposition; accepted names SHALL be implemented as distinct operations where their semantics differ, and rejected legacy names SHALL receive rejection coverage.

#### Scenario: Evaluate numeric boundary cases
- **WHEN** each numeric built-in receives zero, negative, positive, integral, fractional, and domain-boundary inputs applicable to its signature
- **THEN** its value, type, or positioned failure matches the corresponding VisuAlg 3.0.7 oracle rows

#### Scenario: Use trigonometric functions
- **WHEN** `sen`, `cos`, `tan`, `cotan`, an inverse trigonometric function, or an angle-conversion function receives an argument accepted by its recorded signature
- **THEN** the result uses the same angle unit and reference-compatible precision as VisuAlg 3.0.7

#### Scenario: Validate exponentiation arity
- **WHEN** the numeric catalog is qualified against the reference inventory
- **THEN** the recorded acceptance and result of `exp(base, expoente)` are checked explicitly, and a one-argument implementation cannot satisfy that probe

#### Scenario: Skip arguments after an absent numeric result
- **WHEN** `exp` encounters an absent base or a generically absent exponent
- **THEN** the call produces no value and later arguments do not execute their side effects or failures

#### Scenario: Preserve numeric-domain absence while evaluating an exponent
- **WHEN** the second argument of `exp` has numeric-domain absence and extra arguments are present
- **THEN** exactly the next argument is evaluated before the call produces generic absence, and any failure in that argument remains positioned

### Requirement: Text built-ins
The catalog SHALL implement the full oracle-confirmed text set, including `copia`, `maiusc`, `minusc`, `compr`, and `pos`, using the reference's 1-based positions, substring bounds, not-found value, case conversion, accented-character behavior, empty-string behavior, and result types.

#### Scenario: Keep bounded string values
- **WHEN** literals, concatenation, input, or a function result produce more than 255 characters
- **THEN** subsequent operations observe the first 255 characters without splitting decoded Windows-1252 characters
- **AND** input still consumes and echoes the complete entered line before storing the bounded value

#### Scenario: Extract a substring
- **WHEN** `copia` receives a valid character value, 1-based position, and length
- **THEN** it returns exactly the oracle-confirmed substring for ASCII and Windows-1252 text

#### Scenario: Search for absent text
- **WHEN** `pos` cannot find its search value
- **THEN** it returns zero, also returning zero for an empty search value

#### Scenario: Convert accented case and fractional substring bounds
- **WHEN** text contains Windows-1252 characters or `copia` receives real bounds
- **THEN** case conversion preserves `µ` and `ƒ` when uppercasing, real bounds truncate toward zero and narrow to signed 32-bit values, no-value bounds become zero, starts below one clamp to one, and nonpositive lengths or past-end starts yield empty text

### Requirement: Character-code built-ins
Character-code operations, including oracle-confirmed `asc`, `carac`, and aliases, SHALL use the VisuAlg 3.0.7 code domain, accepted string lengths, bounds, recorded character tables, and failure behavior rather than Unicode code-point or UTF-8 byte assumptions. `asc` SHALL use Windows-1252 codes; `carac` SHALL use the separately recorded table, including its drawing-character substitutions. The operations are not required to be inverses.

#### Scenario: Convert extended character codes
- **WHEN** `asc("€")`, `carac(128)`, and `asc(carac(128))` are evaluated
- **THEN** they produce 128, `"Ç"`, and 199 respectively

#### Scenario: Convert an empty value or a control code
- **WHEN** `asc` receives empty text or `carac` receives no argument, a no-value argument, codes 0–31, 127, or 255
- **THEN** `asc` produces no value and `carac` produces one space

#### Scenario: Convert an invalid character code
- **WHEN** a character-conversion built-in receives a code outside the reference domain
- **THEN** `carac` returns no value for out-of-domain integers, rejects non-integer arguments with `E001`, and does not panic

### Requirement: Numeric and character conversion built-ins
The catalog SHALL implement all oracle-confirmed explicit conversion functions and aliases for numeric, logical, and character values. Dynamic numeric conversion such as `caracpnum` SHALL accept and classify integer and real text at runtime according to reference decimal, sign, whitespace, overflow, and result-type rules.

#### Scenario: Convert dynamic numeric text
- **WHEN** `caracpnum` receives text accepted as an integer and text accepted as a real
- **THEN** each call returns the reference value with the reference-selected dynamic numeric type

#### Scenario: Use a conversion in an integer-only context
- **WHEN** the same analyzed `caracpnum` call receives `"1"` and `"1.0"` at runtime
- **THEN** the first result is accepted by `randi` and the second receives positioned `E001`, without a static real-to-integer coercion

#### Scenario: Convert a zero integral prefix
- **WHEN** `caracpnum` receives `"0.5"`, `"-00.5"`, or `"0.5e3"`
- **THEN** each result is integer zero and can be assigned to an integer variable, while `"01.5"` and `"6e-18"` retain their real values and result types

#### Scenario: Apply zero fallback before suffix validation
- **WHEN** decimal conversion receives `"0.5x"`, `"0e9999"`, `"0.0.1"`, or `"0_1"`
- **THEN** each returns integer zero, while the nonzero-prefix control `"01.5x"` receives positioned `R007`

#### Scenario: Preserve conversion failure output
- **WHEN** a character-producing argument emits output and returns malformed numeric text
- **THEN** that output remains visible and conversion receives positioned `R007` at its call

#### Scenario: Reject invalid numeric text
- **WHEN** conversion text contains a malformed number, trailing unsupported characters, or an out-of-range value
- **THEN** it returns the oracle-confirmed fallback or stable positioned built-in diagnostic

#### Scenario: Convert a number to character text
- **WHEN** `numpcarac` receives one numeric argument or an empty parenthesized argument list
- **THEN** it returns the recorded text with 15 significant digits, normalized zero, and uppercase exponents without a plus sign or leading zeros; an empty argument list returns `"0"`

#### Scenario: Preserve an absent conversion result
- **WHEN** `numpcarac` receives character or logical input, or another absent conversion result
- **THEN** it produces no value; an output statement containing that result discards its buffered items and skips later items and format expressions without consuming a pending newline, while a typed return is rejected

### Requirement: Random built-ins
The catalog SHALL implement `rand`, `randi`, and any oracle-confirmed random aliases with the reference arities, argument normalization, inclusive or exclusive bounds, result types, invalid-range behavior, and generator-state interaction. Runs SHALL accept an injected random source for deterministic tests; compatibility SHALL not require matching the reference sequence.

`rand` SHALL be a reserved bare expression returning a real in `[0, 1)`. A parenthesized suffix SHALL be rejected in an output expression with `P001`; assignment and standalone-statement forms SHALL ignore the remaining line tokens without evaluating them. Each evaluated occurrence SHALL consume one injected fraction. A missing source or a non-finite or out-of-range injected fraction SHALL receive positioned `R007` as a project guard and preserve prior output.

#### Scenario: Distinguish random expressions from ignored suffixes
- **WHEN** `x <- rand(1) + 7` and `escreval(rand())` are checked separately
- **THEN** the assignment stores a fraction in `[0, 1)` without evaluating the suffix, and the output expression receives `P001`

#### Scenario: Generate random values at boundaries
- **WHEN** `rand` and `randi` are repeatedly called with valid oracle boundary arguments
- **THEN** every result has the reference type and lies within the reference-confirmed domain

#### Scenario: Reject an invalid random range
- **WHEN** random bounds are empty, reversed, non-finite, or otherwise rejected by VisuAlg 3.0.7
- **THEN** the stable positioned built-in diagnostic is returned without consuming an unintended input value

#### Scenario: Use the recorded integer random bound
- **WHEN** `randi` is called with a signed 32-bit integer
- **THEN** a positive bound gives an integer in `[0, n)`, a negative bound uses its unsigned 32-bit width and a signed 32-bit result, and zero gives zero
- **AND** the empty parenthesized call also gives zero

#### Scenario: Diagnose integer random call syntax and types
- **WHEN** `randi` has omitted parentheses, extra arguments, or an unindexed vector
- **THEN** the recorded `P001` diagnostic is emitted
- **AND** real, character, logical, or absent-value arguments receive `E001`

### Requirement: Built-in failure contract
Invalid built-in arity or statically incompatible arguments SHALL be diagnosed during analysis. Domain, conversion, injected-source, and other runtime-only failures SHALL use the stable positioned built-in failure contract and SHALL never expose a Go panic or silently choose a different overload.

#### Scenario: Call with the wrong arity
- **WHEN** a program calls a known built-in with an unsupported number of arguments
- **THEN** analysis emits a positioned call diagnostic derived from the authoritative catalog

#### Scenario: Trigger a runtime domain failure
- **WHEN** a validly typed built-in call evaluates to an unsupported runtime domain value
- **THEN** execution returns the stable built-in failure diagnostic at the call site
