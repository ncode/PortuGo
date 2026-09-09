# Bundled example progress

The first completed sweep records eleven original programs from the official
VisuAlg 3.0.7 distribution. Their source bytes match the sizes and SHA-256 hashes
in the 73-file catalog. Each observation links its catalog ID, exact source,
reference outcome, and owning tasks in the conformance manifest.

The evidence gate now checks catalog IDs, source hashes and sizes, classifications,
recorded acceptance, and reviewed exclusion reasons. Regression tests reproduce
the previously unchecked mismatches. Pending entries still fail the existing
missing-evidence gate; the new checks do not grant them an exemption.

Eight programs complete successfully in the reference. Six also match the CLI
byte for byte and have permanent tests for original and formatted execution:
`ajustes.alg`, `passo.alg`, `taxaspop.alg`, `Exemplos/TESTE.ALG`, `troca.alg`, and
`vetr2dim.alg`. The accepted `caracfun.alg` and `randomicos.alg` remain pending:
they expose missing string-type/built-in support and `randi`, respectively.
The latter produces no output; that recording establishes successful execution,
not an exact random sequence or a complete randomness contract.

Three bundled files are unusable as supplied. Their ordinary GUI diagnostics
were individually reviewed, and the corpus contains screenshots and labeled
manual transcriptions:

| Original file | Reference rejection |
| --- | --- |
| `interrompa.alg` | `fimrepita` produces a syntax error on line 16. |
| `UDF.ALG` | An invalid function name is reported on line 34. |
| `TESTE.alg` in the distribution root | `ATÉ_QUE` is reported as an unknown variable on line 23. |

These rejection mappings remain pending implementation. The reference may
execute a prefix before discovering malformed syntax; retained partial panel
output documents that behavior. CLI compile-time rejection remains a preflight
check that does not execute an invalid program.

The keyword follow-up adds a completed `PRIMOS.ALG` recording with a longer,
bounded observation window. Its original source hash matches the release
catalog. Execution remains pending in the CLI because the example also uses
bare `escreval`; the keyword aliases alone do not complete its grammar support.

The catalog currently has nine accepted examples, three unusable examples,
and 61 awaiting committed classifications. Incomplete captures are excluded.
This slice does not close the full
example sweep, evidence inventory, or conformance release gate.
