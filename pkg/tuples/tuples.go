package tuples

import (
	"adventofcode2021/pkg/convert"
)

type Pair[T, U convert.Convertable] struct {
	Key   T
	Value U
}
