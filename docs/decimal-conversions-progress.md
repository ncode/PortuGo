# Decimal conversion fallback recordings

This slice adds 21 reference recordings and verifies all of them, together with
the two conversion differences retained by the precision slice. Twenty-one
accepted cases compare original and formatted execution; two positioned
diagnostics check the nonzero-prefix controls.

After hexadecimal handling, `caracpnum` ignores leading zeros when checking
for a nonzero integral prefix. If none remains, it returns integer zero before
parsing a fraction, exponent, or malformed suffix. This covers ordinary
fractions, signed zero, long tiny decimals, comma forms, overflow spellings,
repeated points, underscores, and alphabetic suffixes. Nonzero integral
prefixes continue through normal parsing and validation.

Integer-assignment probes establish the fallback's concrete integer type.
The real-valued `"6e-18"` control still fails integer assignment with `R001`;
the malformed `"01.5x"` control still receives `R007`. The fix adds one prefix
normalization in the shared conversion evaluator. Hexadecimal handling,
numeric source literals, and input parsing retain their existing paths.

The two diagnostic images are reviewed message-window crops with labeled
manual transcripts. Other published evidence contains synthetic sources and
program-only output panels. Raw operational evidence remains private.

The corpus now contains 1155 reference recordings: 1108 verified and 47
pending, plus 16 project-contract records. The conformance gate still needs
10 requirement mappings and 30 bundled-example classifications. Conversion
catalog aliases and the separately recorded signed-exponent source grammar
remain open.
