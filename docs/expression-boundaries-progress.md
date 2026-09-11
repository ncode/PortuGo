# Expression boundary recordings

Sixteen recordings cover logical precedence, comparison grouping, assignment
suffixes, and mixed arithmetic/logical reductions. All sixteen now match the
implementation and have regression tests; the four value-reduction cases were
completed by the [retained operand slice](retained-operands-progress.md).

Logical operators bind more tightly than comparisons. `ou` and `xou` share a
precedence level and associate left to right. Output rejects an ungrouped second
comparison with `P001`; the canonical printer preserves grouped comparisons.
Assignments ignore the remaining tokens after their completed expression,
including a parenthesized suffix or another statement on the same physical line.

The existing ungrouped integer random-range example still matches its recording
after the parser correction. Generator bounds are tested separately; that
example's ignored final comparison is not evidence for its upper-bound check.

The value cases concern propagation when arithmetic consumes an `e`
result. One direct numeric addition produces `10` in the reference; two additions
with retained comparison values produce `P001`; a grouped multiplication retains
`7`. Their original outputs and diagnostics are recorded and verified. These
recordings do not establish complete expression conformance.

Only synthetic source, program-only output panels, reviewed diagnostic-window
crops, and labeled manual transcriptions are published. Full desktop captures and
operational details remain private.
