package matrices

import "constraints"

type IntMatrix[T constraints.Integer] Matrix[T]

func NewIntMatrixFromLines[T constraints.Integer, U any](lines []string, splitter func(string) []U, converter func(U) T) IntMatrix[T] {
	return IntMatrix[T](NewMatrixFromLines(lines, splitter, converter))
}

func (m IntMatrix[T]) ForEach(op func(x, y int, value T)) {
	Matrix[T](m).ForEach(op)
}

func (m IntMatrix[T]) Get(x, y int) T {
	return Matrix[T](m).Get(x, y)
}

func (m IntMatrix[T]) Set(x, y int, val T) {
	Matrix[T](m).Set(x, y, val)
}

func (m IntMatrix[T]) Dimensions() (rows, columns int) {
	return Matrix[T](m).Dimensions()
}

func (m IntMatrix[T]) NumOfElements() int {
	return Matrix[T](m).NumOfElements()
}

func (m IntMatrix[T]) OutOfBounds(x, y int) bool {
	return Matrix[T](m).OutOfBounds(x, y)
}

func (m IntMatrix[T]) ForEachNeighbour(includeDiags bool, x, y int, op func(x1, y1 int)) {
	Matrix[T](m).ForEachNeighbour(includeDiags, x, y, op)
}

// Increment will increment the value at the location provided
func (m IntMatrix[T]) Increment(x, y int) {
	m[y][x]++
}

// Increment all values in the matrix
func (m IntMatrix[T]) IncrementAll() {
	m.ForEach(func(x, y int, value T) {
		m.Increment(x, y)
	})
}
