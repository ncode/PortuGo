## Context

See `proposal.md` for scope. Baseline commit `a57eba8b1aafa83c9e5471ab9bc1b8ea0efd22c3` passes `go test ./...`. Exact CLI replay reproduces all four pending manifest comparisons. Fresh interactive reference runs use the manifest's pinned VisuAlg 3.0.7 executable and verified editor readback:

| Original / probe | Baseline CLI | Fresh reference observation |
| --- | --- | --- |
| `EXEMPLO1.alg` / `bundled-24919c4a2e4c` | Two P001 diagnostics at line 41; no output | One assignment rejection at line 21 after consuming `3`, `1`, `2` and producing an 82-byte prompt/echo prefix. The old empty-output expectation is incomplete. |
| `Tabela_ASCII4.alg` / `bundled-49eb7f47065e` | Exits 0; 14,164 output bytes | Syntax rejection on line 301, with the recorded 152-byte language-output prefix. |
| `ESCOLA1.alg` / `bundled-f251f2cd0c22` | Parsing cascades from line 84; no output | After `SAIR`, missing FIMSE at EOF line 320; the recorded 41-byte prompt/echo prefix is retained. |
| `interrompa.alg` / `bundled-f4fa0c514cf7` | P001 line 16; writes ` 1` | The same line-16 rejection and two output bytes ` 1`, without a newline. The old empty-output expectation is incomplete. |

The original rejection screenshots establish diagnostics, but their associated derived stdout expectations do not all preserve preceding output. The fresh assignment, repeat, school, and table captures were observed on 2026-11-06 at 17:36:37Z, 17:32:58Z, 17:39:22Z, and 17:40:30Z respectively. Their executable hash matches the manifest. Raw captures remain private pending publication review; they are not yet committed corpus evidence and no implementation state has changed.

## Goals / Non-Goals

**Goals:** Qualify four complete original-input contracts, correct incomplete expectations with retained proof, fix the remaining implementation differences, and preserve the baseline's 1,673 verified probes, including 1,227 recorded accepted probes.

**Non-Goals:** General statement-at-a-time interpretation, a parser/runtime rewrite, source-identity special cases, rejecting nested calls generally, adding dependencies, changing releases, or silently expanding this change to every malformed program discovered during reduction.

## Decisions

### 1. Qualify evidence before changing behavior or expectations

Use the existing stage/capture workflow with pinned provenance, exact original source/stdin hashes, fresh timestamps, screenshots and labeled transcriptions. Preserve old observations and derived expectation artifacts unchanged; point any corrected expectation to a separately recorded, reviewed artifact with an explicit disposition.

Keep language-output bytes distinct from GUI presentation. The simulated console can add display rows around screen clearing and loop control, while the results panel can add framing or represent an unfinished line with a display newline. Neither is automatically a native stdout byte stream. For interrupted runs, qualify the complete language-output prefix from the source and both controls; do not apply the completed-run panel normalizer to an unfinished panel, strip significant whitespace, or infer empty output from a diagnostic-only screenshot.

The proposed assignment prefix is exactly:

```text
Digite 3 numeros inteiro!
Primeiro numero: 3
Segundo numero: 1
Terceiro numero: 2
```

It ends in LF and contains 82 UTF-8 bytes. The repeat prefix is exactly two bytes, space followed by `1`, with no LF. Keep the existing 152-byte table and 41-byte school contracts unless reviewed complete evidence changes them. Inconclusive runs remain excluded; conflicts block dependent implementation until their disposition is reviewed.

### 2. Use reduced observations to bound each rule

The following investigation results refine the original hypotheses:

- Scalar and vector incompatible assignments both print `BEFORE` before rejection. Assignment failures are reached during execution; simply selecting a different preflight diagnostic would still lose the reference prefix.
- An unused procedure with a missing FIMSE permits a valid main body to print `MAIN` and complete. Invoking the malformed body prints `BODY` and `MAIN` before rejecting, so immediate rejection at the procedure boundary is not established. Keep that counterexample visible when choosing recovery boundaries.
- Simple sibling/nested string-call controls and a precomputed-argument variant all succeed. A reduced table program retaining seven helpers and two main assignments still reproduces the failure. Inlining the empty-string check, removing the prior call, and simplifying the helper each retain the failure. The precise load-bearing call/body state remains an implementation prerequisite; a general nested-call prohibition is ruled out.
- A reduced repeat marker writes ` 1` before rejecting. An unconditional break escapes the marker successfully. Results-panel stream text preserves concatenation even when the simulated console moves its cursor to a new display row.

