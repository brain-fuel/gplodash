// Package lo — deferred find.go helpers (first/last, index-of, duplicates,
// uniques, earliest/latest, fallible min/max variants). Adapted from samber/lo
// v1.53.0 find.go (MIT), authored in Go+ and differentially tested. Random
// Sample/Samples are excluded.
package lo

import "time"

func Earliest(times ...time.Time) time.Time {
	var mIn time.Time

	if len(times) == 0 {
		return mIn
	}

	mIn = times[0]

	for i := 1; i < len(times); i++ {
		item := times[i]

		if item.Before(mIn) {
			mIn = item
		}
	}

	return mIn
}

func EarliestBy[T any](collection []T, iteratee func(item T) time.Time) T {
	var earliest T

	if len(collection) == 0 {
		return earliest
	}

	earliest = collection[0]
	earliestTime := iteratee(collection[0])

	for i := 1; i < len(collection); i++ {
		itemTime := iteratee(collection[i])

		if itemTime.Before(earliestTime) {
			earliest = collection[i]
			earliestTime = itemTime
		}
	}

	return earliest
}

func EarliestByErr[T any](collection []T, iteratee func(item T) (time.Time, error)) (T, error) {
	var earliest T

	if len(collection) == 0 {
		return earliest, nil
	}

	earliestTime, err := iteratee(collection[0])
	if err != nil {
		return earliest, err
	}
	earliest = collection[0]

	for i := 1; i < len(collection); i++ {
		itemTime, err := iteratee(collection[i])
		if err != nil {
			return earliest, err
		}

		if itemTime.Before(earliestTime) {
			earliest = collection[i]
			earliestTime = itemTime
		}
	}

	return earliest, nil
}

func FindDuplicates[T comparable, Slice ~[]T](collection Slice) Slice {
	isDupl := make(map[T]bool, len(collection))

	duplicates := 0

	for i := range collection {
		duplicated, seen := isDupl[collection[i]]
		if !duplicated {
			isDupl[collection[i]] = seen

			if seen {
				duplicates++
			}
		}
	}

	result := make(Slice, 0, duplicates)

	for i := range collection {
		if duplicated := isDupl[collection[i]]; duplicated {
			result = append(result, collection[i])
			isDupl[collection[i]] = false
		}
	}

	return result
}

func FindDuplicatesBy[T any, U comparable, Slice ~[]T](collection Slice, iteratee func(item T) U) Slice {
	isDupl := make(map[U]bool, len(collection))

	duplicates := 0

	for i := range collection {
		key := iteratee(collection[i])

		duplicated, seen := isDupl[key]
		if !duplicated {
			isDupl[key] = seen

			if seen {
				duplicates++
			}
		}
	}

	result := make(Slice, 0, duplicates)

	for i := range collection {
		key := iteratee(collection[i])

		if duplicated := isDupl[key]; duplicated {
			result = append(result, collection[i])
			isDupl[key] = false
		}
	}

	return result
}

func FindDuplicatesByErr[T any, U comparable, Slice ~[]T](collection Slice, iteratee func(item T) (U, error)) (Slice, error) {
	isDupl := make(map[U]bool, len(collection))

	duplicates := 0

	// First pass: identify duplicates
	for i := range collection {
		key, err := iteratee(collection[i])
		if err != nil {
			var result Slice
			return result, err
		}

		duplicated, seen := isDupl[key]
		if !duplicated {
			isDupl[key] = seen

			if seen {
				duplicates++
			}
		}
	}

	result := make(Slice, 0, duplicates)

	// Second pass: collect first occurrences of duplicates
	for i := range collection {
		key, err := iteratee(collection[i])
		if err != nil {
			var result Slice
			return result, err
		}

		if duplicated := isDupl[key]; duplicated {
			result = append(result, collection[i])
			isDupl[key] = false
		}
	}

	return result, nil
}

func FindErr[T any](collection []T, predicate func(item T) (bool, error)) (T, error) {
	for i := range collection {
		matches, err := predicate(collection[i])
		if err != nil {
			var result T
			return result, err
		}
		if matches {
			return collection[i], nil
		}
	}

	var result T
	return result, nil
}

func FindKey[K, V comparable](object map[K]V, value V) (K, bool) {
	for k, v := range object {
		if v == value {
			return k, true
		}
	}

	return Empty[K](), false
}

func FindKeyBy[K comparable, V any](object map[K]V, predicate func(key K, value V) bool) (K, bool) {
	for k, v := range object {
		if predicate(k, v) {
			return k, true
		}
	}

	return Empty[K](), false
}

func FindOrElse[T any](collection []T, fallback T, predicate func(item T) bool) T {
	for i := range collection {
		if predicate(collection[i]) {
			return collection[i]
		}
	}

	return fallback
}

