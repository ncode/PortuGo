# Recorded keyword forms

Eighteen reduced observations isolate string-type spelling, accented keywords,
name boundaries, choice syntax, and repeat syntax. Twelve complete successfully
in the reference; six have reviewed GUI rejection evidence and labeled manual
transcriptions.
A separate completed recording adds the original bundled `PRIMOS.ALG` program.

This slice verifies nine accepted reduced programs and the `lógico` rejection.
The recognized alternatives are `caracter`, `função`, `então`, `senão`, `faça`,
`até`, and `não`. The keyword table resolves the alternatives without changing
token text. Scalar type parsing uses the token's canonical type name, so the
string alias works in variables, parameters, and function results without a
second alias table in type checking or execution. Formatting emits canonical
unaccented keywords and `caractere`.

Lexer regressions check lower/uppercase spelling and positions. The reference
programs run from their original Windows-1252 bytes and after formatting. The
recorded type spelling `lógico` remains a positioned `P001` rejection.

The reference rejects `caracter` as a variable or function name, and rejects
`caractere` as a variable name. The longer `caracter_extra` identifier executes
successfully and is verified. Exact diagnostic recovery for the reserved
function-name declaration remains pending under task 4.8. The two variable-name
rejections are now verified by the program-header regressions.

The accepted choice-header forms with optional `faca`/`faça` are now verified.
The rejected `ate_que`/`até_que` statements remain pending under task 4.6.
An unconditional break exits
before the `fimrepita` text in one accepted probe; it does not establish that
`fimrepita` is a valid terminator when reached. A separate reachable-terminator
probe confirms a syntax rejection. The bundled prime-number program now passes
with bare-output syntax, while the complete example sweep remains open.

Four further observations verify case-insensitive `div` as the `\` alias,
including operand types, precedence, left-to-right evaluation, and rejection
as a variable name. Formatting emits `\`. An unfinished function-name recording
remains inconclusive and is excluded from the corpus. Four more original bundled
programs now match fixed-input output, including the combination program using
`div`, two factorial programs, and the minimum-of-three program.

Tasks 4.4 and 4.6 remain in progress because complete vocabulary and physical-line
grammar obligations extend beyond these spellings. Three more observations
verify case-insensitive `fimfunção`, canonical formatting as `fimfuncao`, and
rejection as a variable name. The original bundled `ENCRYPT.ALG` program now
matches its fixed-input recording before and after formatting. The evidence
gate retains 63 missing mappings: 27 requirements and 36 bundled examples.
