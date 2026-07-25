// Package lo — positional access helpers (fallback and empty variants).
//
// Exact-compatibility reimplementations of samber/lo v1.53.0 find.go's NthOr
// family, authored in Go+. sliceNth preserves upstream's symmetric negative
// indexing (index -1 is the last element) and out-of-range predicate.
package lo

// sliceNth returns the element at nth, supporting negative indices, and reports
// whether the index was in range. It is the shared helper behind NthOr and
// NthOrEmpty.
func sliceNth[T any, N Integer](collection []T, nth N) (T, bool) {
	n := int(nth)
	l := len(collection)
	if n >= l || -n > l {
		var zero T
		return zero, false
	}
	if n >= 0 {
		return collection[n], true
	}
	return collection[l+n], true
}

// NthOr returns the element at nth, or fallback when nth is out of range.
func NthOr[T any, N Integer](collection []T, nth N, fallback T) T {
	if value, ok := sliceNth(collection, nth); ok {
		return value
	}
	return fallback
}

// NthOrEmpty returns the element at nth, or the zero value when out of range.
func NthOrEmpty[T any, N Integer](collection []T, nth N) T {
	value, _ := sliceNth(collection, nth)
	return value
}

// Min returns the least element, or the zero value for an empty collection.
func Min[T Ordered](collection []T) T {
	var min T
	if len(collection) == 0 {
		return min
	}
	min = collection[0]
	for i := 1; i < len(collection); i++ {
		if collection[i] < min {
			min = collection[i]
		}
	}
	return min
}

// MinIndex returns the least element and its index, or (zero, -1) when empty.
func MinIndex[T Ordered](collection []T) (T, int) {
	var min T
	index := 0
	if len(collection) == 0 {
		return min, -1
	}
	min = collection[0]
	for i := 1; i < len(collection); i++ {
		if collection[i] < min {
			min = collection[i]
			index = i
		}
	}
	return min, index
}

// MinBy returns the least element under the given less relation.
func MinBy[T any](collection []T, less func(a, b T) bool) T {
	var min T
	if len(collection) == 0 {
		return min
	}
	min = collection[0]
	for i := 1; i < len(collection); i++ {
		if less(collection[i], min) {
			min = collection[i]
		}
	}
	return min
}

// MinIndexBy returns the least element and its index under less.
func MinIndexBy[T any](collection []T, less func(a, b T) bool) (T, int) {
	var min T
	index := 0
	if len(collection) == 0 {
		return min, -1
	}
	min = collection[0]
	for i := 1; i < len(collection); i++ {
		if less(collection[i], min) {
			min = collection[i]
			index = i
		}
	}
	return min, index
}

// Max returns the greatest element, or the zero value for an empty collection.
func Max[T Ordered](collection []T) T {
	var max T
	if len(collection) == 0 {
		return max
	}
	max = collection[0]
	for i := 1; i < len(collection); i++ {
		if collection[i] > max {
			max = collection[i]
		}
	}
	return max
}

// MaxIndex returns the greatest element and its index, or (zero, -1) when empty.
func MaxIndex[T Ordered](collection []T) (T, int) {
	var max T
	index := 0
	if len(collection) == 0 {
		return max, -1
	}
	max = collection[0]
	for i := 1; i < len(collection); i++ {
		if collection[i] > max {
			max = collection[i]
			index = i
		}
	}
	return max, index
}

// MaxBy returns the greatest element under the given greater relation.
func MaxBy[T any](collection []T, greater func(a, b T) bool) T {
	var max T
	if len(collection) == 0 {
		return max
	}
	max = collection[0]
	for i := 1; i < len(collection); i++ {
		if greater(collection[i], max) {
			max = collection[i]
		}
	}
	return max
}

// MaxIndexBy returns the greatest element and its index under greater.
func MaxIndexBy[T any](collection []T, greater func(a, b T) bool) (T, int) {
	var max T
	index := 0
	if len(collection) == 0 {
		return max, -1
	}
	max = collection[0]
	for i := 1; i < len(collection); i++ {
		if greater(collection[i], max) {
			max = collection[i]
			index = i
		}
	}
	return max, index
}
