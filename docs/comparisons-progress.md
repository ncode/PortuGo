# Comparison recordings

Ninety-five reduced programs recorded with VisuAlg 3.0.7 cover comparisons and
their consumers: 68 complete successfully and 27 report positioned diagnostics.
Sources, output, diagnostic transcriptions, and reviewed diagnostic crops are
linked from the conformance manifest. Program output is normalized only from
CRLF to LF; diagnostic text is explicitly labeled as a manual transcription.

The recordings establish:

- All six relational operators for numeric, logical, and character pairs,
  including integer/real mixing, false-before-true ordering, case, accents,
  prefix lengths, and Windows-1252 byte ordering.
- Mixed scalar comparisons retain their right operand's concrete value while
  keeping a temporary logical category for assignment and conditions.
- Logical storage can retain a new concrete type in variables, fields, and
  vector elements. Reference parameters copy that type back, including through
  nested calls. Branches and functions reading globals observe the stored type.
- Built-ins and user parameters consume concrete values; logical function
  returns consume the logical payload. Rejected consumers preserve earlier
  completed output and use the recorded diagnostic phase and line.
- Record results stored in logical cells have empty layouts. Whole vectors
  require indexing. Generic absence retains the other operand; left numeric
  domain absence skips the right operand, while right domain absence fails.
- Mixed comparison values passed to `e` retain its right operand. `ou` and
  `xou` evaluate both operands before rejecting these combinations. Unary plus
  rejects a comparison; unary minus suppresses its output statement.
- Mixed division retains a real expression category separately from its output
  value. Unary plus rejects an integer fallback; unary minus interprets its
  unsigned 32-bit payload as real bits. Zero, positive, negative, and maximum
  integer controls pin that representation without accessing host memory.

Two exploratory controls with absent values on both sides produced changing
undefined output. Those controls are excluded from the published oracle; a
positioned `E001` project guard covers that boundary. Characters outside
Windows-1252 use a separately tested Unicode ordering extension.

This slice also corrects two existing logical-width error expectations to
preserve output before the recorded runtime `P001`. Other expression and
conformance tasks remain tracked in the OpenSpec checklist.
