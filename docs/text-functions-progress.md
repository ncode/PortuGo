# Text and character recordings

This slice adds 88 reference recordings: 56 completed programs and 32 positioned
rejections. Seventy-eight match the implementation and have permanent tests;
ten remain pending. Sources, output panels, and reviewed diagnostic-window
crops are linked by hash from the conformance manifest.

The verified cases cover empty and accented text, first and missing matches,
substring clipping, fractional and wrapped bounds, no-value arguments, all 256
character codes, and case conversion of 218 printable Windows-1252 characters.
They also cover argument count/type diagnostics and ordinary left-to-right
evaluation. The shared CP1252 codec preserves source decoding and gives `asc`
ANSI codes. A separate immutable table implements `carac`, including the
reference's control and drawing-character substitutions.

`TestRecordedTextFunctions` checks original and formatted execution;
`TestTextFunctionDiagnostics` pins diagnostic codes and lines. The original
`TABOADA.ALG` now has a matching fixed-input recording and formatter regression.
Source-decoding links cover unambiguous CP1252 input; original-byte position
mapping and BOM position behavior remain separate work.

The ten pending recordings are six initial `caracpnum` conversions and four
argument-order cases. `text-order-bad-bound` and `text-order-extra` retain
reference output before rejection, while current semantic analysis rejects
before execution. `text-order-no-value-bound` and `text-order-no-value-code`
record different no-value/output behavior when later arguments call user
functions. These cases remain visible gaps rather than verified replays.

The corpus contains 963 reference recordings: 912 verified and 51 pending,
plus 16 separate project-contract records. The evidence gate still needs
11 requirement mappings and 30 bundled-example classifications. This slice
does not complete the builtin catalog, dynamic numeric conversion, or the
remaining text/conversion checklist.

## Portable guards

`project.text-code-guards` documents positioned `R007` for a character outside
Windows-1252 and for real substring bounds outside the signed 64-bit conversion
range. Unpositioned reference application faults and complete desktop captures
remain excluded. The later string-limit slice bounds language values at 255
characters. A direct library test retains the case-conversion allocation guard
for oversized Go values that cannot arise from those bounded language strings.
These are portable project contracts, not invented reference diagnostics.
