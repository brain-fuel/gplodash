package lo

import (
	"reflect"
	"testing"

	upstream "github.com/samber/lo"
	"goforge.dev/goplus/std/nonempty"
	"goforge.dev/goplus/std/option"
	"goforge.dev/goplus/std/result"
)

func same[T any](t *testing.T, name string, got, want T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s = %#v, upstream %#v", name, got, want)
	}
}

func TestSliceTransformDifferential(t *testing.T) {
	values := []int{1, 2, 2, 3, 4, 3}
	same(t, "Map", Map(values, func(v, i int) string { return string(rune('a' + v + i)) }), upstream.Map(values, func(v, i int) string { return string(rune('a' + v + i)) }))
	same(t, "Filter", Filter(values, func(v, i int) bool { return (v+i)%2 == 0 }), upstream.Filter(values, func(v, i int) bool { return (v+i)%2 == 0 }))
	same(t, "Reject", Reject(values, func(v, _ int) bool { return v%2 == 0 }), upstream.Reject(values, func(v, _ int) bool { return v%2 == 0 }))
	same(t, "FilterMap", FilterMap(values, func(v, i int) (int, bool) { return v * i, v%2 == 0 }), upstream.FilterMap(values, func(v, i int) (int, bool) { return v * i, v%2 == 0 }))
	same(t, "FlatMap", FlatMap(values, func(v, _ int) []int { return []int{v, -v} }), upstream.FlatMap(values, func(v, _ int) []int { return []int{v, -v} }))
	same(t, "Reduce", Reduce(values, func(a, v, i int) int { return a + v*i }, 7), upstream.Reduce(values, func(a, v, i int) int { return a + v*i }, 7))
	same(t, "ReduceRight", ReduceRight(values, func(a, v, i int) int { return a*2 + v + i }, 0), upstream.ReduceRight(values, func(a, v, i int) int { return a*2 + v + i }, 0))
	same(t, "Times", Times(8, func(i int) int { return i * i }), upstream.Times(8, func(i int) int { return i * i }))
	same(t, "Uniq", Uniq(values), upstream.Uniq(values))
	same(t, "UniqBy", UniqBy(values, func(v int) int { return v % 3 }), upstream.UniqBy(values, func(v int) int { return v % 3 }))
	oursReverse, theirsReverse := append([]int(nil), values...), append([]int(nil), values...)
	same(t, "Reverse", Reverse(oursReverse), upstream.Reverse(theirsReverse))
	same(t, "Drop", Drop(values, 2), upstream.Drop(values, 2))
	same(t, "DropRight", DropRight(values, 2), upstream.DropRight(values, 2))
	same(t, "Take", Take(values, 3), upstream.Take(values, 3))
}

func TestGroupingAndShapeDifferential(t *testing.T) {
	values := []int{1, 2, 3, 4, 5, 6}
	same(t, "GroupBy", GroupBy(values, func(v int) int { return v % 2 }), upstream.GroupBy(values, func(v int) int { return v % 2 }))
	same(t, "GroupByMap", GroupByMap(values, func(v int) (int, string) { return v % 2, string(rune('a' + v)) }), upstream.GroupByMap(values, func(v int) (int, string) { return v % 2, string(rune('a' + v)) }))
	same(t, "KeyBy", KeyBy(values, func(v int) int { return v % 3 }), upstream.KeyBy(values, func(v int) int { return v % 3 }))
	same(t, "Associate", Associate(values, func(v int) (int, int) { return v % 3, v * v }), upstream.Associate(values, func(v int) (int, int) { return v % 3, v * v }))
	same(t, "Chunk", Chunk(values, 4), upstream.Chunk(values, 4))
	same(t, "Window", Window(values, 3), upstream.Window(values, 3))
	same(t, "Sliding", Sliding(values, 3, 2), upstream.Sliding(values, 3, 2))
	same(t, "PartitionBy", PartitionBy(values, func(v int) int { return v % 3 }), upstream.PartitionBy(values, func(v int) int { return v % 3 }))
	same(t, "Flatten", Flatten([][]int{{1, 2}, {}, {3}}), upstream.Flatten([][]int{{1, 2}, {}, {3}}))
	same(t, "Concat", Concat([]int{1}, []int{2, 3}), upstream.Concat([]int{1}, []int{2, 3}))
	same(t, "Interleave", Interleave([]int{1, 2, 3}, []int{4}, []int{5, 6}), upstream.Interleave([]int{1, 2, 3}, []int{4}, []int{5, 6}))
}

