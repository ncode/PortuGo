# Timer and breakpoint recordings

This slice adds 49 synthetic recordings. Twenty-six accepted transcripts and
three positioned failures have exact implementation matches, bringing the
corpus to 1,450 verified recordings and 106 recorded cases pending verification.
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

Ten new elapsed-time recordings remain pending exact replay because wall-clock
values vary. Ten other rejections retain output emitted before a later syntax
or semantic error; the implementation's earlier diagnostics still differ.
Those prefixes are preserved, and the cases remain pending. No clock values or
earlier output are removed to make replay pass.

Subprogram timing is qualified for the recorded ordinary local-variable forms.
Other configuration/declaration combinations, file input, remaining original
examples and the full conformance gate are still unfinished.
