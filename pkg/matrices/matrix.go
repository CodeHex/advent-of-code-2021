package matrices

import "adventofcode2021/pkg/slices"

type Matrix[T any] [][]T

// NewMatrix creates a default matrix with provided dimensions
func NewMatrix[T any](rows, columns int) Matrix[T] {
	m := make([][]T, rows)
	for y := range m {
		m[y] = make([]T, columns)
	}
	return m
}

// NewMatrixFromLines creates a matrix where each line represents a row of the matrix.
// The splitter convertes the line into component entries and the convert converts
// the raw part into the required type
func NewMatrixFromLines[T, U any](lines []string, splitter func(string) []U, converter func(U) T) Matrix[T] {
	m := make([][]T, len(lines))
	for y, line := range lines {
		parts := splitter(line)
		m[y] = slices.Map(parts, converter)
	}
	return m
}

// Dimensions returns the rows and columns of the matrix
func (m Matrix[T]) Dimensions() (rows, columns int) {
	return len(m), len(m[0])
}

// ForEach performs the operation on every element in the matrix,
// referencing the location and value of the element
func (m Matrix[T]) ForEach(op func(x, y int, value T)) {
	rows, columns := m.Dimensions()
	for j := 0; j < rows; j++ {
		for i := 0; i < columns; i++ {
			op(i, j, m[j][i])
		}
	}
}

// Get will return the provided element of the matrix
func (m Matrix[T]) Get(x, y int) T {
	return m[y][x]
}

// Set will set an element of the matrix
func (m Matrix[T]) Set(x, y int, val T) {
	m[y][x] = val
}
