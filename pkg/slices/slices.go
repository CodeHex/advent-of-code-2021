package slices

// Filter will reduce a slices of elements based on the provided predicate
func Filter[T any](source []T, predicate func(T) bool) []T {
	result := []T{}
	for _, entry := range source {
		if predicate(entry) {
			result = append(result, entry)
		}
	}
	return result
}

// IsSingle will check if a slice only has one element and only return that element
func IsSingle[T any](source []T) (T, bool) {
	if len(source) == 1 {
		return source[0], true
	}
	return *new(T), false
}
