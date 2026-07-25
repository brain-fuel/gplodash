package lo_test

import (
	"reflect"
	"sort"
	"testing"

	upstream "github.com/samber/lo"
	forge "goforge.dev/gplodash"
)

func sortedInts(xs []int) []int {
	out := append([]int(nil), xs...)
	sort.Ints(out)
	return out
}

// Map-returning helpers are deterministic (map equality ignores order);
// slice-returning helpers over a map are compared order-insensitively because
// Go randomizes map iteration per range.

func TestMapProjectionProp(t *testing.T) {
	checkProp(t, "Assign", func(a, b map[int]int) bool {
		return reflect.DeepEqual(forge.Assign(a, b), upstream.Assign(a, b))
	})
	checkProp(t, "Invert", func(m map[int]int) bool {
		return reflect.DeepEqual(forge.Invert(m), upstream.Invert(m))
	})
	checkProp(t, "ValueOr", func(m map[int]int, k, fb int) bool {
		return forge.ValueOr(m, k, fb) == upstream.ValueOr(m, k, fb)
	})
	checkProp(t, "HasKey", func(m map[int]int, k int) bool {
		return forge.HasKey(m, k) == upstream.HasKey(m, k)
	})
}

func TestMapFilterProp(t *testing.T) {
	pk := func(k, v int) bool { return (k+v)%2 == 0 }
	checkProp(t, "PickBy/OmitBy", func(m map[int]int) bool {
		return reflect.DeepEqual(forge.PickBy(m, pk), upstream.PickBy(m, pk)) &&
			reflect.DeepEqual(forge.OmitBy(m, pk), upstream.OmitBy(m, pk))
	})
	checkProp(t, "PickByKeys/OmitByKeys", func(m map[int]int, keys []int) bool {
		return reflect.DeepEqual(forge.PickByKeys(m, keys), upstream.PickByKeys(m, keys)) &&
			reflect.DeepEqual(forge.OmitByKeys(m, keys), upstream.OmitByKeys(m, keys))
	})
	checkProp(t, "FilterKeys/FilterValues", func(m map[int]int) bool {
		// FilterKeys/FilterValues return slices drawn from a map, so order is
		// non-deterministic; compare as sorted sequences.
		kp := func(k, v int) bool { return k > 0 }
		vp := func(k, v int) bool { return v > 0 }
		return reflect.DeepEqual(sortedInts(forge.FilterKeys(m, kp)), sortedInts(upstream.FilterKeys(m, kp))) &&
			reflect.DeepEqual(sortedInts(forge.FilterValues(m, vp)), sortedInts(upstream.FilterValues(m, vp)))
	})
}

func TestMapSliceProp(t *testing.T) {
	checkProp(t, "MapToSlice", func(m map[int]int) bool {
		f := forge.MapToSlice(m, func(k, v int) int { return k*1000 + v })
		u := upstream.MapToSlice(m, func(k, v int) int { return k*1000 + v })
		return reflect.DeepEqual(sortedInts(f), sortedInts(u))
	})
	checkProp(t, "UniqKeys/UniqValues", func(a, b map[int]int) bool {
		fk := forge.UniqKeys(a, b)
		uk := upstream.UniqKeys(a, b)
		fv := forge.UniqValues(a, b)
		uv := upstream.UniqValues(a, b)
		return reflect.DeepEqual(sortedInts(fk), sortedInts(uk)) &&
			reflect.DeepEqual(sortedInts(fv), sortedInts(uv))
	})
	checkProp(t, "FromPairs roundtrip", func(m map[int]int) bool {
		// ToPairs order varies; FromPairs(ToPairs(m)) == m is the stable law.
		return reflect.DeepEqual(forge.FromPairs(forge.ToPairs(m)), m) &&
			reflect.DeepEqual(upstream.FromPairs(upstream.ToPairs(m)), m)
	})
}
