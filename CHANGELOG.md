# Changelog

## Unreleased

- Fixed Windows runtime fixture failures by preserving test data bytes through
  Git checkout and staging, regardless of automatic newline conversion settings.
- Added CI build, formatting, vet, pinned lint, race, three-platform Go tests,
  lexer/parser fuzzing, strict OpenSpec validation, and an explicit pending
  oracle gate. Added byte-exact golden helpers and failing adversarial-test
  watchdogs; source/runtime behavior is unchanged.
- Added an initial VisuAlg 3.x implementation with lexer, parser, semantic checks, tree-walking interpreter, CLI, REPL, examples, and regression fixtures.
