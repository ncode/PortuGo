## Purpose

Defines reference-compatible input, output, file-backed execution, random-input mode, encoding, timing, display, debugging, and other environment commands.

## ADDED Requirements

### Requirement: Console input consumption
`leia` SHALL consume user input in the token, line, whitespace, prompt, retry, and end-of-input pattern recorded from VisuAlg 3.0.7 for each destination type. Multiple destinations and consecutive calls SHALL share one input stream without losing buffered data, and conversion failures SHALL identify the responsible source argument.

#### Scenario: Read multiple destinations
- **WHEN** one `leia` statement receives multiple compatible destinations and the input contains the required values across oracle-supported whitespace or line boundaries
- **THEN** each destination receives the corresponding converted value in order and unread input remains available

#### Scenario: Preserve console lines and typed echo
- **WHEN** destinations receive recorded console input lines, including spaces and empty character values
- **THEN** each destination consumes one complete line, character text is preserved, and successful input echoes use plain integers, ten fractional digits for reals, `Verdadeiro`/`Falso` for logical values, or the exact character text, followed by LF

#### Scenario: Reach console end of input
- **WHEN** execution needs another value after the injected console input is exhausted
- **THEN** it returns the positioned input diagnostic or reference-compatible fallback behavior without blocking an explicitly headless run

### Requirement: Input value conversion
Input SHALL reproduce the reference syntax and range rules for integers, reals, logical values, and character values, including accepted decimal separators, boolean spellings and casing, leading signs, surrounding whitespace, and overflow. Completed conversions SHALL replace the destination with the recorded result, including zero or an unsigned mantissa produced from malformed text. Read failures and project range guards SHALL preserve the prior value.

#### Scenario: Read a real value
- **WHEN** input uses a decimal spelling accepted by the reference for a real destination
- **THEN** the destination receives the same numeric value

#### Scenario: Apply recorded scalar conversions
- **WHEN** console input contains malformed numeric text, an integer requiring narrowing, or logical text
- **THEN** malformed numbers retain the unsigned mantissa before scaling or become zero without digits, integers truncate and narrow to signed 32-bit values, and logical input is true exactly when its first character is `v` or `V`

#### Scenario: Guard an unsupported numeric conversion
- **WHEN** the headless runtime encounters input exhaustion, a number outside the finite real range, or integer conversion outside the signed 64-bit range before narrowing
- **THEN** it reports positioned `R004` as a project guard and retains the destination and preceding output

### Requirement: Exact output rendering
`escreva` and `escreval` SHALL reproduce VisuAlg 3.0.7 item separation, implicit spaces, newline placement, logical casing, integer and real rendering, decimal separator, negative zero handling, and string output. Output SHALL be written in argument order without host-language formatting leakage.

#### Scenario: Render mixed output items
- **WHEN** one output statement contains character, integer, real, and logical expressions
- **THEN** its exact bytes, including spaces, casing, decimal characters, and final newline state, match the oracle evidence

#### Scenario: Distinguish write forms
- **WHEN** equivalent items are emitted once with `escreva` and once with `escreval`
- **THEN** the only newline differences are those observed in VisuAlg 3.0.7

#### Scenario: Nested writes during argument evaluation
- **WHEN** a function called by an output item or its format expression performs another output statement
- **THEN** the nested statement emits before the outer statement's buffered items, and the next output statement to finish consumes any pending newline requested by `escreval`, including when that finishing statement is `escreva`

#### Scenario: Replay the recorded portable output profile
- **WHEN** selected 2026-09-07 `en-US` reference observations run through the portable CLI
- **THEN** program output uses decimal dots, numeric and logical leading spaces, and uppercase logical values, with only fixed reference UI notices removed and CRLF converted to LF for comparison; other locale behavior remains pending evidence

### Requirement: Width and precision formatting
Default real rendering SHALL use 15 significant digits, discard the sign of zero, and use uppercase `E` without a plus sign or leading exponent zeros when scientific notation is needed. A nonpositive width SHALL ignore the decimal-count argument.

Output width and precision fields SHALL accept the expression forms and value domains supported by VisuAlg 3.0.7 and SHALL reproduce its alignment, padding, rounding, truncation, sign placement, overflow-width, and non-real precision behavior. Invalid format values SHALL produce positioned diagnostics.

#### Scenario: Cap a positive field width
- **WHEN** a character or numeric output item requests a width above 255
- **THEN** padding uses width 255, preserving character left alignment and numeric right alignment without allocating the requested larger width

#### Scenario: Format a real with width and decimals
- **WHEN** an output item supplies valid width and decimal expressions
- **THEN** the exact padded and rounded result matches the committed oracle bytes

#### Scenario: Value exceeds requested width
- **WHEN** a rendered value is wider than its valid requested width
- **THEN** output expands, truncates, or fails exactly as the reference does

#### Scenario: Replay recorded field formatting
- **WHEN** the recorded width and decimal fields are applied
- **THEN** width zero ignores decimals, positive numeric fields round decimal ties away from zero and expand as necessary, positive string fields left-align and truncate, and logical field widths are rejected before execution

