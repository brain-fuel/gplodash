// Package lo — deferred slice helpers (counts, cuts, sorts, fills, folds,
// fallible variants). Adapted from samber/lo v1.53.0 slice.go (MIT), authored
// in Go+ and differentially/property tested. The random Shuffle is excluded.
package lo

import "sort"

// Clonable is the constraint for values that can clone themselves (used by
// Repeat). Adapted from samber/lo v1.53.0 constraints.go.
type Clonable[T any] interface {
	Clone() T
}

// HasPrefix reports whether collection begins with prefix (element-wise).
func HasPrefix[T comparable](collection, prefix []T) bool {
	if len(collection) < len(prefix) {
		return false
	}
	for i := range prefix {
		if collection[i] != prefix[i] {
			return false
		}
	}
	return true
}

// HasSuffix reports whether collection ends with suffix (element-wise).
func HasSuffix[T comparable](collection, suffix []T) bool {
	if len(collection) < len(suffix) {
		return false
	}
	for i := range suffix {
		if collection[len(collection)-len(suffix)+i] != suffix[i] {
			return false
		}
	}
	return true
}

func AssociateI[T any, K comparable, V any](collection []T, transform func(item T, index int) (K, V)) map[K]V {
	result := make(map[K]V, len(collection))

	for index, item := range collection {
		k, v := transform(item, index)
		result[k] = v
	}

	return result
}

func Clone[T any, Slice ~[]T](collection Slice) Slice {
	// backporting from slices.Clone in Go 1.21
	// when we drop support for Go 1.20, this can be replaced with: return slices.Clone(collection)

	// Preserve nilness in case it matters.
	if collection == nil {
		return nil
	}
	// Avoid s[:0:0] as it leads to unwanted liveness when cloning a
	// zero-length slice of a large array; see https://go.dev/issue/68488.
	return append(Slice{}, collection...)
}

func Count[T comparable](collection []T, value T) int {
	var count int

	for i := range collection {
		if collection[i] == value {
			count++
		}
	}

	return count
}

func CountBy[T any](collection []T, predicate func(item T) bool) int {
	var count int

	for i := range collection {
		if predicate(collection[i]) {
			count++
		}
	}

	return count
}

func CountByErr[T any](collection []T, predicate func(item T) (bool, error)) (int, error) {
	var count int

	for i := range collection {
		ok, err := predicate(collection[i])
		if err != nil {
			return 0, err
		}
		if ok {
			count++
		}
	}

	return count, nil
}

func CountValues[T comparable](collection []T) map[T]int {
	result := make(map[T]int)

	for i := range collection {
		result[collection[i]]++
	}

	return result
}

func CountValuesBy[T any, U comparable](collection []T, transform func(item T) U) map[U]int {
	result := make(map[U]int)

	for i := range collection {
		result[transform(collection[i])]++
	}

	return result
}

func Cut[T comparable, Slice ~[]T](collection, separator Slice) (before, after Slice, found bool) {
	if len(separator) == 0 {
		return make(Slice, 0), collection, true
	}

	for i := 0; i+len(separator) <= len(collection); i++ {
		match := true
		for j := 0; j < len(separator); j++ {
			if collection[i+j] != separator[j] {
				match = false
				break
			}
		}
		if match {
			return collection[:i], collection[i+len(separator):], true
		}
	}

	return collection, make(Slice, 0), false
}

func CutPrefix[T comparable, Slice ~[]T](collection, separator Slice) (after Slice, found bool) {
	if HasPrefix(collection, separator) {
		return collection[len(separator):], true
	}
	return collection, false
}

func CutSuffix[T comparable, Slice ~[]T](collection, separator Slice) (before Slice, found bool) {
	if HasSuffix(collection, separator) {
		return collection[:len(collection)-len(separator)], true
	}
	return collection, false
}

func DropByIndex[T any, Slice ~[]T](collection Slice, indexes ...int) Slice {
	initialSize := len(collection)
	if initialSize == 0 {
		return Slice{}
	}

	// do not change the input
	indexes = append(make([]int, 0, len(indexes)), indexes...)

	for i, index := range indexes {
		if index < 0 {
			indexes[i] += initialSize
		}
	}

	sort.Ints(indexes)

	prev := -1
	indexes = Filter(indexes, func(index int, _ int) bool {
		ok := index != prev && // uniq
			uint(index) < uint(initialSize) // in range

		prev = index
		return ok
	})

	result := make(Slice, 0, initialSize-len(indexes))

	i := 0
	for _, index := range indexes {
		result = append(result, collection[i:index]...)
		i = index + 1
	}

	return append(result, collection[i:]...)
}

