package seqx

import "iter"

// Of returns a sequence over the given values. It is re-iterable.
func Of[T any](vals ...T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range vals {
			if !yield(v) {
				return
			}
		}
	}
}

// Range returns a sequence of integers from start (inclusive) to stop
// (exclusive). If stop <= start it is empty. The sequence is lazy, so a huge
// range costs nothing until consumed.
func Range[I integer](start, stop I) iter.Seq[I] {
	return func(yield func(I) bool) {
		for i := start; i < stop; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

// Repeat returns a sequence that yields v exactly n times. A negative n yields
// nothing.
func Repeat[T any](v T, n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; i < n; i++ {
			if !yield(v) {
				return
			}
		}
	}
}

// FromChan returns a sequence that yields values received from ch until ch is
// closed. The resulting sequence is single-use: each value is delivered once.
// If consumption stops early, the channel is left as-is for the caller to
// drain or close.
func FromChan[T any](ch <-chan T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range ch {
			if !yield(v) {
				return
			}
		}
	}
}
