# Frontend recovery progress

Four existing recordings now match their diagnostics and execution phase:

| Recording | Verified behavior |
| --- | --- |
| `alias-caracter-callable` | One `P001` on line 2; invalid callable names do not cascade into a return-type diagnostic. |
| `header-name-without-keyword` | One `P001` at the following token on line 2 after the leading quoted name. |
| `brace-comment-in-expression` | One `P001` on line 3, where the comment removes the closing output parenthesis. |
| `source-suffix-unterminated-same-line` | `BODY` and a newline precede `P001` on line 4 during execution. |

The first three remain rejected by every formatter mode without replacing the
source. The malformed suffix survives two formatting passes with the same
output and runtime diagnostic. LF, CRLF, repeated CR before LF, EOF, and
original-byte mapping through BOM/Windows-1252 decoding have focused coverage.
Existing malformed-body and EOF-recovery tests remain in place.

`brace-inline-math` and `c-comment-in-expression` remain pending. The former
records successful execution without output; the latter records an undeclared
identifier diagnostic before malformed-operator syntax. Existing observations
do not establish an adjacency rule, a general missing-operand grammar, or the
behavior when that identifier is declared. No evidence bytes or expected
outcomes changed, and no broad frontend task is marked complete.

The manifest contains 1,648 reference recordings: 1,567 verified and 81 pending,
plus 19 verified project contracts and two explicit exclusions (1,669 total).
