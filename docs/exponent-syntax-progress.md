# Source exponent recordings

This slice adds six reference recordings; all six now match the implementation,
including the two positioned real-to-integer assignment diagnostics. Five
accepted cases compare original and formatted execution with exact recorded
output.

Source exponent digits are unsigned. A following sign starts an arithmetic
operator: `1e-2` gives real `-1`, `1e+2` gives real `3`, and `2e-3^2` gives
real `-7`. A bare exponent marker contributes zero while retaining real type.
Recorded controls cover uppercase markers, trailing decimal points, ordinary
positive exponents, spaces around operators, and small decimal literals.

The lexer ends the numeric token before a sign and retains an empty exponent
marker. The parser removes that trailing marker only for value conversion,
after checking the original spelling's integer classification. Token goldens
check operator offsets, and literal-kind tests verify formatting round trips.

The type probes `numeric-source-empty-exponent-type` and
`numeric-source-signed-exponent-type` now defer the recorded real-to-integer
failure to execution and report positioned `R001`. Their reviewed message-window
crops and labeled manual transcripts remain linked to the verified cases.

The corpus now contains 1161 reference recordings: 1113 verified and 48 pending,
plus 16 project-contract records. The evidence gate still needs 10 requirement
mappings and 30 bundled-example classifications. Other literal, comment,
physical-line, and assignment-diagnostic rules remain open.
