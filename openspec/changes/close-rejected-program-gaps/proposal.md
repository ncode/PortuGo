## Why

Four rejected-program comparisons remain pending, although every accepted recording is verified. Closing these differences needs precise error selection, rejection timing, and partial-output behavior without regressing accepted programs or weakening the corpus.

## What Changes

- Reproduce and reduce the four pending rejection paths against the pinned reference application before implementing behavior not established by existing evidence.
- Match the reviewed exit status, diagnostic code and line, and complete standard-output bytes for `EXEMPLO1.alg`, `Tabela_ASCII4.alg`, `ESCOLA1.alg`, and `interrompa.alg`.
- Correct the incomplete output expectations identified by fresh reference runs: the assignment example prints an 82-byte prompt/input prefix, and the repeat example writes ` 1` before failing. Preserve the old observations and qualify new evidence before changing corpus expectations.
- Keep blocking front-end errors non-executable. Classify errors as deferred only where reduced reference controls establish their execution timing, with bounded recovery and no source-name or source-hash special cases.
- Add permanent CLI regressions for all four original inputs and reduced positive/negative controls. Preserve original source, stdin, and evidence bytes; add repeat observations separately.
- Make these four rejected paths required for this change's acceptance, then mark them verified only after original-source replay and traceability pass. Update the language reference, changelog, and current conformance progress from validated results.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `tooling/diagnostics-cli-repl`: Specify the four recorded rejection outcomes and the distinction between blocking front-end diagnostics and evidence-backed deferred diagnostics. Keep `check`, `fmt`, REPL recovery, and resource safeguards explicit.
- `conformance/visualg-3-0-7`: Require replay and test traceability for the four pending rejected examples before this change can be accepted, independently of the accepted-example release sweep.

## Impact

Likely implementation seams are parser recovery, semantic diagnostic classification, interpreter user-call state, CLI integration, and the conformance replay gate. The reference experiment must determine which seams actually need changes. Tests, corpus metadata, `docs/language.md`, `CHANGELOG.md`, and current progress documents change alongside confirmed implementation work. No core dependencies, new language dialects, release publication, or existing release tags are in scope.

The fresh runs distinguish three remaining implementation gaps from an incomplete repeat-output expectation. The character-table trigger still needs further reduction. Inconclusive or conflicting reference observations keep the affected task pending; they do not authorize an invented compatibility rule or an implementation-verification claim.
