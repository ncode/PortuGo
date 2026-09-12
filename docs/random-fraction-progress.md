# Random fraction recordings

Twenty-three recordings pin bare `rand`, assignment and statement suffixes,
reserved-name rejection and comparison grouping. Fourteen accepted cases have
exact original/formatted execution tests, and nine rejected cases have positioned
`P001` regressions. Every recording is verified against the implementation.

Assignments discard the rest of their physical line after the bare random
expression. Thus `x <- rand(1) + 7` stores the draw without evaluating the suffix.
Output expressions require their closing parenthesis and reject `rand()`.
Ungrouped compound comparisons are also rejected, as for other expressions.

An injected source checks exact values and draw counts. Boundary tests include
zero and the largest representable fraction below one, and reject negative,
non-finite and upper-bound draws. Repeated tests across 64 seeds assert the
fraction domain without comparing a reference sequence. Earlier output survives
a positioned `R007` for an invalid injected result.
The bounded legacy random API now shares `randi`'s source guards, covering
missing sources and invalid bounded draws through the same checked operation.

The shared builtin catalog, random-input commands and mode transitions remain
separate pending work. An incomplete nested-call capture is retained privately
and is not used as evidence. Published artifacts contain only synthetic source,
program-only output panels, reviewed diagnostic crops and labeled transcriptions.
