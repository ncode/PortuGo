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

### Requirement: Text built-ins
The catalog SHALL implement the full oracle-confirmed text set, including `copia`, `maiusc`, `minusc`, `compr`, and `pos`, using the reference's 1-based positions, substring bounds, not-found value, case conversion, accented-character behavior, empty-string behavior, and result types.

#### Scenario: Extract a substring
- **WHEN** `copia` receives a valid character value, 1-based position, and length
- **THEN** it returns exactly the oracle-confirmed substring for ASCII and Windows-1252 text

#### Scenario: Search for absent text
- **WHEN** `pos` cannot find its search value
- **THEN** it returns the reference not-found value rather than a Go index convention

### Requirement: Character-code built-ins
Character-code operations, including oracle-confirmed `asc`, `carac`, and aliases, SHALL use the VisuAlg 3.0.7 code domain, accepted string lengths, bounds, CP1252 mapping, and failure behavior rather than Unicode code-point or UTF-8 byte assumptions.

#### Scenario: Round trip an accented character
- **WHEN** an accepted accented character is converted to its code and back
- **THEN** the result and numeric code match the reference CP1252 behavior

#### Scenario: Convert an invalid character code
- **WHEN** a character-conversion built-in receives a code outside the reference domain
- **THEN** it returns the oracle-confirmed value or positioned built-in diagnostic without panicking

### Requirement: Numeric and character conversion built-ins
The catalog SHALL implement all oracle-confirmed explicit conversion functions and aliases for numeric, logical, and character values. Dynamic numeric conversion such as `caracpnum` SHALL accept and classify integer and real text at runtime according to reference decimal, sign, whitespace, overflow, and result-type rules.

#### Scenario: Convert dynamic numeric text
- **WHEN** `caracpnum` receives text accepted as an integer and text accepted as a real
- **THEN** each call returns the reference value with the reference-selected dynamic numeric type

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

#### Scenario: Generate random values at boundaries
- **WHEN** `rand` and `randi` are repeatedly called with valid oracle boundary arguments
- **THEN** every result has the reference type and lies within the reference-confirmed domain

#### Scenario: Reject an invalid random range
- **WHEN** random bounds are empty, reversed, non-finite, or otherwise rejected by VisuAlg 3.0.7
- **THEN** the stable positioned built-in diagnostic is returned without consuming an unintended input value

### Requirement: Built-in failure contract
Invalid built-in arity or statically incompatible arguments SHALL be diagnosed during analysis. Domain, conversion, injected-source, and other runtime-only failures SHALL use the stable positioned built-in failure contract and SHALL never expose a Go panic or silently choose a different overload.

#### Scenario: Call with the wrong arity
- **WHEN** a program calls a known built-in with an unsupported number of arguments
- **THEN** analysis emits a positioned call diagnostic derived from the authoritative catalog

#### Scenario: Trigger a runtime domain failure
- **WHEN** a validly typed built-in call evaluates to an unsupported runtime domain value
- **THEN** execution returns the stable built-in failure diagnostic at the call site
