# Numeric character conversion recordings

This slice adds 114 reference recordings: 78 completed programs and 36
positioned rejections. It verifies 109 of those recordings and the six initial
`caracpnum` cases retained by the preceding text-function slice. Five new
recordings remain pending for the shared string-length behavior.

The 115 verified cases cover decimal and hexadecimal spellings, signs, ASCII
and non-ASCII whitespace, exponent syntax, integer/real classification,
overflow, fallback zero, malformed input, and exact diagnostics. The integer
threshold is recorded across every value from 2147483640 through 2147483680,
with separate negative and hexadecimal boundaries.

Analysis now represents an unresolved numeric result without assigning a
concrete type to its runtime value. Consumers check the actual type after
evaluation. Recorded cases cover arithmetic, integer-only built-ins, vector
indices and constant bounds, loop and format bounds, typed returns, procedure
arguments, assignment failures, and output preceding a conversion failure.
One interpreter test reuses the same analyzed program with integer and real
input to catch accidental static classification.

`TestRecordedCharacterConversions` checks original and formatted execution;
`TestCharacterConversionDiagnostics` checks diagnostic codes, source lines,
and retained output. `character_conversions.alg` supplies a runnable example
and a permanent end-to-end fixture. Diagnostic images contain only reviewed
message windows; transcripts are explicitly labeled manual transcriptions.
The long-message transcript records only visible text and identifies clipping.

The five pending cases are `conversion-enormous`,
`conversion-edge-literal-length`, `conversion-edge-joined-length`,
`conversion-edge-joined-number`, and `conversion-edge-suffix-after-255`.
They establish a 255-character reference limit for literals and concatenation.
That limit belongs in the shared string runtime; truncating only conversion
input would leave other string operations inconsistent.

The corpus now contains 1077 reference recordings: 1027 verified and 50
pending, plus 16 separate project-contract records. The evidence gate still
needs 11 requirement mappings and 30 bundled-example classifications. The
remaining catalog, conversion-alias, string, and evaluation-order work stays
open in the conformance task list.
