# Bundled example progress

The corpus records 60 original programs from the official VisuAlg 3.0.7
distribution. Their source bytes match the sizes and SHA-256 hashes in the
73-file catalog. Each observation links its catalog ID, exact source, reference
outcome, and owning tasks in the conformance manifest.

The evidence gate checks catalog IDs, source hashes and sizes, classifications,
recorded acceptance, and reviewed exclusion reasons. Pending entries still fail
the missing-evidence gate; classifications do not grant them an exemption.

Fifty-two programs complete successfully in the reference. Forty-two match
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

The latest fixed-input recordings add `CALENDARIO.ALG`, `CHECA_CPF.alg`,
`decpoutrasproc.alg`, `ELEMENTO_OCUPADO.alg`, `MEDIA_REGISTRO.ALG`,
`MENU_COM_CASE.alg`, `Numeros_primos.alg`, and `REGISTROS com VETORES.ALG`.
The menu recording exercises all four actions and exits; the vector-registration
recording exercises its exit path. These are recorded input paths, not a claim
that every interactive branch has been tested.

Random-input recordings add `aleatorio1.alg`, `bbsort.alg`, `bbsortp.alg`,
`bbsortreg.alg`, `bubblrec.alg`, `buscabin.alg`, `cntsort.alg`, and `inssort.alg`.
The binary-search recording supplies `-1` to exit after generating its data.
All eight complete in the reference, including the original incomplete-looking
loop header in `aleatorio1.alg`. Their generated values remain in the recorded
transcripts. Implementation verification remains pending; these observations
do not establish exact portable random sequences or complete branch coverage.

Eight bundled files are unusable as supplied. Their ordinary GUI diagnostics
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
| `Tabela_ASCII4.alg` | A syntax diagnostic on line 301 stops the original character-table program. |

The two malformed headers and the unrecognized type now have matching positioned
parser regressions. The other five rejection mappings remain pending. The
reference may execute a prefix before discovering malformed syntax or an invalid
assignment, while CLI analysis reports errors before execution and may collect
more than one diagnostic.

The CLI currently completes `Tabela_ASCII4.alg`, unlike the reference. Its
recording therefore remains pending; the retained diagnostic and preceding output
record the reference outcome without claiming an implementation match.

The catalog has fifty-two accepted examples, eight unusable examples, and 13
awaiting recorded classifications. Incomplete captures are excluded. These
records do not complete the full example sweep, evidence inventory, or
conformance release gate.
