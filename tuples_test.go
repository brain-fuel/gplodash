package lo_test

import (
	"testing"

	upstream "github.com/samber/lo"
	forge "goforge.dev/gplodash"
)

// The tuple/zip/unzip/crossjoin batch and the min/max/nth family are exact
// reimplementations. Because forge.TupleN and upstream.TupleN are distinct
// named types, laws compare the *unpacked* fields rather than the typed tuples.

func eqPairsIS(fs []forge.Tuple2[int, string], us []upstream.Tuple2[int, string]) bool {
	if len(fs) != len(us) {
		return false
	}
	for i := range fs {
		fa, fb := fs[i].Unpack()
		ua, ub := us[i].Unpack()
		if fa != ua || fb != ub {
			return false
		}
	}
	return true
}

func eqPairsII(fs []forge.Tuple2[int, int], us []upstream.Tuple2[int, int]) bool {
	if len(fs) != len(us) {
		return false
	}
	for i := range fs {
		fa, fb := fs[i].Unpack()
		ua, ub := us[i].Unpack()
		if fa != ua || fb != ub {
			return false
		}
	}
	return true
}

func TestTupleConstructAndUnpack(t *testing.T) {
	checkProp(t, "T2/Unpack2", func(a int, b string) bool {
		fa, fb := forge.Unpack2(forge.T2(a, b))
		ua, ub := upstream.Unpack2(upstream.T2(a, b))
		return fa == ua && fb == ub
	})
	checkProp(t, "T3.Unpack", func(a, b, c int) bool {
		x1, x2, x3 := forge.T3(a, b, c).Unpack()
		y1, y2, y3 := upstream.T3(a, b, c).Unpack()
		return x1 == y1 && x2 == y2 && x3 == y3
	})
}

func TestZipUnzip(t *testing.T) {
	checkProp(t, "Zip2", func(a []int, b []string) bool {
		return eqPairsIS(forge.Zip2(a, b), upstream.Zip2(a, b))
	})
	checkProp(t, "ZipBy2", func(a, b []int) bool {
		f := func(x, y int) int { return x + y }
		fr := forge.ZipBy2(a, b, f)
		ur := upstream.ZipBy2(a, b, f)
		if len(fr) != len(ur) {
			return false
		}
		for i := range fr {
			if fr[i] != ur[i] {
				return false
			}
		}
		return true
	})
	checkProp(t, "Unzip2", func(seed []int) bool {
		ft := forge.Map(seed, func(x, _ int) forge.Tuple2[int, int] { return forge.T2(x, x*2) })
		ut := upstream.Map(seed, func(x, _ int) upstream.Tuple2[int, int] { return upstream.T2(x, x*2) })
		fa, fb := forge.Unzip2(ft)
		ua, ub := upstream.Unzip2(ut)
		return intsEq(fa, ua) && intsEq(fb, ub)
	})
}

func TestCrossJoin(t *testing.T) {
	checkProp(t, "CrossJoin2", func(a []int, b []int) bool {
		return eqPairsII(forge.CrossJoin2(a, b), upstream.CrossJoin2(a, b))
	})
	checkProp(t, "CrossJoinBy2", func(a, b []int) bool {
		f := func(x, y int) int { return x*10 + y }
		return intsEq(forge.CrossJoinBy2(a, b, f), upstream.CrossJoinBy2(a, b, f))
	})
}

func intsEq(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestMinMaxNth(t *testing.T) {
	checkProp(t, "Min/Max int", func(xs []int) bool {
		return forge.Min(xs) == upstream.Min(xs) && forge.Max(xs) == upstream.Max(xs)
	})
	checkProp(t, "MinIndex/MaxIndex", func(xs []int) bool {
		fm, fi := forge.MinIndex(xs)
		um, ui := upstream.MinIndex(xs)
		fM, fI := forge.MaxIndex(xs)
		uM, uI := upstream.MaxIndex(xs)
		return fm == um && fi == ui && fM == uM && fI == uI
	})
	checkProp(t, "MinBy/MaxBy", func(xs []int) bool {
		less := func(a, b int) bool { return a < b }
		greater := func(a, b int) bool { return a > b }
		return forge.MinBy(xs, less) == upstream.MinBy(xs, less) &&
			forge.MaxBy(xs, greater) == upstream.MaxBy(xs, greater)
	})
	checkProp(t, "NthOr/NthOrEmpty", func(xs []int, n int8) bool {
		return forge.NthOr(xs, int(n), -999) == upstream.NthOr(xs, int(n), -999) &&
			forge.NthOrEmpty(xs, int(n)) == upstream.NthOrEmpty(xs, int(n))
	})
}
