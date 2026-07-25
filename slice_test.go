package lo_test

import (
	"reflect"
	"testing"

	upstream "github.com/samber/lo"
	forge "goforge.dev/gplodash"
)

func eqI(a, b []int) bool { return reflect.DeepEqual(a, b) }

func TestSliceCountsAndSort(t *testing.T) {
	checkProp(t, "Count/CountBy", func(xs []int, v int) bool {
		by := func(x int) bool { return x > v }
		return forge.Count(xs, v) == upstream.Count(xs, v) &&
			forge.CountBy(xs, by) == upstream.CountBy(xs, by)
	})
	checkProp(t, "CountValues", func(xs []int) bool {
		return reflect.DeepEqual(forge.CountValues(xs), upstream.CountValues(xs))
	})
	checkProp(t, "IsSorted", func(xs []int) bool {
		return forge.IsSorted(xs) == upstream.IsSorted(xs)
	})
}

func TestSliceCutTrim(t *testing.T) {
	checkProp(t, "Cut", func(a, sep []int) bool {
		fb, fa, ff := forge.Cut(a, sep)
		ub, ua, uf := upstream.Cut(a, sep)
		return eqI(fb, ub) && eqI(fa, ua) && ff == uf
	})
	checkProp(t, "CutPrefix/CutSuffix", func(a, sep []int) bool {
		fa, ff := forge.CutPrefix(a, sep)
		ua, uf := upstream.CutPrefix(a, sep)
		fb, fg := forge.CutSuffix(a, sep)
		ub, ug := upstream.CutSuffix(a, sep)
		return eqI(fa, ua) && ff == uf && eqI(fb, ub) && fg == ug
	})
	checkProp(t, "Trim/TrimPrefix/TrimSuffix", func(a, set []int) bool {
		return eqI(forge.Trim(a, set), upstream.Trim(a, set)) &&
			eqI(forge.TrimPrefix(a, set), upstream.TrimPrefix(a, set)) &&
			eqI(forge.TrimSuffix(a, set), upstream.TrimSuffix(a, set))
	})
	checkProp(t, "HasPrefix/HasSuffix", func(a, b []int) bool {
		return forge.HasPrefix(a, b) == upstream.HasPrefix(a, b) &&
			forge.HasSuffix(a, b) == upstream.HasSuffix(a, b)
	})
}

func TestSliceEditsAndTake(t *testing.T) {
	checkProp(t, "DropByIndex", func(xs []int, idx []int) bool {
		return eqI(forge.DropByIndex(xs, idx...), upstream.DropByIndex(xs, idx...))
	})
	checkProp(t, "DropWhile/DropRightWhile", func(xs []int, p int) bool {
		pred := func(x int) bool { return x < p }
		return eqI(forge.DropWhile(xs, pred), upstream.DropWhile(xs, pred)) &&
			eqI(forge.DropRightWhile(xs, pred), upstream.DropRightWhile(xs, pred))
	})
	checkProp(t, "TakeWhile", func(xs []int, p int) bool {
		pred := func(x int) bool { return x < p }
		return eqI(forge.TakeWhile(xs, pred), upstream.TakeWhile(xs, pred))
	})
	checkProp(t, "Replace/ReplaceAll", func(xs []int, old, nEw int, n int8) bool {
		return eqI(forge.Replace(xs, old, nEw, int(n)), upstream.Replace(xs, old, nEw, int(n))) &&
			eqI(forge.ReplaceAll(xs, old, nEw), upstream.ReplaceAll(xs, old, nEw))
	})
	checkProp(t, "Slice/Subset", func(xs []int, a, b int8) bool {
		length := uint(uint8(b))
		return eqI(forge.Slice(xs, int(a), int(b)), upstream.Slice(xs, int(a), int(b))) &&
			eqI(forge.Subset(xs, int(a), length), upstream.Subset(xs, int(a), length))
	})
	checkProp(t, "Splice", func(xs []int, i int8, els []int) bool {
		return eqI(forge.Splice(xs, int(i), els...), upstream.Splice(xs, int(i), els...))
	})
	checkProp(t, "UniqMap/SliceToMap", func(xs []int) bool {
		// UniqMap returns Keys(set) — non-deterministic order; compare sorted.
		um := forge.UniqMap(xs, func(x, _ int) int { return x % 3 })
		uu := upstream.UniqMap(xs, func(x, _ int) int { return x % 3 })
		sm := forge.SliceToMap(xs, func(x int) (int, int) { return x, x * 2 })
		su := upstream.SliceToMap(xs, func(x int) (int, int) { return x, x * 2 })
		return reflect.DeepEqual(sortedInts(um), sortedInts(uu)) && reflect.DeepEqual(sm, su)
	})
}