func DropRightWhile[T any, Slice ~[]T](collection Slice, predicate func(item T) bool) Slice {
	i := len(collection) - 1
	for ; i >= 0; i-- {
		if !predicate(collection[i]) {
			break
		}
	}

	result := make(Slice, 0, i+1)
	return append(result, collection[:i+1]...)
}

func DropWhile[T any, Slice ~[]T](collection Slice, predicate func(item T) bool) Slice {
	i := 0
	for ; i < len(collection); i++ {
		if !predicate(collection[i]) {
			break
		}
	}

	result := make(Slice, 0, len(collection)-i)
	return append(result, collection[i:]...)
}

func Fill[T Clonable[T], Slice ~[]T](collection Slice, initial T) Slice {
	result := make(Slice, 0, len(collection))

	for range collection {
		result = append(result, initial.Clone())
	}

	return result
}

func FilterErr[T any, Slice ~[]T](collection Slice, predicate func(item T, index int) (bool, error)) (Slice, error) {
	result := make(Slice, 0, len(collection))

	for i := range collection {
		ok, err := predicate(collection[i], i)
		if err != nil {
			return nil, err
		}
		if ok {
			result = append(result, collection[i])
		}
	}

	return result, nil
}

func FilterReject[T any, Slice ~[]T](collection Slice, predicate func(T, int) bool) (kept, rejected Slice) {
	kept = make(Slice, 0, len(collection))
	rejected = make(Slice, 0, len(collection))

	for i := range collection {
		if predicate(collection[i], i) {
			kept = append(kept, collection[i])
		} else {
			rejected = append(rejected, collection[i])
		}
	}

	return kept, rejected
}

func FilterSliceToMap[T any, K comparable, V any](collection []T, transform func(item T) (K, V, bool)) map[K]V {
	return FilterSliceToMapI(collection, func(item T, _ int) (K, V, bool) {
		return transform(item)
	})
}

func FilterSliceToMapI[T any, K comparable, V any](collection []T, transform func(item T, index int) (K, V, bool)) map[K]V {
	result := make(map[K]V, len(collection))

	for index, item := range collection {
		k, v, ok := transform(item, index)
		if ok {
			result[k] = v
		}
	}

	return result
}

func FlatMapErr[T, R any](collection []T, transform func(item T, index int) ([]R, error)) ([]R, error) {
	result := make([]R, 0, len(collection))

	for i := range collection {
		r, err := transform(collection[i], i)
		if err != nil {
			return nil, err
		}
		result = append(result, r...)
	}

	return result, nil
}

func ForEachWhile[T any](collection []T, predicate func(item T, index int) bool) {
	for i := range collection {
		if !predicate(collection[i], i) {
			break
		}
	}
}

func GroupByErr[T any, U comparable, Slice ~[]T](collection Slice, iteratee func(item T) (U, error)) (map[U]Slice, error) {
	result := map[U]Slice{}

	for i := range collection {
		key, err := iteratee(collection[i])
		if err != nil {
			return nil, err
		}

		result[key] = append(result[key], collection[i])
	}

	return result, nil
}

func GroupByMapErr[T any, K comparable, V any](collection []T, transform func(item T) (K, V, error)) (map[K][]V, error) {
	result := map[K][]V{}

	for i := range collection {
		k, v, err := transform(collection[i])
		if err != nil {
			return nil, err
		}

		result[k] = append(result[k], v)
	}

	return result, nil
}

func IsSorted[T Ordered](collection []T) bool {
	for i := 1; i < len(collection); i++ {
		if collection[i-1] > collection[i] {
			return false
		}
	}

	return true
}

func IsSortedBy[T any, K Ordered](collection []T, iteratee func(item T) K) bool {
	size := len(collection)

	for i := 0; i < size-1; i++ {
		if iteratee(collection[i]) > iteratee(collection[i+1]) {
			return false
		}
	}

	return true
}

