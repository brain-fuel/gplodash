package it_test

import (
	"iter"
	"reflect"
	"slices"
	"sort"
	"testing"
	"testing/quick"

	upstream "github.com/samber/lo/it"
	forge "goforge.dev/gplodash/it"
)

func check(t *testing.T, name string, f any) {
	t.Helper()
	if err := quick.Check(f, &quick.Config{MaxCount: 1500}); err != nil {
		t.Errorf("%s: %v", name, err)
	}
}

func seq(xs []int) iter.Seq[int] { return slices.Values(xs) }

func sortedI(xs []int) []int {
	out := append([]int(nil), xs...)
	sort.Ints(out)
	return out
}

// Lazy transforms materialize identically to upstream (order-preserving).
func TestItTransforms(t *testing.T) {
	check(t, "Map", func(xs []int) bool {
		f := slices.Collect(forge.Map(seq(xs), func(x int) int { return x*2 + 1 }))
		u := slices.Collect(upstream.Map(seq(xs), func(x int) int { return x*2 + 1 }))
		return reflect.DeepEqual(f, u)
	})
	check(t, "Filter", func(xs []int) bool {
		p := func(x int) bool { return x%2 == 0 }
		f := slices.Collect(forge.Filter(seq(xs), p))
		u := slices.Collect(upstream.Filter(seq(xs), p))
		return reflect.DeepEqual(f, u)
	})
	check(t, "Uniq", func(xs []int) bool {
		f := slices.Collect(forge.Uniq(seq(xs)))
		u := slices.Collect(upstream.Uniq(seq(xs)))
		return reflect.DeepEqual(f, u)
	})
}

// Terminal reductions equal upstream.
func TestItReductions(t *testing.T) {
	check(t, "Reduce", func(xs []int) bool {
		acc := func(a, x int) int { return a + x }
		return forge.Reduce(seq(xs), acc, 0) == upstream.Reduce(seq(xs), acc, 0)
	})
	check(t, "Contains", func(xs []int, v int) bool {
		return forge.Contains(seq(xs), v) == upstream.Contains(seq(xs), v)
	})
	check(t, "Sum", func(xs []int) bool {
		return forge.Sum(seq(xs)) == upstream.Sum(seq(xs))
	})
	check(t, "GroupBy", func(xs []int) bool {
		by := func(x int) int { return x % 3 }
		return reflect.DeepEqual(forge.GroupBy(seq(xs), by), upstream.GroupBy(seq(xs), by))
	})
}

// Keys/Values from maps are order-non-deterministic — compare sorted.
func TestItMap(t *testing.T) {
	check(t, "Keys", func(m map[int]int) bool {
		f := slices.Collect(forge.Keys(m))
		u := slices.Collect(upstream.Keys(m))
		return reflect.DeepEqual(sortedI(f), sortedI(u))
	})
}
