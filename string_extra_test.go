package lo_test

import (
	"errors"
	"testing"

	upstream "github.com/samber/lo"
	forge "goforge.dev/gplodash"
)

func TestSubstringProp(t *testing.T) {
	checkProp(t, "Substring", func(s string, off int, length uint16) bool {
		return forge.Substring(s, off, uint(length)) == upstream.Substring(s, off, uint(length))
	})
}

func TestMathErrProp(t *testing.T) {
	iter := func(fail bool) func(int) (int, error) {
		return func(x int) (int, error) {
			if fail && x == 0 {
				return 0, errors.New("zero")
			}
			return x, nil
		}
	}
	checkProp(t, "SumByErr/ProductByErr", func(xs []int, fail bool) bool {
		fs, fe := forge.SumByErr(xs, iter(fail))
		us, ue := upstream.SumByErr(xs, iter(fail))
		fp, fpe := forge.ProductByErr(xs, iter(fail))
		up, upe := upstream.ProductByErr(xs, iter(fail))
		return fs == us && (fe == nil) == (ue == nil) &&
			fp == up && (fpe == nil) == (upe == nil)
	})
	checkProp(t, "MeanByErr", func(xs []int, fail bool) bool {
		fm, fe := forge.MeanByErr(xs, iter(fail))
		um, ue := upstream.MeanByErr(xs, iter(fail))
		return fm == um && (fe == nil) == (ue == nil)
	})
}

func TestCharsets(t *testing.T) {
	eq := func(a, b []rune) bool {
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
	if !eq(forge.AllCharset, upstream.AllCharset) ||
		!eq(forge.AlphanumericCharset, upstream.AlphanumericCharset) ||
		!eq(forge.LettersCharset, upstream.LettersCharset) ||
		!eq(forge.SpecialCharset, upstream.SpecialCharset) ||
		!eq(forge.NumbersCharset, upstream.NumbersCharset) {
		t.Error("charset mismatch vs upstream")
	}
}
