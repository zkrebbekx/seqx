package seqx

import "iter"

// Keys returns a sequence of the keys of a key/value sequence. Lazy.
func Keys[K, V any](seq iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range seq {
			if !yield(k) {
				return
			}
		}
	}
}

// Vals returns a sequence of the values of a key/value sequence. Lazy.
func Vals[K, V any](seq iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range seq {
			if !yield(v) {
				return
			}
		}
	}
}

// Map2 transforms each key/value pair of a sequence. Lazy.
func Map2[K, V, K2, V2 any](seq iter.Seq2[K, V], f func(K, V) (K2, V2)) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k, v := range seq {
			if !yield(f(k, v)) {
				return
			}
		}
	}
}

// Filter2 keeps the key/value pairs for which keep is true. Lazy.
func Filter2[K, V any](seq iter.Seq2[K, V], keep func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			if keep(k, v) && !yield(k, v) {
				return
			}
		}
	}
}

// Pairs adapts a key/value sequence into a single-value sequence by combining
// each pair, so iter.Seq2 pipelines can rejoin the iter.Seq combinators. Lazy.
func Pairs[K, V, T any](seq iter.Seq2[K, V], combine func(K, V) T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for k, v := range seq {
			if !yield(combine(k, v)) {
				return
			}
		}
	}
}