func IsSortedByKey[T any, K Ordered](collection []T, iteratee func(item T) K) bool {
	return IsSortedBy(collection, iteratee)
}

func KeyByErr[K comparable, V any](collection []V, iteratee func(item V) (K, error)) (map[K]V, error) {
	result := make(map[K]V, len(collection))

	for i := range collection {
		k, err := iteratee(collection[i])
		if err != nil {
			return nil, err
		}
		result[k] = collection[i]
	}

	return result, nil
}

func MapErr[T, R any](collection []T, transform func(item T, index int) (R, error)) ([]R, error) {
	result := make([]R, len(collection))

	for i := range collection {
		r, err := transform(collection[i], i)
		if err != nil {
			return nil, err
		}
		result[i] = r
	}

	return result, nil
}

func PartitionByErr[T any, K comparable, Slice ~[]T](collection Slice, iteratee func(item T) (K, error)) ([]Slice, error) {
	result := []Slice{}
	seen := map[K]int{}

	for i := range collection {
		key, err := iteratee(collection[i])
		if err != nil {
			return nil, err
		}

		resultIndex, ok := seen[key]
		if ok {
			result[resultIndex] = append(result[resultIndex], collection[i])
		} else {
			seen[key] = len(result)
			result = append(result, Slice{collection[i]})
		}
	}

	return result, nil
}

func ReduceErr[T, R any](collection []T, accumulator func(agg R, item T, index int) (R, error), initial R) (R, error) {
	for i := range collection {
		result, err := accumulator(initial, collection[i], i)
		if err != nil {
			var zero R
			return zero, err
		}
		initial = result
	}

	return initial, nil
}

func ReduceRightErr[T, R any](collection []T, accumulator func(agg R, item T, index int) (R, error), initial R) (R, error) {
	for i := len(collection) - 1; i >= 0; i-- {
		result, err := accumulator(initial, collection[i], i)
		if err != nil {
			var zero R
			return zero, err
		}
		initial = result
	}

	return initial, nil
}

func RejectErr[T any, Slice ~[]T](collection Slice, predicate func(item T, index int) (bool, error)) (Slice, error) {
	result := Slice{}

	for i := range collection {
		match, err := predicate(collection[i], i)
		if err != nil {
			return nil, err
		}
		if !match {
			result = append(result, collection[i])
		}
	}

	return result, nil
}

func RejectMap[T, R any](collection []T, callback func(item T, index int) (R, bool)) []R {
	result := []R{}

	for i := range collection {
		if r, ok := callback(collection[i], i); !ok {
			result = append(result, r)
		}
	}

	return result
}

func Repeat[T Clonable[T]](count int, initial T) []T {
	result := make([]T, 0, count)

	for i := 0; i < count; i++ {
		result = append(result, initial.Clone())
	}

	return result
}

func RepeatBy[T any](count int, callback func(index int) T) []T {
	result := make([]T, 0, count)

	for i := 0; i < count; i++ {
		result = append(result, callback(i))
	}

	return result
}

func RepeatByErr[T any](count int, callback func(index int) (T, error)) ([]T, error) {
	result := make([]T, 0, count)

	for i := 0; i < count; i++ {
		r, err := callback(i)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}

	return result, nil
}

func Replace[T comparable, Slice ~[]T](collection Slice, old, nEw T, n int) Slice {
	result := make(Slice, len(collection))
	copy(result, collection)

	for i := range result {
		if result[i] == old && n != 0 {
			result[i] = nEw
			n--
		}
	}

	return result
}

func ReplaceAll[T comparable, Slice ~[]T](collection Slice, old, nEw T) Slice {
	return Replace(collection, old, nEw, -1)
}

func Slice[T any, Slice ~[]T](collection Slice, start, end int) Slice {
	if start >= end {
		return Slice{}
	}

	size := len(collection)
	if start < 0 {
		start = 0
	} else if start > size {
		start = size
	}

	if end < 0 {
		end = 0
	} else if end > size {
		end = size
	}

	return collection[start:end]
}

func SliceToMap[T any, K comparable, V any](collection []T, transform func(item T) (K, V)) map[K]V {
	return Associate(collection, transform)
}

