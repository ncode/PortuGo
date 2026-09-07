# Windows validation — 2026-09-07

This report covers the selected VisuAlg compatibility fixes in PR #2.
Infrastructure names, connection details, machine paths, account names, and
desktop captures are excluded from published evidence.

## Reference evidence

The reference is VisuAlg 3.0.7.0, acquired from the
[official distribution](https://sourceforge.net/projects/visualg30/files/visualg3.0.7.rar/download).

- Distribution SHA-256: `45b8c4787b9af2077aecacf19d4e81b05718d9c2f9f2e6565cef2e6f9ae6c33f`.
- Executable SHA-256: `7b297fb338e2ca7004ed55b42aef215c17a1a1d44b27cff50b6dc9be6dd3becf`.
- Recorded locale: `en-US`; other locales remain unverified.
- The distribution's 73 example programs have not yet been qualified.

The [original observations](validation/visualg-2026-09-07/observations.json)
and [follow-up observations](validation/visualg-2026-09-07/followup-observations.json)
contain 14 reduced sources, source hashes, capture times, reference output,
and bootstrap candidate results. Candidate diagnostic paths are reduced to
filenames; their published hashes describe that redacted text. Desktop captures
and full operational logs are retained privately.

Each source was entered into the actual reference editor, read back to verify
its contents, and executed. Reference panel text was captured through its
Unicode control and encoded as UTF-8. It preserves application notices and CRLF;
it is not a capture of original console bytes.

## Recorded corrections

| Area | Reference behavior now covered |
| --- | --- |
| CLI failure status | Lexical, parse, and semantic rejection exits 1 without executing the program |
| Numeric built-in | `exp(base, exponent)` accepts two numeric arguments and returns a real; one argument is rejected |
| Power associativity | `2^3^2` produces 64 |
| Unary-minus precedence | `-2^2` produces 4 |
| Default output | Numeric/logical leading spaces, decimal dots, uppercase logical values |
| Field formatting | Width zero ignores decimals; numeric ties round away from zero; strings left-align and truncate; logical widths are rejected |
| Loop exit state | Normal, empty, and interrupted loops follow the recorded rules, including mutation and descending steps |

The original CLI defect returned success after `checkedProgram` reported
errors. Regression coverage now exercises lexical, parse, semantic, and valid
execution paths as actual subprocesses.

Loop assignments remain visible within the body without changing progression.
Normal nonempty completion exposes the smaller of the next iteration value and
terminal bound; interruption exposes the smaller of the body value and terminal
bound; an empty loop exposes its initial bound. The reference memory grid can
disagree with subsequent program output; the tests follow program output.

Eager logical evaluation and real-valued `/` were already correct. Invalid
programs are statically rejected before execution, whereas the reference GUI
can print preceding statements before reporting an error.

## Validation

Regression tests were observed failing before the corresponding fixes.

| Check | Result |
| --- | --- |
| macOS/arm64, Go 1.25.0: build, format, vet, staticcheck v0.6.1, golangci-lint v2.4.0 | PASS |
| Uncached ordinary and race tests; module hygiene | PASS |
| Local 30-second lexer and parser fuzz campaigns | PASS; 2,114,724 and 1,020,296 executions |
| Native Windows/amd64, Go 1.26.4: build, uncached tests, vet, CLI build | PASS |
| Native Windows recorded reference replay | PASS, 14/14 |
| Native Windows CLI checks | PASS, 57/57 |
| Native Windows 30-second lexer and parser fuzz campaigns | PASS; 711,411 and 676,715 executions |
| Windows race | Unavailable with CGO disabled; local race testing passed |
| Strict OpenSpec validation | PASS |

The CLI checks cover runtime fixtures, examples, source encoding, paths with
spaces, invalid-source exit statuses, formatter round trips, a REPL transcript,
and all 14 recorded reference programs.

`TestRecordedWindowsProbes` replays 11 accepted programs by removing only fixed
application notices and converting CRLF to the portable CLI's LF. Three rejected
programs assert semantic diagnostic codes and source lines. Runtime fixture
comparisons otherwise preserve exact output bytes.

The [PR checks](https://github.com/ncode/PortuGo/pull/2/checks) cover the current
head on Linux, macOS, and Windows with Go 1.22/1.25, plus independent quality and
fuzz jobs.

## OpenSpec scope

The quality-baseline group and the eight scoped follow-up tasks are complete:
19/216 tasks. All 197 remaining original conformance tasks stay unchecked.
The complete corpus, other locales, integer-boundary/resource guards, and the
73 distribution examples remain unfinished. These results do not establish
full VisuAlg conformance.
