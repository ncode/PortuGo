# Bundled example progress

The corpus records twelve original programs from the official
VisuAlg 3.0.7 distribution. Their source bytes match the sizes and SHA-256 hashes
in the 73-file catalog. Each observation links its catalog ID, exact source,
reference outcome, and owning tasks in the conformance manifest.

The evidence gate now checks catalog IDs, source hashes and sizes, classifications,
recorded acceptance, and reviewed exclusion reasons. Regression tests reproduce
the previously unchecked mismatches. Pending entries still fail the existing
missing-evidence gate; the new checks do not grant them an exemption.

Nine programs complete successfully in the reference. Eight also match the CLI
byte for byte and have permanent tests for original and formatted execution:
`ajustes.alg`, `passo.alg`, `taxaspop.alg`, `Exemplos/TESTE.ALG`, `troca.alg`,
`vetr2dim.alg`, `PRIMOS.ALG`, and `caracfun.alg`. The accepted `randomicos.alg`
remains pending because it needs `randi`.
It produces no output; that recording establishes successful execution,
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
execute a prefix before discovering malformed syntax; private partial panel
observations document that behavior. CLI compile-time rejection remains a preflight
check that does not execute an invalid program.

The `PRIMOS.ALG` prime search has an explicit 5,000,000-step replay allowance for its bounded
nested loops; the five-second replay deadline still applies.

The catalog currently has nine accepted examples, three unusable examples,
and 61 awaiting committed classifications. Incomplete captures are excluded.
These records do not complete the full example sweep, evidence inventory,
or conformance release gate.
