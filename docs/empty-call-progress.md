# Empty-call reference controls

Nine text-output recordings extend call-form coverage. Seven match the
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

Every verified control emits `BEFORE`, enters the body once, and emits the
caller continuation with a global marker equal to 1. The parameterized controls
also print the parameter before and after assignment. Their exact output
includes the recorded numeric and logical leading spaces.

Two accepted recordings remain pending: `call-empty-value-integer` and
`call-empty-value-real`. Both print exactly `BEFORE\nBODY\n`, then show the
reference's normal completion notice. They emit neither the parameter-value
output nor the assignment output or caller continuation. The implementation
currently rejects these calls during analysis, before any output.

The earlier `procedure-empty-arguments` recording prints zero for a simpler
integer procedure. A general zero-default rule would not explain both sets of
observations. Further reduced controls must establish what ends execution and
whether it depends on reading the missing parameter or the surrounding output
expression. No omitted-argument rule is inferred for other scalar types,
reference parameters, vectors, shadowing or repeated-call state.

Error timing for bare and multi-parameter calls with preceding output also
requires qualification. Unpositioned reference application faults do not
establish a language diagnostic or a source position.

All nine published records retain the original source and execution-panel
bytes, including CRLF. `panel-v1` removes only the complete execution notices
and normalizes CRLF to LF in `stdout.txt`; manifest hashes describe the actual
stored bytes. This evidence slice changes no production behavior and leaves
OpenSpec task 5.4 incomplete.
