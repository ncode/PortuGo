# Bundled example progress

The corpus inventories 73 original programs from the official VisuAlg 3.0.7
distribution: 60 accepted, 11 unusable on their recorded paths, and two reviewed
nonterminating exclusions. Their source bytes match the catalog sizes and
SHA-256 hashes. Each observation links its catalog ID, exact source, reference
outcome, and owning tasks in the conformance manifest.

The evidence gate checks catalog IDs, source hashes and sizes, classifications,
recorded acceptance, and reviewed exclusion reasons. Pending entries still fail
the missing-evidence gate; classifications do not grant them an exemption.

All 60 accepted programs have verified original and formatted execution.
Forty-five use byte-for-byte output comparisons; fifteen use reviewed generated-output
contracts that check domains and framing without comparing random sequences.
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
transcripts. Their original and formatted executions now pass the reviewed
generated-output contracts; these observations do not establish exact portable
random sequences or complete branch coverage.

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
The four remaining rejection mappings now match exact original-source CLI
replay. The [November 6 qualification](rejection-paths-progress.md) retains the
previous evidence and separately records the complete 82-byte assignment and
two-byte repeat prefixes. The character table stops at P001 line 301 after its
152-byte banner; the school exit preserves its 41-byte prompt/echo and reports
one P001 at EOF line 320. All eleven rejected example paths are verified.
The explicit `TestRecordedRejectedExamples` gate requires each of these four
mappings to be verified and compares its exit status, complete stdout bytes,
and diagnostic sequence/code/line; accepted-example success cannot satisfy it.

Generated-output examples retain their reviewed contracts. The `randomicos.alg.ALG` transcript (source algorithm `semnome`) now has
a reviewed domain/framing contract for its nine `randi(10)` values and fixed
zero tail; the mixed
integer/text, integer-sort, real-sort, record-sort, record-search,
integer-search, raw integer-search, sorted integer-search, counting-sort,
repeated-value, and relation transcripts now have reviewed shape-specific
contracts. The menu exit now retains its recorded output while deferring the
unexecuted body diagnostic until invocation.
Successful runs can produce different generated values. Local completion alone
does not qualify their output or every branch.

The school example's incomplete initial capture remains excluded. Fresh
complete runs and reduced reached/unreached controls qualify its recovered EOF
path. These results apply to the recorded input, not every interactive branch.

The catalog has sixty accepted examples, eleven unusable examples, and two
reviewed non-goal examples. All 73 catalog entries now have a recorded
disposition; incomplete captures are excluded. The accepted-example sweep and
evidence inventory are complete within their recorded scope. Release acceptance
still requires fresh qualification of the exact candidate; these recordings do
not establish every interactive branch or full-language compatibility.
