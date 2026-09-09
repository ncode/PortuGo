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

## Pending blank-line behavior

`project.repl-blank-lines` remains pending under tasks 16.3 and 16.4. The current
REPL still submits on a blank line instead of preserving it inside incomplete
source. Completing the shared-input requirement therefore remains blocked even
though the common reader and EOF paths are tested.

Automatic submission, recovery after source-size exhaustion, and integration
with future environment input modes also remain pending. This slice does not
close the full REPL or formatter group.
