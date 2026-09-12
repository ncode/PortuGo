# Independent builtin inventory

The function guide shipped in the public 3.0.7 distribution lists 28 names.
`builtin-inventory.json` records that list, the document's archive path and its
SHA-256. The inventory is derived from the guide rather than the implementation.
The document is independently identifiable even though its text mentions an
older release.

A synthetic reference program exercises every listed name. Tests require a
semantic builtin binding for every inventory entry and compare the original and
formatted program's output with the recording. A controlled mutation removing
`arccos` from both semantic recognition and runtime dispatch makes the test fail.
This catches an omission shared by both implementation layers.

Two additional candidates from the command guide have separate recordings.
`pot(2, 3)` is rejected with positioned `E002` and is covered by a regression.
`div(8, 2)` in an output expression is accepted without producing a value; that
expression form remains pending and is not counted as another documented
function. The word `DIV` is already supported as an infix operator.

The shared immutable descriptor registry now covers all 28 names. Descriptors
carry callable forms, value-parameter types, ordinary arity, optional-call and
absence rules, result types, domain descriptions and evaluator dispatch. No
recorded alias or reference parameter is invented. Construction rejects missing
metadata, unsupported modes, contradictory signatures, duplicate names and
aliases, and missing evaluators. Returned metadata cannot mutate the catalog.

`TestBuiltinRegistryAgreement` reports each name's independently specified
signature, semantic binding and runtime output in lowercase and uppercase.
`TestCatalogDomainResults` checks wrapping, conversion, power underflow,
numeric-domain absence and guarded failures through descriptor dispatch.
Existing recorded fixtures retain argument-order, source-position, formatting
and injected-random-state coverage. `examples/numeric_functions.alg` exercises
the numeric set alongside its integration fixture.

Four controlled mutations fail validation: omitting `arccos`, changing the
`caracpnum` result type, changing `arccos` absence metadata, and routing `cos` to
the sine evaluator. The omission affects semantic and runtime recognition
together, so the independent inventory remains necessary.

The broader numeric domain matrix and undocumented callable forms remain
pending. Catalog agreement does not qualify reference behavior that has not
been recorded.

Published evidence contains synthetic source, two program-only output panels,
one reviewed diagnostic crop and a labeled transcription. The original document
hash was checked against the public distribution before publication.
