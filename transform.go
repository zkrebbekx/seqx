package seqx

import "iter"

// Map returns a sequence that applies f to each element of seq. Lazy.
func Map[T, U any](seq iter.Seq[T], f func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range seq {
			if !yield(f(v)) {
				return
			}
		}
	}
}

// Filter returns a sequence of the elements of seq for which keep is true. Lazy.
func Filter[T any](seq iter.Seq[T], keep func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			if keep(v) && !yield(v) {
				return
			}
		}
	}
}

// FlatMap maps each element to a sequence and concatenates the results. Lazy.
func FlatMap[T, U any](seq iter.Seq[T], f func(T) iter.Seq[U]) iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range seq {
			for u := range f(v) {
				if !yield(u) {
					return
				}
			}
		}
	}
}

// Take returns a sequence of at most the first n elements of seq. It stops
// pulling from seq once n have been yielded, so it is safe on unbounded
// sequences. A non-positive n yields nothing.
func Take[T any](seq iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		if n <= 0 {
			return
		}
		count := 0
		for v := range seq {
			if !yield(v) {
				return
			}
			count++
			if count >= n {
				return
			}
		}
	}
}

// Drop returns a sequence that skips the first n elements of seq. Lazy.
func Drop[T any](seq iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		skipped := 0
		for v := range seq {
			if skipped < n {
				skipped++
				continue
			}
			if !yield(v) {
				return
			}
		}
	}
}

// TakeWhile yields elements while pred holds, stopping at the first element for
// which it does not (and pulling no further). Lazy.
func TakeWhile[T any](seq iter.Seq[T], pred func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			if !pred(v) {
				return
			}
			if !yield(v) {
				return
			}
		}
	}
}

// DropWhile skips leading elements while pred holds, then yields the rest. Lazy.
func DropWhile[T any](seq iter.Seq[T], pred func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		dropping := true
		for v := range seq {
			if dropping {
				if pred(v) {
					continue
				}
				dropping = false
			}
			if !yield(v) {
				return
			}
		}
	}
}

// Chunk groups consecutive elements into slices of at most size. The final
// chunk may be shorter. It panics if size <= 0. Each emitted slice is freshly
// allocated and owned by the consumer.
func Chunk[T any](seq iter.Seq[T], size int) iter.Seq[[]T] {
	if size <= 0 {
		panic("seqx: Chunk size must be positive")
	}
	return func(yield func([]T) bool) {
		buf := make([]T, 0, size)
		for v := range seq {
			buf = append(buf, v)
			if len(buf) == size {
				if !yield(buf) {
					return
				}
				buf = make([]T, 0, size)
			}
		}
		if len(buf) > 0 {
			yield(buf)
		}
	}
}

// Distinct yields the elements of seq with duplicates removed, keeping first
// occurrence order. It holds the set of seen values, so memory grows with the
// number of distinct elements. Lazy in output.
func Distinct[T comparable](seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		seen := make(map[T]struct{})
		for v := range seq {
			if _, ok := seen[v]; ok {
				continue
			}
			seen[v] = struct{}{}
			if !yield(v) {
				return
			}
		}
	}
}

// Enumerate pairs each element with its zero-based index, as an iter.Seq2.
func Enumerate[T any](seq iter.Seq[T]) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		i := 0
		for v := range seq {
			if !yield(i, v) {
				return
			}
			i++
		}
	}
}
