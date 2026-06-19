package seqx_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/zkrebbekx/seqx"

	. "github.com/smartystreets/goconvey/convey"
)

// naturals is an unbounded sequence 0,1,2,... used to prove laziness and
// early termination: any test that consumes it without a bound would hang.
func naturals() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; ; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

func TestTransforms(t *testing.T) {
	Convey("Given a finite source 1..6", t, func() {
		src := func() iter.Seq[int] { return seqx.Range(1, 7) }

		Convey("When Filter then Map are composed", func() {
			out := seqx.Collect(seqx.Map(
				seqx.Filter(src(), func(n int) bool { return n%2 == 0 }),
				func(n int) int { return n * n },
			))
			Convey("Then only the matching, transformed elements remain", func() {
				So(out, ShouldResemble, []int{4, 16, 36})
			})
		})

		Convey("When FlatMap expands each element", func() {
			out := seqx.Collect(seqx.FlatMap(seqx.Of(1, 2), func(n int) iter.Seq[int] {
				return seqx.Of(n, n*10)
			}))
			Convey("Then the sub-sequences are concatenated", func() {
				So(out, ShouldResemble, []int{1, 10, 2, 20})
			})
		})

		Convey("When Distinct removes duplicates", func() {
			out := seqx.Collect(seqx.Distinct(seqx.Of(1, 1, 2, 3, 3, 3, 2)))
			Convey("Then first-occurrence order is kept", func() {
				So(out, ShouldResemble, []int{1, 2, 3})
			})
		})

		Convey("When Chunk groups by 2", func() {
			out := seqx.Collect(seqx.Chunk(src(), 2))
			Convey("Then it yields equal groups with a possibly-short tail", func() {
				So(out, ShouldResemble, [][]int{{1, 2}, {3, 4}, {5, 6}})
			})
		})

		Convey("When Chunk size is not positive", func() {
			Convey("Then it panics", func() {
				So(func() { seqx.Chunk(src(), 0) }, ShouldPanic)
			})
		})
	})
}

func TestLazinessAndEarlyTermination(t *testing.T) {
	Convey("Given an unbounded source", t, func() {
		Convey("When Take(3) is collected", func() {
			out := seqx.Collect(seqx.Take(naturals(), 3))
			Convey("Then only the first three are produced (no hang)", func() {
				So(out, ShouldResemble, []int{0, 1, 2})
			})
		})

		Convey("When Find stops at the first match", func() {
			v, ok := seqx.Find(naturals(), func(n int) bool { return n == 5 })
			Convey("Then it returns without draining the infinite source", func() {
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 5)
			})
		})

		Convey("When TakeWhile bounds an unbounded source", func() {
			out := seqx.Collect(seqx.TakeWhile(naturals(), func(n int) bool { return n < 4 }))
			Convey("Then it stops at the predicate boundary", func() {
				So(out, ShouldResemble, []int{0, 1, 2, 3})
			})
		})

		Convey("When a pipeline counts how many source elements it pulls", func() {
			pulled := 0
			counted := func(yield func(int) bool) {
				for i := 0; ; i++ {
					pulled++
					if !yield(i) {
						return
					}
				}
			}
			_ = seqx.Collect(seqx.Take(seqx.Map(counted, func(n int) int { return n }), 3))
			Convey("Then it pulls exactly as many as needed (single-pass, lazy)", func() {
				So(pulled, ShouldEqual, 3)
			})
		})
	})
}

func TestDropAndWhile(t *testing.T) {
	Convey("Given 0..9", t, func() {
		src := func() iter.Seq[int] { return seqx.Range(0, 10) }

		Convey("When Drop(3)", func() {
			So(seqx.Collect(seqx.Drop(src(), 3)), ShouldResemble, []int{3, 4, 5, 6, 7, 8, 9})
		})
		Convey("When DropWhile < 5", func() {
			out := seqx.Collect(seqx.DropWhile(src(), func(n int) bool { return n < 5 }))
			So(out, ShouldResemble, []int{5, 6, 7, 8, 9})
		})
		Convey("When Take exceeds the source length", func() {
			So(seqx.Collect(seqx.Take(seqx.Of(1, 2), 9)), ShouldResemble, []int{1, 2})
		})
	})
}

