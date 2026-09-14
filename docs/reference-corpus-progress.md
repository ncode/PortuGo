# Reference corpus progress

OpenSpec group 2 is complete. The corpus inventories 229 stable records across
73 requirements, 73 bundled example filenames, sizes and hashes, 25 checklist
rows, 48 assumptions, 9 feature records, and one audited defect. All 73 examples now have
recorded dispositions: 60 accepted, 11 unusable, and 2 reviewed non-goals. The
manifest has 1,657 recorded probes and 22 reviewed non-applicable probes; 19
implementation checks remain pending. Every one of the 33 legacy `[VERIFICAR]`
markers has a stable inventory link. The original compatibility checklist is
present in `especificacao-visualg-3.md`; its 25 numbered rows now have stable
manifest links to existing probe coverage. Twenty-six evidenced source
assumptions, the builtin catalog, and one recorded CLI correction are also
linked. Stable feature entries now cover every documented environment-command
family through existing probe coverage. The latest source-obligation recordings
close six previously unrecorded questions; the reference diagnostic behavior
for nonassignable `var` arguments and maximum recursion or stack-display behavior
remain unclaimed, along with unqualified
timing, GUI, and reference-sequence behavior.

The evidence-recording portion of group 2 is complete: all 73 bundled examples
have a recorded disposition, and every manifest probe is either recorded or
reviewed not-applicable. Evidence validation reports 1,657 recorded probes, 22
reviewed not-applicable probes, and no unrecorded entries. The 19 pending
implementation checks remain visible under their implementation tasks.

The September 14 source-obligation slice adds eight reviewed recordings from
the Windows reference application. `const` and `dos` are reserved, a second
main `var` section is rejected, and builtin names cannot be declared as user
procedures. A finite recursive call at depth 32 completes, while `interrompa`
outside a loop is accepted and execution continues. The reference application
raised an internal fault for expression and literal arguments passed to a
`var` parameter; those two observations are retained as labeled GUI
transcriptions without treating the fault as a language-level diagnostic.
The four deterministic rejection cases are now linked to local implementation
tests; the two nonassignable `var` argument guards now have project regression
coverage, while their reference application faults remain separate from any
language-level diagnostic claim. Published GUI evidence contains only reviewed,
sanitized images and text; raw operational captures remain private.

The follow-up zero-value recording confirms that declared integer, real, text,
and logical variables read before assignment produce their type zero values.
The probe is implementation-verified and closes the corresponding legacy
checklist and source-obligation trace.

The legacy marker audit now classifies seven remaining occurrences explicitly:
the declaration-placement question is linked to its existing accepted and
rejected recordings, and the zero-value question is resolved by its new
recording. The prose instructions and grammar comments at their recorded lines
are contextual
references (the grammar aliases point to their existing assumptions). Trace
validation anchors all seven classifications and accepts no unclassified
qualified marker.

The environment inventory now includes stable links for random input, file
input, and console configuration alongside the existing display, timer,
breakpoint, echo, and chronometer entries. Every linked probe is recorded and
implementation-verified; unresolved timing, GUI, and reference-sequence
qualifications remain explicitly pending.

Expression and control-flow inventory links now cover precedence, numeric and
 logical evaluation, comparisons, evaluation order, case ranges, loops, and
 interruption with 299 unique recorded and implementation-verified probes.
Overlapping requirement links share one probe where a recording serves more
 than one category; unqualified operand combinations remain pending.

The traceability handoff is complete: a fresh audit resolves all 229 inventory
entries, 1,679 probes, and 1,813 implementation-test links without stale or
untraced IDs. Shared links are retained where one probe serves multiple
requirements. Pending implementation probes keep their owning tasks and omit
nonexistent tests, while reviewed project-only probes keep their test links.

