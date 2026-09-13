# Safeguard acceptance

The assembled defensive profile is exercised by the ordinary test suite. The
tests use bounded child processes for cases that could otherwise recurse,
allocate, or loop indefinitely. A watchdog expiration is a test failure; it is
never treated as a language diagnostic.

| Area | Acceptance coverage |
| --- | --- |
| Source size, decoding, and positions | `internal/source/limits_test.go#TestSourceLimit`, `internal/source/source_test.go#TestDecodeCP1252`, `internal/parser/adversarial_test.go#TestEncodedSourceLimitSubprocess` |
| Syntax and traversal depth | `internal/parser/limits_test.go#TestStructuralLimits`, `internal/parser/limits_test.go#TestTypeDepthLimit`, `internal/parser/limits_test.go#TestConstantDepthLimit`, `internal/parser/limits_test.go#TestLimitKeepsEarlierDiagnostics` |
| Steps, calls, flat operands, and evaluation cleanup | `internal/interp/execution_test.go#TestStepBudgetAndReuse`, `internal/interp/execution_test.go#TestStepBudgetSharedAcrossCalls`, `internal/interp/execution_test.go#TestCallAndValueLimits`, `internal/interp/execution_test.go#TestRetainedOperandStorage`, `internal/interp/execution_test.go#TestDefensiveEvaluationDepth`, `scripts/conformance/execute_test.go#TestReplayBudgetExhaustion`, `scripts/conformance/replay_test.go#TestReplay` |
| Text, input, and output allocation bounds | `internal/interp/io_test.go#TestInputStringLimitPreservesEcho`, `internal/interp/io_test.go#TestWriteBufferLimit`, `internal/interp/execution_test.go#TestFormattedItemBoundary`, `internal/stdlib/stdlib_test.go#TestCaseConversionAllocationGuard`, `internal/stdlib/text_conversion_test.go#TestCopiaChecksResultSizeBeforeMaterializing` |
| Input failure/reuse and run cleanup | `internal/interp/io_test.go#TestInputFailurePreservesValue`, `internal/interp/execution_test.go#TestInputBufferOwnership`, `internal/interp/random_input_test.go#TestRandomInputRunReset`, `internal/interp/file_input_test.go#TestFileInputCleanupAndReset`, `internal/repl/input_modes_test.go#TestInputModesPreserveNextSubmission` |
| Storage and filesystem failures | `internal/interp/storage_test.go#TestCorruptedVectorStorage`, `internal/interp/bounds_test.go#TestConstantBoundAllocationGuard`, `internal/interp/file_input_test.go#TestFileInputFailures`, `internal/interp/file_input_test.go#TestMissingParentClearsPreviousFile`, `internal/runtime/vector_test.go#TestVectorCellRejectsInvalidStorage` |
| Host failures and state restoration | `internal/interp/display_test.go#TestDisplayHostFailures`, `internal/interp/environment_test.go#TestEchoHost`, `internal/interp/environment_test.go#TestChronometerResetAndFailures`, `internal/interp/execution_controls_test.go#TestTimerResetAndFailures`, `internal/interp/execution_controls_test.go#TestTimerCallFailureRestoresFrame` |
| REPL recovery and bounded subprocess behavior | `internal/repl/recovery_test.go#TestSourceLimitRecovery`, `internal/repl/recovery_test.go#TestHostDiagnosticRecovery`, `internal/repl/submission_test.go#TestAutomaticSubmissionDiagnostics`, `internal/testprocess/process.go` |

The replay runner also supplies every accepted probe with a finite execution
budget and a bounded subprocess timeout. `TestReplayBudgetExhaustion` verifies
that an empty infinite loop reports `R006` through the observation adapter
without producing output or relying on the watchdog. `TestReplay` verifies that
an accepted probe which exits with that diagnostic is rejected by the replay
contract. The remaining replay tests verify deadline and output bounds, so an
accepted example cannot pass by timing out or exhausting its output capture.

These are implementation safeguards and test-process contracts. They do not
change the reference-language evidence state of any probe or bundled example.
