package matrices

import (
	"constraints"
	"fmt"
)

type IntMatrix[T constraints.Integer] struct {
	Matrix[T]
}

func NewIntMatrixFromData[T constraints.Integer](data [][]T) IntMatrix[T] {
	return IntMatrix[T]{Matrix: NewMatrixFromData(data)}
}

// Increment will increment the value at the location provided
func (m IntMatrix[T]) Increment(x, y int) {
	m.data[y][x]++
}

// Increment all values in the matrix
func (m IntMatrix[T]) IncrementAll() {
	m.ForEach(func(x, y int, value T) {
		m.Increment(x, y)
	})
}

// CompactString produces a compact representation of the matrix, assuming single digit entries
func (m IntMatrix[T]) CompactString() string {
	out := ""
	for _, line := range m.data {
		for _, val := range line {
			out += fmt.Sprintf("%d", val)
		}
		out += "\n"
	}
	return out
}
