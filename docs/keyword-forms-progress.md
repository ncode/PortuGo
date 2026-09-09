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

Other newly recorded grammar forms remain pending under task 4.6. The accepted
choice-header forms and the rejected `ate_que`/`até_que` statements are retained
without claiming their parser work is complete. An unconditional break exits
before the `fimrepita` text in one accepted probe; it does not establish that
`fimrepita` is a valid terminator when reached. The bundled prime-number program
also needs bare-output syntax and remains pending under task 17.3.

Tasks 4.4 and 4.6 remain in progress because complete vocabulary and physical-line
grammar obligations extend beyond these spellings. The evidence gate retains
93 missing mappings: 32 requirements and 61 bundled examples.