The reduced lexical and grammar batch is complete across nine requirement
traces and 152 unique implementation-verified probes (149 recorded cases and
three reviewed project-only contracts). The reduced type and declaration batch
is complete across fourteen requirement traces and 273 unique recorded and
implementation-verified probes. Broader unqualified forms remain outside this
recorded scope; the reviewed nonassignable-argument faults remain separate from
the language-level diagnostic contract.

The group 2 documentation reconciliation is complete. The proposal, delta
specifications, design, task ledger, `AGENTS.md`, and language reference now
use the recorded scope: scalar constants and aliases, named records and scalar
fields, assignment aliases, and case ranges are positive only in their
confirmed forms. Rejected variants retain positioned negative coverage;
unrecorded forms keep their pending owning tasks, and no dependent layout or
example obligation is inferred. Manifest dispositions and retired-ID history
were preserved.

Fresh recorder, normalizer, replay, manifest, traceability, quality, and strict
OpenSpec checks are captured in the [current quality report](quality-acceptance-2026-09-14.md).
Evidence validation has zero verified regressions; its published inventory
keeps the current 19 pending implementation probes visible; the two recorded
nonassignable `var` guards are now covered by project regressions.

The clean-checkout rehearsal in
[the release rehearsal report](release-rehearsal-2026-09-14.md) passes the
build, CLI and REPL smoke commands, ordinary and race suites, vet, formatter,
incremental replay, and strict OpenSpec checks. Pending implementation and
handoff tasks continue to block release acceptance.

The clean-checkout rehearsal in
[the release rehearsal report](release-rehearsal-2026-09-14.md) passes the
build, CLI and REPL smoke commands, ordinary and race suites, vet, formatter,
incremental replay, and strict OpenSpec checks. Pending implementation and
handoff tasks continue to block release acceptance.

The initial manifest contained 18 reduced probes:

- Eleven accepted observations from the earlier Windows validation have complete
  source, raw output and normalized output artifacts, with existing replay tests.
- Three earlier rejection tests are implemented and passing. Their retained
  original GUI screenshots now supply the missing error evidence, with labeled
  manual transcriptions and source hashes checked against the original captures.
- Four fresh declaration observations have complete recorded evidence and pending
  implementation, as detailed below.

| Probe | Recorded result | Implementation owner |
| --- | --- | --- |
| `constant-declaration` | `const n = 3` followed directly by `inicio` is rejected, with an error reported on line 4. | 6.3 |
| `constant-var` | The same constant followed by a `var` section with an integer variable executes and prints 3. | 6.3 |
| `named-type` | `tipo numero = inteiro`, with a variable of that type, executes and prints 3. | 6.4 |
| `record-field` | A `registro` type, record variable, field assignment and field read execute and print 3. | 6.5 |

These observations establish the demonstrated syntax, not the complete semantics
of constants, aliases or records. The newer batch below covers additional empty
sections, copying, reference coercion, bounds, and literal forms; layout, limits,
errors and other interactions still require broader coverage.
The rejected constant form does not justify removing constant support: the next
probe demonstrates an accepted form.

The [corpus workflow and schema](../testdata/conformance/visualg-3.0.7/README.md)
describe staging, capture, normalization, validation and bounded CLI replay.
The evidence gate deliberately fails while required recordings or mappings are
missing. A successful tooling test suite is not evidence-mode acceptance or a
claim of complete VisuAlg compatibility.

The tooling review corrections now enforce reference acceptance/rejection,
unique generated-file coverage, mandatory history comparison in validation and
CI, and auxiliary input integrity during capture. Those earlier corrections did
not fill the outstanding inventory mappings; the current audit above closes
them. All registered reference probes have evidence, and inventory validation
discovers the entire change's spec tree independently of the manifest's source
list.

Five additional [project tooling entries](conformance-project-evidence.md) link
existing provenance, recording, normalization, traceability, and phase-validation
tests with explicit reference non-applicability reasons. Existing recorded
observations also now trace reference authority and feature dispositions. These
close seven missing requirement mappings without adding language observations
or relaxing the evidence gate.