func SliceToMapI[T any, K comparable, V any](collection []T, transform func(item T, index int) (K, V)) map[K]V {
	return AssociateI(collection, transform)
}

func Splice[T any, Slice ~[]T](collection Slice, i int, elements ...T) Slice {
	sizeCollection := len(collection)
	sizeElements := len(elements)
	output := make(Slice, 0, sizeCollection+sizeElements) // preallocate memory for the output slice

	switch {
	case sizeElements == 0:
		return append(output, collection...) // simple copy
	case i > sizeCollection:
		// positive overflow
		return append(append(output, collection...), elements...)
	case i < -sizeCollection:
		// negative overflow
		return append(append(output, elements...), collection...)
	case i < 0:
		// backward
		i = sizeCollection + i
	}

	return append(append(append(output, collection[:i]...), elements...), collection[i:]...)
}

func Subset[T any, Slice ~[]T](collection Slice, offset int, length uint) Slice {
	size := len(collection)

	if offset < 0 {
		offset = size + offset
		if offset < 0 {
			offset = 0
		}
	}

	if offset > size {
		return Slice{}
	}

	if length > uint(size)-uint(offset) {
		length = uint(size - offset)
	}

	return collection[offset : offset+int(length)]
}

func TakeFilter[T any, Slice ~[]T](collection Slice, n int, predicate func(item T, index int) bool) Slice {
	if n < 0 {
		panic("lo.TakeFilter: n must not be negative")
	}

	if n == 0 {
		return make(Slice, 0)
	}

	result := make(Slice, 0, n)
	count := 0

	for i := range collection {
		if predicate(collection[i], i) {
			result = append(result, collection[i])
			count++
			if count >= n {
				break
			}
		}
	}

	return result
}

func TakeWhile[T any, Slice ~[]T](collection Slice, predicate func(item T) bool) Slice {
	i := 0
	for ; i < len(collection); i++ {
		if !predicate(collection[i]) {
			break
		}
	}

	result := make(Slice, i)
	copy(result, collection[:i])
	return result
}

func Trim[T comparable, Slice ~[]T](collection, cutset Slice) Slice {
	set := Keyify(cutset)

	i := 0
	for ; i < len(collection); i++ {
		if _, ok := set[collection[i]]; !ok {
			break
		}
	}

	if i >= len(collection) {
		return Slice{}
	}

	j := len(collection) - 1
	for ; j >= 0; j-- {
		if _, ok := set[collection[j]]; !ok {
			break
		}
	}

	result := make(Slice, 0, j+1-i)
	return append(result, collection[i:j+1]...)
}

func TrimLeft[T comparable, Slice ~[]T](collection, cutset Slice) Slice {
	set := Keyify(cutset)

	return DropWhile(collection, func(item T) bool {
		_, ok := set[item]
		return ok
	})
}

func TrimPrefix[T comparable, Slice ~[]T](collection, prefix Slice) Slice {
	if len(prefix) == 0 {
		return collection
	}

	for HasPrefix(collection, prefix) {
		collection = collection[len(prefix):]
	}

	return collection
}

func TrimRight[T comparable, Slice ~[]T](collection, cutset Slice) Slice {
	set := Keyify(cutset)

	return DropRightWhile(collection, func(item T) bool {
		_, ok := set[item]
		return ok
	})
}

func TrimSuffix[T comparable, Slice ~[]T](collection, suffix Slice) Slice {
	if len(suffix) == 0 {
		return collection
	}

	for HasSuffix(collection, suffix) {
		collection = collection[:len(collection)-len(suffix)]
	}

	return collection
}

func UniqByErr[T any, U comparable, Slice ~[]T](collection Slice, iteratee func(item T) (U, error)) (Slice, error) {
	result := make(Slice, 0, len(collection))
	seen := make(map[U]struct{}, len(collection))

	for i := range collection {
		key, err := iteratee(collection[i])
		if err != nil {
			return nil, err
		}

		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}
		result = append(result, collection[i])
	}

	return result, nil
}

func UniqMap[T any, R comparable](collection []T, transform func(item T, index int) R) []R {
	seen := make(map[R]struct{}, len(collection))

	for i := range collection {
		r := transform(collection[i], i)
		if _, ok := seen[r]; !ok {
			seen[r] = struct{}{}
		}
	}

	return Keys(seen)
}
