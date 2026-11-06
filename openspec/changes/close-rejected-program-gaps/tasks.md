## 1. Establish the reference contract

- [ ] 1.1 Run a fresh clean test baseline and add bounded subprocess regressions for the four unchanged corpus sources, asserting exact stdout, diagnostic sequence/code/line, and exit status; observe the current mismatches, and preserve the red implementation cases after reference expectations are qualified.
- [ ] 1.2 Record reduced assignment-versus-later-syntax controls and a fresh original `EXEMPLO1.alg` run, establishing reached-error timing, target/RHS ordering, and the complete prompt/input prefix.
- [ ] 1.3 Record reduced sibling/nested user-call controls and a fresh original `Tabela_ASCII4.alg` run, reducing the reproduced line-301 trigger without rejecting the successful simple-call and precomputed-argument controls.
- [ ] 1.4 Record unused-versus-called malformed subprogram controls and a fresh `ESCOLA1.alg` run with `SAIR`, establishing the prompt/echo prefix and reached EOF error without assuming that invoked malformed bodies fail at their declaration boundary.
- [ ] 1.5 Record reached-versus-escaped repeat markers, newline-versus-no-newline writes, successful final flushing, and a fresh original `interrompa.alg` run.
- [ ] 1.6 Review and privacy-screen all new text, screenshots, and metadata; add pinned provenance and traceability through the existing recorder. Preserve old evidence. Qualify the fresh 82-byte assignment and two-byte repeat prefixes, add separately reviewed expectation artifacts, and resolve conflicts explicitly before changing mappings. Update/review this proposal if its observable contracts change. Block dependent implementation for every unresolved case.

## 2. Correct assignment diagnostic precedence

- [ ] 2.1 Add the reduced assignment-order and neighboring valid/error controls from 1.2; observe the applicable regressions fail before production changes.
- [ ] 2.2 Repair the shared analysis/recovery seam established by 1.2, preserve preceding I/O and evaluation order, unchanged static check diagnostics, unrelated diagnostic collection, and non-execution after blocking errors, and pass the original EXEMPLO1 regression plus its controls.

## 3. Correct the character-table rejection

- [ ] 3.1 Add the reduced trigger and neighboring successful call controls established by 1.3; observe the applicable regressions fail before production changes.
- [ ] 3.2 Apply the smallest confirmed shared call/evaluation fix and pass the original table regression with its exact 152-byte prefix, P001 line 301, and neighboring controls; add no source-identity special cases.

## 4. Correct the school exit-path rejection

- [ ] 4.1 Add reached/unreached malformed-body and EOF recovery controls established by 1.4, including nested boundary and resource-limit cases; observe the applicable regressions fail before production changes.
- [ ] 4.2 Repair only the evidenced recovery/deferred-diagnostic seam, preserving progress and positions, then pass the original school regression with its 41-byte prompt/echo and single EOF diagnostic.

## 5. Correct the incomplete repeat output expectation

- [ ] 5.1 Add complete-output evidence/data regressions from 1.5 and the known accepted escape control; demonstrate that the old empty-output expectation conflicts with the reference while the baseline CLI preserves the observed two bytes.
- [ ] 5.2 Correct the repeat expected-output mapping to the separately reviewed two-byte artifact, retain old evidence, and pass its exact original-source regression. Preserve the verified unconditional-break escape; make no interpreter/output change without a separately confirmed mismatch.

## 6. Integrate and qualify all four changes

- [ ] 6.1 Verify analysis-only commands perform no runtime side effects, formatting preserves deferred behavior or rejects in-place changes without modifying bytes, and a valid REPL submission after each rejection starts cleanly. Audit affected CLI, REPL, printer, and conformance-adapter callers.
- [ ] 6.2 Add the explicit four-probe completion gate with a regression proving a pending or mismatching named probe fails despite a green accepted-example sweep. Preserve historical task IDs, add supported new test/task links, and mark each probe verified only after its exact original replay passes.
- [ ] 6.3 Update docs/language.md, CHANGELOG.md, and current corpus/example progress from validated outcomes. Verify all existing source/input/evidence hashes remain unchanged and no pending result or excluded run is relabeled as success.
- [ ] 6.4 Run strict OpenSpec validation, build, tests, race tests, vet, pinned staticcheck and golangci-lint, 30-second lexer/parser fuzzing, full conformance evidence/trace/replay and the explicit four-case gate. Privacy-screen and run the required OCR review on the exact implementation diff, address confirmed findings, disclose exclusions, and validate the implementation PR's required CI checks before handoff.
