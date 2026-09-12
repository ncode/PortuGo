# Random input recordings

This slice adds 56 synthetic recordings and eight original bundled programs.
Thirty-seven synthetic cases have exact implementation matches: 26 accepted
outputs and 11 positioned rejections. The complete corpus now contains 1,403
verified reference recordings and 80 recorded cases still pending implementation
verification. The missing-evidence gate has 18 unmapped entries: five requirements
and 13 bundled examples.

The command, range, expression and repeated-read recordings qualify the numeric
and text input behavior described in [the language reference](language.md).
Synthetic sources, exact reference versions, input bytes, output panels and
diagnostic crops are linked by SHA-256 in the conformance manifest. Diagnostic
text is labeled as a manual transcription; only reviewed diagnostic windows are
retained as images.

The reference differs from its bundled manual: bare `aleatorio` is rejected,
numeric bounds accept expressions, and a single bound sets the lower endpoint
with an upper endpoint of 100. `on` restores defaults and `off` returns to console
reads. Character input contains five uppercase ASCII letters; logical reads
continue using console input. The console transcript retains input echo even
around `eco off`; echo commands themselves remain a separate implementation task.

Repeated real reads establish zero default precision, truncation and clamping of
the precision argument to 0 through 5, and fractional draws added to the base
value. In particular, the first `aleatorio 2,2,3` sample happened to be 2; later
repeated samples include fractions. That first recording is not promoted as a
deterministic exact-output test.

Fixed-output recordings test original and canonical source execution. Scripted
sources pin the implementation's draw bounds, shared-source use and failures;
property checks cover numeric and character domains across many seeds. These
tests do not assert equality with the reference generator's seeds or sequences.
Random sample transcripts remain pending exact replay qualification; no random
values are removed or normalized away to make replay pass.

Eight additional bundled programs complete in the reference, as described in
[the example catalog progress](bundled-examples-progress.md). Their implementation
status remains pending because their output is nondeterministic. File-input
interactions and extreme numeric bounds also remain unqualified; the current
slice does not complete the full conformance change.
