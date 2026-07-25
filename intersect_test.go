package lo_test

import (
	"reflect"
	"testing"

	upstream "github.com/samber/lo"
	forge "goforge.dev/gplodash"
)

func TestIntersectBatch(t *testing.T) {
	pred := func(x int) bool { return x%2 == 0 }
	id := func(x int) int { return x }

	checkProp(t, "ContainsBy/EveryBy/SomeBy/NoneBy", func(xs []int) bool {
		return forge.ContainsBy(xs, pred) == upstream.ContainsBy(xs, pred) &&
			forge.EveryBy(xs, pred) == upstream.EveryBy(xs, pred) &&
			forge.SomeBy(xs, pred) == upstream.SomeBy(xs, pred) &&
			forge.NoneBy(xs, pred) == upstream.NoneBy(xs, pred)
	})
	checkProp(t, "IntersectBy", func(a, b, c []int) bool {
		return reflect.DeepEqual(
			forge.IntersectBy(id, a, b, c),
			upstream.IntersectBy(id, a, b, c),
		)
	})
	checkProp(t, "WithoutBy", func(xs, ex []int) bool {
		return reflect.DeepEqual(
			forge.WithoutBy(xs, id, ex...),
			upstream.WithoutBy(xs, id, ex...),
		)
	})
	checkProp(t, "WithoutEmpty/Compact", func(xs []int) bool {
		return reflect.DeepEqual(forge.WithoutEmpty(xs), upstream.WithoutEmpty(xs)) &&
			reflect.DeepEqual(forge.Compact(xs), upstream.Compact(xs))
	})
	checkProp(t, "WithoutNth", func(xs []int, nths []int) bool {
		return reflect.DeepEqual(
			forge.WithoutNth(xs, nths...),
			upstream.WithoutNth(xs, nths...),
		)
	})
	checkProp(t, "ElementsMatch", func(a, b []int) bool {
		return forge.ElementsMatch(a, b) == upstream.ElementsMatch(a, b)
	})
	checkProp(t, "Keyify", func(xs []int) bool {
		return reflect.DeepEqual(forge.Keyify(xs), upstream.Keyify(xs))
	})
}
