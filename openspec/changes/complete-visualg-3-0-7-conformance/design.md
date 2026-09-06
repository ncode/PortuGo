## Context

See `proposal.md` for motivation and the eight delta specifications for observable requirements. The current repository is a small Go 1.22 tree walker whose semantic entry point returns diagnostics only, whose interpreter reparses AST type syntax and returns an unpositioned `error`, whose built-in signatures and implementations are duplicated, and whose CLI and REPL construct input streams independently. The AST separates globals from subprograms, stores vector bounds as integers, and loses comments during lexing. Declaration ordering, bound expressions, optional forms, and environment commands need evidence before their syntax is expanded. In particular, `especificacao-visualg-3.md` identifies records and named types as third-party extensions, not established VisuAlg 3.0.7 features.

The reference application is an interactive Windows program. Deterministic behavior can be captured, but some commands require a person to drive the GUI and record screenshots. Core language packages must remain dependency-free, source positions must survive decoding and every pipeline stage, malformed input must not panic, and implementation work must stay reviewable as focused stacked pull requests.

## Goals / Non-Goals

**Goals:**

- Make reference evidence a repeatable engineering input rather than informal notes, with a validator that closes every traceability gap.
- Pass resolved semantic facts to execution so the interpreter never guesses types, bindings, declaration values, or layouts from syntax.
- Isolate input, output, working-directory, clock/display host behavior, and randomness for headless use and deterministic tests.
- Represent the complete oracle-confirmed grammar without unordered side tables or sentinel expressions.
- Make every runtime failure a stable, positioned diagnostic and keep semantic and runtime built-in behavior derived from one catalog.
- Allow each conformance slice to land independently with a regression, focused checks, documentation, and an explicit OpenSpec checkpoint.

**Non-Goals:**

- Preserve permissive legacy behavior behind a compatibility flag; conformance changes intentionally replace it.
- Automate or redistribute the proprietary VisuAlg executable, installer, or archive.
- Model GUI pixels or terminal escape sequences as language semantics; only typed command effects and recorded observable behavior are modeled.
- Introduce a public Go API, third-party core dependency, native compiler, VM, LSP, or exact reference RNG algorithm.

## Decisions

### 1. Record an oracle corpus before disputed implementation

The repository will store evidence under `testdata/conformance/visualg-3.0.7/`. A checked-in JSON manifest will identify the reference hashes and environment and map stable probe IDs to source, input, normalized output/error, raw-observation hashes, generated-file hashes, screenshots, specification requirements, implementation tests, checklist items, defects, and bundled examples. Probe directories will contain reduced `.alg` programs and only redistributable evidence; the reference executable and archive remain outside Git.

A Windows recorder will guide an interactive run, collect console text and generated files, accept screenshot paths for GUI-only results, hash raw material, and apply a versioned normalizer. The validator will recalculate committed hashes, reject missing files and prohibited binary extensions or known reference hashes, and validate traceability according to the current phase.

Each manifest entry has an `evidenceState` (`unrecorded`, `recorded`, or `not-applicable`), an `implementationState` (`pending`, `verified`, or `not-applicable`), and an owning implementation group. A recorded probe also states whether the reference accepts or rejects its source. `pending` entries identify their future task IDs but need no nonexistent test paths. `verified` entries require existing test IDs and a successful replay against the current checkout. `not-applicable` requires a reviewed reason; project-specific CLI, formatter, and safeguard requirements use project evidence and tests with an explicit rationale for having no reference probe. Reference rejection remains a required negative test, not a reason to waive conformance.

| Gate | Evidence required | Implementation required |
|---|---|---|
| Group 2: evidence | All required probes recorded, hashes valid, inventory and owner/task links complete; explicit reviewed exceptions | Later behavior may remain `pending`; existing verified baseline probes must replay successfully |
| Groups 3–16: incremental | Evidence gate continues to pass | All probes whose owning group has completed are `verified`; replay every verified probe and report later pending entries |
| Group 17: implementation acceptance | No unresolved required evidence or stale links | Every required probe and accepted example is verified, no pending behavior or deterministic mismatch |

