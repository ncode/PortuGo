# File-input recordings

This slice retains 42 synthetic reference recordings: 29 accepted programs and
13 positioned rejections. Forty have matching implementation tests. Two accepted
missing-parent-directory cases remain pending: the reference silently continues
without creating a file, while the project's positioned file-error guard reports
`R008`. No successful reference outcome is relabeled as a rejection.

The controls cover literal configuration syntax, repeated and local directives,
relative and nested paths, Windows-1252 text, LF/CRLF and unterminated input,
empty and exhausted files, console fallback, converted recording bytes, echo,
random input, replacement and failure buffering. Ordinary and formatted programs
share the same isolated filesystem expectations. The example reads supplied
synthetic data and leaves it unchanged.

New recordings use 128-byte output blocks. Exact-fill controls distinguish an
immediate flush from a flush deferred until the next byte. A failed program or
replacement leaves the file and its flushed prefix, including a possible trailing
CR; successful completion writes the unfinished block. The headless runtime also
closes handles retained by the reference application after a failure.

The ten buffer and replacement controls record file presence and sizes while
the reference process is still running. The process is then stopped to release
locked handles before reading the file bytes. Each published byte length matches
its preceding size observation. The evidence entries explicitly label this
capture method. Locked or unreadable files are never treated as absent.

Sources, program-only output, initial and final file bytes, and thirteen reviewed
diagnostic-window crops are linked by SHA-256. Diagnostic text is labeled as a
manual transcription. Explicit absence expectations are checked by recording,
manifest validation and CLI replay, with regressions for unexpected files and
conflicting inventories.

Undefined Windows-1252 bytes, unrepresentable recording characters, host file
errors and resource limits have separate project guards. Unreadable-path
compatibility and the missing-parent difference remain unfinished; these tests
do not establish full file-system or full-language conformance.

Validation includes local build, tests, race detection, vet, staticcheck, lint,
both fuzz targets, strict OpenSpec validation and CLI replay of all 1,490 verified
recordings. Native Windows build, vet, both fuzz targets and 3,671 tests pass,
including the same 1,490 replays. The evidence gate still reports three unmapped
requirements and thirteen original examples; its requirements are unchanged.
