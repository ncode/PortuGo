# Text argument-order recordings

The six previously pending text-call controls now replay successfully in both
original and formatted sources. They cover left-to-right output before a
deferred rejection, arity rejection without evaluating an extra argument, and
the no-value state carried through direct and stored character conversions.

| Probe | Recorded behavior | Implementation |
| --- | --- | --- |
| `text-order-no-value-bound` | `Texto(1)` and `Numero(3)` write before the absent bound suppresses the outer substring. | Verified |
| `text-order-bad-bound` | `Texto(1)` writes before the invalid bound receives `E001`. | Verified |
| `text-order-extra` | The first three arguments write before the extra argument receives `P001`; the extra call is not evaluated. | Verified |
| `text-order-no-value-code` | `Texto(9)` writes, then the pending output and later `DONE` are skipped. | Verified |
| `text-state-code-zero-call` | Explicit `carac(0)` retains pending text and later `DONE` around a nested function write. | Verified |
| `text-state-copy-absent-earlier-call` | A function call before an absent substring bound retains the substring and later `DONE`. | Verified |
| `text-state-code-absent-stored` | Assigning `carac(abs())` preserves the state for later output with `Texto(9)`. | Verified |
| `text-state-code-absent-tail-error` | The state skips later `DONE` and the deliberate `raizq(-1)` failure, completing normally. | Verified |

`TestRecordedTextArgumentOrderGaps` checks the six formerly pending controls in
both original and formatted sources; the two earlier controls remain covered by
`TestRecordedTextFunctions`. The state is per interpreter run, and pending
output remains bounded by the existing output limits.
