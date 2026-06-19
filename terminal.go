package seqx

import "iter"

// ForEach calls fn for every element of seq, draining it.
func ForEach[T any](seq iter.Seq[T], fn func(T)) {
	for v := range seq {
		fn(v)
	}
}

// Reduce folds seq into a single value, starting from init.
func Reduce[T, A any](seq iter.Seq[T], init A, fn func(acc A, v T) A) A {
	acc := init
	for v := range seq {
		acc = fn(acc, v)
	}
	return acc
}

// Collect drains seq into a slice. It is slices.Collect, provided for symmetry.
func Collect[T any](seq iter.Seq[T]) []T {
	var out []T
	for v := range seq {
		out = append(out, v)
	}
	return out
}

// Count returns the number of elements in seq, draining it.
func Count[T any](seq iter.Seq[T]) int {
	n := 0
	for range seq {
		n++
	}
	return n
}

// Find returns the first element satisfying pred, and true; or the zero value
// and false. It stops pulling once found.
func Find[T any](seq iter.Seq[T], pred func(T) bool) (T, bool) {
	for v := range seq {
		if pred(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// First returns the first element of seq and true, or the zero value and
// false if seq is empty. It pulls a single element.
func First[T any](seq iter.Seq[T]) (T, bool) {
	for v := range seq {
		return v, true
	}
	var zero T
	return zero, false
}

// Any reports whether any element satisfies pred, stopping at the first match.
func Any[T any](seq iter.Seq[T], pred func(T) bool) bool {
	for v := range seq {
		if pred(v) {
			return true
		}
	}
	return false
}

// All reports whether every element satisfies pred, stopping at the first
// counter-example. It is vacuously true for an empty sequence.
func All[T any](seq iter.Seq[T], pred func(T) bool) bool {
	for v := range seq {
		if !pred(v) {
			return false
		}
	}
	return true
}

// Sum totals the elements of seq. An empty sequence sums to zero.
func Sum[T Number](seq iter.Seq[T]) T {
	var total T
	for v := range seq {
		total += v
	}
	return total
}

// Min returns the smallest element and true, or the zero value and false if
// seq is empty.
func Min[T Ordered](seq iter.Seq[T]) (T, bool) {
	var best T
	found := false
	for v := range seq {
		if !found || v < best {
			best, found = v, true
		}
	}
	return best, found
}

// Max returns the largest element and true, or the zero value and false if
// seq is empty.
func Max[T Ordered](seq iter.Seq[T]) (T, bool) {
	var best T
	found := false
	for v := range seq {
		if !found || v > best {
			best, found = v, true
		}
	}
	return best, found
}

// GroupBy drains seq into a map from key(v) to the elements with that key,
// preserving encounter order within each group.
func GroupBy[T any, K comparable](seq iter.Seq[T], key func(T) K) map[K][]T {
	out := make(map[K][]T)
	for v := range seq {
		k := key(v)
		out[k] = append(out[k], v)
	}
	return out
}

// ToMap drains seq into a map, deriving a key and value from each element. On
// duplicate keys, the later element wins.
func ToMap[T any, K comparable, V any](seq iter.Seq[T], kv func(T) (K, V)) map[K]V {
	out := make(map[K]V)
	for v := range seq {
		k, val := kv(v)
		out[k] = val
	}
	return out
}

// Partition splits seq into the elements for which pred is true and those for
// which it is false, draining seq.
func Partition[T any](seq iter.Seq[T], pred func(T) bool) (yes, no []T) {
	for v := range seq {
		if pred(v) {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}
	return yes, no
}
