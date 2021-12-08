package slices

import (
	"constraints"
)

// Filter will reduce a slice of elements based on the provided predicate
func Filter[T any](source []T, predicate func(T) bool) []T {
	result := []T{}
	for _, entry := range source {
		if predicate(entry) {
			result = append(result, entry)
		}
	}
	return result
}

// First will select the first matching element based on the provided predicate
func First[T any](source []T, predicate func(T) bool) (T, bool) {
	for _, entry := range source {
		if predicate(entry) {
			return entry, true
		}
	}
	var blank T
	return blank, false
}

// FirstOrDefault will select the first matching element based on the provided predicate, or the default value
// if one is not found
func FirstOrDefault[T any](source []T, predicate func(T) bool) T {
	result, _ := First(source, predicate)
	return result
}

// Contains returns if the provided element is in the slice
func Contains[T comparable](source []T, element T) bool {
	for _, entry := range source {
		if entry == element {
			return true
		}
	}
	return false
}

// Divide will split a slice of elements into 2 slices, with the first slice elements matching
// the predicate, the second slice elements do not
func Divide[T any](source []T, predicate func(T) bool) ([]T, []T) {
	resultMatch := []T{}
	resultNotMatch := []T{}
	for _, entry := range source {
		if predicate(entry) {
			resultMatch = append(resultMatch, entry)
		} else {
			resultNotMatch = append(resultNotMatch, entry)
		}
	}
	return resultMatch, resultNotMatch
}

// IsSingle will check if a slice only has one element and only return that element
func IsSingle[T any](source []T) (T, bool) {
	if len(source) == 1 {
		return source[0], true
	}
	return *new(T), false
}

// InitGrid will initialse a 2D set of slices with default values
func InitGrid[T any](numX, numY int) [][]T {
	grid := make([][]T, numY)
	for i := range grid {
		grid[i] = make([]T, numX)
	}
	return grid
}

// Map converted one slice to another slice
func Map[T, U any](source []U, selector func(U) T) []T {
	result := make([]T, len(source))
	for i, val := range source {
		result[i] = selector(val)
	}
	return result
}

// Max determines the maximum value in a slice of ordered values
func Max[T constraints.Ordered](source []T) T {
	var maxVal T
	for _, val := range source {
		if val > maxVal {
			maxVal = val
		}
	}
	return maxVal
}

// Min determines the minimum value in a slice of ordered values
func Min[T constraints.Ordered](source []T) T {
	var minVal T
	if len(source) != 0 {
		minVal = source[0]
	}
	for _, val := range source {
		if val < minVal {
			minVal = val
		}
	}
	return minVal
}

// MinMax determines the minimum and maximum value in a slice of ordered values efficiently
func MinMax[T constraints.Ordered](source []T) (T, T) {
	var minVal T
	var maxVal T
	if len(source) != 0 {
		minVal = source[0]
	}
	for _, val := range source {
		if val < minVal {
			minVal = val
		}
		if val > maxVal {
			maxVal = val
		}
	}
	return minVal, maxVal
}

// Sum will sum all values in the slice
func Sum[T constraints.Integer](source []T) T {
	return SumWeighted(source, func(x T) T { return x })
}

// SumWeighted will sum all values in the slice using the provided weighting function
func SumWeighted[T any, U constraints.Integer](source []T, weightFunc func(T) U) U {
	var total U
	for _, entry := range source {
		total += weightFunc(entry)
	}
	return total
}

// CountIf counts the number of elements that match the predicate
func CountIf[T any](source []T, predicate func(x T) bool) int {
	count := 0
	for _, entry := range source {
		if predicate(entry) {
			count++
		}
	}
	return count
}
