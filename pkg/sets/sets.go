package sets

type Set[T comparable] map[T]struct{}

// NewEmptySet generates an empty set
func NewEmptySet[T comparable]() Set[T] {
	return make(Set[T])
}

// Add will add an element to a set
func (s Set[T]) Add(entry T) {
	s[entry] = struct{}{}
}

// Filter will generate a new set containing elements that match the predicate
func (s Set[T]) Filter(predicate func(val T) bool) Set[T] {
	result := NewEmptySet[T]()
	for k := range s {
		if predicate(k) {
			result.Add(k)
		}
	}
	return result
}

// ToSlice will generate a slice will all the set elements (undefined order) for iteration
func (s Set[T]) ToSlice() []T {
	result := make([]T, len(s))
	i := 0
	for k := range s {
		result[i] = k
		i++
	}
	return result
}
