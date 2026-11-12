# Performance baseline

The first reproducible baseline identifies allocation work in semantic symbol
maps/catalog copies, vector type queries, builtin argument buffers, and recursive
call frames. The initial baseline includes measurements and correctness-checked
workloads and makes no production optimization or speedup claim. The
matched follow-up below records the changes retained after measurement.

## Scope and method

Measured on November 9, 2026 with Go 1.27.3, macOS arm64, `GOMAXPROCS=1`, and
`-cpu=1`. The implementation baseline is `041d71a`; the exact benchmark harness
is commit `c2631b6`. The hardware model is omitted from published samples for
privacy. `GOGC`, `GOMEMLIMIT`, and `GODEBUG` were unset. Compare future
versions on the same machine and toolchain.

`BenchmarkPortuGo` has six workloads and six stages: scan, parse, analyze, run,
format, and pipeline. The pipeline starts with decoded source already in memory
and ends with interpreter output to `io.Discard`. Run includes a new interpreter
and reuses checked syntax/information; format is the canonical AST printer.
Working directory is preset. File reads/decoding, process startup, real terminal
I/O, and host delays are outside these measurements.

Workloads cover a small scalar program, 1,000 integer loop iterations,
Fibonacci(12) with 465 calls, initialization/summing of 256 vector elements,
200 string rotations/conversions, and 500 scalar declarations/assignments.
Each fixture has exact expected output. `TestBenchmarkPrograms` checks all six
with an execution budget, and benchmark setup validates output before timing.
Fixtures are deterministic, but do not represent every language feature or a
production workload distribution.

## Repeated measurements

The complete campaign uses six 1-second measurements per stage/workload, 216
samples in total. Values below are benchstat medians; timing ranges are its
reported 95% confidence intervals. Stages run separately, so their medians are
not additive and should not be subtracted from pipeline medians.

| Pipeline workload | µs/op and confidence range | KiB/op | allocations/op |
| --- | ---: | ---: | ---: |
| small | 8.62 ± 1% | 38.45 | 73 |
| arithmetic | 543.28 ± 9% | 39.20 | 82 |
| calls | 752.38 ± 39% | 291.53 | 3,393 |
| vectors | 614.08 ± 5% | 252.21 | 3,198 |
| strings | 827.20 ± 14% | 164.23 | 1,134 |
| declarations | 1,936.11 ± 14% | 3,027.17 | 2,695 |

Several timings were noisy. A separate campaign repeats the main allocation
candidates six times with 2-second measurements. The implementation is identical:
this is confirmation, not a before/after optimization comparison.

| Confirmation stage | µs/op and confidence range | KiB/op | allocations/op |
| --- | ---: | ---: | ---: |
| calls/run | 727.52 ± 8% | 247.72 | 3,278 |
| vectors/run | 486.09 ± 4% | 201.74 | 3,096 |
| strings/run | 565.15 ± 2% | 108.51 | 1,018 |
| declarations/analyze | 471.58 ± 7% | 1,315.38 | 91 |

Allocation ranges in the confirmation table round to 0%. Timing uncertainty remains material
for calls and larger-source analysis; use longer controlled measurements before
accepting a small latency improvement. The original campaign and confirmation
samples are both retained, including noisy measurements.

## Profile findings and next changes

CPU profiles used a separate 5-second benchmark campaign. Allocation profiles
used separate 2-second campaigns with `-memprofilerate=1` and the `alloc_space`
view. Profiled timings are excluded from the timing tables above. Allocation
percentages describe cumulative allocated bytes, not resident heap or RSS.

1. **Semantic maps and builtin catalog copies.** Small-program analysis allocates
   25,920 bytes/op. `scope.declare` accounts for 71.0% of allocated bytes and the
   catalog slice copy for 18.8% in that profile. With 500 declarations,
   `recordBinding` is 46.9%, `writable` 29.1%, and scope declarations 23.5% of flat
   allocation bytes. First investigate capacity hints for known symbol/binding
   counts and a read-only builtin-catalog traversal. Keep collision checks,
   per-analysis isolation, and the catalog's copy contract intact. See
   [scope construction](../internal/sema/sema.go),
   [binding records](../internal/sema/info.go), and
   [catalog access](../internal/stdlib/registry.go).
