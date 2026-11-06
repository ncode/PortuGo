## Context

See `proposal.md` for scope. Baseline commit `a57eba8b1aafa83c9e5471ab9bc1b8ea0efd22c3` passes `go test ./...`. A bounded CLI comparison against the four committed probes reproduces every pending difference:

| Original program / probe | Recorded outcome | Current difference |
| --- | --- | --- |
| `EXEMPLO1.alg` / `bundled-24919c4a2e4c` | Exit 1; one E001 at line 21; 0 output bytes | Two P001 diagnostics at line 41 hide the earlier assignment error. |
| `Tabela_ASCII4.alg` / `bundled-49eb7f47065e` | Exit 1; one P001 at line 301; 152 output bytes | Exits 0 and prints 14,164 bytes. |
| `ESCOLA1.alg` / `bundled-f251f2cd0c22` | Exit 1; one P001 at EOF line 320; 41 output bytes | Parsing cascades from line 84; no prompt or input echo. |
| `interrompa.alg` / `bundled-f4fa0c514cf7` | Exit 1; one P001 at line 16; 0 output bytes | Diagnostic matches, but prints ` 1`. |

The pipeline currently prevents execution after blocking front-end diagnostics. It already represents selected malformed constructs as deferred diagnostic statements and defers certain uncalled-subprogram lookup failures. Existing repeat tests require a break before `fimrepita` to succeed. These are constraints on a fix, not proof that their mechanisms explain all four cases.

## Goals / Non-Goals

**Goals:** Four exact original-input comparisons, reusable rules justified by reduced controls, and no regression in the existing 1,673 verified probes (including 1,227 recorded accepted probes).

**Non-Goals:** General statement-at-a-time interpretation, rewriting the parser or callable runtime, filename-specific compatibility, accepting malformed compound comparisons, adding dependencies, or publishing a release.

## Decisions

### 1. Make reference reproduction the first implementation gate

Use the existing conformance stage/capture workflow and the executable hash pinned by the manifest. Record trigger and neighboring control programs, plus a fresh full-source run for each disputed output prefix. Preserve captures privately until text, images, metadata, and hashes have been reviewed for publication.

The gate must establish diagnostic order, trigger point, complete preceding output, and whether uncalled malformed bodies affect main execution. If a recording cannot be reproduced or complete output cannot be observed, keep that case pending and report the blocker. Do not infer empty output from a partial panel or a modal diagnostic screenshot.

Ranked, falsifiable investigation hypotheses (unconfirmed):

1. **Earlier assignment errors dominate later syntax errors.** Moving the incompatible assignment after the malformed condition, or correcting its type alone, will change the first reference error if ordering explains `EXEMPLO1`.
2. **Unused malformed subprogram bodies are skipped.** Removing the unused malformed procedure will preserve the school's prompt and EOF error, while invoking it will expose an earlier diagnostic if lazy body validation explains the prefix.
3. **Sibling user calls share expression state.** Replacing `Zeros(4)+Dec2Oct(xs)` with an equivalent precomputed scalar, while retaining `Direita`, will remove the table failure if sibling-call state is responsible. A single nested call and built-in-only controls distinguish this from nesting generally.
4. **An unterminated output line is withheld at the repeat failure.** Changing only `escreva` to `escreval` will expose preceding output if line buffering explains the extra two bytes. A successful no-newline run checks final flushing.
5. **The original output capture is incomplete or depends on application state.** Fresh isolated full-source runs will disagree with the old transcript, or differ after an editor reset, if capture/state rather than a language rule is responsible.

These are experiment predictions, not implementation rules. The four task branches below remain fixed; the gate selects their smallest confirmed repair. If evidence changes the proposed observable contract, revise and review this proposal before code changes.

### 2. Fix shared seams only after the evidence gate

- **Assignment precedence:** retain a usable partial syntax tree only where safe, and select the evidenced analysis failure without dropping unrelated collected diagnostics. Never execute a partial tree with blocking errors or accept the invalid later comparison to hide it.
- **Character-table call:** repair or model the confirmed call-state rule in the shared callable/evaluation path. Do not special-case the builtin names, example text, or line 301 unless an independently recorded language rule actually depends on a built-in.
- **School recovery:** synchronize at confirmed block/subprogram boundaries, retaining source positions and deferred diagnostics where justified. Keep malformed uncalled-body handling distinct from the reached EOF error. Every recovery loop must advance or terminate under existing structural limits.
- **Repeat output:** change the narrowest confirmed output/rejection seam. Keep successful no-newline output, input prompts, host write errors, and the accepted break-before-marker control intact. A global output discard on runtime error is not an acceptable fix.

Prefer existing deferred-diagnostic AST/semantic information and output handling. Add a representation only if tests demonstrate that existing structures cannot preserve the confirmed trigger and location. Caller audit includes CLI `run`/`check`/`fmt`, REPL submissions, canonical printing, and the conformance execution adapter.

### 3. Test originals and rules independently

Use existing subprocess test helpers with explicit stdin, timeout, finite execution budget, and bounded output to exercise the real CLI. Assert exact stdout bytes, exit status, diagnostic count/code/line, and empty stderr on successful controls. Do not assert reference columns that were never recorded. Keep each regression red before its corresponding production change.

For every new deferred construct, test reached and unreached paths, analysis-only commands without side effects, in-place formatting safety or equivalent original/formatted replay, and interpreter/REPL state reset. Preserve all verified accepted and rejected controls through full corpus replay.

### 4. Keep evidence state separate from implementation acceptance

Add test/task traceability without changing the old source or observation bytes. Keep a probe pending until its original-source CLI regression and full replay pass. Existing release acceptance treats these rejected examples as non-blocking, so this change needs an explicit four-probe completion check; a generic green accepted-example report is insufficient.

At completion, derive counts and current documentation from validation. Keep historical reports and archived task lists unchanged except for a required reviewed traceability migration. Do not repoint the manifest's historical task catalog to a new change that lacks its original task IDs; link the new regressions through supported traceability, or add the smallest validated extension if current metadata cannot represent both owners.

## Risks / Trade-offs

- **Parser recovery hides real errors or loops** → Retain blocking diagnostics outside recorded rules; add boundary/resource cases and run lexer/parser fuzzing.
- **Output changes break prompts or successful writes** → Test success, newline/no-newline, input echo, runtime failure, writer failure, and interpreter reuse at the affected seam.
- **A reference application quirk has no general reduced trigger** → Keep its probe pending rather than introducing a source-identity workaround.
- **Evidence conflict changes acceptance** → Preserve both observations and require an explicit reviewed disposition before implementation or count changes.
- **GUI access requires changes to the remote desktop environment** → Use an existing authorized interactive session or a specifically approved temporary test launcher; no persistent privileged task is part of this design.

## Migration Plan

No data migration or release/tag operation is required. Land implementation, tests, corpus status, and current documentation together after the evidence and quality gates. Rollback reverts their focused commits together; never rewrite the original reference observations or published release history.

## Implementation Gate Status

The four current failures and clean local test baseline are reproduced. The reference executable version and hash match the manifest. No new reduced observation was recorded: a temporary process could not access interactive controls, and a proposed persistent elevated launcher was not authorized. Implementation must wait for the reduced-reference gate; the planning artifacts do not claim that gate passed.
