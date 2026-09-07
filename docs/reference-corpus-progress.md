# Reference corpus progress

OpenSpec group 2 is in progress. The corpus currently inventories 73 requirements
and 73 bundled example filenames, sizes and hashes. Example classifications and
the remaining reduced recordings are pending. The original compatibility
checklist and audit referenced by the plan have not been located in the checkout.
Their inventories must be reconciled before this group can close.

The initial manifest contains 18 reduced probes:

- Eleven accepted observations from the earlier Windows validation have complete
  source, raw output and normalized output artifacts, with existing replay tests.
- Three earlier rejection tests are implemented and passing. Their retained
  original GUI screenshots now supply the missing error evidence, with labeled
  manual transcriptions and source hashes checked against the original captures.
- Four fresh declaration observations have complete recorded evidence and pending
  implementation, as detailed below.

| Probe | Recorded result | Implementation owner |
| --- | --- | --- |
| `constant-declaration` | `const n = 3` followed directly by `inicio` is rejected, with an error reported on line 4. | 6.3 |
| `constant-var` | The same constant followed by a `var` section with an integer variable executes and prints 3. | 6.3 |
| `named-type` | `tipo numero = inteiro`, with a variable of that type, executes and prints 3. | 6.4 |
| `record-field` | A `registro` type, record variable, field assignment and field read execute and print 3. | 6.5 |

These observations establish the demonstrated syntax, not the complete semantics
of constants, aliases or records. Empty sections, copying, reference aliasing,
layout, limits, errors and other interactions still need their own observations.
The rejected constant form does not justify removing constant support: the next
probe demonstrates an accepted form.

The [corpus workflow and schema](../testdata/conformance/visualg-3.0.7/README.md)
describe staging, capture, normalization, validation and bounded CLI replay.
The evidence gate deliberately fails while required recordings or mappings are
missing. A successful tooling test suite is not evidence-mode acceptance or a
claim of complete VisuAlg compatibility.

The tooling review corrections now enforce reference acceptance/rejection,
unique generated-file coverage, mandatory history comparison in validation and
CI, and auxiliary input integrity during capture. These corrections do not fill
the outstanding inventory mappings; group 2 remains in progress. All 18 currently
registered probes have evidence; additional required probes still need to be
prepared and recorded. Inventory validation discovers the entire change's spec
tree independently of the manifest's source list.
