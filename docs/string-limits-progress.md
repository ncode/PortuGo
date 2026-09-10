# String value and width recordings

This slice adds 25 reference recordings and verifies 21 of them, together with
the five string-boundary cases retained by numeric conversion. Twenty-six
cases now have original-source and formatted-source execution regressions.

The recordings cover 254, 255, 256, and 511-character literals; accented
characters at the cutoff; concatenation; search, slicing, and comparison;
constants; assignments; parameters and returns; vector elements; and input.
Input echoes the complete line while storing the first 255 characters. A
separate interpreter test checks both the stored value and the complete echo.

Output widths are capped at 255 before padding, including a recorded request
of 16777217. Numeric padding keeps its right alignment; character padding
keeps its left alignment. A shared character-limit helper preserves UTF-8
boundaries and copies truncated prefixes so stored values do not retain a
large input buffer.

Three large-precision recordings were retained for the subsequent
[precision slice](format-precision-progress.md): `string-field-real-precision`,
`string-field-integer-precision`, and `string-field-wide-precision`. The positioned
`string-size-format-values` rejection also remains pending because reference
output precedes a later invalid logical format, while current analysis rejects
the program before execution. Its reviewed diagnostic crop and labeled manual
transcript establish the previously missing I/O-validation evidence mapping.

An interrupted batch contributed only completed, source-hash-verified cases.
A malformed constant declaration and an input request above the recorder's
limit were excluded and replaced by corrected, completed probes. Raw machine
diagnostics and desktop captures remain outside the repository.

Allocation checks remain covered at their actual boundaries: bounded values
still count toward the aggregate pending-output budget, oversized input lines
are rejected, and direct library casing rejects excessive expansion before
allocation. An indefinitely growing concatenation now reaches the execution
step limit because its string value stays bounded.

The corpus contains 1102 reference recordings: 1053 verified and 49 pending,
plus 16 project-contract records. The evidence gate still needs 10 requirement
mappings and 30 bundled-example classifications. Remaining catalog,
and evaluation-order work stays open in the conformance task list.
