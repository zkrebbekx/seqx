// Command examples shows a realistic seqx pipeline: parse, filter, transform,
// and aggregate a stream of records lazily.
package main

import (
	"fmt"
	"slices"

	"github.com/zkrebbekx/seqx"
)

type order struct {
	Customer string
	Total    int // cents
}

func main() {
	orders := []order{
		{"alice", 1200}, {"bob", 50}, {"alice", 800},
		{"carol", 3000}, {"bob", 1500}, {"carol", 250},
	}

	// Lazy pipeline: keep orders over $1, total cents per customer.
	big := seqx.Filter(slices.Values(orders), func(o order) bool { return o.Total >= 100 })
	perCustomer := seqx.GroupBy(big, func(o order) string { return o.Customer })

	for _, name := range slices.Sorted(maps(perCustomer)) {
		total := seqx.Sum(seqx.Map(slices.Values(perCustomer[name]),
			func(o order) int { return o.Total }))
		fmt.Printf("%-6s $%.2f\n", name, float64(total)/100)
	}

	// Take the first 3 squares of an unbounded sequence.
	naturals := func(yield func(int) bool) {
		for i := 1; ; i++ {
			if !yield(i) {
				return
			}
		}
	}
	squares := seqx.Map(naturals, func(n int) int { return n * n })
	fmt.Println("first squares:", seqx.Collect(seqx.Take(squares, 3)))
}

// maps yields the keys of m as a sequence (tiny local helper to avoid importing
// the experimental maps.Keys signature differences across versions).
func maps[K comparable, V any](m map[K]V) func(func(K) bool) {
	return func(yield func(K) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}
