# Changelog

## Unreleased

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
