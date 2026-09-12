# Console configuration recordings

Sixteen VisuAlg 3.0.7 recordings cover `dos`: 13 accepted programs and three
positioned rejections. The accepted cases cover the program header, repetition,
mixed case, ignored line tails (including `on`, `off`, and an unmatched quote),
a procedure header, an expression, and a control without the directive.

The rejections distinguish a declaration-section error from executable-body
errors. The latter preserve any output written before the rejected directive.
Sources, program-only output panels, diagnostic-window crops, and labeled manual
transcriptions are retained in the conformance corpus. Raw desktop captures and
operational details are excluded.

`TestRecordedDisplayCommands` compares original and formatted execution;
`TestDisplayCommandDiagnostics` checks the code, source line, and prior output.
Interpreter tests separately pin the injected host contract: configuration order,
per-call local configuration, repeated runs, and sanitized positioned failures.
The default host keeps output in the supplied writer. These host-call details are
implementation contracts; the recordings do not establish a GUI mode transition
or an `on`/`off` option distinction.

Other environment commands and the remaining conformance inventory stay pending.
