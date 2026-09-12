# Evidence foundation

This branch brings already recorded observations into the evidence foundation
so each implementation slice can inherit the required coverage. It preserves
the earlier probes and their verification states. Newly introduced observations
remain **pending implementation** here, even where later branches contain their
implementation and tests. Recorded sources and observation bytes are copied
without modification and checked against their existing hashes.

This is evidence coverage, not a claim of complete language support. The
incremental and final acceptance gates retain their verification requirements.
No verified behavior is downgraded and no missing observation is reported as a
successful run.

## Project tooling

AST source spans, comment anchors, canonical printing, CLI/REPL contracts and
defensive Go storage or host failures are project interfaces that cannot be
observed by running a program in the reference editor. Each such entry retains
its specific non-applicability rationale and owning implementation tasks.
New project entries remain pending on this foundation branch; their tests and
verification arrive with the corresponding implementation slices. Language
syntax and runtime results still require separate recorded programs.

## Bundled rejections

Each unusable bundled example retains its original catalog hash and byte count,
the recorded source, a reviewed diagnostic crop and a labeled transcription.
The individual catalog rationale identifies the observed error and line.
These classifications describe recorded reference outcomes; they do not claim
that this branch already produces matching diagnostics or execution prefixes.
The implementation state remains pending until a later slice supplies tests.

The two [externally stopped examples](bundled-external-stop.md) have an explicit
finite-completion exclusion. Their original sources and partial observations
are retained; incomplete runs are not treated as successful or as syntax errors.
