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

The complete shared descriptor registry and signature/domain agreement checks
remain pending. This inventory establishes a required set, without claiming
that all callable forms or undocumented aliases have been exhausted.

Published evidence contains synthetic source, two program-only output panels,
one reviewed diagnostic crop and a labeled transcription. The original document
hash was checked against the public distribution before publication.
