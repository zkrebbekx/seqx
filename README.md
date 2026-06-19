# seqx

Lazy iterator combinators for Go 1.23 — the `Map`/`Filter`/`Reduce` the standard
library left out, built directly on `iter.Seq`.

Go 1.23 added range-over-func and the `iter` package, but the standard library
deliberately ships only producers and collectors (`slices.Values`,
`slices.Collect`, `maps.All`, …). `seqx` adds the transforms — lazily,
single-pass, and **on the standard `iter.Seq` type itself**, so it composes with
everything.

- **Lazy & single-pass.** Nothing runs until a terminal pulls. `Take`, `Find`,
  `TakeWhile`, `Any`, `All` stop early and that stop propagates upstream — so
  they work on **unbounded** sequences.
- **No wrapper type.** Every function takes and returns `iter.Seq[T]` /
  `iter.Seq2[K,V]`. Drop it into existing iterator code; mix freely with the
  standard library.
- **Zero dependencies.** Go 1.23+.

```go
import "github.com/zkrebbekx/seqx"
```

## Example

```go
src := slices.Values([]int{1, 2, 3, 4, 5, 6})
evens := seqx.Filter(src, func(n int) bool { return n%2 == 0 })
squares := seqx.Map(evens, func(n int) int { return n * n })

out := slices.Collect(squares) // [4 16 36]
```

Unbounded sources are fine — laziness means only what's needed is produced:

```go
naturals := func(yield func(int) bool) { for i := 0; ; i++ { if !yield(i) { return } } }

first5 := seqx.Collect(seqx.Take(naturals, 5)) // [0 1 2 3 4] — pulls exactly 5
```

## Why free functions, not method chaining

A fluent `.Map().Filter()` chain is impossible in Go when a step changes the
element type: a method cannot introduce a new type parameter (`Map[T,U]` needs
`U`). So `seqx` uses free functions, which read naturally with intermediate
variables and — crucially — keep the type as the plain `iter.Seq`, so there's no
lock-in and full interop with the standard library and other iterator code.

## API

**Producers** — `Of`, `Range`, `Repeat`, `FromChan` (and use `slices.Values`,
`maps.All` from the stdlib).

**Lazy transforms** — `Map`, `Filter`, `FlatMap`, `Take`, `Drop`, `TakeWhile`,
`DropWhile`, `Chunk`, `Distinct`, `Enumerate`.

**Terminals** — `ForEach`, `Reduce`, `Collect`, `Count`, `Find`, `First`,
`Any`, `All`, `Sum`, `Min`, `Max`, `GroupBy`, `ToMap`, `Partition`.

**Key/value (`iter.Seq2`)** — `Keys`, `Vals`, `Map2`, `Filter2`, `Pairs`.

## Re-iteration & concurrency

A pipeline is as re-iterable as its source: sequences from `Of`, `Range`,
`Repeat`, or `slices.Values` can be ranged many times; one from `FromChan` is
single-use. A single `iter.Seq` is not safe for concurrent iteration, but
distinct pipelines over re-iterable sources are independent.

## How it compares

- **samber/lo** operates eagerly over slices; `seqx` is lazy over `iter.Seq`,
  so it streams, stops early, and handles unbounded sources.
- Versus other `iter.Seq` helpers: `seqx` stays on the bare `iter.Seq` type
  (no wrapper), covers the full core set, and is careful about early-termination
  propagation.

## License

MIT
