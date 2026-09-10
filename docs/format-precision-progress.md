# Numeric field precision recordings

This slice adds 32 reference recordings and verifies 29 of them, together with
the three precision cases retained by the string-width slice. All 32 verified
cases compare original-source and formatted-source execution with exact
reference output bytes.

Requested decimal counts are capped at 216 before expansion. Fixed real fields
round the stored binary value, preserve negative zero, and pad decimal places
beyond the conversion's digit budget with zeros. Recordings distinguish false
decimal ties from exact halves and cover binary-exponent transitions, small
values produced by arithmetic, and large whole-number rounding. At absolute
value `2^120`, fields switch to the recorded scientific format, including its
minimum width, sign column, fractional-digit limit, and four-digit exponent.

The previous precision-allocation rejection is replaced by the measured cap.
Input-line and direct-library allocation safeguards remain, as does the shared
pending-output limit across nested writes. The example and interpreter fixture
show rounding, negative zero, trailing-zero padding, and scientific notation.

Three new observations were retained outside formatting:

- `format-precision-small-values` exposes signed exponent tokens being parsed
  differently by the reference source parser.
- `format-precision-small-rounding` and `format-identity-tiny-decimals` expose
  long decimal strings converting to zero where equivalent exponent strings
  retain their numeric value. The subsequent
  [decimal-conversion slice](decimal-conversions-progress.md) verifies these
  as the zero-integral-prefix fallback, which also applies to ordinary fractions.

Arithmetic and exponent-string control cases show that these differences
precede formatting. No magnitude-based zeroing was added to the formatter.
Raw machine diagnostics and desktop captures remain outside the repository;
published evidence contains synthetic sources and program-only output panels.

The corpus now contains 1134 reference recordings: 1085 verified and 49 pending,
plus 16 project-contract records. The evidence gate still needs 10 requirement
mappings and 30 bundled-example classifications. I/O diagnostic timing,
numeric parsing, and the remaining conformance work stay open.
