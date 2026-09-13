# Scalar type alias recordings

This slice adds 34 reference recordings and verifies 33 of them, together with
the earlier named-type case. Fifteen accepted programs compare original and
formatted execution with exact reference output. Nineteen rejections pin the
diagnostic code and source line.

Scalar aliases support earlier-alias chains, case-insensitive names, separate
type and variable namespaces, local shadowing, vector element declarations,
and compatible scalar arguments. Duplicate declarations retain the first
definition while validating later definitions. Empty type sections and a
semicolon after the section header are accepted. A following `var` is required.

Forward, cyclic, missing, and sibling-local types receive positioned syntax
diagnostics. Vector aliases, named parameter and result types, misplaced
sections, and semicolons after alias definitions are rejected. A long flat
alias chain and AST type-depth checks cover bounded processing independently
of the reference observations.

The `type-alias-duplicate` narrowing-assignment case now reports the recorded
`R001` at execution. Record aliases retain separate identity without scalar
compatibility: the recorded incompatible alias-value assignments now report
`R001` during execution. Broader aggregate alias behavior remains separate
work.

The declaration-scope requirement now links accepted constant, alias, and
record observations together with rejected vector aliases and named callable
headers. The remaining record observations stay scoped to their own declarations
and field behavior. No evidence gate or implementation requirement was relaxed.

The corpus contains 1195 reference recordings: 1147 verified and 48 pending,
plus 16 project-contract records. The evidence gate still needs nine
requirement mappings and 30 bundled-example classifications.
