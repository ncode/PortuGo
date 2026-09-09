# Bundled example progress

The corpus records 43 original programs from the official VisuAlg 3.0.7
distribution. Their source bytes match the sizes and SHA-256 hashes in the
73-file catalog. Each observation links its catalog ID, exact source, reference
outcome, and owning tasks in the conformance manifest.

The evidence gate checks catalog IDs, source hashes and sizes, classifications,
recorded acceptance, and reviewed exclusion reasons. Pending entries still fail
the missing-evidence gate; classifications do not grant them an exemption.

Thirty-six programs complete successfully in the reference. Thirty-four match
the CLI byte for byte and have permanent tests for original and formatted execution.
The fixed-input cases cover combinations, factorials, minimum selection, means,
prime decomposition, base conversion, vector sorting, remainders, reversed text,
perfect numbers, square-root approximation, and text choices. Each input is
retained alongside the original source and recorded output.

The original `ENCRYPT.ALG` and `mediaar.alg` programs now match their fixed-input
recordings with the accented `fimfunção` terminator and trailing declaration
semicolon. The accepted `escolha.alg` program now also matches its recording
with inclusive `ate` case ranges. A bundled-file classification alone does not
claim CLI support.

The display slice also verifies `CARACOL.ALG`, `Caracol2.ALG`, and `graus.alg`
with their original screen-clear and color commands. `randomicos.alg.ALG` and
`RELACIONAR.ALG` now have completed reference recordings, but their random output
remains pending a domain-based replay contract. Their recorded samples do not
establish portable exact sequences.

The text slice adds `TABOADA.ALG`, which uses screen clearing, character input,
and a counted multiplication-table loop. Its original and formatted executions
match the reference with the recorded fixed input.

The `randomicos.alg` example produces no output; its recording establishes
successful execution, not an exact random sequence or a complete randomness
contract. The `PRIMOS.ALG` prime search has an explicit 5,000,000-step replay
allowance for its bounded nested loops; the five-second replay deadline still
applies.

Seven bundled files are unusable as supplied. Their ordinary GUI diagnostics
were individually reviewed; the corpus retains diagnostic evidence and labeled
manual transcriptions:

| Original file | Reference rejection |
| --- | --- |
| `interrompa.alg` | `fimrepita` produces a syntax error on line 16. |
| `UDF.ALG` | An invalid function name is reported on line 34. |
| `TESTE.alg` in the distribution root | `ATÉ_QUE` is reported as an unknown variable on line 23. |
| `Calculo_media2.alg` | The malformed `Agoritmo` header is rejected on line 1. |
| `Calculo_media2.alg.ALG` | The malformed `Agoritmo` header is rejected on line 1. |
| `estcivil.alg` | The `NUMERICO` type is not recognized on line 13. |
| `EXEMPLO1.alg` | Assignment from integer to a character vector element fails on line 21. |

The two malformed headers and the unrecognized type now have matching positioned
parser regressions. The other four rejection mappings remain pending. The
reference may execute a prefix before discovering malformed syntax or an invalid
assignment, while CLI analysis reports errors before execution and may collect
more than one diagnostic.

The catalog has thirty-six accepted examples, seven unusable examples, and 30
awaiting recorded classifications. Incomplete captures are excluded. These
records do not complete the full example sweep, evidence inventory, or
conformance release gate, which still has 41 missing mappings: 11 requirements
and 30 bundled examples.
