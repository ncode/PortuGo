# File-input recordings

This slice retains 51 synthetic reference recordings: 38 accepted programs and
13 positioned rejections. Forty-five have verified implementation coverage;
six access-constrained recordings remain pending generic CLI replay. The two
missing-parent-directory recordings now continue without creating a file,
including a program with no reads and one that consumes and echoes console
input. Original and formatted sources retain the exact recorded output and
file absence. No successful reference outcome is relabeled as a rejection.

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
errors and resource limits have separate project guards. Missing parents leave
file input inactive; temporary-filesystem regressions also check random input
followed by console input, exhausted-console diagnostics, both path separators,
and closure of a replaced file without flushing an unfinished recording block.
REPL transcripts verify that these modes leave the next submission intact.
These interaction checks are project regressions, not additional reference
recordings. These tests do not establish full filesystem or full-language
conformance.

Nine further controls compare readable and empty files, explicit read-data
denial, exclusive read/write sharing locks, and logical input while random mode is enabled.
The restricted fixtures are existing regular files containing `41` followed by
CRLF. The setup independently verifies a failed read with native error 5 for
access denial or error 32 for the sharing lock. Both restrictions allow the
program to continue without reading, consume console values `7` and `9`, or
generate `7` under fixed random bounds. No case changes the existing file bytes.
The reference process is stopped and fixture access is restored before the
final byte read; the original and final hashes match. No new desktop captures
are published. These cases qualify only the recorded existing-file restrictions;
other directory, permission, ACL and concurrent-access arrangements remain
unqualified.

The `file-random-logical-source` control selects a file containing
`Verdadeiro` followed by CRLF while `aleatorio on` is active. It echoes
`Verdadeiro` and prints `VALUE= VERDADEIRO`; the contrasting optional console
response `Falso` is unused, and the file remains unchanged. This control and
the readable/empty controls replay normally and execute before and after
canonical formatting.

`TestRecordedFileInput` recreates real mode-based read denial where the host
enforces it and a Windows read/write handle with sharing disabled. The read/write
access matches the recorded fixture: a read-only handle does not reliably deny
Go's backup-semantics read open. The test independently verifies sharing error 32
before execution. It compares original and
formatted output and unchanged file bytes. Tests explicitly skip unsupported
restrictions or an account that bypasses mode permissions. Windows ACL denial
has native reference evidence and error-classification coverage, but does not
have an automated candidate ACL fixture. The corpus runner creates readable byte
fixtures, so the six `file-access-denied-*` and `file-access-locked-*` probes
remain pending rather than treating an ordinary readable file as a faithful
replay. No interpreter file-open injection or fixture-access schema is added.

Earlier file-input validation included local build, tests, race detection, vet,
staticcheck, lint, both fuzz targets, strict OpenSpec validation and CLI replay
of all 1,490 verified recordings. Native Windows build, vet, both fuzz targets
and 3,671 tests passed,
including the same 1,490 replays. At that stage, the evidence gate reported three
unmapped requirements and thirteen original examples; its requirements are
unchanged.
