## MODIFIED Requirements

### Requirement: CLI command behavior
The CLI SHALL provide `run`, `check`, `fmt`, and `repl` with documented arguments, standard-input behavior, output stream separation, and working-directory handling. `run` SHALL execute only after decoding, parsing, and analysis establish a safe program with no blocking diagnostics. A failure SHALL be classified as deferred only when reduced reference evidence establishes its execution timing; `check` SHALL not execute; and `fmt` SHALL not run the program. `run` and `repl` SHALL accept `--max-steps N` as a nonnegative execution budget, defaulting to zero (unlimited); a REPL budget SHALL reset for each submitted program, never for a subcall or loop.

#### Scenario: Check a valid file
- **WHEN** `portugol check` receives a valid accepted program
- **THEN** it exits successfully without program output

#### Scenario: Prevent execution after analysis errors
- **WHEN** `portugol run` receives a program with blocking front-end diagnostics
- **THEN** it reports those blocking diagnostics and produces none of the program's runtime side effects

#### Scenario: Exhaust a CLI execution budget
- **WHEN** `portugol run --max-steps 10000` reaches its budget while executing a valid program
- **THEN** it preserves preceding output, emits positioned `R006` on standard error, and exits with status `1`

#### Scenario: Distinguish deferred diagnostics from blocking failure
- **WHEN** a recorded recoverable construct carries a deferred diagnostic and decoding, parsing, and semantic analysis establish that the remaining program can execute safely
- **THEN** `run` preserves the recorded execution prefix and reports that diagnostic only at its recorded trigger
- **AND** blocking lexical, parsing, semantic, and resource-limit diagnostics continue to prevent execution

#### Scenario: Preserve static checking without runtime effects
- **WHEN** an incompatible assignment is statically diagnosable and its `run` failure is classified as deferred by recorded evidence
- **THEN** `check` retains its static diagnostic behavior without consuming program input or producing runtime output
- **AND** `run` follows the separately recorded reached-execution behavior

## ADDED Requirements

### Requirement: Qualified bundled rejection fidelity
For the four pending rejected bundled probes, the CLI SHALL match the reviewed exit status, diagnostic sequence, established locations, and complete language-output bytes with unchanged original source and stdin. Corrected output expectations SHALL require separately retained, complete reference evidence and an explicit reviewed disposition. Columns absent from reference evidence SHALL NOT become reference assertions. Implementation rules SHALL be independent of source filename, algorithm name, source hash, and fixed line-number checks.

#### Scenario: Preserve input and output before the assignment rejection
- **WHEN** `run` executes `bundled-24919c4a2e4c` (`EXEMPLO1.alg`) after its corrected prefix has passed reference qualification
- **THEN** it exits with status `1`, reports exactly one `E001` on line 21, and preserves the reviewed 82-byte prompt/input prefix ending in LF

#### Scenario: Stop the character table at its recorded failure
- **WHEN** `run` executes `bundled-49eb7f47065e` (`Tabela_ASCII4.alg`)
- **THEN** it exits with status `1`, reports exactly one `P001` on line 301, and writes exactly the reviewed 152-byte language-output prefix

#### Scenario: Preserve the school program's exit-path prefix
- **WHEN** `run` executes `bundled-f251f2cd0c22` (`ESCOLA1.alg`) with its committed stdin
- **THEN** it exits with status `1`, reports exactly one `P001` at EOF on line 320, and writes exactly the reviewed 41-byte prompt/input prefix

#### Scenario: Preserve output before the repeat rejection
- **WHEN** `run` executes `bundled-f4fa0c514cf7` (`interrompa.alg`) after its corrected prefix has passed reference qualification
- **THEN** it exits with status `1`, reports exactly one `P001` on line 16, and writes exactly space followed by `1`, with no newline

### Requirement: Rejection compatibility safeguards
Evidence-backed rejection handling SHALL preserve existing accepted programs, diagnostic collection outside the confirmed rule, positions, bounded traversal, execution budgets, and host-error handling. It SHALL NOT execute an unchecked AST after blocking front-end failures or leak state into subsequent runs/submissions. `check` and `fmt` SHALL NOT execute deferred constructs. Formatter output SHALL preserve their observable behavior or reject formatting before changing the input file.

#### Scenario: Preserve the accepted repeat escape
- **WHEN** a recorded `repita` body executes `interrompa` before its `fimrepita` marker
- **THEN** execution preserves the already verified successful escape and continuation behavior

#### Scenario: Preserve valid nested calls
- **WHEN** a neighboring string-call form is recorded as successful
- **THEN** implementing the character-table rejection preserves its recorded result without introducing a general sibling/nested-call ban

#### Scenario: Keep analysis-only commands free of runtime effects
- **WHEN** `check` or `fmt` processes a source with a deferred rejection
- **THEN** it consumes no program input and performs no program writes or host operations
- **AND** in-place formatting either preserves deferred behavior on replay or rejects formatting with the original bytes intact

#### Scenario: Recover after a rejected submission
- **WHEN** a REPL submission reaches a deferred rejection and a valid program follows
- **THEN** the first submission's diagnostics and pending output do not contaminate the next program
