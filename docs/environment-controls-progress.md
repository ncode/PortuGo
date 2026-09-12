# Echo and chronometer recordings

This is a historical snapshot of the original slice. Its counts and completion
statements describe that point in development; current behavior is documented
in the [language reference](language.md).

This slice adds 24 synthetic recordings: ten echo controls, ten chronometer
controls and four elapsed-time samples. Eighteen recordings have exact
implementation matches: twelve accepted outputs, including two earlier echo
controls, and six positioned rejections. The complete corpus contains 1,421
verified reference recordings and 86 recorded cases pending verification.

The command sources, program-only output panels and reviewed diagnostic-window
crops are linked by SHA-256 in the conformance manifest. Diagnostic text is
labeled as a manual transcription. Original and canonical source execution
must produce identical output for each verified accepted recording.

Echo settings use typed host calls. The recorded console and random-input
transcripts retain converted input echo around `eco off`; the headless profile
preserves it. Bare, numeric, quoted and unknown echo tails are accepted, but their
GUI state and file-input effects remain unqualified.

The chronometer records start, repeated start, stop, inactive stop, ignored
tails, rejected modes and expression-form absence. Fake-clock tests reproduce
the recorded zero, 16-millisecond, 2.047-second, 2.422-second and 72.141-second
messages, including restart, per-run reset, output failures and clock regressions.
Whole-second formatting without a remainder is the implementation's boundary
rule; no exact whole-second reference sample is claimed.

Eight accepted elapsed-time recordings remain pending exact replay because
their wall-clock values vary. Four of them use numeric timer commands to create
bounded intervals; timer implementation is a following slice. No elapsed values
are removed or normalized to make the conformance replay pass. Pause, debug,
timer, file input and the complete conformance gate remain unfinished.

Later slices implemented the recorded timer, pause/debug and file-input
behavior. `file-detail-echo-off` now verifies retained file-input echo and
unchanged input bytes through original and formatted execution. The GUI state
of ignored echo tails remains unqualified. Three zero-elapsed-time recordings
now have [deterministic replay](deterministic-qualification.md); variable
elapsed-time recordings and rejection-order differences remain pending.
The existing [host-failure coverage](display-commands-progress.md#injected-host-failures)
now traces the complete typed host error boundary separately from those gaps.
