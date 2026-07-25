# GoForge collection algebra

`goforge.dev/gplodash` is a deliberately bounded compatibility facade for
[`github.com/samber/lo`](https://github.com/samber/lo), pinned to v1.53.0
(`cf6fb4f9b08c1d3d6e309581316f106dc30b458e`, MIT).

The module has two layers:

- 47 compatible declarations cover the high-use collection transformation,
  fold, grouping, shape, search/set, and map tier.
- New semantic operations make partiality and allocation explicit:
  `FirstOption`, `LastOption`, `FindOption`, `NthResult`, `ReduceNonEmpty`,
  `MapInto`, `FilterInto`, `FilterMapInto`, and `ReverseCopy`.
- The Go+ `Search[T]` enum and `Locate`/`LocateLast` require exhaustive
  handling of `Located(value, index)` versus `Absent`, eliminating the
  compatibility API's zero-value, boolean, and `-1` sentinel combination.

The production implementation is authored in `lo.gp` and
`find_set_map.gp`. `lo_gp.go` and `find_set_map_gp.go` are reproducible
generated artifacts, not implementation sources. The Go+ compiler is pinned
as a Go tool in `go.mod`.

```sh
go generate ./...
go tool goplus gen -check ./...
go test ./...
```

The complete 651-declaration upstream-module inventory is recorded in
[`API_MANIFEST.csv`](API_MANIFEST.csv); no unimplemented symbol is silently
claimed. See [`COMPATIBILITY.md`](COMPATIBILITY.md) for behavior and aliasing
contracts and [`PERFORMANCE.md`](PERFORMANCE.md) for the benchmark gate.

The proof-oriented primitives live in Go+'s standard library:

- `std/nonempty` seals owned, non-empty sequences and provides total head,
  last, map, append, and reduction.
- `std/vec` provides `Vec[T,n]`, shape-preserving `Map` and `Zip`, and
  bounds evidence `Fin[n]` for total `At`.

The project and both supporting packages are authored in Go+ and distributed
as generated ordinary Go.

## Reproduce the API inventory

```sh
upstream_root=$(go env GOMODCACHE)/github.com/samber/lo@v1.53.0
go run ./internal/cmd/apimanifest "$upstream_root" API_MANIFEST.csv
```
