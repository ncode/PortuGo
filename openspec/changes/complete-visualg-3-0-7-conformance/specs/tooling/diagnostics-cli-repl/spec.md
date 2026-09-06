## Purpose

Defines stable positioned diagnostics and predictable command-line, formatting, and interactive behavior for conformance-oriented Portugol workflows.

## ADDED Requirements

### Requirement: Stable compile-time diagnostics
Source decoding, lexical, syntax, and semantic failures SHALL be returned as ordered diagnostic values with stable public codes, severity, source name, byte or rune span, and rendered line and column. Tests and automation SHALL be able to rely on codes independently of human-readable wording.

#### Scenario: Report multiple analysis failures
- **WHEN** a source contains recoverable failures from more than one front-end stage
- **THEN** the command reports deterministic ordered diagnostics with stable codes and accurate source positions

### Requirement: Stable runtime diagnostics
Runtime failures SHALL use these stable categories: `R001` type or coercion, `R002` arithmetic, `R003` storage or indexing, `R004` input, `R005` call or return, `R006` loop state, `R007` built-in failure, and `R008` host or file I/O. Every runtime diagnostic SHALL identify the source construct responsible for the failure and SHALL retain an inspectable underlying cause without making host-specific text the stable contract.

#### Scenario: Report an out-of-bounds index
- **WHEN** a vector index is outside a resolved declared bound
- **THEN** execution returns `R003` positioned at the indexing expression and the CLI renders that position

#### Scenario: Report a missing input file
- **WHEN** an `arquivo` operation fails for a missing path and the reference semantics require failure
- **THEN** execution returns `R008` positioned at that statement with the filesystem cause available for debugging

### Requirement: No user-input panics
All CLI and REPL operations SHALL handle malformed, adversarial, truncated, deeply nested, and oversized source or input with diagnostics or controlled resource-limit failures. They SHALL not expose a Go panic, deadlock, unbounded recovery loop, or partial process crash for user-controlled data.

#### Scenario: Process fuzz-generated input
- **WHEN** lexer, parser, semantic, runtime, formatter, or REPL entry points receive arbitrary bytes with structural limits, a finite runtime execution budget when executing, and bounded nonblocking test input/output and host operations
- **THEN** the operation terminates with a result or diagnostic and no user-visible panic

### Requirement: Resource limits and bounded execution
The implementation SHALL document and enforce source-size, syntax/traversal-depth, call-depth, and text/format-allocation limits before the corresponding allocation or descent. Front-end limit failures SHALL use stable positioned `E900` diagnostics; runtime call-depth failures SHALL use `R005`, text/allocation or evaluation-depth failures SHALL use `R003`, and execution-budget exhaustion SHALL use `R006`. Execution SHALL offer an opt-in finite work budget shared across the entire run, including recursive calls, input retries, and loops with empty bodies. Exhaustion SHALL stop before the next operation's side effects and release resources normally. Ordinary execution with no work budget SHALL allow intentional infinite loops and interactive input waits; universal termination SHALL NOT be a language guarantee.

#### Scenario: Exceed nesting or source limits
- **WHEN** source size or syntax/traversal depth would exceed a documented front-end limit
- **THEN** processing returns a positioned `E900` before unbounded buffering or recursive descent and in-place formatting preserves the original file

#### Scenario: Exceed call depth
- **WHEN** another language invocation would exceed the active-call limit
- **THEN** execution returns `R005` at that call before allocating its frame and unwinds existing frames normally

#### Scenario: Bound an empty infinite loop in a test
- **WHEN** a valid loop with an empty body runs under a finite execution budget
- **THEN** its iterations consume that shared budget and execution returns positioned `R006` without a Go panic or watchdog kill

#### Scenario: Exceed a text or formatting allocation limit
- **WHEN** concatenation, conversion, input buffering, copying, or a format width/precision would require an allocation beyond its documented bound
- **THEN** the operation returns positioned `R003` before the oversized allocation

#### Scenario: Distinguish safeguards from oracle results
- **WHEN** a known accepted reference example exhausts a configured budget or resource limit
- **THEN** the run fails acceptance until the bound is deliberately reviewed and retested; the example is not silently excluded and exhaustion is not counted as matching reference output

