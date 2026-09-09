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
