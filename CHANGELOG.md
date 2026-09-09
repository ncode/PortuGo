# Changelog

## Unreleased

- Preserve physical LF/CRLF tokens and match recorded program-header rules:
  require the name on the header line, ignore text after its closing quote,
  and avoid cascading errors for invalid headers and main-body starts.
  Add six reference observations and verify eleven header/declaration cases.

- Accept the recorded `caracter` type alias and accented `função`, `então`,
  `senão`, `faça`, `até`, and `não` keywords. Preserve token spelling/positions
  and normalize parsed scalar types for declarations, parameters, and returns.
  Add nineteen reference observations, verify ten keyword cases, and retain
  the remaining grammar and bundled-program gaps.
- Record eleven hash-verified bundled examples: eight accepted by the reference
  and three rejected as supplied. Add byte-exact original/formatted regressions
  for six matching programs. Validate catalog IDs, source hashes and sizes,
  acceptance classifications, and exclusion reviews against the recorded probes.
  Retain the remaining implementation gaps.
- Recover after oversized REPL submissions: report `E900` once, discard the
  rejected input with bounded reads, and resume at the next program header.
  Preserve subsequent `leia` input and the session's failure status. Add exact
  CLI recovery coverage, transient output-failure tests, and nil diagnostic
  writer handling without panics.
- Run REPL programs immediately at a real `fimalgoritmo` token and preserve
  blank lines in unfinished source. Ignore terminator lookalikes in strings,
  comments, and longer identifiers. Keep shared input, diagnostic positions,
  consecutive program execution, and per-program budget resets consistent;
  update the prompt banner and add an exact CLI transcript.
- Validate vector layout products and backing lengths before computing storage
  offsets. Corrupted storage now returns positioned `R003` diagnostics for
  reads, writes, input, and reference arguments instead of panicking. Add
  overflow, boundary, and corruption regressions that preserve prior output,
  elements, and unread input.
- Match recorded vector declaration and indexing rules: accept one or two
  dimensions and unsigned literal bounds; reject signed, fractional, expression,
  reversed, and third-dimension forms. Default an omitted second index to that
  dimension's lower bound and reject whole-vector assignment. Add 26 reviewed
  observations and regressions for 34 vector cases, including accepted sizes
  beyond 500 elements and positioned bounds failures with preceding output.
- Verify 12 existing frontend observations with permanent execution, formatting,
  and positioned-diagnostic regressions: empty declarations, ignored suffixes,
  underscores, mixed-case identifiers, exponent and decimal literals, and
  malformed inline operators. Comment/suffix retention and broader grammar
  rules remain pending.
- Make `retorne` update the active function result while execution continues to
  `fimfuncao`. Accept fallthrough paths, retain same-type results at each call
  depth, and preserve independent recursive frames. Match return-type and
  function-name assignment diagnostics; accept the recorded colonless `caso`
  label. Add 45 reviewed observations and verify 40 additional cases. Bare-return
  diagnostics and nested-write formatting remain pending; cross-type fallthrough
  uses a documented fresh-zero guard.
- Separate callable and variable names while preserving recorded function and
  procedure priority. Reject assignments beginning with a procedure name and
  report duplicate callables at the second declaration without secondary body
  errors. Add 15 reviewed observations and verify 16 additional name-resolution
  cases, including global collisions and indexed assignment targets.
- Match recorded numeric argument conversions and scalar `var` copy-in/copy-out,
  including ordered copy-back, nested calls, and caller type changes. Preserve
  lexical scope and argument order; report procedure argument errors at their
  declaration lines. Add 35 reviewed observations and verify 33 additional call
  cases. Empty-argument edge cases and broader return behavior remain pending.
- Accept parameterless declarations without parentheses, bare procedure calls,
  and bare function values. Match recorded function-name priority over local
  variables and parameters. Reject function results as reference storage and
  retain call-depth safeguards. Add 11 reviewed observations and verify 14
  additional call cases; broader procedure diagnostics remain pending.
- Match recorded comment lines: brace prefixes and slash/star prefixes at the
  start of a line stop at that line's end, without opening multiline blocks.
  Preserve arithmetic operators and quoted delimiters; reject `//` inside a
  quoted value as the reference does. Add 24 reviewed observations and verify
  29 comment-related cases. Comment retention and expression recovery remain
  pending.
- Accept the recorded `:=` assignment spelling through the existing assignment
  token and AST. Token positions and source spelling are retained; formatting
  consistently emits `<-`. The reference assignment probe is now verified.
- Preserve backslashes literally in strings, including a final backslash, and
  reject backslash-escaped quotes as recorded by the reference. Formatting now
  preserves string values and program names without introducing Go escapes.
  Six additional reference probes have verified implementation coverage.
- Document the finite-sweep exclusion for two nonterminating bundled examples,
  retaining their original sources and partial reference observations. Validate
  retained artifact hashes even when a reviewed exclusion applies.
- Added 61 recorded syntax, literal, declaration, call, and vector probes with
  exact source/output hashes and reviewed GUI rejection evidence. Implementation
  remains pending. Corrected the draft storage specification after the reference
  accepted vectors with 501, 5000, and 5001 elements.
- Fixed the REPL dropping buffered source at EOF and bypassing source decoding.
  Complete input now executes, incomplete input reports diagnostics, and UTF-8
  BOM and Windows-1252 source use the file decoder. Earlier submission failures
  still determine the session exit status without preventing later execution.
- Linked existing corpus-tooling tests and recorded feature dispositions to their
  specification requirements, with reviewed reasons for project-only evidence.
  Language recordings, bundled examples, and full acceptance remain incomplete.
- Fixed formatting of deep expressions by omitting unnecessary parentheses.
  The CLI rejects formatted output that exceeds source or syntax limits before
  writing any source bytes, so successful output can be parsed again.
- Added 4 MiB source/submission limits and 256-level syntax/AST traversal limits
  with positioned `E900` diagnostics. Windows-1252 decoding checks original
  bytes before expansion; analysis and formatting reject deep trees before
  traversal or output. Grammar, vocabulary, and comment compatibility remain
  pending additional reference evidence.
- Replaced the internal check/run APIs with immutable semantic information,
  interpreter options, injectable host/random sources, and positioned runtime
  diagnostics (`R001`–`R008`). Added call, expression, text, and formatting guards;
  `run` and `repl` accept `--max-steps`. CLI usage errors exit 2, handled failures
  exit 1, and preceding language output is preserved. REPL input is shared with
  `leia`. Added finite corpus budgets and deterministic state/host observations.
  Fixed overflowing substring lengths and large valid output widths. These
  project safeguards do not establish additional reference compatibility.
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