func TestTerminals(t *testing.T) {
	Convey("Given 1..5", t, func() {
		src := func() iter.Seq[int] { return seqx.Range(1, 6) }

		Convey("Reduce / Sum / Count", func() {
			So(seqx.Reduce(src(), 0, func(a, v int) int { return a + v }), ShouldEqual, 15)
			So(seqx.Sum(src()), ShouldEqual, 15)
			So(seqx.Count(src()), ShouldEqual, 5)
		})

		Convey("Min / Max over an empty and non-empty sequence", func() {
			mn, ok1 := seqx.Min(src())
			mx, ok2 := seqx.Max(src())
			So(ok1, ShouldBeTrue)
			So(mn, ShouldEqual, 1)
			So(ok2, ShouldBeTrue)
			So(mx, ShouldEqual, 5)

			_, okEmpty := seqx.Max(seqx.Of[int]())
			So(okEmpty, ShouldBeFalse)
		})

		Convey("Any / All", func() {
			So(seqx.Any(src(), func(n int) bool { return n == 3 }), ShouldBeTrue)
			So(seqx.All(src(), func(n int) bool { return n > 0 }), ShouldBeTrue)
			So(seqx.All(src(), func(n int) bool { return n > 2 }), ShouldBeFalse)
			So(seqx.All(seqx.Of[int](), func(int) bool { return false }), ShouldBeTrue) // vacuous
		})

		Convey("First on empty", func() {
			_, ok := seqx.First(seqx.Of[int]())
			So(ok, ShouldBeFalse)
		})

		Convey("GroupBy / ToMap / Partition", func() {
			groups := seqx.GroupBy(src(), func(n int) string {
				if n%2 == 0 {
					return "even"
				}
				return "odd"
			})
			So(groups["even"], ShouldResemble, []int{2, 4})
			So(groups["odd"], ShouldResemble, []int{1, 3, 5})

			m := seqx.ToMap(src(), func(n int) (int, int) { return n, n * n })
			So(m[4], ShouldEqual, 16)

			yes, no := seqx.Partition(src(), func(n int) bool { return n > 3 })
			So(yes, ShouldResemble, []int{4, 5})
			So(no, ShouldResemble, []int{1, 2, 3})
		})
	})
}

func TestSeq2(t *testing.T) {
	Convey("Given Enumerate over three letters", t, func() {
		pairs := seqx.Enumerate(seqx.Of("a", "b", "c"))

		Convey("When taking keys and values", func() {
			So(seqx.Collect(seqx.Keys(pairs)), ShouldResemble, []int{0, 1, 2})
			So(seqx.Collect(seqx.Vals(pairs)), ShouldResemble, []string{"a", "b", "c"})
		})

		Convey("When Filter2 then Pairs recombine", func() {
			out := seqx.Collect(seqx.Pairs(
				seqx.Filter2(pairs, func(i int, _ string) bool { return i > 0 }),
				func(_ int, s string) string { return s },
			))
			So(out, ShouldResemble, []string{"b", "c"})
		})
	})
}

func TestInteropWithStdlib(t *testing.T) {
	Convey("Given a slice via slices.Values", t, func() {
		src := slices.Values([]int{5, 3, 1, 4, 2})

		Convey("When run through seqx and collected with slices.Collect", func() {
			out := slices.Collect(seqx.Filter(src, func(n int) bool { return n > 2 }))
			slices.Sort(out)
			Convey("Then seqx composes with the standard library", func() {
				So(out, ShouldResemble, []int{3, 4, 5})
			})
		})
	})
}

func TestReIteration(t *testing.T) {
	Convey("Given a pipeline over a re-iterable source", t, func() {
		evens := seqx.Filter(seqx.Range(0, 10), func(n int) bool { return n%2 == 0 })

		Convey("When iterated twice", func() {
			first := seqx.Collect(evens)
			second := seqx.Collect(evens)
			Convey("Then both passes yield the same result", func() {
				So(first, ShouldResemble, []int{0, 2, 4, 6, 8})
				So(second, ShouldResemble, first)
			})
		})
	})
}
