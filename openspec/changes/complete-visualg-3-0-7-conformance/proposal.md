## Why

The current interpreter implements a useful VisuAlg-like subset, but its grammar, semantics, runtime behavior, built-ins, diagnostics, and tools are not yet demonstrably compatible with the official VisuAlg 3.0.7 release. A strict, evidence-backed conformance target is needed now so that further work converges on one observable language rather than preserving assumptions that conflict with the reference implementation.

## What Changes

- Establish official VisuAlg 3.0.7 as the behavioral authority and require disputed behavior to be recorded on an interactive Windows reference host before implementation.
- Add a reproducible oracle corpus containing probes, inputs, normalized output and errors, raw hashes, generated file bytes, and screenshots for GUI-only observations, with separate evidence and implementation states and incremental replay gates; do not vendor the VisuAlg executable or archive.
- **BREAKING** Tighten lexical, syntactic, semantic, runtime, formatting, and diagnostic behavior wherever the current implementation differs from the oracle, even when the current behavior is more permissive.
- Complete the oracle-confirmed language surface, including physical-newline rules, declarations, vectors, subprogram call variants, returns, expressions, and control flow. Treat constants, named types, records, fields, assignment aliases, and case ranges as disputed candidates: implement accepted forms and preserve rejection tests for unsupported forms, revising dependent plans after recording.
- Reproduce recorded storage restrictions, checked allocation and indexing, zero values, copying and aliasing rules, and overflow behavior. Fresh recordings accept vectors beyond the formerly assumed 500-slot limit.
- Complete console I/O, CP1252 handling, `arquivo`, random-input mode, timers, echo, chronometer, pause/debug, screen/color integration, and any additional oracle-confirmed environment commands through injectable headless host hooks.
- Replace duplicated built-in knowledge with one authoritative descriptor registry shared by semantic analysis and runtime, and implement the oracle-confirmed numeric, text, conversion, `rand`, and `randi` behavior.
- **BREAKING** Introduce stable positioned runtime diagnostic codes, correct CLI exit statuses, a formatter that retains comments through lexing and parsing, and REPL submission when `fimalgoritmo` completes a program.
- Define documented source, nesting, call-depth, and value-size safeguards plus opt-in execution budgets for bounded tests, while allowing intentional infinite loops in ordinary execution.
- Match random value domains and command semantics while explicitly excluding byte-for-byte random sequences.
- Deliver the work as focused stacked pull requests with test-first regressions, immediate task tracking, documentation and changelog updates, full oracle traceability, and final implementation acceptance. Archive and `v0.1.0` tagging follow the separate post-merge checklist in `release.md`.
- Exclude a vendored reference executable, native compilation or a bytecode VM, an LSP, recreation of the VisuAlg GUI, and exact RNG sequence reproduction.

## Capabilities

### New Capabilities

- `conformance/visualg-3-0-7`: Defines the reference release, oracle evidence, compatibility authority, traceability, and conformance acceptance for bundled examples.
- `language/lexical-grammar`: Defines source encoding, physical newlines, tokens, identifiers, keywords, comments, literals, declaration grammar, recovery, and canonical printing.
- `language/types-declarations`: Defines confirmed declarations, vectors, assignments, coercion, storage semantics, and recorded storage restrictions, with evidence-dependent acceptance or rejection of constants, named types, records, and fields.
- `language/expressions-control-flow`: Defines operators and precedence, comparisons, cases and ranges, loops, interruption, I/O statement validation, and overflow behavior.
- `language/subprograms`: Defines declaration and call forms, parameters, exact reference typing, evaluation order, lexical scope, recursion, and return behavior.
- `runtime/io-environment`: Defines console and file input, output formatting, encoding, random-input mode, timers, echo, chronometer, and host-backed environment commands.
- `runtime/standard-library`: Defines the authoritative built-in catalog and the numeric, trigonometric, text, character-code, conversion, and random functions.
- `tooling/diagnostics-cli-repl`: Defines stable diagnostics, source positions, CLI statuses, formatting, and interactive input and submission behavior.

### Modified Capabilities

None. This repository has no existing OpenSpec capability specifications; the new capability set formalizes and, where necessary, supersedes behavior currently described only in project documentation and code.

## Impact

This program affects every pipeline stage (`source`, `token`, `lexer`, `ast`, `parser`, `sema`, `runtime`, `interp`, and `stdlib`), the CLI and REPL, diagnostics, examples and fixtures, language documentation, and CI. Public internal boundaries will change to pass immutable semantic information into an option-configured interpreter and to route environment-dependent behavior through typed host, clock, filesystem, input/output, and randomness seams. The test suite gains an oracle manifest and conformance corpus, platform jobs, fuzzing and malformed-input coverage, filesystem and CP1252 cases, deterministic fakes, subprocess checks, and traceability validation. Core language packages remain free of third-party dependencies.
