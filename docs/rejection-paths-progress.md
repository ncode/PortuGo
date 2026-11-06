# Rejected example qualification

The four original rejected examples have fresh VisuAlg 3.0.7.0 observations from
November 6, using the pinned executable on Windows amd64 with en-US culture.
Original source and stdin bytes are unchanged. The corpus retains the original
September diagnostic screenshots, transcriptions, and derived stdout artifacts;
`reobserved-2026-11-06/previous-evidence.json` records their former mappings.

| Original | Qualified diagnostic | Complete language-output prefix |
| --- | --- | --- |
| EXEMPLO1.alg | E001, line 21 | 82 bytes: prompts and three input echoes, ending in LF. |
| Tabela_ASCII4.alg | P001, line 301 | The retained 152-byte waiting banner. |
| ESCOLA1.alg | P001, EOF line 320 | The retained 41-byte prompt and SAIR echo. |
| interrompa.alg | P001, line 16 | Two bytes, space then 1, without LF. |

The earlier empty assignment and repeat output expectations were incomplete.
The replacement artifacts are separately captured and reviewed; old artifacts
are retained. The baseline CLI already matches the complete repeat prefix.
All four mappings are verified by exact original-source CLI replay and
`TestRecordedRejectedExamples`, which also rejects incomplete mappings.

The recorder verifies restaged source/stdin hashes and preserves the actual
observation timestamps. Source was delivered as decoded text and editor
readback was verified before execution; this does not establish reference
file-open decoding. Published raw text selects unchanged Unicode text-control
contents. Operational handles, private automation metadata, and connection
details are omitted and labeled as publication redaction. Private originals
are retained. The application diagnostic PNGs are unchanged and contain no
private identifiers or textual metadata.

The simulated console is a display state, not a native stdout byte stream.
It can include rows created by screen clearing or cursor movement. The results
panel also adds framing and represents partial lines with display line endings.
The reviewed language-output artifacts distinguish those presentation effects
from program writes; a completed-run normalizer is not applied to interrupted
panels. Diagnostic columns absent from the reference are not asserted.

## Follow-up ownership and controls

The manifest preserves the original example-qualification owner, task 17.3.
Follow-up implementation belongs to OpenSpec change
`close-rejected-program-gaps`: assignment controls map to tasks 1.2, 2.1 and
2.2; malformed-condition, EOF and unused-body controls map to tasks 1.4, 4.1
and 4.2; repeat controls map to tasks 1.5, 5.1 and 5.2; neighboring call controls
map to tasks 1.3, 3.1 and 3.2. The new command regressions and exact corpus replay
qualify these controls independently of their historical owner.

The called malformed-body control prints `BODY` and `MAIN` before a missing
`fimse` rejection at EOF. Its separately recorded contract prevents recovery
from forcing a premature procedure-boundary failure. The corresponding unused
body completes successfully. Resource-limit and other closing-boundary guards
remain project regressions and are not new reference claims.

The character-table trigger reduces to a numeric remainder in a text-returning
user helper while a caller text operand is pending. `%` and `mod` reproduce the
same caller P001. Direct and precomputed calls, reversed call order, skipped
remainders, ordinary callee arithmetic/comparisons, and numeric caller additions
remain successful. The shared evaluator records the call's operand boundary
and invalidates pending caller text when the remainder executes; it does not
inspect filenames, names, hashes, or line numbers.

Thirty-nine new controls are recorded and verified. Full replay validates all
1,718 probes: 1,716 verified implementations and two unchanged nonterminating
exclusions. There are 1,696 reference recordings, 22 reviewed evidence
exclusions/project contracts, and 1,877 implementation-test links. None of the
1,673 previously verified probes regressed. All previous source, stdin, and
evidence artifacts retain their original bytes and hashes.

Exploratory reference-application faults found during reduction are retained
privately and excluded from language-diagnostic claims. Wider unrecorded
malformed-program and remainder compositions remain unqualified. No general
compatibility claim follows from the completed four input paths.

## Review coverage

OCR reviewed all 14 changed production Go files. Its two parser findings were
reproduced with failing regressions and fixed: incomplete `se` headers retain
closing boundaries, and a declared callable named `fimrepita` retains ordinary
call resolution after a bounded repeat, including forward declarations. A
follow-up review of the parser correction reported no findings.

OCR's default path exclusions omit `cmd/portugo/deferred_rejections_test.go`,
`cmd/portugo/recorded_rejections_test.go`, `internal/parser/recovery_test.go`,
the conformance manifest, and the four `previous-evidence.json` records.
Unsupported extensions omit the changed Markdown documents/specifications and
all new `.alg`/`.txt` fixtures and observations. The 27 new diagnostic PNGs are
binary exclusions. These files were reviewed locally; exact replay, trace
validation, historical hash verification, and individual image inspection cover
their evidence and privacy obligations. No selected production file was omitted
or failed model review.