func FindUniques[T comparable, Slice ~[]T](collection Slice) Slice {
	isDupl := make(map[T]bool, len(collection))

	duplicates := 0

	for i := range collection {
		duplicated, seen := isDupl[collection[i]]
		if !duplicated {
			isDupl[collection[i]] = seen

			if seen {
				duplicates++
			}
		}
	}

	result := make(Slice, 0, len(isDupl)-duplicates)

	for i := range collection {
		if duplicated := isDupl[collection[i]]; !duplicated {
			result = append(result, collection[i])
		}
	}

	return result
}

func FindUniquesBy[T any, U comparable, Slice ~[]T](collection Slice, iteratee func(item T) U) Slice {
	isDupl := make(map[U]bool, len(collection))

	duplicates := 0

	for i := range collection {
		key := iteratee(collection[i])

		duplicated, seen := isDupl[key]
		if !duplicated {
			isDupl[key] = seen

			if seen {
				duplicates++
			}
		}
	}

	result := make(Slice, 0, len(isDupl)-duplicates)

	for i := range collection {
		key := iteratee(collection[i])

		if duplicated := isDupl[key]; !duplicated {
			result = append(result, collection[i])
		}
	}

	return result
}

func FirstOr[T any](collection []T, fallback T) T {
	i, ok := First(collection)
	if !ok {
		return fallback
	}

	return i
}

func FirstOrEmpty[T any](collection []T) T {
	i, _ := First(collection)
	return i
}

func IndexOf[T comparable](collection []T, element T) int {
	for i := range collection {
		if collection[i] == element {
			return i
		}
	}

	return -1
}

func LastIndexOf[T comparable](collection []T, element T) int {
	length := len(collection)

	for i := length - 1; i >= 0; i-- {
		if collection[i] == element {
			return i
		}
	}

	return -1
}

func LastOr[T any](collection []T, fallback T) T {
	i, ok := Last(collection)
	if !ok {
		return fallback
	}

	return i
}

func LastOrEmpty[T any](collection []T) T {
	i, _ := Last(collection)
	return i
}

func Latest(times ...time.Time) time.Time {
	var mAx time.Time

	if len(times) == 0 {
		return mAx
	}

	mAx = times[0]

	for i := 1; i < len(times); i++ {
		item := times[i]

		if item.After(mAx) {
			mAx = item
		}
	}

	return mAx
}

func LatestBy[T any](collection []T, iteratee func(item T) time.Time) T {
	var latest T

	if len(collection) == 0 {
		return latest
	}

	latest = collection[0]
	latestTime := iteratee(collection[0])

	for i := 1; i < len(collection); i++ {
		itemTime := iteratee(collection[i])

		if itemTime.After(latestTime) {
			latest = collection[i]
			latestTime = itemTime
		}
	}

	return latest
}

func LatestByErr[T any](collection []T, iteratee func(item T) (time.Time, error)) (T, error) {
	var latest T

	if len(collection) == 0 {
		return latest, nil
	}

	latestTime, err := iteratee(collection[0])
	if err != nil {
		return latest, err
	}
	latest = collection[0]

	for i := 1; i < len(collection); i++ {
		itemTime, err := iteratee(collection[i])
		if err != nil {
			return latest, err
		}

		if itemTime.After(latestTime) {
			latest = collection[i]
			latestTime = itemTime
		}
	}

	return latest, nil
}

func MaxByErr[T any](collection []T, greater func(a, b T) (bool, error)) (T, error) {
	var mAx T

	if len(collection) == 0 {
		return mAx, nil
	}

	mAx = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		isGreater, err := greater(item, mAx)
		if err != nil {
			return mAx, err
		}
		if isGreater {
			mAx = item
		}
	}

	return mAx, nil
}

func MaxIndexByErr[T any](collection []T, greater func(a, b T) (bool, error)) (T, int, error) {
	var (
		mAx   T
		index int
	)

	if len(collection) == 0 {
		return mAx, -1, nil
	}

	mAx = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		isGreater, err := greater(item, mAx)
		if err != nil {
			var zero T
			return zero, -1, err
		}
		if isGreater {
			mAx = item
			index = i
		}
	}

	return mAx, index, nil
}

func MinByErr[T any](collection []T, less func(a, b T) (bool, error)) (T, error) {
	var mIn T

	if len(collection) == 0 {
		return mIn, nil
	}

	mIn = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		isLess, err := less(item, mIn)
		if err != nil {
			var zero T
			return zero, err
		}
		if isLess {
			mIn = item
		}
	}

	return mIn, nil
}

func MinIndexByErr[T any](collection []T, less func(a, b T) (bool, error)) (T, int, error) {
	var (
		mIn   T
		index int
	)

	if len(collection) == 0 {
		return mIn, -1, nil
	}

	mIn = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		isLess, err := less(item, mIn)
		if err != nil {
			var zero T
			return zero, -1, err
		}

		if isLess {
			mIn = item
			index = i
		}
	}

	return mIn, index, nil
}
