# Qualification of recorded deterministic cases

Thirty-four previously pending probes now have focused implementation tests and
incremental replay checks. Original sources, reference screenshots, partial
outputs, normalized transcripts, and their recorded hashes are unchanged.

| Cases | Qualified behavior |
| --- | --- |
| `post-terminator-broken-string` | An unterminated string on a later physical line after the terminator is ignored. Original and formatted source both print `PASS` followed by LF. |
| `reference-widening-followup` | A real reference parameter changes the caller's integer variable to a real value; subsequent arithmetic and assignments reproduce the complete recorded transcript. Original and formatted source agree. |
| `pascal-comment`, `doubled-string-quote`, `assignment-equals` | Syntax rejection occurs on the recorded line before any program output. The retained partial-output files contain only the reference execution wrapper. |
| `c-comment-after-statement`, `c-inline-without-statement` | A C-like opener after an output statement is rejected on that line before printing the statement's text. The reference partial output is empty after removing its execution wrapper. |
| `chronometer-on-off`, `chronometer-repeat-start`, `chronometer-tail` | Original and formatted source reproduce their recorded zero-elapsed-time transcripts using an injected clock that does not advance. |
| `chronometer-repeat-stop`, `environment-chronometer-milliseconds`, `environment-chronometer-seconds`, `environment-chronometer-fractional-seconds` | Original and formatted source reproduce their recorded elapsed-time transcripts using explicit clock schedules and timer delays through the deterministic host adapter. |
| Eight accepted timed-subprogram recordings | Procedure/function calls, interrupted loops, a choice branch, and local declaration layouts reproduce their recorded elapsed text and half-second delay order through explicit clock schedules. |
| `environment-timer-clock-start`, `environment-timer-clock-stop`, `environment-timer-upper-bound`, `environment-timer-real-delay`, `environment-timer-large-delay`, `environment-timer-enormous-delay`, `environment-timer-fraction-loop`, `environment-timer-rounding-zero`, `environment-timer-rounding-one`, `environment-timer-rounding-one-half`, `environment-timer-rounding-two-half`, `environment-chronometer-minute` | Original and formatted source reproduce recorded timer ordering, the ten-second cap, whole-millisecond truncation, repeated-loop delays, and elapsed output beyond one minute through explicit clock schedules and host traces. |

The five rejected programs receive one `P001` parser diagnostic. This code is
the project's syntax category; reference screenshots establish the reported
line, and separate retained partial-output files establish the empty program
output. Malformed programs are rejected by formatting. No formatted executable
is manufactured from a rejected source, and no broader claim is made about
errors following earlier successful statements in other programs.

Chronometer replay uses the existing deterministic observation adapter. The
`expected-clock.json` and `expected-host.json` files are **project fixtures**,
not captured reference traces. The clock fixtures provide the recorded elapsed
sample while the host trace checks the clock-read and timer-delay order. Exact
stdout still compares against the unchanged reference transcripts; no timing
digits are removed or rewritten. This qualifies the listed elapsed-time
transcripts together with the zero-elapsed start, restart, stop, and ignored-tail
cases. It does not promise that an arbitrary CLI wall-clock run will reproduce
these durations, or establish other elapsed-time boundaries.

`TestRecordedDeterministicFollowups` checks execution before and after
formatting and formatting idempotence. `TestRecordedSyntaxRejectionsBeforeOutput`
pins the rejection stage, code, line, and retained reference partial output.
`TestRecordedZeroElapsedChronometers` checks exact transcript bytes, clock-read
counts, and source/formatted execution. Incremental replay additionally checks
CLI exit status and stdout for these deterministic programs, and adapter stdout,
clock schedules and host traces for their timing cases.
`TestRecordedTimedChronometers` applies the same checks to the four elapsed-time
transcripts and verifies their timer-delay calls.
`TestRecordedTimedSubprograms` applies the same checks to the eight accepted
timed-subprogram transcripts and verifies their timer-delay calls.
`TestRecordedTimerTimingBoundaries` applies the same checks to the twelve timer
and chronometer boundary transcripts and verifies their exact host traces.
