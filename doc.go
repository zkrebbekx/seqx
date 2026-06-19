// Package seqx is a small, lazy combinator library for Go 1.23 iterators.
//
// Go 1.23 added range-over-func and the [iter] package: an iter.Seq[T] is just
// a func(yield func(T) bool). The standard library gives you a few producers
// and collectors (slices.Values, slices.Collect, maps.All, ...) but
// deliberately omits Map/Filter/Reduce and friends. seqx fills that gap.
//
// # Design
//
// Every function operates on the standard [iter.Seq] and [iter.Seq2] types —
// seqx introduces no wrapper type. That means seqx composes directly with the
// standard library and any other iterator code:
//
//	import (
//		"iter"
//		"slices"
//
//		"github.com/zkrebbekx/seqx"
//	)
//
//	src := slices.Values([]int{1, 2, 3, 4, 5, 6})
//	evens := seqx.Filter(src, func(n int) bool { return n%2 == 0 })
//	squares := seqx.Map(evens, func(n int) int { return n * n })
//	out := slices.Collect(squares) // [4, 16, 36]
//
// Transforms are lazy and single-pass: nothing runs until a terminal (or a
// range loop) pulls values through, and only as many source elements as needed
// are produced. [Take], [Find], [First], [Any] and [All] stop early, and that
// stop propagates upstream — so they work on unbounded sequences:
//
//	first3 := seqx.Take(seqx.Range(0, 1_000_000_000), 3) // pulls only 3
//
// Because Go methods cannot introduce new type parameters, a type-changing
// step like Map cannot be a fluent method; seqx uses free functions and reads
// naturally with intermediate variables, as above.
//
// # Re-iteration and concurrency
//
// A seqx pipeline is as re-iterable as its source. Sequences built from
// [Of], [Range], [Repeat] or slices.Values can be ranged over many times;
// one built from a channel ([FromChan]) is single-use. A single iter.Seq is
// not safe for concurrent iteration, but distinct pipelines over re-iterable
// sources are independent.
package seqx
