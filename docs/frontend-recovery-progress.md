# Frontend recovery progress

Six existing recordings now match their diagnostics and execution phase:

| Recording | Verified behavior |
| --- | --- |
| `alias-caracter-callable` | One `P001` on line 2; invalid callable names do not cascade into a return-type diagnostic. |
| `header-name-without-keyword` | One `P001` at the following token on line 2 after the leading quoted name. |
| `brace-comment-in-expression` | One `P001` on line 3 when the truncated write executes. |
| `source-suffix-unterminated-same-line` | `BODY` and a newline precede `P001` on line 4 during execution. |
| `brace-inline-math` | The adjacent numeric brace suppresses the current write without stopping the program. |
| `c-comment-in-expression` | The name between the adjacent operator pairs receives `E002` during execution. |

Qualification controls separate numeric adjacency from spaced and parenthesized
braces, and distinguish undeclared names from declared values, literals, and
function calls. Earlier output survives the two diagnostic forms. Operand calls
retain their effects before an absent result suppresses the current write;
subsequent statements continue after successful suppression.

The two header/name errors remain rejected by every formatter mode without
replacing the source. Recovered expressions, incomplete writes, and the malformed
terminator suffix retain their spelling and execution behavior across two
formatting passes. Tests cover LF, CRLF, repeated CR before LF, and original-byte
positions through BOM/Windows-1252 decoding. Recovery operands participate in
AST traversal limits, and the two newly qualified forms seed parser fuzzing.

Other adjacent identifier forms and malformed operator combinations remain
under qualification. No generic missing-operand rule is claimed. The two
previously pending implementation entries are promoted; reference sources,
recordings, expected outcomes, and task checkboxes remain unchanged. Raw
qualification captures are kept outside the repository.

The manifest contains 1,648 reference recordings: 1,573 verified and 75 pending,
plus 19 verified project contracts and two explicit exclusions (1,669 total).
