# Generated-input replay contracts

Seventeen recorded `aleatorio` programs print a generated value twice: first
the `leia` echo, then `escreval(value)`. Their values vary between executions.
The original source, raw observation, normalized observation and hashes remain
unchanged. Replay qualifies the documented domains without promising matching
reference seeds, sequences, distributions or generator-consumption counts.

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
| `random-value-single-bound` | 7..100 integers | 0 | 8 |
| `random-value-single-negative` | -3..100 integers | 0 | 8 |
| `random-range-text` | Five uppercase ASCII letters | — | 1 |

`TestRecordedRandomInputDomains` checks all seventeen original and formatted
programs using both forced generator endpoints and sixteen deterministic seeds.
Contract tests reject wrong types, bounds, precision, repeated values, spacing,
line counts, invalid specifications, and attempts to replace recorded evidence.
The ordinary subprocess replay uses the production CLI and its own generator.
