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

Two further rejection recordings now match before execution:
`accented-identifier` receives one `L001` on declaration line 3, and
`accented-keywords` receives one `L001` at the unsupported `início` on line 2.
The parser applies these restrictions at their grammatical boundaries, retaining
accepted accented keywords and the existing `até_que`/`lógico` diagnostics.
Tests cover original Windows-1252 bytes, UTF-8 and an uppercase UTF-8 BOM form,
including original-byte positions and unchanged token spelling. Every formatter
mode rejects the recorded programs without changing their source. Broader
identifier and keyword tasks remain open; the reference evidence is unchanged.

The manifest contains 1,657 recorded probes: 1,647 verified and 10 pending,
plus 20 verified project contracts and two explicit exclusions (1,679 total).
