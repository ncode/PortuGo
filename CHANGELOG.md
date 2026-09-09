# Changelog

## Unreleased

- Linked existing corpus-tooling tests and recorded feature dispositions to their
  specification requirements, with reviewed reasons for project-only evidence.
  Language recordings, bundled examples, and full acceptance remain incomplete.
- Preserve conformance history for pending probes, inventory entries, linked
  tasks, and prior retirements. Removed IDs require reviewed dispositions;
  replacements may identify an active probe, inventory entry, or task.
- Made corpus traceability discover every specification in the change's spec tree.
  Added retained GUI error evidence and labeled transcriptions for the three
  earlier rejection probes; the broader evidence inventory remains incomplete.
- Fixed corpus gates to enforce recorded acceptance/rejection and complete
  generated-file inventories. Validation now requires a history comparison, and
  CI supplies the base commit. Recording rejects altered auxiliary input files
  while permitting declared outputs to change.
- Added reference-corpus validation, bounded CLI replay, and a private staging
  and capture workflow. The initial corpus includes an inventory of requirements
  and bundled examples, earlier observations, and fresh declaration probes.
  Evidence collection and complete conformance remain unfinished.
- Fixed the Windows-recorded VisuAlg gaps: failed `run` validation now exits 1;
  `exp` takes two numeric arguments; power is left-associative and binds below
  unary minus; `para` exposes the recorded exit values, including empty and
  interrupted loops. Output now uses the recorded deterministic dot-decimal
  profile, numeric/logical spacing, uppercase logical values, decimal rounding,
  width-zero behavior, and string alignment/truncation; logical field widths
  are rejected. Added regression fixtures and replay of 14 preserved reference
  observations. Other locales and full conformance remain unverified.
- Fixed Windows runtime fixture failures by preserving test data bytes through
  Git checkout and staging, regardless of automatic newline conversion settings.
- Added CI build, formatting, vet, pinned lint, race, three-platform Go tests,
  lexer/parser fuzzing, strict OpenSpec validation, and an explicit pending
  oracle gate. Added byte-exact golden helpers and failing adversarial-test
  watchdogs; source/runtime behavior is unchanged.
- Added an initial VisuAlg 3.x implementation with lexer, parser, semantic checks, tree-walking interpreter, CLI, REPL, examples, and regression fixtures.
