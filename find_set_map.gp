// Go+ collection, set, and map compatibility operations.
package lo

import (
	"fmt"

	"goforge.dev/goplus/std/option"
	"goforge.dev/goplus/std/result"
)

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Search is the exhaustive Go+ alternative to the compatibility API's
// zero-value/bool and sentinel-index triples.
type Search[T any] enum {
	Located(Value T, Index int)
	Absent
}

// Locate returns both the matching value and its index without sentinel
// values. Go+ callers must handle both variants.
func Locate[T any](collection []T, predicate func(T) bool) Search[T] {
	for index, item := range collection {
		if predicate(item) {
			return Located(item, index)
		}
	}
	return Absent
}

// LocateLast is the reverse-search counterpart to Locate.
func LocateLast[T any](collection []T, predicate func(T) bool) Search[T] {
	for index := len(collection) - 1; index >= 0; index-- {
		if predicate(collection[index]) {
			return Located(collection[index], index)
		}
	}
	return Absent
}

// SearchGet is the ordinary-Go erasure boundary for Search.
func SearchGet[T any](search Search[T]) (T, int, bool) {
	match search {
	case Located(value, index):
		return value, index, true
	case Absent:
		var zero T
		return zero, -1, false
	}
}

func Find[T any](collection []T, predicate func(T) bool) (T, bool) {
	for _, item := range collection {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

func FindIndexOf[T any](collection []T, predicate func(T) bool) (T, int, bool) {
	for i, item := range collection {
		if predicate(item) {
			return item, i, true
		}
	}
	var zero T
	return zero, -1, false
}

func FindLastIndexOf[T any](collection []T, predicate func(T) bool) (T, int, bool) {
	for i := len(collection) - 1; i >= 0; i-- {
		if predicate(collection[i]) {
			return collection[i], i, true
		}
	}
	var zero T
	return zero, -1, false
}

func First[T any](collection []T) (T, bool) {
	if len(collection) != 0 {
		return collection[0], true
	}
	var zero T
	return zero, false
}

func Last[T any](collection []T) (T, bool) {
	if len(collection) != 0 {
		return collection[len(collection)-1], true
	}
	var zero T
	return zero, false
}

func Nth[T any, N Integer](collection []T, nth N) (T, error) {
	index := int(nth)
	if index < 0 {
		index = len(collection) + index
	}
	if index < 0 || index >= len(collection) {
		var zero T
		return zero, fmt.Errorf("nth: %d out of slice bounds", nth)
	}
	return collection[index], nil
}

func FindOption[T any](collection []T, predicate func(T) bool) option.Option[T] {
	if value, ok := Find(collection, predicate); ok {
		return option.Some[T]{Value: value}
	}
	return option.None[T]{}
}

func FirstOption[T any](collection []T) option.Option[T] {
	if value, ok := First(collection); ok {
		return option.Some[T]{Value: value}
	}
	return option.None[T]{}
}

func LastOption[T any](collection []T) option.Option[T] {
	if value, ok := Last(collection); ok {
		return option.Some[T]{Value: value}
	}
	return option.None[T]{}
}

func NthResult[T any, N Integer](collection []T, nth N) result.Result[T, error] {
	value, err := Nth(collection, nth)
	if err != nil {
		return result.Err[T, error]{Err: err}
	}
	return result.Ok[T, error]{Value: value}
}

func Contains[T comparable](collection []T, element T) bool {
	for _, item := range collection {
		if item == element {
			return true
		}
	}
	return false
}

func Every[T comparable](collection, subset []T) bool {
	set := make(map[T]struct{}, len(collection))
	for _, item := range collection {
		set[item] = struct{}{}
	}
	for _, item := range subset {
		if _, ok := set[item]; !ok {
			return false
		}
	}
	return true
}

func Some[T comparable](collection, subset []T) bool {
	set := make(map[T]struct{}, len(collection))
	for _, item := range collection {
		set[item] = struct{}{}
	}
	for _, item := range subset {
		if _, ok := set[item]; ok {
			return true
		}
	}
	return false
}

func None[T comparable](collection, subset []T) bool { return !Some(collection, subset) }

func Without[T comparable, Slice ~[]T](collection Slice, exclude ...T) Slice {
	set := make(map[T]struct{}, len(exclude))
	for _, item := range exclude {
		set[item] = struct{}{}
	}
	return Filter(collection, func(item T, _ int) bool { _, found := set[item]; return !found })
}

func Union[T comparable, Slice ~[]T](lists ...Slice) Slice { return Uniq(Concat(lists...)) }

func Intersect[T comparable, Slice ~[]T](lists ...Slice) Slice {
	if len(lists) == 0 {
		return Slice{}
	}
	counts := make(map[T]int)
	for i, list := range lists {
		seen := make(map[T]struct{}, len(list))
		for _, item := range list {
			if _, duplicate := seen[item]; duplicate {
				continue
			}
			seen[item] = struct{}{}
			if i == 0 || counts[item] == i {
				counts[item]++
			}
		}
	}
	out := make(Slice, 0)
	seen := make(map[T]struct{})
	for _, item := range lists[0] {
		if counts[item] == len(lists) {
			if _, duplicate := seen[item]; !duplicate {
				seen[item] = struct{}{}
				out = append(out, item)
			}
		}
	}
	return out
}

func Difference[T comparable, Slice ~[]T](left, right Slice) (Slice, Slice) {
	return Without(left, right...), Without(right, left...)
}

func Keys[K comparable, V any](maps ...map[K]V) []K {
	total := 0
	for _, in := range maps {
		total += len(in)
	}
	out := make([]K, 0, total)
	for _, in := range maps {
		for key := range in {
			out = append(out, key)
		}
	}
	return out
}

func Values[K comparable, V any](maps ...map[K]V) []V {
	total := 0
	for _, in := range maps {
		total += len(in)
	}
	out := make([]V, 0, total)
	for _, in := range maps {
		for _, value := range in {
			out = append(out, value)
		}
	}
	return out
}

func Entries[K comparable, V any](in map[K]V) []Entry[K, V] {
	out := make([]Entry[K, V], 0, len(in))
	for key, value := range in {
		out = append(out, Entry[K, V]{Key: key, Value: value})
	}
	return out
}

func FromEntries[K comparable, V any](entries []Entry[K, V]) map[K]V {
	out := make(map[K]V, len(entries))
	for _, entry := range entries {
		out[entry.Key] = entry.Value
	}
	return out
}

func MapKeys[K comparable, V any, R comparable](in map[K]V, transform func(V, K) R) map[R]V {
	out := make(map[R]V, len(in))
	for key, value := range in {
		out[transform(value, key)] = value
	}
	return out
}

func MapValues[K comparable, V, R any](in map[K]V, transform func(V, K) R) map[K]R {
	out := make(map[K]R, len(in))
	for key, value := range in {
		out[key] = transform(value, key)
	}
	return out
}
