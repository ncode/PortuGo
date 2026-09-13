# Empty-call reference controls

Twelve text-output recordings extend call-form coverage. All twelve match the
implementation in original and formatted execution through
`TestRecordedCallArguments` in `call_arguments_test.go`:

| Probe | Recorded result |
| --- | --- |
| `call-supplied-value-integer` | The parameter changes from 7 to 9; the caller remains 7. |
| `call-supplied-value-real` | The parameter changes from 2.5 to 9.5; the caller remains 2.5. |
| `call-supplied-value-logical` | The parameter changes from true to false; the caller remains true. |
| `call-supplied-value-character` | The parameter changes from `seed` to `set`; the caller remains `seed`. |
| `call-supplied-reference-integer` | The parameter changes from 7 to 9; the caller receives 9. |
| `call-supplied-reference-real` | The parameter changes from 2.5 to 9.5; the caller receives 9.5. |
| `call-empty-zero-parameters` | `P()` enters a parameterless procedure and returns to the caller. |
| `call-empty-value-integer` | A parenthesized empty call enters the body; printing the absent integer parameter ends execution after `BEFORE` and `BODY`. |
| `call-empty-value-real` | A parenthesized empty call enters the body; printing the absent real parameter ends execution after `BEFORE` and `BODY`. |
| `procedure-empty-arguments` | A parenthesized empty call enters a one-argument integer procedure; its direct output exposes the typed zero. |
| `function-empty-arguments` | A parenthesized empty call enters a one-argument integer function; returning its absent parameter produces no enclosing output. |
| `function-empty-arguments-continuation` | The function body prints `BODY` after caller output `BEFORE`; its absent result suppresses the enclosing write and `AFTER`. |

The six supplied-argument controls emit `BEFORE`, enter the body once, and
emit the caller continuation with a global marker equal to 1. The parameterized
controls also print the parameter before and after assignment. Their exact
output includes the recorded numeric and logical leading spaces.

The two multi-argument empty numeric procedure controls print exactly
`BEFORE\nBODY\n`, then complete normally. They emit neither the
parameter-value output nor the assignment output or caller continuation. The
simple one-argument procedure control prints its typed zero and returns. The
function controls propagate no value when the absent parameter is returned: the
function body output remains visible, while the enclosing write and following
caller statement are suppressed.

A general zero-default rule would not explain both procedure output forms or
the function no-value result. This implementation therefore qualifies only a
parenthesized call with one by-value numeric parameter and no arguments for
the recorded procedure and function forms. A direct one-argument procedure
write exposes typed zero; reaching the absent value in composed procedure
output or in a function result ends the enclosing output as recorded. No
omitted-argument rule is inferred for bare calls, multi-parameter calls, other
scalar types, reference parameters, vectors, shadowing or repeated-call state.

Error timing for bare and multi-parameter calls with preceding output also
requires qualification. Unpositioned reference application faults do not
establish a language diagnostic or a source position.

All twelve published records retain the original source and execution-panel
bytes, including CRLF. `panel-v1` removes only the complete execution notices
and normalizes CRLF to LF in `stdout.txt`; manifest hashes describe the actual
stored bytes. Broader call-frame and numeric-reference cases remain under
OpenSpec task 5.4.
