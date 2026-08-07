## Context

See `proposal.md` for motivation and the eight delta specifications for observable requirements. The current repository is a small Go 1.22 tree walker whose semantic entry point returns diagnostics only, whose interpreter reparses AST type syntax and returns an unpositioned `error`, whose built-in signatures and implementations are duplicated, and whose CLI and REPL construct input streams independently. The AST separates globals from subprograms and stores vector bounds as integers, which cannot represent the declaration ordering, constant expressions, records, optional forms, or environment statements required by the target.

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

A Windows recorder will guide an interactive run, collect console text and generated files, accept screenshot paths for GUI-only results, hash raw material, and apply a versioned normalizer. The validator will recalculate committed hashes, reject missing files and prohibited binary extensions or known reference hashes, and require bidirectional traceability. `[VERIFICAR]` and other disputed cases become blocking manifest states until recorded. Oracle task group 2 may correct the delta specs and design when evidence disproves an assumption, before affected implementation starts.

Alternative considered: encode behavior directly from manuals and the current implementation. This is insufficient because manuals omit edge cases and the plan explicitly makes the executable's behavior authoritative.

### 2. Make semantic analysis produce the execution contract

The semantic boundary is:

```go
func Analyze(program *ast.Program) (*sema.Info, []diag.Diagnostic)
```

`sema.Info` owns resolved expression and designator types, constant values, named-type identities, record layouts, vector bounds and slot counts, and declaration/use bindings. Its maps remain unexported; query methods return immutable values or defensive copies. `Analyze` may return partial information alongside diagnostics for tooling, but callers MUST execute only when the diagnostic slice has no errors. The existing `Check` function is removed after callers migrate rather than maintained as a second semantic path.

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
}

type RandomSource interface {
    Float64() float64
    Uint64N(n uint64) uint64
}