2. **Vector type metadata copies.** `runtime.Type.Clone` and its range-slice
   copies account for 77.1% of allocation bytes in vector execution. The workload
   allocates about 202 KiB and 3,096 objects per run. Investigate avoiding repeated
   deep copies in interpreter type/binding queries, while preserving independent
   copies returned by the semantic information API. See
   [type cloning](../internal/runtime/type.go) and
   [semantic information queries](../internal/sema/info.go).
3. **Builtin argument buffers.** In the string workload, 92.0% of flat allocated
   bytes are at `args := make([]runtime.Value, len(exprs))` in `callFunction`.
   Investigate bounded-arity storage or a narrower call path, with escape-analysis
   evidence. Preserve left-to-right effects, nested calls, no-value propagation,
   and halt behavior. See [builtin calls](../internal/interp/call.go).
4. **Recursive call frames.** Parameter-cell/environment creation and direct
   `callSub` allocations account for about 90.5% of flat allocated bytes in the
   recursive workload. Investigate reducing per-call environment/cell/result
   churn, preserving recursion isolation, retained results, and reference
   parameter copy-back. See [call setup](../internal/interp/call.go) and
   [environment cells](../internal/interp/env.go).

Start with semantic map/catalog setup: it is a narrow target, affects short
programs, and has stable allocation evidence. Use the vector and builtin-buffer
benchmarks for subsequent changes. These are candidate optimizations, not
promised percentage improvements.

CPU captures were dominated by `runtime.kevent` (61.3–98.4% of samples), often
with GC/runtime frames. They do not establish a reliable ranking of application
CPU costs. Keep CPU-specific conclusions pending a clearer capture on another
supported host or independent sampling method. No CPU-profile percentages above
are used to predict a speedup.

## Reproduce and validate

```sh
bench_dir=$(mktemp -d)
GOBIN="$bench_dir" go install golang.org/x/perf/cmd/benchstat@v0.0.0-20261009192801-be2c69fb417e
go test . -run '^TestBenchmarkPrograms$' -count=1
GOMAXPROCS=1 go test . -run '^$' -bench '^BenchmarkPortuGo$' -benchmem -benchtime=1s -count=6 -cpu=1 > "$bench_dir/before.txt"
"$bench_dir/benchstat" "$bench_dir/before.txt"
```

After one targeted change, repeat the same command and environment:

```sh
GOMAXPROCS=1 go test . -run '^$' -bench '^BenchmarkPortuGo$' -benchmem -benchtime=1s -count=6 -cpu=1 > "$bench_dir/after.txt"
"$bench_dir/benchstat" "$bench_dir/before.txt" "$bench_dir/after.txt"
```

The checked-in samples support reproducing this baseline; new claims require a
matched A/B run. Keep raw benchmark output private until metadata is screened.
Benchstat used here is `golang.org/x/perf` version
`v0.0.0-20261009192801-be2c69fb417e`, installed outside the project module.

Keep profile binaries and raw profiles outside the checkout:

```sh
profile_dir=$(mktemp -d)
GOMAXPROCS=1 go test . -o "$profile_dir/portugo.test" -run '^$' \
  -bench '^BenchmarkPortuGo$/^vectors$/^run$' -cpu=1 -benchtime=5s \
  -cpuprofile="$profile_dir/cpu.pprof"
go tool pprof -top "$profile_dir/portugo.test" "$profile_dir/cpu.pprof"
GOMAXPROCS=1 go test . -o "$profile_dir/portugo.test" -run '^$' \
  -bench '^BenchmarkPortuGo$/^vectors$/^run$' -cpu=1 -benchtime=2s \
  -memprofile="$profile_dir/alloc.pprof" -memprofilerate=1
go tool pprof -alloc_space -top "$profile_dir/portugo.test" "$profile_dir/alloc.pprof"
```

