# Changelog

## Unreleased

- Add numeric timer commands, conditional debug breakpoints and pause commands.
  Route delays and pauses through the host, preserve recorded loop and call
  timing, reset timer state per run and report positioned host failures.

- Add typed echo settings and chronometer commands with per-run state, injected
  clocks, recorded elapsed-time text, canonical formatting and positioned host
  failures. Preserve console input echo in the deterministic headless profile.

- Add command-form random input with evaluated bounds, bounded real precision,
  generated text, console transitions, per-run reset and guarded source failures.
  Replace the legacy callable `aleatorio` API with recorded no-value behavior.
  Record eight more original bundled examples without claiming exact random
  sequence compatibility.

- Pin an independent inventory of 28 documented builtins to a recorded reference
  program and check semantic binding, runtime results and formatting. Record
  rejection of the `pot` candidate and pending `div(...)` value behavior.

- Preserve intermediate operands in recorded mixed logical/comparison reductions,
  including enclosing arithmetic and function calls. Bound retained storage and
  release it after each outer evaluation or failure.

- Add bare `rand` expressions, recorded statement and assignment-suffix forms,
  deterministic injection, and guarded random fractions in `[0, 1)`.
  Reject missing sources and invalid bounded draws in the legacy random API.

- Correct logical precedence and single-comparison syntax, preserve grouped
  comparisons during formatting, and ignore suffixes after assignment expressions.

- Verify eight more bundled examples with fixed input and original/formatted
  execution. Retain the reference's character-table rejection separately from
  implementation support, reducing the unrecorded example backlog to 21.

- Add `dos` console configuration with ignored line tails, repeated directives,
  subprogram configuration, typed host calls, headless output, and positioned
  failures that preserve preceding output. Preserve directives through formatting.

- Match logical ordering and Windows-1252 character ordering. Preserve mixed
  comparisons' concrete values and logical assignment/condition category,
  including changed logical storage, reference parameters, positioned consumer
  errors, and no-value handling. Preserve prior output before rejected logical
  field formatting.

- Add named records with scalar fields, local layouts, zero initialization,
  independent copies, vector elements, and scalar reference arguments. Match
  empty output, record alias identity, non-addressable nested fields, and
  positioned field/declaration errors; guard storage before access or copying
  and preserve captured field references across whole-record assignment.

- Add global and local scalar type aliases with earlier-alias resolution,
  first-definition precedence, compatible vector elements and scalar arguments,
  positioned declaration errors, and preserved formatter output.

- Accept an empty source exponent as zero and treat a following sign as an
  arithmetic operator. Preserve real literal types and behavior through formatting.

- Preserve integer-zero fallback for numeric text with a zero integral
  prefix, including fractional and exponent forms. Keep nonzero-prefix real
  conversions and their runtime type checks intact.

- Match recorded fixed and scientific field formatting, including the
  216-place cap, binary-value rounding, negative zero, digit padding, and
  the large-value scientific-notation threshold.

- Limit string values and output widths to the recorded 255-character bound.
  Preserve complete input echo while storing the bounded value, keep accented
  characters intact, and apply the same limit before text conversion.

- Add `caracpnum` with recorded integer/real result selection, numeric spellings,
  fallback values, and positioned conversion errors. Preserve the runtime type
  through expressions, assignments, calls, vector indices, and format bounds.

- Match recorded text slicing, empty searches, accented casing, ANSI character
  codes, and the reference character table. Accept fractional substring bounds,
  preserve no-value character results, and pin text-call diagnostics. Verify
  the original bundled multiplication-table example with fixed input.

- Add recorded screen-clear and color commands through the typed host interface,
  preserving headless output and reporting positioned host failures. Verify
  original and formatted execution of three more bundled display examples.

- Resolve integer constant vector bounds before allocation, including negative
  bounds and per-call local layouts. Preserve prior output on a later invalid
  declaration, reject recorded inline vector parameter/result types, and guard
  individual aggregate allocations with a documented project slot limit.

- Add global and local scalar constants with ordered initialization, current
  parameter/global values, immutable names, and preserved formatter output.
  Match recorded declaration errors, dependencies, built-ins, and overflow;
  retain constant-based vector layouts as pending conformance work.

