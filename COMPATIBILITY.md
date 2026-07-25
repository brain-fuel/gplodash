# samber/lo v1.53.0 compatibility

## Tier 1

| Area | Compatible declarations |
|---|---|
| Transform/fold | `Map`, `Filter`, `Reject`, `FilterMap`, `FlatMap`, `Reduce`, `ReduceRight`, `ForEach`, `Times` |
| Uniqueness/grouping | `Uniq`, `UniqBy`, `GroupBy`, `GroupByMap`, `KeyBy`, `Associate`, `PartitionBy` |
| Shape/order | `Chunk`, `Window`, `Sliding`, `Flatten`, `Concat`, `Interleave`, `Reverse`, `Drop`, `DropRight`, `Take` |
| Search | `Find`, `FindIndexOf`, `FindLastIndexOf`, `First`, `Last`, `Nth` |
| Sets | `Contains`, `Every`, `Some`, `None`, `Without`, `Union`, `Intersect`, `Difference` |
| Maps | `Entry`, `Keys`, `Values`, `Entries`, `FromEntries`, `MapKeys`, `MapValues` |

Results preserve upstream encounter order. Fresh-output operations do not alias
their input. `Chunk` and `Sliding` copy each child slice, matching upstream's
retention/aliasing contract. The deprecated compatibility function `Reverse`
mutates in place exactly as upstream does; `ReverseCopy` is the explicit
alias-safe semantic alternative. Map iteration order remains unspecified.

Invalid `Chunk`, `Window`, `Sliding`, `Drop`, `DropRight`, and `Take` arguments
panic where the pinned compatibility API panics. New code can avoid partial
collection operations with `Option`, `Result`, `std/nonempty`, and `std/vec`.

## Explicitly deferred

Every declaration has a row in `API_MANIFEST.csv`. The remaining root surface
is marked `deferred`; `it`, `mutable`, and `parallel` have separate statuses
because they require decisions about Go's `iter.Seq`, aliasing, cancellation,
and ordered concurrency rather than mechanical copying. Experimental SIMD is a
separate upstream module and is not part of the v1.53.0 root-module inventory.

This package does not claim drop-in coverage for deferred declarations.
