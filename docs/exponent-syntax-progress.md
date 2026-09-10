# Source exponent recordings

This slice adds six reference recordings and verifies four of them, together
with the source-exponent case retained by numeric formatting. Five accepted
cases compare original and formatted execution with exact recorded output.

Source exponent digits are unsigned. A following sign starts an arithmetic
operator: `1e-2` gives real `-1`, `1e+2` gives real `3`, and `2e-3^2` gives
real `-7`. A bare exponent marker contributes zero while retaining real type.
Recorded controls cover uppercase markers, trailing decimal points, ordinary
positive exponents, spaces around operators, and small decimal literals.

The lexer ends the numeric token before a sign and retains an empty exponent
marker. The parser removes that trailing marker only for value conversion,
after checking the original spelling's integer classification. Token goldens
check operator offsets, and literal-kind tests verify formatting round trips.

Two type probes, `numeric-source-empty-exponent-type` and
`numeric-source-signed-exponent-type`, remain pending for diagnostic timing.
The reference reports `R001` when real values are assigned to an integer;
current analysis reports `E001` before execution. Both reviewed message-window
crops and labeled manual transcripts remain linked to those pending cases.

The corpus now contains 1161 reference recordings: 1113 verified and 48 pending,
plus 16 project-contract records. The evidence gate still needs 10 requirement
mappings and 30 bundled-example classifications. Other literal, comment,
physical-line, and assignment-diagnostic rules remain open.
