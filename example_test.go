package seqx_test

import (
	"fmt"
	"iter"
	"slices"

	"github.com/zkrebbekx/seqx"
)

// Example shows a lazy Filter→Map pipeline collected with the standard library.
func Example() {
	src := slices.Values([]int{1, 2, 3, 4, 5, 6})
	evens := seqx.Filter(src, func(n int) bool { return n%2 == 0 })
	squares := seqx.Map(evens, func(n int) int { return n * n })

	fmt.Println(slices.Collect(squares))
	// Output: [4 16 36]
}

// Example_unbounded shows that Take stops an infinite sequence after n elements.
func Example_unbounded() {
	naturals := func(yield func(int) bool) {
		for i := 0; ; i++ {
			if !yield(i) {
				return
			}
		}
	}

	first5 := seqx.Take(iter.Seq[int](naturals), 5)
	fmt.Println(seqx.Collect(first5))
	// Output: [0 1 2 3 4]
}

// Example_groupBy buckets values by a derived key.
func Example_groupBy() {
	words := seqx.Of("ant", "bee", "bear", "cat", "crow")
	byInitial := seqx.GroupBy(words, func(s string) byte { return s[0] })

	fmt.Println(byInitial['b'])
	fmt.Println(byInitial['c'])
	// Output:
	// [bee bear]
	// [cat crow]
}