### Requirement: CLI command behavior
The CLI SHALL provide `run`, `check`, `fmt`, and `repl` with documented arguments, standard-input behavior, output stream separation, and working-directory handling. `run` SHALL execute only after successful decoding, parsing, and analysis; `check` SHALL not execute; and `fmt` SHALL not run the program. `run` and `repl` SHALL accept `--max-steps N` as a nonnegative execution budget, defaulting to zero (unlimited); a REPL budget SHALL reset for each submitted program, never for a subcall or loop.

#### Scenario: Check a valid file
- **WHEN** `portugol check` receives a valid accepted program
- **THEN** it exits successfully without program output

#### Scenario: Prevent execution after analysis errors
- **WHEN** `portugol run` receives a program with front-end diagnostics
- **THEN** it reports those diagnostics and produces none of the program's runtime side effects

#### Scenario: Exhaust a CLI execution budget
- **WHEN** `portugol run --max-steps 10000` reaches its budget while executing a valid program
- **THEN** it preserves preceding output, emits positioned `R006` on standard error, and exits with status `1`

### Requirement: CLI exit statuses
CLI commands SHALL exit with status `0` on success, status `1` when valid command usage encounters source, semantic, runtime, formatting, or operational diagnostics, and status `2` for command-line usage errors. Diagnostics SHALL be written to standard error, while language output and requested formatted source SHALL be written to standard output unless an explicit output file is selected.

#### Scenario: Runtime failure exit status
- **WHEN** `portugol run` encounters a positioned runtime diagnostic after producing any preceding language output
- **THEN** it preserves that standard output, writes the diagnostic to standard error, and exits with status `1`

#### Scenario: Invalid command usage
- **WHEN** a command receives an unknown option or invalid argument count
- **THEN** it writes usage guidance to standard error and exits with status `2`

### Requirement: Formatter contract
`portugol fmt` SHALL use the canonical printer, preserve source behavior and all accepted comments and reference-ignored suffixes, be idempotent, and reject malformed or resource-limited input without overwriting the original. Check, standard-output, and in-place modes SHALL have documented and testable exit statuses.

#### Scenario: Format a valid file in place
- **WHEN** in-place formatting succeeds
- **THEN** the file contains canonical source that parses equivalently and a second formatting pass makes no change

#### Scenario: Format a malformed file
- **WHEN** in-place formatting encounters syntax diagnostics
- **THEN** it exits with status `1`, reports the diagnostics, and leaves the original bytes unchanged

### Requirement: Shared REPL input
The REPL SHALL use a single input abstraction for prompts, multi-line program entry, and program `leia` operations so buffered text is neither lost nor consumed by the wrong layer. Blank lines inside an incomplete program SHALL be preserved or ignored according to grammar needs rather than unconditionally submitting or discarding the buffer.

#### Scenario: Program reads after multi-line entry
- **WHEN** a completed REPL program executes `leia` and more input follows in the shared stream
- **THEN** the program receives that input once and the next REPL prompt begins after it

#### Scenario: Enter a blank line in an incomplete program
- **WHEN** a blank line occurs before the program is syntactically complete
- **THEN** the REPL remains in continuation mode without losing the accumulated source

### Requirement: Automatic REPL submission
The REPL SHALL automatically submit an accumulated complete program when `fimalgoritmo` terminates it according to the parser, without requiring a second blank line or end-of-file. Text inside strings or comments that resembles `fimalgoritmo` SHALL not trigger submission.

#### Scenario: Finish a complete program
- **WHEN** the parser recognizes `fimalgoritmo` as the terminating token of the accumulated REPL program
- **THEN** the REPL immediately analyzes and runs it, prints diagnostics or output, clears the submitted buffer, and displays the next primary prompt

### Requirement: REPL recovery
After lexical, syntax, semantic, or runtime diagnostics, the REPL SHALL return to a usable primary prompt with program state reset or retained according to documented behavior. End-of-input SHALL terminate cleanly, and a failed submission SHALL not contaminate the next program buffer.

#### Scenario: Recover after malformed submission
- **WHEN** an automatically submitted program contains diagnostics
- **THEN** the REPL reports them once and accepts a subsequent valid program in a fresh buffer
