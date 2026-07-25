package lo_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	upstream "github.com/samber/lo"
	forge "goforge.dev/gplodash"
)

func TestFindIndexAndFirstLast(t *testing.T) {
	checkProp(t, "IndexOf/LastIndexOf", func(xs []int, v int) bool {
		return forge.IndexOf(xs, v) == upstream.IndexOf(xs, v) &&
			forge.LastIndexOf(xs, v) == upstream.LastIndexOf(xs, v)
	})
	checkProp(t, "FirstOr/FirstOrEmpty/LastOr/LastOrEmpty", func(xs []int, fb int) bool {
		return forge.FirstOr(xs, fb) == upstream.FirstOr(xs, fb) &&
			forge.FirstOrEmpty(xs) == upstream.FirstOrEmpty(xs) &&
			forge.LastOr(xs, fb) == upstream.LastOr(xs, fb) &&
			forge.LastOrEmpty(xs) == upstream.LastOrEmpty(xs)
	})
}

func TestFindKeyAndOrElse(t *testing.T) {
	checkProp(t, "FindKey", func(m map[int]int, v int) bool {
		fk, fok := forge.FindKey(m, v)
		uk, uok := upstream.FindKey(m, v)
		// key search over a map may find any matching key; agree on existence.
		return fok == uok && (!fok || m[fk] == v && m[uk] == v)
	})
	checkProp(t, "FindOrElse", func(xs []int, fb, p int) bool {
		pred := func(x int) bool { return x > p }
		return forge.FindOrElse(xs, fb, pred) == upstream.FindOrElse(xs, fb, pred)
	})
	checkProp(t, "FindErr", func(xs []int, p int) bool {
		pred := func(x int) (bool, error) {
			if x == 0 {
				return false, errors.New("zero")
			}
			return x > p, nil
		}
		fv, fe := forge.FindErr(xs, pred)
		uv, ue := upstream.FindErr(xs, pred)
		return fv == uv && (fe == nil) == (ue == nil)
	})
}

func TestFindDuplicatesUniques(t *testing.T) {
	checkProp(t, "FindDuplicates", func(xs []int) bool {
		return reflect.DeepEqual(forge.FindDuplicates(xs), upstream.FindDuplicates(xs))
	})
	checkProp(t, "FindUniques", func(xs []int) bool {
		return reflect.DeepEqual(forge.FindUniques(xs), upstream.FindUniques(xs))
	})
	checkProp(t, "MinByErr/MaxByErr", func(xs []int) bool {
		less := func(a, b int) bool { return a < b }
		greater := func(a, b int) bool { return a > b }
		fmn, femn := forge.MinByErr(xs, func(a, b int) (bool, error) { return less(a, b), nil })
		umn, uemn := upstream.MinByErr(xs, func(a, b int) (bool, error) { return less(a, b), nil })
		fmx, femx := forge.MaxByErr(xs, func(a, b int) (bool, error) { return greater(a, b), nil })
		umx, uemx := upstream.MaxByErr(xs, func(a, b int) (bool, error) { return greater(a, b), nil })
		return fmn == umn && (femn == nil) == (uemn == nil) &&
			fmx == umx && (femx == nil) == (uemx == nil)
	})
}

func TestEarliestLatest(t *testing.T) {
	checkProp(t, "Earliest/Latest", func(secs []int32) bool {
		ts := make([]time.Time, len(secs))
		for i, s := range secs {
			ts[i] = time.Unix(int64(s), 0)
		}
		return forge.Earliest(ts...).Equal(upstream.Earliest(ts...)) &&
			forge.Latest(ts...).Equal(upstream.Latest(ts...))
	})
}
