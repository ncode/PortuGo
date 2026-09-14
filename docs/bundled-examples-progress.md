# Bundled example progress

The corpus records 71 original programs from the official VisuAlg 3.0.7
distribution. Their source bytes match the sizes and SHA-256 hashes in the
73-file catalog. Each observation links its catalog ID, exact source, reference
outcome, and owning tasks in the conformance manifest.

The evidence gate checks catalog IDs, source hashes and sizes, classifications,
recorded acceptance, and reviewed exclusion reasons. Pending entries still fail
the missing-evidence gate; classifications do not grant them an exemption.

Sixty-one programs complete successfully in the reference. Forty-five match
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
`RELACIONAR.ALG` now have completed reference recordings. The former has a
reviewed domain/framing contract; the latter now has a reviewed relation
branch/framing contract without a portable exact sequence claim.

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

The latest recordings add `bin2dec.alg`, `buscaseq.alg`, `buscaseqreg.alg`,
`buscbinr.alg`, `DESTAQUES.ALG`, `SEMNOME.ALG`, `MENU_PRINCIPAL.alg`, and
`REGISTROS.ALG`. The original `REGISTROS.ALG` exit path and the menu exit path
now have exact original and formatted replay coverage. No generated-data
transcript remains pending a random-output comparison contract. The menu
diagnostic is deferred for uncalled subprogram bodies and remains available
when an affected subprogram is invoked.

Eleven bundled files are unusable on their recorded paths. Their GUI diagnostics
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
| `Jogo_velha.alg` | The `LITERAL` type is not recognized on line 9. |
| `EXTENSO4.ALG` | The `LITERAL` vector element type is not recognized on line 8. |
| `ESCOLA1.alg` | After input `SAIR`, a missing `FIMSE` is reported at EOF on line 320. The preceding prompt and echoed input are retained. |

The two malformed headers, the three unrecognized-type forms, and the invalid
function name now have matching positioned parser regressions. The game example
now also has a matching type rejection.
The bundled compound-repeat spelling now also has a matching `E002` on line 23,
before any output. Original and formatted source retain the unknown-identifier
rejection; formatting does not make the spelling an accepted terminator.
The other four rejection mappings remain pending. The reference may execute a
prefix before discovering malformed syntax or an invalid assignment, while CLI
analysis reports errors before execution and may collect more than one
diagnostic.

The latest bounded grammar audit leaves four rejected recordings pending; no
accepted recording remains pending. These observations do not justify adding a
default loop bound or generally suppressing diagnostics in unexecuted
subprograms.

A replay audit after the lazy body diagnostic fix found no additional byte-exact
match. The `randomicos.alg.ALG` transcript (source algorithm `semnome`) now has
a reviewed domain/framing contract for its nine `randi(10)` values and fixed
zero tail; the mixed
integer/text, integer-sort, real-sort, record-sort, record-search,
integer-search, raw integer-search, sorted integer-search, counting-sort,
repeated-value, and relation transcripts now have reviewed shape-specific
contracts. The menu exit now retains its recorded output while deferring the
unexecuted body diagnostic until invocation.
Successful runs still produce different generated values, while a small number
of paths stop during analysis or execution. Local completion alone does not
qualify their output or every branch.

The four rejected paths retain their recorded diagnostics and remain pending;
the accepted menu path now completes in the CLI with the recorded exit input.

The CLI currently completes `Tabela_ASCII4.alg`, unlike the reference. Its
recording therefore remains pending; the retained diagnostic and preceding output
record the reference outcome without claiming an implementation match.

The school example's initial capture was waiting for input and was excluded;
the later capture records the displayed EOF diagnostic. Parser recovery now
keeps the source's EOF position after consuming the final token, instead of
moving subsequent missing-delimiter diagnostics to line 1. This corrects
position loss without claiming to match the school's earlier parsing behavior.

The catalog has sixty accepted examples, eleven unusable examples, and two
reviewed non-goal examples. All 73 catalog entries now have a recorded
disposition; incomplete captures are excluded. These records do not complete
the full example sweep, evidence inventory, or conformance release gate.