The [September 9 batch](reference-observations-2026-09-09.md) adds 61 recorded
probes: 41 accepted and 20 rejected. A focused comment follow-up adds another 24
(15 accepted, 9 rejected). A call follow-up adds 11 more usable observations
(4 accepted, 7 rejected). The argument follow-ups add 35 observations
(28 accepted, 7 rejected). The naming follow-ups add 15 observations
(7 accepted, 8 rejected). The return follow-ups add 45 observations
(36 accepted, 9 rejected). The vector follow-ups add 26 observations
(10 accepted, 16 rejected). The first bundled-example sweep adds eleven recorded
programs (8 accepted, 3 unusable). Keyword follow-ups add eighteen reduced
observations (12 accepted, 6 rejected), plus one completed bundled-program run.
The corpus now contains 265 reference probes (215 verified, 50 pending
implementation), plus the five project-tooling entries.
The missing mappings at the corpus layer cover 44 requirements and all 73 bundled
examples. The stable storage requirement ID is retained while its linked heading
and text now reflect the recorded acceptance beyond the former 500-slot claim.

The [REPL input slice](repl-input-progress.md) also traces the verified common
reader and EOF paths separately from pending blank-line behavior. The remaining
missing mappings cover 43 requirements and all 73 bundled examples. Pending
project behavior remains visible and still fails implementation acceptance.

The string-literal slice verifies five accepted literal-backslash forms and the
rejected backslash-quote form. Original and formatted programs produce the same
recorded output, and repeated formatting preserves both values and program names.
Local build, formatting, vet, lint, ordinary/race tests, and 30-second lexer and
parser fuzz runs pass. Native Windows build, vet, ordinary tests, both 30-second
fuzz runs, and all 20 verified CLI reference replays pass. Raw environment logs
remain outside the repository.

The assignment slice also verifies `:=` through the existing assignment token,
semantic checks, and runtime. Token spelling and positions are retained, while
canonical formatting uses `<-`.
Native Windows build, vet, ordinary tests, both 30-second fuzz checks, and all
21 currently verified CLI reference cases pass for this slice.

The comment slice verifies 26 accepted programs and three rejected strings
containing `//`. Regression tests also check that comments preserve the next
token's line and column with LF and CRLF. Comment retention in formatted source
and the reference's incomplete-expression behavior remain pending.

Local build, formatting, vet, static analysis, lint, ordinary/race tests, both
30-second fuzz runs, and strict OpenSpec validation pass. Native Windows build,
vet, ordinary tests, both 30-second fuzz runs, and all 50 verified CLI reference
cases pass. The full evidence gate still reports only the 116 missing mappings.

The parameterless-call slice verifies eight accepted programs and six rejected
call contexts. This includes the reference's function-name priority over local
variables and parameters, and unchanged execution after formatting. Three
temporary-reference-argument probes caused internal application faults and were
excluded from language acceptance/rejection evidence. Fresh processes isolate
each follow-up probe. The declaration-line error positions in
`bare-procedure-missing-argument` and `procedure-name-priority` remain pending.
The new scope observation leaves 42 requirements and 73 examples unmapped.

Local build, formatting, vet, static analysis, lint, ordinary/race tests, both
30-second fuzz runs, and strict OpenSpec validation pass for the call slice.
Native Windows build, vet, ordinary tests, both 30-second fuzz runs, and all
64 verified CLI reference cases pass on the same production source and manifest.

The argument slice verifies 25 accepted programs, seven static call rejections,
and one runtime type-change failure with its preceding output. Scalar `var`
parameters use the recorded copy-in/copy-out behavior, including conversion,
ordered copy-back to repeated destinations, and independent nested/recursive
parameter storage. The new recordings also verify argument evaluation order,
reference element capture, lexical globals, local shadowing, and mutual calls.
Seven GUI diagnostics were individually reviewed and transcribed; raw window
inventories and operational logs remain private. Guard tests cover failed setup,
interpreter reuse, and out-of-range numeric arguments.