After a measured improvement, run `go test ./...`, `go test -race -count=1 ./...`,
`go vet ./...`, pinned linters, and the complete conformance replay. Benchmark
results are not noisy CI pass/fail thresholds; CI checks workload correctness.

## Matched optimization follow-up

On November 12, 2026, a fresh baseline at `cb78635` and the review fixes were
measured on the same host with Go 1.27.4, macOS arm64, `GOMAXPROCS=1`, and
`-cpu=1`. Each version ran all six workloads and six stages with `-benchmem`,
`-benchtime=1s`, and `-count=6`: 216 samples per version, 432 total. The benchmark
harness and checked outputs were unchanged. The benchstat version above was
used for comparison. Raw output remains private; the table contains screened
medians without hardware-model or local-path metadata.

| Stage | Before B/op | After B/op | Before allocations/op | After allocations/op |
| --- | ---: | ---: | ---: | ---: |
| small/analyze | 25,920 | 12,160 | 23 | 17 |
| declarations/analyze | 1,346,952 | 669,856 | 91 | 29 |
| vectors/run | 206,585 | 46,898 | 3,096 | 24 |
| strings/run | 111,116 | 111,180 | 1,018 | 1,018 |
| calls/run | 253,666 | 253,730 | 3,278 | 3,278 |

Retained changes avoid full layout copies for immutable semantic facts, size
semantic maps from cheap known counts, and traverse builtin names without
copying the catalog. Small and declaration analysis use 53.09% and 50.27% fewer
bytes respectively; vector execution uses 77.30% fewer bytes and 99.22% fewer
allocations. Full-layout semantic queries still return independent copies.
The new live-storage accounting adds 64 B/op to non-vector run workloads without
adding allocations. Recursive frame allocation remains a separate candidate.
A separate two-second vector-run allocation capture with `-memprofilerate=1`
attributes 87.49% of allocated bytes to vector creation; raw profiles and their
local build metadata remain private.

The invocation-local builtin argument-array candidate was rejected. Escape
analysis showed the three-value array moving to the heap through the existing
call interface. Six matched string-run samples used 162,382 B/op and 1,018
allocations/op, compared with 111,180 B/op and the same allocation count for the
retained slice implementation. No reusable buffer or pooling was introduced.

The targeted analysis and execution stages showed no timing regression. Vector
run medians were 450.2 and 324.0 µs/op. Small-program formatting rose 3.74%
(720.0 to 746.9 ns/op); unchanged small scan/parse stages rose 2.23% and 2.81%,
and declaration parsing rose 3.00%. Other unchanged stages shifted down slightly.
Those campaign shifts limit causal timing claims; the allocation reductions
above are the acceptance evidence. These synthetic workloads and this single
host do not establish a general CPU or production-latency improvement.

## Published evidence

Hardware-model/binary/capture metadata redactions are labeled in each artifact;
benchmark numeric samples and profile accounting values are retained. Raw
profile binaries remain private because they contain local build paths.

- [216 baseline samples](benchmarks/2026-11-09-baseline.txt)
- [24 longer confirmation samples](benchmarks/2026-11-09-confirmation.txt)
- [sanitized CPU/allocation summaries](benchmarks/2026-11-09-profiles.txt)

Published SHA-256 values:

```text
2026-11-09-baseline.txt c9138144634bc784cf425842d89f6dbcf779acedfb4f061fb29946bba4f11e66
2026-11-09-confirmation.txt f8c7e92832907b3d9016bbe73a1af71176bc57478a065ab4c4d4a341ac720166
2026-11-09-profiles.txt 4531b7f5a319c535c12d129a14eb58caf87c024a6e0780e705bed78508848056
```
