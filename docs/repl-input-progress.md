# REPL input progress

The REPL is a project interface, separate from the reference application's
editor. Its stream-sharing, submission, and recovery rules therefore use project
tests with explicit reference non-applicability. Language programs entered
through the REPL still require the same reference evidence as file execution.

## Shared input

`project.repl-shared-input` verifies that source entry and interpreter input use
one buffered reader. `TestSharedInputPreservesFinalProgram` reads a value exactly
once and then executes another program from the remaining stream.
`TestBufferedSourceAtEOF` verifies final-buffer submission, diagnostics, explicit
cancellation, and preservation of a prior failure's session status.
`TestCommandContracts` checks the process exit statuses and output ordering.

`TestSourceEncoding` verifies that submitted UTF-8, UTF-8 BOM, and Windows-1252
source follows the existing file decoder. This is input-path consistency, not a
new reference encoding observation.

## Blank-line preservation

`project.repl-blank-lines` now verifies that blank lines remain inside unfinished
source. `TestIncompleteBlankLinesAndTerminatorLookalikes` checks blank lines in
declarations and statements, strings, comments, and longer identifiers.
`TestAutomaticSubmissionDiagnostics` verifies that source positions still count
preserved blank lines.

## Automatic submission and recovery

`project.repl-submission` covers the project REPL protocol. Its token completion
result recognizes a real `fimalgoritmo`; lexer tokens distinguish
strings, comments, and longer identifiers. Complete invalid input still goes
through the ordinary decoder, parser, and analyzer and reports its diagnostics.
Each new physical line is scanned once, and the full source is analyzed only
when submitted, avoiding repeated parsing of the growing buffer.

`TestAutomaticSubmission` checks immediate execution, consecutive programs,
shared `leia` input, uppercase terminators, and LF/CRLF. The CLI transcript
`testdata/cli/repl_auto.in` also checks exact prompt/output order and exit status
through `TestCommandContracts`. Lexical, syntax, semantic, runtime, and budget
failures report once, retain the session's failure status, and permit a fresh
program with its own execution budget.

Line framing uses decoded text, preventing Windows-1252 bytes in longer names
from being mistaken for a terminator. `TestCP1252TerminatorLookalike` verifies
cancellation of that unfinished input; it does not claim acceptance of its
identifier form as a reference language feature.

Source-size and host-diagnostic recovery are covered below. Integration with
future environment input modes and formatter retention remain pending.

Local build, formatting, vet, staticcheck, lint, ordinary/race tests, both
30-second fuzz checks, strict specification validation, and all 199 verified
reference CLI replays pass. Native Windows build, vet, 766 tests, both fuzz
checks, and all 199 CLI replays pass. All 11 transported files match the local
snapshot before the final validation notes. The evidence gate reports only
106 missing mappings: 33 requirements and 73 bundled examples.

## Limit and host recovery

`TestSourceLimitRecovery` covers oversized first and continuation lines, a
terminator crossing the cap, a short line exceeding the remaining allowance,
long discarded lines, header lookalikes, EOF, and cancellation. The rejected
submission produces one `E900`. Recovery discards the rest of its physical line
with bounded reads and resumes at a line whose first token is `algoritmo`.
The next program receives its `leia` input once, and the session retains a
failure status. The CLI contract test checks exact prompts, output, diagnostic
position, and exit status for this transition.

`TestHostDiagnosticRecovery` verifies that a transient program output failure
produces positioned `R008`, preserves earlier output, and permits a later
submission. `TestNilDiagnosticWriter` verifies controlled lexical and source-limit
failures when diagnostic output is discarded. Invalid decoded bytes also retain
their lexical diagnostic and allow a fresh next program. Permanently unusable
input or prompt streams remain operational failures that end the session.

Local build, formatting, vet, staticcheck, lint, ordinary/race tests, both
30-second fuzz checks, strict specification validation, and all 199 verified
reference CLI replays pass for recovery. Native Windows build, vet, 777 tests,
both fuzz checks, and all 199 CLI replays pass. All 11 transported files match
the local snapshot before final validation prose. The evidence gate still
reports only the 106 missing inventory mappings.