The procedure declaration-line argument errors are now covered. Empty-argument
edge cases, procedure-name collisions, integer operators after reference type
changes, and broader return behavior remain pending. The new mappings leave
40 requirements and all 73 bundled examples unmapped.

Local build, formatting, vet, static analysis, lint, ordinary/race tests, both
30-second fuzz checks, and strict OpenSpec validation pass. Native Windows build,
vet, ordinary tests, both 30-second fuzz checks, and all 97 verified CLI reference
cases pass on the same production source and manifest. The evidence gate reports
the 113 missing mappings and no other failures.

The naming slice verifies seven accepted programs and nine rejections, including
the earlier procedure-name collision. Variables and callables have separate
namespaces; functions and procedures share a namespace. Calls follow the
recorded contextual priority, and assignments beginning with procedure names
report the declaration-line error for scalar and indexed targets. Duplicate
callable declarations are rejected before ambiguous bodies produce additional
diagnostics. Eight new GUI diagnostics were individually reviewed and their
transcriptions are labeled manual.

Empty-call repeat observations remain private while their early-termination
behavior is investigated; this slice does not infer default arguments from
those incomplete executions. Existing published empty-call observations remain
pending. The missing inventory count is unchanged at 113.

Local build, formatting, vet, static analysis, lint, ordinary/race tests, both
30-second fuzz checks, strict OpenSpec validation, and all 113 verified CLI cases
pass. Native Windows build, vet, 549 tests, both 30-second fuzz checks, and all
113 verified CLI cases pass on the same production source and manifest.

The return slice verifies 33 accepted programs and seven rejection diagnostics.
`retorne` updates the active result while statements and finite loop iterations
continue; the function terminator ends the call. Completed calls retain
same-type results at each depth, including across different names and frame
sizes. Recursive results, parameters, and locals remain independent. Argument
expressions run before the callee selects its retained result. Error guards
cover failures after a result assignment, skipped parameter copy-back, and reuse
of the interpreter for a new program. The recorded colonless `caso` form is
accepted and formatting preserves execution.

Nine GUI rejection diagnostics were individually reviewed and manually
transcribed. An unfinished infinite-loop observation and six cross-type
fallthrough captures were excluded. The latter expose internal storage; their
raw output and images remain private. Cross-type fallthrough uses a documented
fresh-zero project guard rather than claiming reference equivalence.

Bare-return diagnostics and nested-write formatting remain pending, including
the new observations that depend on them. New mappings for function results,
call-frame safety, and nested result updates leave 37 requirements and all
73 bundled examples unmapped. The evidence gate remains enabled.

Local build, formatting, vet, staticcheck, lint, ordinary/race tests, both
30-second fuzz checks, strict OpenSpec validation, and all 153 verified CLI
replays pass. Native Windows build, vet, 634 tests, both 30-second fuzz checks,
and all 153 verified CLI replays pass on the same production source and manifest.
The native snapshot also confirms removal of the obsolete return-path checker.
The evidence gate reports the 110 missing mappings and no other failures.

The frontend regression slice verifies 12 previously recorded observations:
eight accepted programs and four positioned rejections. Permanent tests cover
empty `var` sections, ignored post-termination words and statements, underscore
and mixed-case identifiers, exponent and decimal-point forms, unterminated
strings, and invalid inline slash/star usage. Accepted programs preserve their
output after canonical formatting. This establishes the recorded execution
behavior; comment and suffix text retention remain pending.

No recordings or production behavior changed in this slice. The corpus contains
165 verified reference cases and 44 pending cases. The missing mapping count
remains 110, and no broader frontend task is marked complete by this subset.

Local build, formatting, vet, staticcheck, lint, ordinary/race tests, both
30-second fuzz checks, strict OpenSpec validation, and all 165 verified CLI
replays pass. Native Windows build, vet, 660 tests, both 30-second fuzz checks,
and all 165 CLI replays pass with the same regression source and manifest.