Group 2 builds a corpus runner as well as a metadata validator, initially replaying the existing CLI in bounded child processes. Group 3 migrates it to the new Go pipeline with recorded input, budgeted execution, deterministic fakes, and observation adapters. The runner compares output, mapped diagnostic categories/positions, generated bytes, and observable state/host events, and reports pending mismatches without treating them as passing results. Fixture limits and subprocess deadlines apply as described in decision 10. A mismatch in a verified entry always fails CI; moving a verified entry back to pending requires an explicit reviewed scope correction. Group 17 and the release gate reject all pending entries.

Before group 3, group 2 corrects `proposal.md`, the delta specs, `design.md`, **and `tasks.md`**, along with affected project documentation, using the recorded dispositions. Constants, named types, records, fields, assignment aliases, and case ranges have no assumed positive support. For rejected features, replace positive requirements and implementation subtasks with rejection coverage, remove dependent layouts/examples, and retain the original probe and requirement/task ID disposition in the manifest. Preserve existing task IDs when rewording; explicitly map retired IDs to their replacement or reviewed disposition. Groups 3, 5–9, and 16 consume that revised scope; they must not construct unused record/type machinery just to satisfy the original plan.

Alternative considered: encode behavior directly from manuals and the current implementation. This is insufficient because manuals omit edge cases and the plan explicitly makes the executable's behavior authoritative.

### 2. Make semantic analysis produce the execution contract

The semantic boundary is:

```go
func Analyze(program *ast.Program) (*sema.Info, []diag.Diagnostic)
```

`sema.Info` owns resolved expression and designator types, vector layouts, and declaration/use bindings. Group 3 implements facts for the existing supported syntax. Constant values, named-type identities, record layouts, and extended bounds/slot accounting are added with their owning groups only when reference evidence requires them. Its maps remain unexported; query methods return immutable values or defensive copies. `Analyze` may return partial information alongside diagnostics for tooling, but callers MUST execute only when the diagnostic slice has no errors. The existing `Check` function is removed after callers migrate rather than maintained as a second semantic path.

The interpreter boundary becomes:

```go
func (i *Interpreter) Run(program *ast.Program, info *sema.Info) []diag.Diagnostic
```

Execution consumes bindings and layouts from `Info`, validates impossible or corrupted state defensively, and does not call a syntax-to-runtime type converter. A slice result aligns runtime with the other pipeline stages; execution normally stops after the first runtime diagnostic but can preserve cleanup diagnostics without inventing a separate error channel.

Alternative considered: annotate mutable fields directly onto AST nodes. A separate immutable result keeps syntax reusable by the formatter and avoids analysis-order mutation, hidden coupling, and stale annotations after a reparse.

### 3. Configure execution through explicit options and typed host seams

Construction is locked to:

```go
type Options struct {
    Input      io.Reader
    Output     io.Writer
    WorkingDir string
    Host       Host
    Random     RandomSource
    MaxSteps   uint64
}

type RandomSource interface {
    Float64() float64
    Uint64N(n uint64) uint64
}

func New(options Options) *Interpreter
```

Nil input and output use empty input and `io.Discard`; an empty working directory resolves once to the process working directory; a nil random source receives a per-interpreter standard source. The interpreter remains single-goroutine and owns all buffered input state so the REPL can pass a shared reader without competing scanners.

`MaxSteps` is an optional execution budget: zero means no step limit for ordinary programs; a positive value bounds interpreter work as defined in decision 10. Call-depth and value-size safeguards remain active regardless of this setting.

The host boundary is a typed interface rather than a generic command bag:

```go
type Host interface {
    Delay(duration time.Duration) error
    Breakpoint(event Breakpoint) error
    ClearScreen() error
    SetDisplay(state DisplayState) error
    Now() time.Time
}
```

`Breakpoint` and `DisplayState` are closed typed values that carry only oracle-confirmed options, including foreground/background color and visibility state. Newly discovered host commands extend the typed interface or these closed values only after an oracle recording. The default headless host uses the system clock and delay, makes UI-only breakpoint, clear-screen, and display operations non-blocking no-ops, and emits no terminal escape bytes. Fakes record calls and control time. Host errors become `R008` at the originating statement.

Alternative considered: call `time`, terminal functions, and global randomness directly. That would make conformance tests slow or nondeterministic and would entangle GUI assumptions with the language runtime.

### 4. Expand the AST around ordered declarations and explicit optionality

