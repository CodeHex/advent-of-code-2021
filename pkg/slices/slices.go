package slices

import "constraints"

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