func New(options Options) *Interpreter
```

Nil input and output use empty input and `io.Discard`; an empty working directory resolves once to the process working directory; a nil random source receives a per-interpreter standard source. The interpreter remains single-goroutine and owns all buffered input state so the REPL can pass a shared reader without competing scanners.

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

`ast.Program` will hold an ordered `[]ast.Decl` containing constants, named types, variables, procedures, and functions. Local declaration regions use the same declaration node family where the reference permits them. Constant initializers and vector bounds remain expressions until semantic resolution. Type syntax gains named references and record specifications with ordered positioned field declarations. Designators gain explicit field selection in addition to vector indexing.

Case clauses use positioned labels with `Low` and optional `High` expressions rather than treating ranges as ordinary binary expressions. `ReturnStmt` stores an optional value explicitly. `RepeatStmt` stores an optional termination condition so infinite and conditional forms are distinct. Environment operations use typed statement nodes with positioned, command-specific operands; they are not encoded as magic procedure names. Assignment aliases accepted by the reference share the ordinary assignment node after lexing records the source token.

The canonical printer covers every new node in the same group that introduces it. Parser goldens and parse-print-parse tests make ordering and optional forms visible.

Alternative considered: add more fields to the current `Globals`, `Subs`, and `TypeSpec` structures. That cannot preserve interleaved declaration order and encourages sentinel values for missing conditions, bounds, and returns.

### 5. Separate resolved type identity from storage locations

Resolved types are immutable descriptors for scalar, named, record, and vector types. Record layouts contain ordered fields with offsets and resolved field types. Vector layouts contain every lower and upper bound plus an overflow-checked flattened size. Slot accounting is performed during analysis using oracle-recorded scope and aggregate rules, and execution checks the resolved sizes again before allocation.

Runtime environments bind semantic symbols to typed locations. A location can select a scalar, record field, or vector element without copying its container. Value assignment and value parameters use a single recursive copy operation; `var` parameters carry locations and require exact semantic type compatibility. This centralizes deep-copy versus alias behavior and ensures index or field failure occurs before mutation.

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

Conversion operations with value-dependent result types, notably `caracpnum`, use a descriptor result rule that records a numeric union during analysis and a concrete integer or real value at runtime. Call-site operations requiring a single concrete type must narrow or diagnose that union according to oracle evidence.

Alternative considered: generate one switch from another. A shared descriptor is easier to table-test and prevents signature metadata and runtime dispatch from drifting again.

### 8. Model input sources as one state machine

The interpreter owns one buffered input controller whose active mode is console, `arquivo`, or random input. Mode transitions retain the injected console reader and encode file exhaustion, fallback recording, and echo as states derived from oracle probes. Paths resolve against `Options.WorkingDir`; tests use temporary directories. CP1252 decoding and encoding live in shared source/text helpers so source files, `arquivo`, generated files, character-code built-ins, and output agree.

The REPL reads program text and program input through the same buffered abstraction. It uses lexer/parser completeness, not substring matching, to submit immediately when the terminating `fimalgoritmo` token completes a program. Blank lines remain ordinary input to the incomplete-program state, and every submitted program gets a fresh semantic/runtime state while the underlying reader remains shared.

Alternative considered: keep separate `bufio.Scanner` instances for REPL entry and `leia`. Competing buffering can consume each other's bytes and cannot safely hand unread input from one layer to the other.

### 9. Treat each task group as a verifiable stacked change

The 17 task groups in `tasks.md` are pull-request boundaries. Each group starts with a committed failing regression or manifest validation, implements only that slice, runs focused and full checks, updates `docs/language.md` and `CHANGELOG.md` for visible behavior, checks its OpenSpec items immediately, and publishes a draft PR stacked on the preceding branch. Oracle corrections land before dependent behavior. The last group closes traceability, archives the OpenSpec change, and tags only after all platform, fuzz, race, lint, and deterministic oracle gates pass.

Alternative considered: one implementation pull request. The change spans every stage and would be too large to review, bisect, or correct when reference evidence changes.

## Risks / Trade-offs

- [Interactive Windows evidence is slow or unavailable] → Keep unresolved probes explicitly blocking, prioritize reduced probe batches early, and never substitute undocumented assumptions.
- [Oracle evidence changes a foundational grammar or type assumption] → Finish the recorder/spec-correction group before dependent implementation and keep later PRs stacked so they can be rebased in order.
- [The 500-slot rule has context-dependent accounting] → Record exact boundary probes for scalars, records, vectors, globals, locals, and recursion; keep accounting policy separate from overflow-safe arithmetic.
- [A broad `Host` interface makes simple hosts cumbersome] → Provide a complete default headless implementation and small reusable recording fake; add methods only for oracle-confirmed operations.
- [Dynamic conversion result types complicate analysis] → Limit union-like semantic types to catalog operations proven to need them and require explicit narrowing at ordinary language boundaries.
- [Exact output and CP1252 behavior is platform-sensitive] → Compare byte fixtures, centralize encoding and newline conversion, and run Windows, macOS, and Linux jobs.
- [Stacked PR bookkeeping diverges from OpenSpec tasks] → Require checkbox updates in each group before publishing its draft PR and validate that branch/PR evidence is linked in the final traceability report.
- [Strict compatibility breaks existing example behavior] → Treat the break as intentional, update fixtures and documentation in the responsible group, and retain oracle evidence explaining the change.

## Migration Plan

1. Land this planning change on top of the clean bootstrap `main`; it changes no interpreter behavior.
2. Create each implementation branch from the preceding accepted group, beginning with quality infrastructure and then the Windows oracle recorder and spec corrections.
3. Introduce the semantic result and interpreter options behind a repository-wide caller migration in one group; remove `sema.Check` and the old `interp.New(reader, writer)` / `Run(program) error` entry points in that same group so there is no split execution path.
4. Migrate AST and runtime capabilities in dependency order defined by `tasks.md`, keeping the tree buildable and all prior conformance fixtures green at every group boundary.
5. Update language documentation and changelog with each visible behavior change; do not postpone compatibility notes to the final group.
6. Run the complete accepted-example, oracle, fuzz, race, lint, platform, CLI, REPL, filesystem, and traceability gates before archiving the OpenSpec change and creating `v0.1.0`.

Before a group is merged, rollback is a normal revert of that focused group plus rebasing any later stacked branches. No persistent data migration or compatibility mode is required. If final conformance cannot be demonstrated, do not archive or tag; the implemented groups remain independently testable without claiming the release target.
