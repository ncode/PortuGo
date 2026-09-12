# Display command recordings

This slice adds 72 VisuAlg 3.0.7 recordings: 57 completed programs and 15
positioned rejections. Sixty-seven match the implementation and have permanent
tests; five remain pending. Source bytes, program-output text, and reviewed
diagnostic transcriptions are linked by hash from the corpus manifest.

The observations cover clear-screen syntax, ignored arguments, expression
contexts, color argument order, physical line boundaries, invalid values,
case handling, whitespace, and five unchanged bundled programs. The supported
RGB colors and the plural `fundos` target were checked against the displayed
program region. The reference's display is separate from its captured text
panel; clearing or recoloring it does not erase or decorate that text stream.
Only program text and diagnostic-window crops are retained in this repository.

`TestRecordedDisplayCommands` checks original and canonical formatted execution.
`TestDisplayCommandDiagnostics` pins rejection codes and lines.
`TestDisplayHostOrder`, `TestDisplayPalette`, and
`TestDisplayKeywordsHaveNoValue` check typed host events, both color targets,
ignored options, call order, interpreter reuse, and expression contexts.
The parser depth tests cover both color expressions at and beyond the limit.

The five pending records are:

- `display-color-bad-second`: the reference prints the first argument's effect
  before rejecting the second; current semantic analysis rejects before execution.
- `display-color-no-value-first`: mixed integer/string addition has a different
  diagnostic from the reference before the display operation is reached.
- `display-color-no-value-second`: the same expression difference combines with
  the reference's retained first-argument output.
- `bundled-991ec2bd1566` and `bundled-7f44aa03fcc5`: accepted original examples
  print random values. Their samples do not establish exact portable sequences.

The corpus now has 875 reference recordings: 834 verified and 41 pending,
plus 15 explicitly separate project-contract records. The evidence gate still
requires mappings for 16 requirements and 31 bundled examples. These recordings
do not close the remaining environment-command or overall conformance tasks.

## Injected host failures

`project.display-host-failures` covers errors returned by the embedding Go host.
They cannot be injected as equivalent source-level operations in the reference
editor. `TestDisplayHostFailures` checks `R008` at the original command,
preserved output and cause, suppressed operational error details, stopped later
effects, and successful interpreter reuse after failure. This project contract
does not substitute for file, clock, encoding, or other reference recordings.

The same project record now links the existing tests for the other fallible
typed host operations: `TestConsoleConfiguration`, `TestEchoHost`,
`TestTimerResetAndFailures`, and `TestBreakpointHost`. They check console,
echo, delay, and pause/debug failures at their source positions while stopping
later effects and retaining preceding output and inspectable causes.
`TestTimerCallFailureRestoresFrame` checks cleanup when a delay fails during
subprogram entry. `TestChronometerResetAndFailures` covers clock regression,
per-run reset, and output failure; `TestDiagnosticContract` checks that
rendered diagnostics omit underlying error details.

Headless behavior is also exercised through bounded CLI replay of recorded
programs, including `console-header`, `echo-command-on`, `echo-command-off`,
`execution-pause-bare`, `execution-debug-true`, `display-clear-bare`,
`display-color-foreground`, and `environment-timer-straight`. Each uses a
five-second failure deadline and exact output comparison. UI operations add no
terminal escapes and require no interactive continuation; requested timer
delays still wait. This closes task 15.9 without changing recorded expectations
or qualifying pending elapsed-time and source-level rejection-order cases.
