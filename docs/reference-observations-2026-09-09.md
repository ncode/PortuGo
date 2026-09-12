# Reference observations, September 9

This batch records 61 reduced programs against the same hashed VisuAlg 3.0.7.0
executable and `en-US` profile as the existing corpus. Forty-one completed
successfully and twenty produced a visible rejection diagnostic. Recording
alone does not verify implementation; current verification states and tests
are tracked in the manifest. These observations guide subsequent changes and
do not certify complete conformance.

Exact CP1252 source bytes were decoded into the editor control and checked by
readback before execution. Capture timestamps, source hashes, and executable
hashes were verified. Successful output uses the existing `panel-v1` normalizer.
Rejected programs retain individually reviewed, unchanged screenshots, labeled
manual transcriptions, and separate partial output. Published artifacts contain
synthetic programs and reference results only. One run without a final output
notice or visible diagnostic was excluded as inconclusive.

| Area | Recorded result | Representative probe IDs |
| --- | --- | --- |
| Program lines | A one-line complete program, a header name on the next line, two statements on one line, and an expression split across lines are rejected. An empty `var` section is accepted. | `header-same-line`, `header-name-next-line`, `same-line-statements`, `semicolon-statements`, `expression-next-line`, `empty-var-section` |
| Post-program text | Words, another statement, and an unterminated quoted suffix after `fimalgoritmo` are ignored. | `post-terminator-words`, `post-terminator-statement`, `post-terminator-broken-string` |
| Identifiers | Leading/internal underscores and mixed ASCII case are accepted. The tested accented identifier and accented keywords are rejected. | `leading-underscore`, `internal-underscore`, `mixed-case-identifier`, `accented-identifier`, `accented-keywords` |
| Strings | Backslashes remain literal, including `\n`, `\t`, `\\`, an unknown sequence, and a final backslash. Backslash-quoted and doubled-quote forms are rejected. | `backslash-string-newline`, `backslash-string-tab`, `backslash-string-pair`, `backslash-string-unknown`, `backslash-string-final`, `backslash-string-quote`, `doubled-string-quote` |
| Comments | The tested brace and C-like openers suppress the rest of their physical line, including text after a closing delimiter. They do not suppress statements on following lines. A Pascal-like opener is rejected. | `brace-comment-inline`, `brace-comment-cross-line`, `c-comment-inline`, `c-comment-cross-line`, `c-comment-unclosed`, `pascal-comment` |
| Numbers | `1e2` and `5.` are accepted; `.5` is rejected. | `numeric-exponent`, `trailing-decimal-point`, `leading-decimal-point` |
| Calls | Parameterless declarations and calls can omit parentheses. A function used as a statement is rejected. The tested integer passed to a real reference parameter is accepted and receives the assigned value 3. | `procedure-no-parentheses`, `procedure-bare-call`, `function-no-parentheses`, `function-bare-call`, `function-used-as-statement`, `reference-integer-to-real` |
| Declarations and assignment | A constant followed by an empty `var` section, a constant vector bound, record copying, and `:=` assignment are accepted. Whole-vector assignment and `=` assignment are rejected. | `constant-empty-var`, `constant-vector-bound`, `record-copy`, `vector-copy`, `assignment-colon-equals`, `assignment-equals` |
| Bounds | Lower bounds 0 and 2 are accepted. The tested literal negative bound is rejected. | `vector-zero-bound`, `vector-positive-bound`, `vector-non-one-bound` |
| Storage | Vectors with 500, 501, 5000, and 5001 elements can assign and read their final element. | `vector-500-slots`, `vector-501-slots`, `vector-5000-slots`, `vector-5001-slots` |
| Choice | The tested `caso 1..5` range is rejected. | `case-range` |

The bundled manual's 500-variable description does not match the tested vector
behavior of this executable. The draft specification now preserves recorded
accepted sizes and requires evidence for any further compatibility restriction.
These probes establish no universal upper limit and do not resolve accounting
for every scalar, aggregate, local, or recursive declaration context.

Diagnostic codes in the manifest are the intended project mappings, not codes
displayed by the reference application. GUI line numbers are preserved; columns
are omitted where they were not established. Full evidence and implementation
acceptance remain blocked by the outstanding inventory and pending behavior.

## Focused comment follow-up

Twenty-four additional reduced programs use the same executable, profile, and
capture method: fifteen complete and nine report visible errors. Each rejection
has an individually reviewed screenshot and a labeled manual transcription;
all source and artifact hashes describe the published bytes.

| Context | Recorded result | Representative probe IDs |
| --- | --- | --- |
| First non-whitespace character | A single `/`, `*`, `}`, or `*/` ignores the rest of its line; indentation does not change the slash/star cases. | `single-slash-comment`, `single-star-comment`, `closing-brace-comment`, `closing-c-comment`, `indented-slash-line`, `indented-star-line` |
| After a complete statement | Braces ignore the remaining line, including another statement. A single slash/star or `/* ... */` is rejected there. | `brace-inline-trailing-statement`, `closing-brace-inline`, `single-slash-inline`, `single-star-inline`, `c-inline-without-statement`, `c-comment-after-statement` |
| Quoted delimiters | `/`, `/*`, `*/`, `{`, and `}` remain literal. `//` truncates the line even inside a quoted value and produces an error. | `string-single-slash`, `string-c-opening`, `string-c-closing`, `string-brace-opening`, `string-brace-closing`, `string-brace-content`, `string-double-slash`, `string-comment-after-content`, `comment-symbols-in-string` |
| Incomplete expressions | `escreval(1 { note } + 2)` reports a missing `)`, while `escreval(1{ note })` completes with no program output. A C-like opener within the tested expression reports an undeclared `NOTE`. These exact observations remain pending implementation and do not establish a general recovery rule. | `brace-comment-in-expression`, `brace-inline-math`, `c-comment-in-expression` |

## Isolated call follow-up

Fourteen call and name-resolution probes were run with a fresh reference
process for each program. Eleven supply usable evidence: four accepted programs
and seven visible language rejections. Three temporary-reference-argument cases
caused internal application faults and remain inconclusive; their raw captures
are not published or treated as language rejection evidence.

| Context | Recorded result | Probe IDs |
| --- | --- | --- |
| Repeated bare calls | Two occurrences of a counter function execute independently and print successive values. | `bare-function-side-effects` |
| Function names and variables | Bare and parenthesized function values take priority over a colliding local variable. Bare function values also take priority over a colliding parameter; assignments to the local variable remain accepted. | `local-variable-shadows-function`, `function-name-priority`, `function-name-priority-over-parameter` |
| Missing function arguments | A bare function that needs arguments is rejected at its use, including when a local variable has the same name. | `bare-function-missing-argument`, `function-name-priority-with-arguments` |
| Invalid call context | A variable alone or followed by `()` is rejected. A procedure used as a value is rejected. The bare-variable diagnostic has a visible line number but no message text. | `bare-variable-statement`, `variable-call-statement`, `procedure-value-expression` |
| Procedure errors | A bare call missing required arguments and a procedure name colliding with a local variable are rejected with line 2 reported in these programs. Those error positions remain pending implementation. | `bare-procedure-missing-argument`, `procedure-name-priority` |

The declaration and call observations already recorded earlier are now verified
alongside these focused cases where implementation is marked verified. Expected
call-context errors use the project's `E004` category; GUI line numbers remain
unchanged. Each published rejection screenshot was inspected individually, and
all artifact hashes describe the unchanged published bytes.
