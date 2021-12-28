package slices

import (
	"constraints"
	"sort"
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

// ContainsAny returns if any of the elements provided in the slice exist in the source
func ContainsAny[T comparable](source []T, elements []T) bool {
	for _, ele := range elements {
		if Contains(source, ele) {
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

// Reverse will reverse the order of the slice
func Reverse[T any](source []T) []T {
	reversed := make([]T, len(source))
	for i, val := range source {
		reversed[len(source)-1-i] = val
	}
	return reversed
}

// Last will return the final element in the slice
func Last[T any](source []T) T {
	size := len(source)
	if size == 0 {
		var blank T
		return blank
	}
	return source[size-1]
}

// TrimEnd will trim the final n elements from the end of the slice
func TrimEnd[T any](source []T, n int) []T {
	if n > len(source) {
		return make([]T, 0)
	}
	return source[0 : len(source)-n]
}

// Median will return the median value of odd length slices
func Median[T constraints.Integer](source []T) T {
	n := len(source)
	if n%2 == 0 {
		panic("median of even slices is not supported")
	}
	result := make([]T, n)
	copy(result, source)
	sort.SliceStable(result, func(i, j int) bool {
		return result[i] < result[j]
	})
	return result[n/2]
}

// IndexOf returns the index where the first occurence of val is, otherwise -1 if not found
func IndexOf[T comparable](source []T, val T) int {
	for i, v := range source {
		if v == val {
			return i
		}
	}
	return -1
}
