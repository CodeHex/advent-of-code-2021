package slices

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
