package lo_test

import (
	"sort"
	"testing"

	upstream "github.com/samber/lo"
	forge "goforge.dev/gplodash"
)

func TestCaseConversions(t *testing.T) {
	checkProp(t, "CamelCase", func(s string) bool {
		return forge.CamelCase(s) == upstream.CamelCase(s)
	})
	checkProp(t, "PascalCase", func(s string) bool {
		return forge.PascalCase(s) == upstream.PascalCase(s)
	})
	checkProp(t, "Capitalize", func(s string) bool {
		return forge.Capitalize(s) == upstream.Capitalize(s)
	})
}

func TestSampleByDeterministic(t *testing.T) {
	// With a fixed generator, SampleBy/SamplesBy are deterministic and must
	// match upstream exactly.
	gen := func(n int) int { return 0 } // always pick index 0
	checkProp(t, "SampleBy", func(xs []int) bool {
		if len(xs) == 0 {
			return forge.SampleBy(xs, gen) == upstream.SampleBy(xs, gen)
		}
		return forge.SampleBy(xs, gen) == upstream.SampleBy(xs, gen)
	})
	checkProp(t, "SamplesBy", func(xs []int, c uint8) bool {
		count := int(c) % 8
		f := forge.SamplesBy(xs, count, gen)
		u := upstream.SamplesBy(xs, count, gen)
		if len(f) != len(u) {
			return false
		}
		for i := range f {
			if f[i] != u[i] {
				return false
			}
		}
		return true
	})
}

func TestRandomInvariants(t *testing.T) {
	checkProp(t, "Shuffle preserves multiset", func(xs []int) bool {
		cp := append([]int(nil), xs...)
		forge.Shuffle(cp)
		a := append([]int(nil), xs...)
		b := append([]int(nil), cp...)
		sort.Ints(a)
		sort.Ints(b)
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return len(a) == len(b)
	})
	checkProp(t, "Sample from collection", func(xs []int) bool {
		if len(xs) == 0 {
			return true
		}
		got := forge.Sample(xs)
		for _, x := range xs {
			if x == got {
				return true
			}
		}
		return false
	})
	checkProp(t, "RandomString length+charset", func(sz uint8) bool {
		size := int(sz)%20 + 1
		charset := []rune("abcXYZ123")
		s := forge.RandomString(size, charset)
		if len([]rune(s)) != size {
			return false
		}
		set := make(map[rune]bool)
		for _, r := range charset {
			set[r] = true
		}
		for _, r := range s {
			if !set[r] {
				return false
			}
		}
		return true
	})
}