`ast.Program` will hold ordered declarations for the section kinds accepted by the oracle. Variables, procedures, and functions are the baseline; constants, named types, record specifications, and field designators are added only if group 2 confirms them. Local declaration regions use the same declaration node family where permitted. Accepted declaration expressions remain expressions until semantic resolution. Rejected forms receive positioned rejection coverage without adding unused AST variants.

For oracle-accepted ranges, case clauses use positioned labels with `Low` and optional `High` expressions rather than ordinary binary expressions. `ReturnStmt` stores an optional value only if bare returns are accepted. `RepeatStmt` stores an optional termination condition if the infinite form is accepted, keeping it distinct from conditional repetition. Rejected variants do not require optional fields. Environment operations use typed statement nodes with positioned, command-specific operands; they are not encoded as magic procedure names. Assignment aliases accepted by the reference share the ordinary assignment node after lexing records the source token.

The canonical printer covers every new node in the same group that introduces it. Parser goldens and parse-print-parse tests make ordering and optional forms visible.

Group 4 retains comments end to end. The lexer emits positioned comment tokens without consuming their physical newline. The parser stores ordered comment groups on the program, preserving decoded comment text and source spans, and anchors each group before the following construct, after a construct on the same line, before a closing delimiter, or at end of file. Comment-only and empty blocks retain their groups. The printer emits each group exactly once at its anchor; only indentation and line endings are normalized. If the oracle ignores a suffix after `fimalgoritmo`, retain that suffix as opaque source and re-emit it after the terminator so formatting cannot discard notes or reinterpret ignored text. Round-trip checks compare comment contents, ordering, and anchors as well as syntax. Later AST additions preserve this representation rather than introducing a second comment path.

Alternative considered: add more fields to the current `Globals`, `Subs`, and `TypeSpec` structures. That cannot preserve interleaved declaration order and encourages sentinel values for missing conditions, bounds, and returns.

### 5. Separate resolved type identity from storage locations

Resolved types are immutable descriptors for the scalar and vector types accepted by the oracle. Named-type identities and record layouts are conditional on group 2 evidence; accepted records contain ordered fields with offsets and resolved field types. Vector layouts contain every lower and upper bound plus an overflow-checked flattened size. Slot accounting is performed during analysis using oracle-recorded scope and accepted aggregate rules, and execution checks the resolved sizes again before allocation.

Runtime environments bind semantic symbols to typed locations. A location can select a scalar or vector element, and a record field only if supported, without copying its container. Accepted value assignment and value parameters use one copy operation following the recorded depth rules; `var` parameters carry locations and require exact semantic type compatibility. This centralizes copying versus alias behavior and ensures index or field failure occurs before mutation.

Alternative considered: represent aliases as copied `runtime.Value` instances plus write-back. Write-back fails for early returns and nested designators, evaluates indices at the wrong time, and cannot reproduce aliasing between arguments.

### 6. Use structured runtime diagnostics end to end

`diag.Diagnostic` will gain severity, source span, and an optional wrapped cause while retaining stable rendering. Runtime categories are fixed as:

- `R001`: type or coercion
- `R002`: arithmetic
- `R003`: storage or indexing
- `R004`: input
- `R005`: call or return
- `R006`: loop state
- `R007`: built-in failure
- `R008`: host or file I/O

Interpreter helpers accept the triggering AST position and return diagnostics rather than formatting anonymous errors. Internal invariants may still panic with an `internal:` prefix, but every state reachable from source, input, file contents, host failure, or size arithmetic must use a diagnostic. CLI rendering sends diagnostics to standard error and returns status 1; usage returns 2; success returns 0.

Alternative considered: wrap the current `error` at the CLI. The source node and failure category have already been lost by then, so stable positioning cannot be reconstructed reliably.

### 7. Drive semantic and runtime built-ins from one descriptor registry

`internal/stdlib` will expose immutable descriptors containing canonical name and aliases, callable forms, parameter modes and accepted types, result-type selection, domain metadata, and an evaluator. Semantic analysis uses descriptor metadata for binding, arity, reference requirements, and result types; execution dispatches through the same descriptor. Registry construction validates duplicate aliases and incomplete descriptors once without mutable global state. Stateful library instances receive the interpreter's random source.

