# Qualification of recorded deterministic cases

Ten previously pending probes now have focused implementation tests and
incremental replay checks. Original sources, reference screenshots, partial
outputs, normalized transcripts, and their recorded hashes are unchanged.

| Cases | Qualified behavior |
| --- | --- |
| `post-terminator-broken-string` | An unterminated string on a later physical line after the terminator is ignored. Original and formatted source both print `PASS` followed by LF. |
| `reference-widening-followup` | A real reference parameter changes the caller's integer variable to a real value; subsequent arithmetic and assignments reproduce the complete recorded transcript. Original and formatted source agree. |
| `pascal-comment`, `doubled-string-quote`, `assignment-equals` | Syntax rejection occurs on the recorded line before any program output. The retained partial-output files contain only the reference execution wrapper. |
| `c-comment-after-statement`, `c-inline-without-statement` | A C-like opener after an output statement is rejected on that line before printing the statement's text. The reference partial output is empty after removing its execution wrapper. |
| `chronometer-on-off`, `chronometer-repeat-start`, `chronometer-tail` | Original and formatted source reproduce their recorded zero-elapsed-time transcripts using an injected clock that does not advance. |

The five rejected programs receive one `P001` parser diagnostic. This code is
the project's syntax category; reference screenshots establish the reported
line, and separate retained partial-output files establish the empty program
output. Malformed programs are rejected by formatting. No formatted executable
is manufactured from a rejected source, and no broader claim is made about
errors following earlier successful statements in other programs.

Chronometer replay uses the existing deterministic observation adapter. The
new `expected-host.json` files are **project fixtures**, not captured reference
traces: start/stop requires two clock reads, and restart/stop requires three.
The adapter returns the same instant for each read in these programs. Exact
stdout still compares against the unchanged reference transcripts; no timing
digits are removed or rewritten. This qualifies zero elapsed time, start,
restart, stop, and ignored tails. It does not promise that an arbitrary CLI
wall-clock run will measure zero, or establish other elapsed-time boundaries.

`TestRecordedDeterministicFollowups` checks execution before and after
formatting and formatting idempotence. `TestRecordedSyntaxRejectionsBeforeOutput`
pins the rejection stage, code, line, and retained reference partial output.
`TestRecordedZeroElapsedChronometers` checks exact transcript bytes, clock-read
counts, and source/formatted execution. Incremental replay additionally checks
CLI exit status and stdout for the seven deterministic programs, and adapter
stdout and host traces for the three qualified chronometer programs.
