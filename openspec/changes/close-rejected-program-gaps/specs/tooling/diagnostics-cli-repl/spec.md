## MODIFIED Requirements

### Requirement: CLI command behavior
The CLI SHALL provide `run`, `check`, `fmt`, and `repl` with documented arguments, standard-input behavior, output stream separation, and working-directory handling. `run` SHALL execute only after successful decoding, parsing, and analysis. A malformed construct SHALL be deferred until execution only when reduced reference evidence establishes that timing and the front end can retain its diagnostic without a blocking analysis failure; `check` SHALL not execute; and `fmt` SHALL not run the program. `run` and `repl` SHALL accept `--max-steps N` as a nonnegative execution budget, defaulting to zero (unlimited); a REPL budget SHALL reset for each submitted program, never for a subcall or loop.

#### Scenario: Check a valid file
- **WHEN** `portugol check` receives a valid accepted program
- **THEN** it exits successfully without program output

#### Scenario: Prevent execution after analysis errors
- **WHEN** `portugol run` receives a program with front-end diagnostics
- **THEN** it reports those diagnostics and produces none of the program's runtime side effects

#### Scenario: Exhaust a CLI execution budget
- **WHEN** `portugol run --max-steps 10000` reaches its budget while executing a valid program
- **THEN** it preserves preceding output, emits positioned `R006` on standard error, and exits with status `1`

#### Scenario: Distinguish deferred diagnostics from analysis failure
- **WHEN** a recorded recoverable construct carries a deferred diagnostic and decoding, parsing, and semantic analysis complete without blocking diagnostics
- **THEN** `run` preserves the recorded execution prefix and reports the deferred diagnostic only at its recorded trigger
- **AND** ordinary lexical, parsing, semantic, and resource-limit failures continue to prevent execution

## ADDED Requirements

### Requirement: Recorded bundled rejection fidelity
For each of the four pending rejected bundled probes, the CLI SHALL match the recorded exit status, diagnostic sequence, established source locations, and exact standard-output bytes with the committed source and stdin. Diagnostic columns absent from reference evidence SHALL NOT be invented as reference assertions. Matching SHALL follow a rule established by reduced reference controls, independent of source filename, algorithm name, source hash, or fixed line-number checks.

#### Scenario: Report the earlier assignment rejection
- **WHEN** `run` executes the unchanged source and stdin of `bundled-24919c4a2e4c` (`EXEMPLO1.alg`)
- **THEN** it exits with status `1`, reports exactly one `E001` on line 21, and writes no standard output

#### Scenario: Stop the character table at its recorded failure
- **WHEN** `run` executes the unchanged source and stdin of `bundled-49eb7f47065e` (`Tabela_ASCII4.alg`)
- **THEN** it exits with status `1`, reports exactly one `P001` on line 301, and writes exactly the 152 bytes in its committed expected-output artifact

#### Scenario: Preserve the school program's exit-path prefix
- **WHEN** `run` executes the unchanged source of `bundled-f251f2cd0c22` (`ESCOLA1.alg`) with its committed `SAIR` input
- **THEN** it exits with status `1`, reports exactly one `P001` at EOF on line 320, and writes exactly the 41 bytes in its committed expected-output artifact

#### Scenario: Reject the repeat marker without extra output
- **WHEN** `run` executes the unchanged source and stdin of `bundled-f4fa0c514cf7` (`interrompa.alg`)
- **THEN** it exits with status `1`, reports exactly one `P001` on line 16, and writes no standard output

### Requirement: Rejection compatibility safeguards
Evidence-backed rejection recovery SHALL preserve existing accepted programs, diagnostic collection outside the recorded rule, source positions, bounded traversal, execution budgets, and host-error handling. It SHALL NOT execute an AST after a blocking front-end diagnostic. Rejected-program handling SHALL NOT leak state into subsequent interpreter or REPL submissions. `check` and `fmt` SHALL NOT execute deferred constructs; formatter output SHALL preserve their observable behavior or reject formatting before changing the input file.

#### Scenario: Preserve the accepted repeat escape
- **WHEN** a recorded `repita` body unconditionally executes `interrompa` before its `fimrepita` marker
- **THEN** execution keeps the already verified successful escape and continuation behavior

#### Scenario: Keep analysis-only commands free of runtime effects
- **WHEN** `check` or `fmt` processes a source containing a deferred rejection
- **THEN** it consumes no program input and performs no program writes or host operations
- **AND** in-place formatting either preserves the deferred behavior on replay or rejects formatting with the original bytes intact

#### Scenario: Recover after a rejected submission
- **WHEN** a REPL submission reaches a deferred rejection and is followed by a valid program
- **THEN** the rejected submission's diagnostic and pending output do not contaminate the subsequent program