func TestSearchAndSetDifferential(t *testing.T) {
	values := []int{1, 2, 2, 3, 4, 3}
	gotV, gotOK := Find(values, func(v int) bool { return v > 2 })
	wantV, wantOK := upstream.Find(values, func(v int) bool { return v > 2 })
	if gotV != wantV || gotOK != wantOK {
		t.Fatalf("Find = %v/%v, upstream %v/%v", gotV, gotOK, wantV, wantOK)
	}
	gotV, gotI, gotOK := FindIndexOf(values, func(v int) bool { return v == 3 })
	wantV, wantI, wantOK := upstream.FindIndexOf(values, func(v int) bool { return v == 3 })
	if gotV != wantV || gotI != wantI || gotOK != wantOK {
		t.Fatalf("FindIndexOf differs")
	}
	gotV, gotI, gotOK = FindLastIndexOf(values, func(v int) bool { return v == 3 })
	wantV, wantI, wantOK = upstream.FindLastIndexOf(values, func(v int) bool { return v == 3 })
	if gotV != wantV || gotI != wantI || gotOK != wantOK {
		t.Fatalf("FindLastIndexOf differs")
	}
	for _, index := range []int{-7, -1, 0, 5, 6} {
		got, gotErr := Nth(values, index)
		want, wantErr := upstream.Nth(values, index)
		if got != want || (gotErr == nil) != (wantErr == nil) {
			t.Fatalf("Nth(%d) = %v/%v, upstream %v/%v", index, got, gotErr, want, wantErr)
		}
	}
	same(t, "Contains", Contains(values, 4), upstream.Contains(values, 4))
	same(t, "Every", Every(values, []int{2, 4}), upstream.Every(values, []int{2, 4}))
	same(t, "Some", Some(values, []int{9, 4}), upstream.Some(values, []int{9, 4}))
	same(t, "None", None(values, []int{8, 9}), upstream.None(values, []int{8, 9}))
	same(t, "Without", Without(values, 2, 4), upstream.Without(values, 2, 4))
	same(t, "Union", Union(values, []int{5, 1}), upstream.Union(values, []int{5, 1}))
	same(t, "Intersect", Intersect(values, []int{3, 2, 8}, []int{2, 3}), upstream.Intersect(values, []int{3, 2, 8}, []int{2, 3}))
	gotLeft, gotRight := Difference(values, []int{2, 8})
	wantLeft, wantRight := upstream.Difference(values, []int{2, 8})
	same(t, "Difference-left", gotLeft, wantLeft)
	same(t, "Difference-right", gotRight, wantRight)
}

func TestTotalOptionAndResultVariants(t *testing.T) {
	if got := FirstOption([]int{7}); got != (option.Some[int]{Value: 7}) {
		t.Fatalf("FirstOption = %#v", got)
	}
	if _, ok := FindOption([]int{1, 2}, func(v int) bool { return v > 9 }).(option.None[int]); !ok {
		t.Fatal("FindOption miss is not None")
	}
	if got := NthResult([]string{"a", "b"}, -1); got != (result.Ok[string, error]{Value: "b"}) {
		t.Fatalf("NthResult = %#v", got)
	}
	if _, ok := NthResult([]string{"a"}, 2).(result.Err[string, error]); !ok {
		t.Fatal("NthResult bounds failure is not Err")
	}
}

func TestMapDifferential(t *testing.T) {
	input := map[string]int{"a": 1, "b": 2}
	same(t, "FromEntries", FromEntries(Entries(input)), input)
	same(t, "MapKeys", MapKeys(input, func(v int, k string) int { return v + len(k) }), upstream.MapKeys(input, func(v int, k string) int { return v + len(k) }))
	same(t, "MapValues", MapValues(input, func(v int, k string) string { return k + string(rune('0'+v)) }), upstream.MapValues(input, func(v int, k string) string { return k + string(rune('0'+v)) }))
	if len(Keys(input)) != len(upstream.Keys(input)) || len(Values(input)) != len(upstream.Values(input)) {
		t.Fatal("map cardinality differs")
	}
}

func TestFreshAllocationAndNonEmptySemantics(t *testing.T) {
	source := []int{1, 2, 3, 4}
	for name, result := range map[string][]int{
		"Map":    Map(source, func(v, _ int) int { return v }),
		"Filter": Filter(source, func(int, int) bool { return true }),
		"Take":   Take(source, len(source)),
		"Drop":   Drop(source, 0),
	} {
		result[0] = 99
		if source[0] != 1 {
			t.Fatalf("%s aliases its input", name)
		}
	}
	chunks := Chunk(source, 2)
	chunks[0][0] = 99
	if source[0] != 1 {
		t.Fatal("Chunk aliases its input")
	}
	if got := ReduceNonEmpty(nonempty.Of(1, 2, 3), func(a, b int) int { return a + b }); got != 6 {
		t.Fatalf("ReduceNonEmpty = %d", got)
	}
}

func TestInvalidSizesPanicLikeUpstream(t *testing.T) {
	for name, operation := range map[string]func(){
		"Chunk": func() { Chunk([]int{1}, 0) }, "Window": func() { Window([]int{1}, 0) },
		"Sliding-size": func() { Sliding([]int{1}, 0, 1) }, "Sliding-step": func() { Sliding([]int{1}, 1, 0) },
		"Drop": func() { Drop([]int{1}, -1) }, "DropRight": func() { DropRight([]int{1}, -1) }, "Take": func() { Take([]int{1}, -1) },
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("did not panic")
				}
			}()
			operation()
		})
	}
}

func TestExhaustiveSearch(t *testing.T) {
	value, index, found := SearchGet(Locate([]int{1, 4, 6}, func(value int) bool {
		return value%2 == 0
	}))
	if !found || value != 4 || index != 1 {
		t.Fatalf("Locate = %d/%d/%v", value, index, found)
	}
	value, index, found = SearchGet(LocateLast([]int{1, 4, 6}, func(value int) bool {
		return value%2 == 0
	}))
	if !found || value != 6 || index != 2 {
		t.Fatalf("LocateLast = %d/%d/%v", value, index, found)
	}
	value, index, found = SearchGet(Locate([]int{1, 3}, func(value int) bool {
		return value%2 == 0
	}))
	if found || value != 0 || index != -1 {
		t.Fatalf("missing Locate = %d/%d/%v", value, index, found)
	}
}
