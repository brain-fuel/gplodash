# Performance contract

The primary paired operation transforms 1,024 integers while retaining and
squaring the even values. `FilterMapInto` is compared directly with upstream
`FilterMap`; both execute the same callback and produce the same ordered
result. The GoForge operation makes buffer reuse explicit.

```sh
go test -run '^TestFusedIntoAllocationBudget$' \
  -bench 'Benchmark(CoreFusedInto|UpstreamFilterMap|UpstreamMapFilterPipeline)$' \
  -benchmem -count=5
```

Completion requires the slowest GoForge run to be at least twice as fast as
the fastest upstream `FilterMap` run, and at least 50% fewer allocations.
`TestFusedIntoAllocationBudget` enforces zero allocations for the reused-buffer
path and verifies the equivalent upstream pipeline allocates at least twice.

## Completion measurement

Measured 2026-07-23 from generated Go+ output on Apple M5 Max
(`darwin/arm64`):

| Operation | Five-run range | Bytes/op | Allocations/op |
|---|---:|---:|---:|
| GoForge `FilterMapInto` | 370.0–371.3 ns | 0 | 0 |
| upstream `FilterMap` | 985.4–994.7 ns | 8,192 | 1 |
| upstream `Map` + `Filter` | 1,804–1,832 ns | 16,384 | 2 |

Using the slowest GoForge run and fastest upstream run, the direct fused pairing
is **2.65x faster** with **100% fewer allocations**. Against the equivalent
two-helper pipeline it is 4.86x faster, also with 100% fewer allocations.

The exact-compatibility `Map` path is intentionally ordinary rather than
buffer-reusing: three runs measured 1,035–1,077 ns and one 8,192-byte
allocation, versus upstream's 1,010–1,049 ns and the same single allocation.
`TestCompatibilityAllocationParity` enforces that `Map`, `Filter`, and
`FilterMap` do not allocate more than their upstream counterparts. The speedup
claim belongs to the explicit semantic `Into` API, not to compatibility names.
