# Retained operand recordings

Sixteen new recordings and four previously pending cases cover intermediate
operands retained by mixed logical and comparison reductions. All twenty are
verified: sixteen accepted original/formatted executions and four positioned
syntax failures.

Examples include `2 + (3 e 7) = 10`, `20 - (3 e 7) = -4`, and
`2 * (3 e 7) = 21`. Two nested reductions in `2 + ((3 e 7) + (4 e 9))` produce
`20`. The effect persists through a numeric builtin and a language function;
ordinary arithmetic and existing random-domain recordings still match.

Comparisons retain their own left value, while enclosing arithmetic can consume
an operand left by a nested mixed reduction. Incompatible retained addition
operands produce `P001`; the recorded mixed multiplication retains its right
value. These are compatibility behaviors, not general arithmetic identities.

An enclosing evaluation retains at most 65,536 operands. This is a project
storage safeguard with positioned `R003`, not a measured reference quota.
Regression tests cover the limit, independent expressions, reuse after failure
and release of references held by the operand buffer.

The complete operand-type matrix remains pending. Published evidence contains
synthetic source, program-only output panels, two reviewed diagnostic crops and
labeled manual transcriptions. Operational metadata stays private.
