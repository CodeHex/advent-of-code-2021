package matrices

type Matrix[T any] struct {
	data          [][]T
	Rows, Columns int
	Size          int
}

// NewMatrix creates a default matrix with provided dimensions
func NewMatrix[T any](rows, columns int) Matrix[T] {
	m := make([][]T, rows)
	for y := range m {
		m[y] = make([]T, columns)
	}
	return NewMatrixFromData(m)
}

// NewMatrixFromLines creates a matrix where each line represents a row of the matrix.
// The splitter convertes the line into component entries and the convert converts
// the raw part into the required type
func NewMatrixFromData[T any](data [][]T) Matrix[T] {
	rows := len(data)
	columns := len(data[0])

	for _, row := range data {
		if len(row) != columns {
			panic("unable to create matrix, mismatching row lengths")
		}
	}
	return Matrix[T]{
		data:    data,
		Rows:    rows,
		Columns: columns,
		Size:    rows * columns,
	}
}

// ForEach performs the operation on every element in the matrix,
// referencing the location and value of the element
func (m *Matrix[T]) ForEach(op func(x, y int, value T)) {
	for j := 0; j < m.Rows; j++ {
		for i := 0; i < m.Columns; i++ {
			op(i, j, m.data[j][i])
		}
	}
}

// Get will return the provided element of the matrix
func (m *Matrix[T]) Get(x, y int) T {
	return m.data[y][x]
}

// Set will set an element of the matrix
func (m Matrix[T]) Set(x, y int, val T) {
	m.data[y][x] = val
}

// OutOfBounds indicates if the provided location exists in the matrix
func (m Matrix[T]) OutOfBounds(x, y int) bool {
	return x < 0 || x > m.Columns-1 || y < 0 || y > m.Rows-1
}

// ForEachNeighbour performs the operation on itself and its closest neighbours
func (m Matrix[T]) ForEachNeighbour(includeDiags bool, originX, originY int, op func(x, y int)) {
	opIfNotOutOfBounds := func(i, j int) {
		if !m.OutOfBounds(i, j) {
			op(i, j)
		}
	}
	opIfNotOutOfBounds(originX, originY+1)
	opIfNotOutOfBounds(originX, originY-1)
	opIfNotOutOfBounds(originX-1, originY)
	opIfNotOutOfBounds(originX+1, originY)

	if includeDiags {
		opIfNotOutOfBounds(originX+1, originY+1)
		opIfNotOutOfBounds(originX+1, originY-1)
		opIfNotOutOfBounds(originX-1, originY+1)
		opIfNotOutOfBounds(originX-1, originY-1)
	}
}

// Expand will increase the size the matrix. The value of these new entries are provided
func (m Matrix[T]) Expand(sizeLeft, sizeRight, sizeTop, sizeBottom int, val T) Matrix[T] {
	newMatrix := NewMatrix[T](m.Rows+(sizeTop+sizeBottom), m.Columns+(sizeLeft+sizeRight))
	for j := 0; j < newMatrix.Rows; j++ {
		for i := 0; i < newMatrix.Columns; i++ {
			if i < sizeLeft || i > newMatrix.Columns-sizeRight-1 || j < sizeTop || j > newMatrix.Rows-sizeBottom-1 {
				newMatrix.Set(i, j, val)
			} else {
				newMatrix.Set(i, j, m.Get(i-sizeLeft, j-sizeTop))
			}
		}
	}
	return newMatrix
}
