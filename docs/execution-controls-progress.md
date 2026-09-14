# Timer and breakpoint recordings

This slice adds 58 synthetic recordings. Twenty-six accepted transcripts and
three positioned failures have exact implementation matches, bringing the
corpus to 1,450 verified recordings and 115 recorded cases pending verification.
Sources, program-only panels and thirteen reviewed diagnostic-window crops are
linked by SHA-256. Diagnostic text is labeled as a manual transcription.

The pause and true-debug recordings include an explicit operator continuation
before completion. Headless execution uses a nonblocking host, with injected
host tests asserting request order, source positions, false conditions and
failure handling. No GUI action is inferred from expression-form commands.

Clock probes distinguish sampling from the current command's delay. With a
1,000-millisecond timer, a start followed by timer disable and stop records
1.031 seconds; starting before timer enable and stopping before disable records
1.015 seconds. Fake-clock tests require exactly one second in both cases.
Half-second controls separately qualify calls, local declaration grouping,
interrupted loops and choice branches. Timer state, argument evaluation,
host errors, duration guards and frame cleanup have deterministic tests.

The twelve remaining elapsed-time recordings now replay through explicit project
clock schedules, preserving their reference timing text and timer-delay order.
Ten other rejections now retain output emitted before a later syntax or type
error. Missing timer/debug arguments, unresolved timer modes, and nonlogical
debug conditions are diagnosed when their commands execute. The unchanged
recordings compare exact `BEFORE`/`STEP` prefixes, diagnostic code and line,
and exit status. No clock values or earlier output are removed to make replay
pass.

`TestRecordedDeferredExecutionCommands` runs all ten original and formatted
programs through the typed host harness, verifies that no delay or breakpoint
is requested, and checks formatter idempotence. Additional ordering regressions
start with an active timer and preserve function effects before an unresolved
timer operand; the failing command adds no host action. The generic CLI replay
needs no clock fixture for these programs because their delays remain zero.

Additional controls pin the 10,000-millisecond cap, including inputs of 30,000,
2,147,483,648 and 9,223,372,036,855. Repeated `0.9`-millisecond commands track the
zero-delay control, supporting whole-millisecond truncation. Integer and
half-integer controls also show why short wall-clock measurements cannot be
treated as exact host-call traces: scheduling quantizes their observed delays.

Subprogram timing is now qualified for eight accepted recorded forms spanning
procedure and function calls, interrupted loops, a choice branch, and one or
two local declaration lines. Their original and formatted sources replay with
fixed host-clock schedules, exact elapsed output, and the expected delay order.
Other configuration/declaration combinations, file input, remaining original
examples and the full conformance gate are still unfinished.

The newly qualified elapsed-time controls are `environment-timer-clock-start`,
`environment-timer-clock-stop`, `environment-timer-upper-bound`,
`environment-timer-real-delay`, `environment-timer-large-delay`,
`environment-timer-enormous-delay`, `environment-timer-fraction-loop`,
`environment-timer-rounding-zero`, `environment-timer-rounding-one`,
`environment-timer-rounding-one-half`, `environment-timer-rounding-two-half`,
and `environment-chronometer-minute`. Their reference elapsed values are
supplied to the deterministic host adapter as project clock fixtures;
arbitrary wall-clock measurements remain pending.

The manifest contains 1,657 recorded probes: 1,645 verified and 12 pending,
plus 20 verified project contracts and two explicit exclusions (1,679 total).