The vector slice verifies 34 cases, including eight earlier observations and
26 new recordings. One- and two-dimensional vectors retain zero and positive
non-one bounds; an omitted second index selects that dimension's lower bound.
Recorded unsigned literal bounds are accepted, while signed, fractional,
parenthesized, arithmetic, reversed, and third-dimension declarations are
rejected. Whole-vector assignment is rejected. Regression tests cover accepted
500/501/5000/5001-element vectors, zero initialization, extra indices, and runtime
bounds failures with their preceding output. Accepted programs retain their
output after formatting, and the example combines recorded dimension and
index forms.

Sixteen new GUI diagnostics were individually inspected and manually transcribed.
Only reviewed synthetic programs, output, and language diagnostics are retained;
raw operational inventories remain private. New zero-initialization and storage
access mappings leave 35 requirements and all 73 bundled examples unmapped.
Constant bounds, aggregate parameters, and corrupted-layout defenses remain
pending; this slice completes only task 7.9.

Local build, formatting, vet, staticcheck, lint, ordinary/race tests, both
30-second fuzz checks, strict specification validation, and all 199 verified CLI
replays pass. Native Windows build, vet, 732 tests, both 30-second fuzz checks,
and all 199 CLI replays pass. All 112 transported branch files match the local
snapshot before this validation note. The evidence gate reports the 108 missing
mappings and no other errors.

The defensive storage slice validates vector metadata, element kinds, dimension
products, and exact backing lengths before offset arithmetic. Direct regression
tests cover malformed layouts, both integer extremes, bounds, and every element
of small matrices. Interpreter regressions inject truncated storage before a
read, assignment, input, or reference argument; all return positioned `R003`
without a panic, element changes, lost preceding output, or consumed input.

These synthetic corruption cases are linked as `project.vector-storage`, with
an explicit reference non-applicability rationale. All reference recordings and
expectations remain unchanged. Tasks 7.6 and 7.7 are complete; allocation and
slot-accounting work remains pending. The corpus retains 199 verified reference
cases and 36 pending cases, with 108 inventory mappings still missing.

Local build, formatting, vet, staticcheck, lint, ordinary/race tests, both
30-second fuzz checks, strict specification validation, and all 199 verified CLI
replays pass. Native Windows build, vet, both fuzz checks, 756 tests, and all 199
CLI replays pass on the final production source. The final manifest correction
also passes a fresh native test/replay run, with all nine transported files
matching locally before this validation note. No extra evidence-gate errors
remain beyond the 108 missing inventory mappings.

The automatic REPL slice preserves blank lines inside unfinished programs and
submits immediately on a real terminator token. New transcripts cover LF/CRLF,
consecutive programs, shared input, terminator lookalikes, exact CLI prompt
ordering, source positions, and recovery after lexical, syntax, semantic,
runtime, and execution-budget failures. Decoded line framing prevents an encoded
longer identifier from prematurely submitting an unfinished program.

The pending project blank-line entry is now verified. A new project submission
entry traces automatic submission and recovery without changing reference
recordings or expectations. Tasks 16.3, 16.4, and 16.5 are complete. Source-size
and host-failure recovery, environment input modes, and formatter retention
remain pending. The missing inventory now covers 33 requirements and all 73
bundled examples; the reference corpus remains at 199 verified and 36 pending.

Local build, formatting, vet, staticcheck, lint, ordinary/race tests, both
30-second fuzz checks, strict specification validation, and all 199 verified CLI
replays pass. Native Windows build, vet, 766 tests, both fuzz checks, and all 199
CLI replays pass with all 11 transported files matching locally before these
validation notes. The required-evidence gate reports the 106 missing mappings
and no other errors.

The REPL recovery slice adds bounded recovery after source-size failures and
regressions for transient host output failures, invalid decoded bytes, and a
missing diagnostic writer. Rejected oversized source reports one `E900` and is
discarded through the next program header without losing its buffered input.
Header lookalikes, a terminator crossing the remaining allowance, EOF, and
cancellation are covered. Exact CLI output and exit status remain tested.