Catalog completeness is checked against the independently recorded reference inventory, not just against descriptors already present in Go. Numeric probe and implementation coverage explicitly includes `arccos`, `arcsen`, `arctan`, `cotan`, `grauprad`, `radpgrau`, and `quad` in addition to the existing names. These are distinct functions, not aliases. Every candidate, including existing `frac`, gets a recorded disposition; accepted names receive descriptors and rejected legacy names receive rejection regressions. The two-argument `exp(base, expoente)` form in the local reference gets an explicit arity/value probe. Removing a required function from both sema and runtime must still fail the inventory comparison.

Conversion operations with value-dependent result types, notably `caracpnum`, use a descriptor result rule that records a numeric union during analysis and a concrete integer or real value at runtime. Call-site operations requiring a single concrete type must narrow or diagnose that union according to oracle evidence.

Alternative considered: generate one switch from another. A shared descriptor is easier to table-test and prevents signature metadata and runtime dispatch from drifting again.

### 8. Model input sources as one state machine

The interpreter owns one buffered input controller whose active mode is console, `arquivo`, or random input. Mode transitions retain the injected console reader and encode file exhaustion, fallback recording, and echo as states derived from oracle probes. Paths resolve against `Options.WorkingDir`; tests use temporary directories. CP1252 decoding and encoding live in shared source/text helpers so source files, `arquivo`, generated files, character-code built-ins, and output agree.

The REPL reads program text and program input through the same buffered abstraction. It uses lexer/parser completeness, not substring matching, to submit immediately when the terminating `fimalgoritmo` token completes a program. Blank lines remain ordinary input to the incomplete-program state, and every submitted program gets a fresh semantic/runtime state while the underlying reader remains shared.

Alternative considered: keep separate `bufio.Scanner` instances for REPL entry and `leia`. Competing buffering can consume each other's bytes and cannot safely hand unread input from one layer to the other.

### 9. Treat each task group as a verifiable stacked change

The 17 task groups in `tasks.md` are implementation pull-request boundaries. Each group starts with a committed failing regression or manifest validation, implements only that slice, runs focused and full checks under the applicable corpus phase, updates `docs/language.md` and `CHANGELOG.md` for visible behavior, checks its OpenSpec items immediately, and publishes a draft PR stacked on the preceding branch. Oracle corrections land before dependent behavior. Group 17 closes implementation traceability and produces the acceptance report and release handoff.

Archive and tag operations live in `release.md`, outside the implementation task count. During group 17, acceptance checks validate behavior and report outstanding handoff tasks without requiring their own future completion. After the final handoff task is evidenced and checked, the post-merge release gate independently requires all implementation tasks complete and all acceptance checks green on the release candidate. Only then does it archive and validate the specs, followed by tagging. It never requires an archive or tag to already exist as a precondition for creating it. OpenSpec planning status alone is not proof of task completion.

Alternative considered: one implementation pull request. The change spans every stage and would be too large to review, bisect, or correct when reference evidence changes.

### 10. Bound structural resources and test execution explicitly

These are project safeguards, with project-owned tests and a documented oracle non-applicability rationale; they are not claims about reference limits. The initial limits are 4 MiB of source per file or REPL submission, 256 levels of syntax/AST traversal or expression evaluation within a call frame, 256 active language call frames, and 16 MiB for a single text value, input token/line, or formatted item. Check source/input size during accumulation, depth before recursive descent or traversal, call depth before frame creation, and text/format sizes with checked arithmetic before allocation. Cover long flat expressions whose resulting AST is deep, declaration dependency chains, copies, concatenation, and width/precision expansion as well as visibly nested syntax. Resource limits must not turn a known accepted reference example into a silently excluded example: a hit blocks acceptance until the safeguard is deliberately revised and retested.

Front-end resource failures use a new stable `E900` diagnostic at the first construct or byte exceeding the limit; reserve that code without renumbering existing codes. Runtime call depth uses `R005` at the attempted call, text/format allocation uses `R003` at the producing expression, and step exhaustion uses `R006` at the next operation. Evaluation-depth protection also uses `R003`. Stop before that operation's side effects and unwind resources normally. A size cap is checked before expanding CP1252 to decoded text, using a bounded decoded allocation; the source cap measures original bytes.

