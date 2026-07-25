package lo_test

import (
	"reflect"
	"testing"
	"testing/quick"

	upstream "github.com/samber/lo"
	forge "goforge.dev/gplodash"
)

// The math batch is an exact reimplementation of samber/lo's math.go. Each law
// below is a property: for arbitrary inputs, the GoForge result must equal the
// pinned upstream result. testing/quick generates and (on failure) reports the
// counterexample.

func checkProp(t *testing.T, name string, f any) {
	t.Helper()
	if err := quick.Check(f, &quick.Config{MaxCount: 2000}); err != nil {
		t.Errorf("%s: %v", name, err)
	}
}

func TestSumProductProp(t *testing.T) {
	checkProp(t, "Sum/int", func(xs []int) bool {
		return forge.Sum(xs) == upstream.Sum(xs)
	})
	checkProp(t, "Sum/float64", func(xs []float64) bool {
		return forge.Sum(xs) == upstream.Sum(xs)
	})
	checkProp(t, "Product/int", func(xs []int) bool {
		return forge.Product(xs) == upstream.Product(xs)
	})
	checkProp(t, "SumBy/int", func(xs []int) bool {
		by := func(x int) int { return x * 2 }
		return forge.SumBy(xs, by) == upstream.SumBy(xs, by)
	})
	checkProp(t, "ProductBy/int", func(xs []int8) bool {
		by := func(x int8) int64 { return int64(x) }
		return forge.ProductBy(xs, by) == upstream.ProductBy(xs, by)
	})
}

func TestMeanProp(t *testing.T) {
	checkProp(t, "Mean/int", func(xs []int) bool {
		return forge.Mean(xs) == upstream.Mean(xs) // truncating integer division
	})
	checkProp(t, "Mean/float64", func(xs []float64) bool {
		return forge.Mean(xs) == upstream.Mean(xs)
	})
	checkProp(t, "MeanBy/int", func(xs []int) bool {
		by := func(x int) int { return x }
		return forge.MeanBy(xs, by) == upstream.MeanBy(xs, by)
	})
}

func TestModeProp(t *testing.T) {
	checkProp(t, "Mode/int", func(xs []int) bool {
		return reflect.DeepEqual(forge.Mode(xs), upstream.Mode(xs))
	})
}

func TestClampProp(t *testing.T) {
	checkProp(t, "Clamp/int", func(v, a, b int) bool {
		lo, hi := a, b
		if lo > hi {
			lo, hi = hi, lo
		}
		return forge.Clamp(v, lo, hi) == upstream.Clamp(v, lo, hi)
	})
	checkProp(t, "Clamp/string", func(v, a, b string) bool {
		lo, hi := a, b
		if lo > hi {
			lo, hi = hi, lo
		}
		return forge.Clamp(v, lo, hi) == upstream.Clamp(v, lo, hi)
	})
}

func TestRangeProp(t *testing.T) {
	checkProp(t, "Range", func(n int16) bool { // bounded so allocation stays sane
		return reflect.DeepEqual(forge.Range(int(n)), upstream.Range(int(n)))
	})
	checkProp(t, "RangeFrom/int", func(start int, n int8) bool {
		return reflect.DeepEqual(forge.RangeFrom(start, int(n)), upstream.RangeFrom(start, int(n)))
	})
	checkProp(t, "RangeWithSteps/int", func(start, end, step int8) bool {
		a := forge.RangeWithSteps(int(start), int(end), int(step))
		b := upstream.RangeWithSteps(int(start), int(end), int(step))
		return reflect.DeepEqual(a, b)
	})
	checkProp(t, "RangeWithSteps/float64", func(s, e, st int8) bool {
		// derive small floats from int8 to keep steps meaningful and finite
		start, end, step := float64(s)/4, float64(e)/4, float64(st)/4
		return reflect.DeepEqual(
			forge.RangeWithSteps(start, end, step),
			upstream.RangeWithSteps(start, end, step),
		)
	})
}
