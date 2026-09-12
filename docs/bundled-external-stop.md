# Bundled examples requiring external interruption

The finite automated example sweep excludes only the complete runs of the two
original programs below. This is an explicit scope exception for their lack of
normal completion, not a syntax rejection or a passing execution result. Their
unchanged Windows-1252 sources, original catalog hashes, entered input and
partial output panels remain available for inspection. Language features used
by these programs remain required and retain their independent probes and tests.

The catalog uses `non-goal` for this narrow finite-completion exception, with
reviewed `not-applicable` evidence and implementation states. Retained artifact
hashes are still validated. Reproducing an externally stopped interactive session
is outside the finite sweep; a timeout never counts as successful execution.

## Clock

`example.f322405e9d84` increments its clock indefinitely: its loop tests
`Hora < 25`, but resets the hour to zero when it exceeds 24. The reference run
on 2026-09-12 displayed six successive updates, from `00:00:01` to `00:00:06`,
without a completion notice before it was externally stopped. No input was
entered. Timer, loop, string-formatting and display behavior remain required.

- [Original source](../testdata/conformance/visualg-3.0.7/probes/bundled-f322405e9d84/source.alg)
- [Partial output panel](../testdata/conformance/visualg-3.0.7/probes/bundled-f322405e9d84/panel.txt)

## Conversion menu

`example.80addd4fac6d` reads a local `opcao` in `menu` but does not assign the
function result to the global variable used by its outer exit condition.
The reference run on 2026-09-12 accepted `F`, requested Enter and displayed the
menu again. It was waiting for another choice when externally stopped, with
no completion notice. Input consumption, lexical scope, function results and
choice behavior remain required.

- [Original source](../testdata/conformance/visualg-3.0.7/probes/bundled-80addd4fac6d/source.alg)
- [Entered input](../testdata/conformance/visualg-3.0.7/probes/bundled-80addd4fac6d/input.txt)
- [Partial output panel](../testdata/conformance/visualg-3.0.7/probes/bundled-80addd4fac6d/panel.txt)

Panel files preserve the captured UTF-8 program-output bytes, CRLF line endings
and start notice. They contain no execution-completion notice, desktop
capture, process details or diagnostic logs. The source and panel hashes in the
manifest describe those published bytes.