A step is charged before each statement dispatch, expression evaluation, and loop iteration/back-edge, including an empty body. One counter belongs to the whole run, including subcalls and input-retry attempts; it never resets on a call and cannot overflow when unlimited. A positive `MaxSteps` permits exactly that many charged operations. Normal execution leaves it zero, so a valid infinite loop is not promised to terminate. This is a work budget, not a wall-clock timeout and cannot interrupt a blocking injected reader or host method.

The fuzz/adversarial profile uses at most 64 KiB of generated source, the structural limits above, and 10,000 runtime steps. Tests inject finite input, output capture capped at 1 MiB (returning a writer error on overflow), scripted randomness, and a nonblocking fake host/clock. File probes use isolated temporary directories. Potentially blocking or crashing cases run in child processes with a five-second watchdog; a watchdog kill is a failing test, never a successful language diagnostic. Ordinary interactive reads may wait for the user. Deterministic corpus entries record sufficient per-probe budgets; budget exhaustion or timeout in an accepted example fails acceptance. No goroutines are added to the interpreter to enforce these limits.

Group 3 introduces runtime budget/call/value guards and bounded corpus execution; group 4 adds source, parser, AST-traversal, and formatter guards. Groups 6–16 extend the guards with each new allocation, traversal, loop, and input mode. Group 17 tests the assembled safeguards rather than deferring their implementation until release.

Alternative considered: claim termination for every source or rely only on a process timeout. Valid infinite programs disprove the former; the latter cannot provide positioned language diagnostics and is useful only as a failing-test backstop.

## Risks / Trade-offs

- [Interactive Windows evidence is slow or unavailable] → Keep unresolved probes explicitly blocking, prioritize reduced probe batches early, and never substitute undocumented assumptions.
- [Oracle evidence changes a foundational grammar or type assumption] → Finish the recorder/spec-correction group before dependent implementation and keep later PRs stacked so they can be rebased in order.
- [The 500-slot rule has context-dependent accounting] → Record exact boundary probes for scalars, records, vectors, globals, locals, and recursion; keep accounting policy separate from overflow-safe arithmetic.
- [A broad `Host` interface makes simple hosts cumbersome] → Provide a complete default headless implementation and small reusable recording fake; add methods only for oracle-confirmed operations.
- [Dynamic conversion result types complicate analysis] → Limit union-like semantic types to catalog operations proven to need them and require explicit narrowing at ordinary language boundaries.
- [Exact output and CP1252 behavior is platform-sensitive] → Compare byte fixtures, centralize encoding and newline conversion, and run Windows, macOS, and Linux jobs.
- [Stacked PR bookkeeping diverges from OpenSpec tasks] → Require checkbox updates in each group before publishing its draft PR and validate that branch/PR evidence is linked in the final traceability report.
- [Strict compatibility breaks existing example behavior] → Treat the break as intentional, update fixtures and documentation in the responsible group, and retain oracle evidence explaining the change.
- [Resource safeguards reject an accepted reference example] → Keep the example in acceptance, review the specific bound, and rerun boundary and adversarial tests after any change; never report a timeout as a match.

## Migration Plan

1. Land this planning change on top of the clean bootstrap `main`; it changes no interpreter behavior.
2. Create each implementation branch from the preceding accepted group, beginning with quality infrastructure and then the Windows oracle recorder and spec corrections.
3. Introduce the semantic result and interpreter options behind a repository-wide caller migration in one group; remove `sema.Check` and the old `interp.New(reader, writer)` / `Run(program) error` entry points in that same group so there is no split execution path.
4. Migrate AST and runtime capabilities in dependency order defined by `tasks.md`, keeping the tree buildable and all prior conformance fixtures green at every group boundary.
5. Update language documentation and changelog with each visible behavior change; do not postpone compatibility notes to the final group.
6. Complete the implementation acceptance report and final PR handoff, then use `release.md` after merge to check task completion, rerun the accepted-example, oracle, fuzz, race, lint, platform, CLI, REPL, filesystem, and traceability gates, archive, and create `v0.1.0` in that order.

Before a group is merged, rollback is a normal revert of that focused group plus rebasing any later stacked branches. No persistent data migration or compatibility mode is required. If final conformance cannot be demonstrated, do not archive or tag; the implemented groups remain independently testable without claiming the release target.
