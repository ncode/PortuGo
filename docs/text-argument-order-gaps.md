# Text argument-order gaps

Two existing successful reference recordings still differ from execution:
`text-order-no-value-bound` omits the final substring after its nested calls,
and `text-order-no-value-code` omits both the pending output and a later `DONE`.
Independent reruns reproduce both results.

Four further recordings separate ordinary zero conversion from the missing
execution-state behavior. Each contains the complete reference execution
notices, uses `panel-v1` normalization, and retains unchanged source/output
bytes with their hashes. Only source and output text are published.

| Probe | Recorded behavior | Implementation |
| --- | --- | --- |
| `text-state-code-zero-call` | Explicit `carac(0)` retains pending text and later `DONE` around a nested function write. | Verified |
| `text-state-copy-absent-earlier-call` | A function call before an absent substring bound retains the substring and later `DONE`. | Verified |
| `text-state-code-absent-stored` | Assigning `carac(abs())` before the output statement still ends after the nested function's write. | Pending |
| `text-state-code-absent-tail-error` | The reference reports normal completion without reaching later `DONE` or the deliberate `raizq(-1)` domain error. | Pending |

`TestRecordedTextFunctions` checks original and formatted execution of the two
controls. The termination cases remain pending: discarding only the current
output buffer would not account for the missing later output and error. The
recordings do not yet establish the lifetime or reset rules of the reference
state, so they do not justify a runtime special case.

Two separate rejection recordings also remain pending. `text-order-bad-bound`
writes `1` before rejecting the second argument with `E001`;
`text-order-extra` writes `1`, `2`, and `3` before rejecting the extra argument
with `P001`, without evaluating that argument. Semantic analysis currently
reports both diagnostics before execution. Matching these prefixes requires
an explicit design for preserving bindings and deferring the relevant errors
until their execution boundary. Suppressing semantic diagnostics alone would
leave incomplete analysis metadata.

These findings do not change the completed indexing/case-conversion task or
complete any additional group 12 task. Earlier text-progress counts describe
their original recording slice; the conformance manifest is the current
source of implementation states.