Tasks 16.1 and 16.6 are complete. The existing project submission entry gains
only additive task and test links; all reference evidence, expected outcomes,
and inventory mappings remain intact. The required-evidence gate still has
106 missing mappings. Environment input modes and formatter retention remain
pending under the corresponding tasks.

Local build, formatting, vet, staticcheck, lint, ordinary/race tests, both
30-second fuzz checks, strict specification validation, and all 199 verified CLI
replays pass. Native Windows build, vet, 777 tests, both fuzz checks, and all 199
CLI replays pass with all 11 transported files matching locally before these
validation notes. The evidence gate reports the same 106 missing mappings and
no other failures.

The first bundled-example slice records eleven original sources against their
catalog hashes: eight accepted programs and three unusable files with reviewed
GUI diagnostics. Six accepted programs have byte-exact original/formatted
regressions. The new catalog check rejects stale IDs, source hash/size changes,
false classifications, and unreviewed exclusions. Its regression test is linked
from the existing project traceability entry.

Local build, formatting, vet, staticcheck, lint, ordinary/race tests, both
30-second fuzz checks, strict specification validation, and all 205 verified CLI
replays pass. Native Windows build, vet, 803 tests, both fuzz checks, and all 205
CLI replays pass. The final trace-link metadata also passes a fresh native
test/replay run, with all 46 transported files matching before this validation
note. The evidence gate reports 94 missing mappings (32 requirements and 62
examples), with no other errors. The full inventory and example tasks remain
in progress.

The keyword slice adds the recorded `caracter` type spelling and six accented
keyword aliases. Canonical scalar type parsing keeps declarations, parameters,
and function results consistent while lexer tests retain spelling and positions.
Ten newly recorded keyword cases are verified, including the rejected accented
logical type and an accepted longer identifier. Reserved-name recovery, choice
and repeat grammar, and the newly recorded bundled prime-number program remain
pending. The corpus contains 265 recorded probes, of which 215 are verified.

Local build, formatting, vet, staticcheck, lint, ordinary/race tests, both
30-second fuzz checks, strict specification validation, and all 215 verified CLI
replays pass. Native Windows build, vet, and both fuzz checks pass on the final
production source; the completed name-boundary evidence then passes a fresh
826-test run with all 215 reference replays. All 79 transported files match
locally before this validation note. The evidence gate reports 93 missing
mappings (32 requirements and 61 examples), with no other errors.

The empty-argument subprogram follow-up now verifies three additional accepted
records through the original and formatted execution regression. A simple
one-argument numeric procedure exposes its typed zero; a one-argument numeric
function propagates no value when that absent parameter is returned, preserving
the recorded body output and caller stop. Bare, multi-parameter, and other
omitted-argument forms remain pending.

The record-rejection follow-up now verifies six previously pending recorded
controls. One keyword-shaped alias use reports `P001` at its variable
declaration, while five incompatible or duplicate-layout assignments report
positioned `R001` during execution. The evidence remains recorded and broader
aggregate semantics remain pending.

The timed-subprogram follow-up now verifies eight accepted elapsed-time
recordings through original and formatted execution. Fixed host-clock schedules
preserve each recorded elapsed message and half-second delay order across
procedure/function calls, interrupted loops, a choice branch, and local
declaration layouts. Other variable timing combinations remain pending.

The replay adapter now rejects a run when a deterministic clock schedule is
exhausted, preventing an incomplete fixture from silently falling back to
elapsed time.

Conformance artifact reads now enforce the repository size limit while reading,
so files that grow after metadata inspection are rejected before their bytes are
used.

Retained replay-output artifacts now use the smaller observation limit during
the read, avoiding a larger allocation before an oversized expectation is
rejected.

Replay source artifacts now apply the 64 KiB execution profile while reading,
before a candidate process is staged.
