package lo

import (
	"encoding/csv"
	"os"
	"reflect"
	"testing"
	"testing/quick"

	upstream "github.com/samber/lo"
)

func TestCollectionLaws(t *testing.T) {
	properties := map[string]func([]int) bool{
		"map identity": func(values []int) bool {
			return reflect.DeepEqual(Map(values, func(value, _ int) int { return value }), upstream.Map(values, func(value, _ int) int { return value }))
		},
		"map composition": func(values []int) bool {
			left := Map(Map(values, func(value, _ int) int { return value + 1 }), func(value, _ int) int { return value * 2 })
			right := Map(values, func(value, _ int) int { return (value + 1) * 2 })
			return reflect.DeepEqual(left, right)
		},
		"filter idempotence": func(values []int) bool {
			predicate := func(value, _ int) bool { return value%2 == 0 }
			once := Filter(values, predicate)
			return reflect.DeepEqual(Filter(once, predicate), once)
		},
		"uniq idempotence":   func(values []int) bool { return reflect.DeepEqual(Uniq(Uniq(values)), Uniq(values)) },
		"reverse involution": func(values []int) bool { return reflect.DeepEqual(ReverseCopy(ReverseCopy(values)), values) },
	}
	for name, property := range properties {
		t.Run(name, func(t *testing.T) {
			if err := quick.Check(property, &quick.Config{MaxCount: 1000}); err != nil {
				t.Fatal(err)
			}
		})
	}
	if !reflect.DeepEqual(Concat(Concat([]int{1}, []int{2}), []int{3}), Concat([]int{1}, Concat([]int{2}, []int{3}))) {
		t.Fatal("Concat is not associative")
	}
}

func FuzzCollectionDifferential(f *testing.F) {
	for _, seed := range [][]byte{nil, {}, {1}, {1, 2, 2, 3}, {255, 0, 127}} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, raw []byte) {
		values := make([]int, len(raw))
		for i, value := range raw {
			values[i] = int(int8(value))
		}
		same(t, "Map", Map(values, func(v, i int) int { return v + i }), upstream.Map(values, func(v, i int) int { return v + i }))
		same(t, "Filter", Filter(values, func(v, _ int) bool { return v%3 == 0 }), upstream.Filter(values, func(v, _ int) bool { return v%3 == 0 }))
		same(t, "Uniq", Uniq(values), upstream.Uniq(values))
		same(t, "PartitionBy", PartitionBy(values, func(v int) int { return v % 5 }), upstream.PartitionBy(values, func(v int) int { return v % 5 }))
		reversed := ReverseCopy(values)
		theirsReversed := upstream.Reverse(append([]int(nil), values...))
		same(t, "Intersect", Intersect(values, reversed), upstream.Intersect(values, theirsReversed))
	})
}

func TestAPIManifestIsCompleteAndExplicit(t *testing.T) {
	file, err := os.Open("API_MANIFEST.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 652 {
		t.Fatalf("manifest rows = %d, want header + 651 declarations", len(records))
	}
	compatible := 0
	for row, record := range records[1:] {
		if len(record) != 6 || record[4] == "" || record[5] == "" {
			t.Fatalf("manifest row %d is incomplete: %#v", row+2, record)
		}
		if record[4] == "compatible" {
			compatible++
		}
	}
	if compatible != 651 {
		t.Fatalf("compatible declarations = %d, want 651", compatible)
	}
}
