# Random replay contracts

Nineteen recorded `aleatorio` programs print a generated value twice: first
the `leia` echo, then `escreval(value)`. Their values vary between executions.
The original source, raw observation, normalized observation and hashes remain
unchanged. Replay qualifies the documented domains without promising matching
reference seeds, sequences, distributions or generator-consumption counts.

The bundled `randomicos.alg.ALG` example (source algorithm `semnome`, probe
`bundled-991ec2bd1566`) has a separate `randomOutput` contract. It emits ten
newline-terminated integer lines with one leading space; the first nine values
are `randi(10)` results in `[0, 10)`, and the final unassigned vector slot remains
zero. Replay checks that framing and domain for the original and formatted
source while continuing to ignore the reference generator sequence.

One bundled mixed-input example has a separate `randomOutput` contract. It
emits ten lines in five alternating pairs: a canonical integer in `[0, 101)`
followed by five uppercase ASCII letters. Replay checks the line shape,
integer domain and text alphabet for the original and formatted source without
promising the generated sequence.

One bundled integer-sort example has a separate `randomOutput` contract. It
emits an unconstrained sequence of bounded integers followed by the same values
in nondecreasing order, using the documented output spacing. Replay checks the
input domain, sorted permutation and line framing for the original and formatted
source without promising the generated sequence.

Two bundled real-sort examples have separate `randomOutput` contracts. Each
emits bounded real input with a fixed decimal grid followed by the same values
in nondecreasing order, rendered in the documented numbered-row format. Replay
checks the input grid, sorted permutation and row framing for the original and
formatted source without promising the generated sequence.

The optional `implementation.expected.randomInput` contract is restricted to
accepted group-13 recordings and requires a review reason and document link.
`minimum` and `maximum` are inclusive integer ticks at `10^-decimals`; for
example, 2000 through 2999 with three decimals means 2 through 2.999. These
bounds follow the recorded language rules, rather than the sampled minimum and
maximum. Integer contracts require zero decimals; text contracts require zero
numeric bounds and exactly five uppercase ASCII letters. The numeric contract
profile is limited to signed-32-bit values and at most five fractional places.

Every qualified transcript must contain exactly `samples` echo/output pairs
and a final LF. Numeric echoes have no leading space; real echoes have exactly
ten fractional digits. Numeric `escreval` lines have one leading space and the
documented default number representation. Both lines must contain the same
value. Text pairs must match exactly. Real values must lie on the exact decimal
grid; no floating-point tolerance or whitespace normalization is used.

The validator checks each contract against the retained reference observation
and still requires the expected stdout artifact to equal that observation.
Replay checks the actual output against the same contract. Exit status,
diagnostics, files, state and host observations retain their existing exact
comparisons. All probes without this explicit contract retain byte-exact stdout
comparison. An arbitrary regular expression or ignored output channel is not
supported.

| Recordings | Inclusive domain | Fractional places | Pairs |
|---|---|---:|---:|
| `random-input-fixed-real`, `random-value-real-grid` | 2..2.999 | 3 | 1, 8 |
| `random-range-fraction-precision` | 2.75..3.749 | 3 | 1 |
| `random-range-fractional-precision` | 2..2.9 | 1 | 1 |
| `random-range-large-precision`, `random-value-precision-six`, `random-value-precision-seven`, `random-value-precision-large` | 2..2.99999 | 5 | 1, 8, 8, 8 |
| `random-range-zero-bound`, `random-range-on-reset`, `random-value-on-integer` | 0..100 integers | 0 | 1, 1, 8 |
| `random-value-on-real` | 0..100 whole reals | 0 | 8 |
| `random-range-expression-bound` | 7..8 integers | 0 | 1 |
| `random-range-reversed` | 3..9 integers | 0 | 1 |
| `random-boundary-full-signed32`, `random-boundary-full-signed32-reversed` | -2147483648..2147483647 integers | 0 | 1, 1 |
| `random-value-single-bound` | 7..100 integers | 0 | 8 |
| `random-value-single-negative` | -3..100 integers | 0 | 8 |
| `random-range-text` | Five uppercase ASCII letters | — | 1 |

`TestRecordedRandomInputDomains` checks all nineteen original and formatted
programs using both forced generator endpoints and sixteen deterministic seeds.
Contract tests reject wrong types, bounds, precision, echo/output disagreement,
spacing, line counts, invalid specifications, and attempts to replace recorded
evidence. The bundled `randomicos.alg.ALG` contract is replayed in
`TestRecordedBundledRandiOutput` for both source forms, forced endpoints, and
sixteen deterministic seeds. The mixed-input contract is replayed by
`TestRecordedBundledMixedRandomOutput`, and the integer-sort contract by
`TestRecordedBundledRandomSortOutput`, with the same source and generator
coverage. The real-sort contracts are replayed by
`TestRecordedBundledRandomRealSortOutput`. Bundled programs with generated arrays, records, sorting, branching,
or other mixed input/output remain pending until their own observable contracts
exist.
The ordinary subprocess replay uses the production CLI and its own generator.
