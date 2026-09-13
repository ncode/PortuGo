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

The six previously pending rejection controls are now verified. The
`record-layout-alias` control accepts the keyword-shaped alias declaration but
reports `P001` when the alias is used for a variable type. The five assignment
controls report positioned `R001` during execution: incompatible whole-record
copies, record-alias value assignments, and the two duplicate-layout cases.
Their evidence remains recorded, and no broader aggregate semantics are inferred.

The corpus contains 1241 reference recordings: 1195 verified and 46 pending,
plus 16 project-contract records. Nine requirement mappings and 30
bundled-example classifications remain missing. No evidence gate was relaxed.
