package maps

import "constraints"

// SumValues will return the total of all the maps values
func SumValues[T comparable, U constraints.Integer](source map[T]U) U {
	var total U
	for _, val := range source {
		total += val
	}
	return total
}

// SumValuesFor will return the total of all the maps values where the keys match the predicate
func SumValuesFor[T comparable, U constraints.Integer](source map[T]U, predicate func(T) bool) U {
	var total U
	for key, val := range source {
		if predicate(key) {
			total += val
		}
	}
	return total
}

// MinValue will return the entry with the smallest value
func MinValue[T constraints.Integer, U constraints.Integer](source map[T]U) (T, U) {
	started := false
	var minKey T
	var minVal U
	for key, val := range source {
		if !started {
			minKey = key
			minVal = val
			started = true
		}
		if val < minVal {
			minKey = key
			minVal = val
		}
	}
	return minKey, minVal
}