### Requirement: CP1252 file and stream behavior
Source-independent text read from or written to VisuAlg-compatible files SHALL use the oracle-confirmed Windows-1252 and newline behavior. Unsupported Unicode output SHALL follow a documented reference-compatible substitution or positioned failure rule; conversion SHALL never silently corrupt unrelated bytes.

#### Scenario: Read accented file input
- **WHEN** an `arquivo` input contains Windows-1252 accented text
- **THEN** `leia` observes the intended character value and byte offsets remain traceable

#### Scenario: Write a generated text file
- **WHEN** a reference command produces a text file
- **THEN** its encoding and newline bytes match the recorded generated-file evidence

### Requirement: Arquivo input lifecycle
The oracle-confirmed `arquivo` command SHALL resolve relative paths against the configured program working directory and SHALL reproduce VisuAlg 3.0.7 behavior for existing, missing, unreadable, empty, and exhausted files. Any transition to interactive input, fallback recording, created file, or echo mode SHALL match the committed reference state machine and generated bytes.

#### Scenario: Consume an existing arquivo
- **WHEN** `arquivo` names a readable file containing enough accepted values
- **THEN** subsequent input consumes those values in reference order without consulting console input

#### Scenario: Handle a missing arquivo
- **WHEN** `arquivo` names a missing path
- **THEN** the recorded reference fallback, recording, or positioned file diagnostic occurs relative to the configured working directory

#### Scenario: Exhaust arquivo input
- **WHEN** the configured file contains fewer values than subsequent reads require
- **THEN** input changes source, records fallback data, or fails exactly as the oracle evidence specifies

### Requirement: Random-input mode
The command-form random-input facility, including `aleatorio` and any paired range or disable commands confirmed by the oracle, SHALL reproduce its activation, bounds, destination-type conversion, echo, interaction with console and `arquivo` input, and reset behavior. It SHALL match the reference value domains but need not reproduce exact sequences.

#### Scenario: Read while random-input mode is active
- **WHEN** random-input mode is enabled with valid bounds and `leia` requests a supported destination
- **THEN** a value within the oracle-confirmed domain is supplied without requiring console bytes and all observable echo behavior matches the reference

#### Scenario: Disable random-input mode
- **WHEN** the reference-confirmed disable form executes
- **THEN** later reads return to the correct prior or default input source

### Requirement: Echo, timer, and chronometer commands
Echo, delay or timer, and chronometer commands SHALL reproduce the state transitions, units, output text, clock sampling, reset behavior, and interaction with input recorded from VisuAlg 3.0.7. Runs with an injected clock SHALL be deterministic.

#### Scenario: Measure elapsed time
- **WHEN** a program starts, queries, resets, or stops the chronometer using accepted forms while the injected clock advances
- **THEN** the observed value and output match the reference semantics for that clock sequence

#### Scenario: Toggle input echo
- **WHEN** echo is enabled or disabled around console, file, or random input
- **THEN** exactly the reference-selected input representations are emitted

### Requirement: Host-backed UI commands
Pause, debug or breakpoint, clear-screen, color or display, clock, and every additional oracle-confirmed environment command SHALL expose the reference arguments and observable sequencing through a host abstraction. The default headless host SHALL make UI-only effects no-ops while preserving validation, non-UI output, state transitions, and errors; it SHALL not block waiting for a GUI.

#### Scenario: Clear the screen in a headless run
- **WHEN** a valid clear-screen command executes with the default headless host
- **THEN** execution continues without terminal escape leakage or failure and subsequent ordinary output is preserved

#### Scenario: Select a recorded display color
- **WHEN** `mudacor` receives two character expressions on the command line
- **THEN** they are evaluated left to right, the seven recorded color names and `frente`/`fundos` targets are matched without case distinctions or whitespace trimming, and recognized pairs produce one typed host event
- **AND** unknown names preserve display settings, trailing syntax is ignored, and host failures return positioned `R008` without rendering the underlying error details

#### Scenario: Use a display keyword as a value
- **WHEN** `limpatela` or `mudacor` occurs in an expression with or without apparent call arguments
- **THEN** it produces no value without a host event or argument evaluation, and a containing output statement follows the recorded no-value discard rule

#### Scenario: Execute a breakpoint with a fake host
- **WHEN** a debug or pause command executes with a recording host
- **THEN** the host receives one typed operation at the correct point between surrounding language effects

### Requirement: Environment failure diagnostics
File path, encoding, injected input/output, clock, random-source, and host operation failures SHALL return positioned runtime diagnostics and SHALL not panic, terminate the process directly, or expose raw host errors as the only explanation. Resources SHALL be released when execution succeeds or fails.

#### Scenario: Host operation fails
- **WHEN** an injected host reports a failure for an operation that can fail
- **THEN** execution returns the stable host/file diagnostic at the triggering statement and closes any opened resources