If a new neighboring control reveals an independent mismatch, retain its observation and report its scope explicitly. Do not misclassify it as passing, force an early diagnostic the reference does not produce, or broaden the implementation silently.

### 3. Fix shared seams after classification is established

- **Assignment path:** distinguish evidence-backed deferred failures from blocking front-end errors. Preserve evaluation order and preceding I/O, and let the earlier reached assignment error win over the later malformed condition. Keep `check`'s existing static diagnostics and lack of runtime effects. Do not execute an unchecked partial tree or suppress errors by code/message/line alone.
- **Character table:** repair or model only the confirmed shared call/evaluation rule. Preserve the successful simple-call and precomputed controls; do not special-case function names, example text, or line 301.
- **School path:** retain positions and bounded recovery for the recorded uncalled-body and reached-EOF paths. Every recovery loop must advance or terminate under existing structural limits. Reached malformed-body behavior needs its own reviewed contract before a proposed recovery rule depends on it.
- **Repeat path:** correct the derived expectation after evidence qualification. The baseline CLI already matches the freshly observed two-byte prefix; no output-buffering or interpreter change is justified for this case.

Prefer existing deferred-diagnostic information and statement representations. Add a representation only if a regression demonstrates that current structures cannot preserve the confirmed trigger and position. Audit CLI `run`/`check`/`fmt`, REPL submissions, canonical printing, and the conformance execution adapter.

### 4. Test original programs and rules independently

Use existing subprocess helpers with exact stdin, timeout, finite execution budget, and bounded output. Assert the reviewed stdout artifact bytes, exit status, diagnostic count/code/line, and empty stderr on successful controls. Do not assert reference columns that were never recorded. Keep each implementation regression red before its production fix; use a data-only regression for the incomplete repeat expectation.

Test reached/unreached paths, RHS and target evaluation order where assignment classification changes, unchanged static `check` behavior, formatting safety or equivalent original/formatted replay, interpreter/REPL state reset, and existing resource limits. Preserve all previously verified accepted and rejected probes through full replay.

### 5. Make four-probe completion explicit

Existing release acceptance treats these rejected examples as non-blocking, so add an explicit completion check for all four reviewed original-source contracts. A green accepted-example sweep is insufficient. Keep each mapping pending until its qualified reference expectation, exact original-source CLI regression, and test/task traceability pass.

Do not repoint the manifest's historical task catalog to a new change lacking its original IDs. Use supported traceability or the smallest validated extension if both historical and new owners need representation. Derive final counts and current documentation from validation; retain historical reports and unchanged source/stdin/evidence bytes.

## Risks / Trade-offs

- **Incomplete GUI captures become false language rules** → Preserve both channels and old artifacts, publish reviewed complete evidence, and correct expectations before coding toward them.
- **Recovery masks real errors or loops** → Keep blocking diagnostics outside confirmed deferred rules, preserve static `check` coverage, and exercise boundary/resource cases plus lexer/parser fuzzing.
- **A reduced counterexample expands the feature** → Report it explicitly and review scope before changing contracts or task completion claims.
- **Call-state cause remains uncertain** → Keep the character-table implementation task gated on a reduced trigger and neighboring controls; no source-specific workaround.

## Migration Plan

No release/tag or user-data migration is required. After evidence qualification, land focused implementation/data corrections, tests, traceability, and current docs together. Rollback their commits together; preserve old and new reference observations and published release history.

## Implementation Gate Status

All four original reference paths have been rerun, and reduced controls in all four areas have been observed privately. The proposal now accounts for the two incomplete output expectations. Publication/traceability of that evidence and final reduction of the table rule remain required before dependent implementation; all implementation tasks remain unchecked.
