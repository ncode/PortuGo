# Record declaration recordings

This slice adds 46 reference recordings and verifies 40 of them, together with
the earlier field and copy cases. Twenty-eight accepted programs compare
original and formatted execution with exact recorded output. Fourteen
rejections pin the diagnostic code and source line.

Four final controls caught detached field references after whole-record
assignment. Record copying now preserves existing field locations, including
vector elements, self-assignment, and replacement while evaluating a later
argument. Original and formatted execution match all four reference results.

Records contain scalar fields, including scalar aliases and grouped names.
Field names ignore case, duplicate definitions retain the first field, and
the first type definition remains effective. Local record definitions shadow
global definitions. Records can be empty, and an optional semicolon may follow
`registro` or `fimregistro`; field-declaration semicolons and a reserved field
name are rejected.

Fields support input, assignment, output, and compatible scalar reference
parameters. Records used as vector elements initialize independently and can
be copied without sharing mutable fields. Whole-record output is empty text.
Record aliases retain separate record identity without copied fields. Named
nested-record declarations create no addressable nested field; inline record
and vector fields are rejected. The specification now distinguishes these
recorded limitations from supported scalar fields.

The implementation stores ordered field layouts in semantic facts, preserves
declaration spelling during formatting, and resolves fields through checked
storage locations. Tests cover case-insensitive first-field lookup, independent
copies and metadata, invalid backing storage, vector slot boundaries, and AST
traversal limits. Empty records count as one cell for the project allocation
guard; this is not a measured reference storage quota.

Five cases retain recorded assignment errors while current analysis reports
a type mismatch before execution: `record-layout-copy-identity`,
`record-boundary-duplicate-real`, `record-alias-value-alias-integer`,
`record-alias-value-alias-copy`, and `record-name-duplicate-integer-first`.
The keyword type-name case `record-layout-alias` also remains pending because
its rejection occurs at a different source position. Their evidence is retained.

The corpus contains 1241 reference recordings: 1189 verified and 52 pending,
plus 16 project-contract records. Nine requirement mappings and 30
bundled-example classifications remain missing. No evidence gate was relaxed.