- Add dynamic `ate` case ranges with inclusive numeric matching, first-match
  selection, and positioned label diagnostics. Match numeric selector
  truncation and recorded text-label casing, preserve ranges through
  formatting, and verify the original bundled choice example.

- Accept one optional semicolon at the end of a `var` or variable-declaration
  line, including local and vector declarations. Reject repeated or misplaced
  semicolons with a positioned diagnostic and verify the bundled mean example.

- Accept case-insensitive `fimfunção` as a function terminator, reserve it as
  a variable name, and format it as `fimfuncao`. Verify the original bundled
  string-processing example that uses the accented ending.

- Record 21 more original bundled examples. Verify fourteen fixed-input
  programs and three source rejections; retain the remaining accepted-program
  and diagnostic differences as explicit conformance work.

- Recognize case-insensitive `DIV` as the integer-division operator alias and
  preserve its behavior through formatting. Verify four more bundled examples
  with fixed input, including the combination program that uses this alias.

- Preserve finite large-angle conversions, report real squaring overflow at
  the call site, and match numeric-domain absence through nested calls and
  exponentiation argument evaluation.

- Correct numeric empty calls, logarithm signatures and domains, and argument
  evaluation when a call produces no value. Reject unsupported `frac` calls
  and report invalid exponentiation domains at the call site.

- Add recorded inverse trigonometric functions, cotangent, angle conversions,
  and `quad`, including optional arguments, result types, and no-value domains.
  Support bare `pi` and reject its parenthesized call form.

- Report invalid numeric powers and nonfinite arithmetic results at their
  operators, preserve recorded nonnumeric division and no-value behavior,
  and accept numeric unary plus through parsing, analysis, and formatting.

- Match recorded division and remainder operand types, signed conversion,
  divisor handling, and left-to-right evaluation. Preserve positioned guards
  for unsupported conversions and integer division failures.

- Match recorded signed 32-bit wrapping in integer addition, subtraction,
  multiplication, negation, `abs`, and `int`, including assignments and function
  returns. Guard unsupported integer conversions and division overflow with
  positioned diagnostics; retain the remaining division and remainder cases.

- Classify whole-number literals above the signed 32-bit limit as real values,
  preserve real literal types through formatting, and match the recorded
  15-significant-digit default output. Share numeric rendering with `numpcarac`,
  diagnose large random bounds during analysis, and retain overflow and
  assignment-timing observations as pending conformance work.

- Implement recorded `randi` calls, positive and negative domains, empty/zero
  bounds, argument evaluation, and positioned guards for invalid injected
  sources. Add deterministic and seeded domain tests, verify the original
  bundled random example, and retain large literal typing as pending evidence.

- Implement recorded `numpcarac` conversion, including precision, exponent
  spelling, zero arguments, signed zero, and no-value output behavior. Add
  twenty-four reference programs and verify the original bundled string
  example, with positioned syntax and conversion guards.

- Read console input one line per destination, preserve spaces and empty
  character values, and reproduce the recorded scalar input echo. Keep unread
  lines across calls and interpreter reuse; update CLI and REPL regressions
  for the shared input stream. Match numeric mantissa fallback, integer
  narrowing, and first-character logical conversion, with positioned guards
  for input exhaustion and numeric range failures.

- Match recorded nested-output ordering and newline consumption. Evaluate and
  format a statement's items before emitting them, preserve nested output on
  evaluation failure, bound the statement buffer, and clear pending newline
  state before interpreter reuse. Add nine reference programs and verify
  twelve nested-output and function-return cases.

- Diagnose missing return values on the return line without consuming the next
  statement or function terminator. Verify all four scalar result types and
  bare returns outside functions; retain mixed syntax/semantic error ordering
  as pending reference work.

- Accept standalone `escreva`/`escreval` and optional `faca`/`faça` in choice
  headers. Recover malformed statements at physical line boundaries, verify
  matching LF/CRLF behavior, and run the original bundled prime-number example.
  Add thirteen reference observations and verify fifteen more corpus cases.

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
