package fileparser

import (
	"adventofcode2021/pkg/convert"
	"adventofcode2021/pkg/matrices"
	"adventofcode2021/pkg/slices"
	"adventofcode2021/pkg/tuples"
	"fmt"
	"os"
	"strings"
)

func ReadSingles[T convert.Convertable](filename string) []T {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	dataString := strings.TrimSpace(string(data))
	dataParts := strings.Split(dataString, "\n")

	resultParts := make([]T, len(dataParts))
	converter := convert.FuncFor[T]()
	for i, part := range dataParts {
		resultParts[i] = converter(part)
	}
	return resultParts
}

func ReadPairs[T, U convert.Convertable](filename string, separator string) []tuples.Pair[T, U] {
	strParts := ReadSingles[string](filename)
	return ReadPairsFromStrings[T, U](strParts, separator)
}

func ReadPairsFromStrings[T, U convert.Convertable](data []string, separator string) []tuples.Pair[T, U] {
	result := make([]tuples.Pair[T, U], len(data))

	convertKey := convert.FuncFor[T]()
	convertValue := convert.FuncFor[U]()

	for i, part := range data {
		vals := strings.Split(part, separator)
		if len(vals) != 2 {
			panic(fmt.Sprintf("expecting 2 parts, '%s'", part))
		}
		result[i] = tuples.Pair[T, U]{
			Key:   convertKey(strings.TrimSpace(vals[0])),
			Value: convertValue(strings.TrimSpace(vals[1])),
		}
	}
	return result
}

func ReadLines(filename string) []string {
	return ReadSingles[string](filename)
}

func ReadTypedLines[T any](filename string, constructor func(string) T) []T {
	lines := ReadLines(filename)
	result := make([]T, len(lines))
	for i, data := range lines {
		result[i] = constructor(data)
	}
	return result
}

func ReadCSVLine[T convert.Convertable](filename string) []T {
	lines := ReadLines(filename)
	return Split[T](lines[0], ",")
}

func ReadCharMatrix[T convert.Convertable](filename string) matrices.Matrix[T] {
	return ReadCharMatrixFromLines[T](ReadLines(filename))
}

func ReadCharMatrixFromLines[T convert.Convertable](lines []string) matrices.Matrix[T] {
	m := make([][]T, len(lines))
	converter := convert.FuncFor[T]()
	for y, line := range lines {
		m[y] = slices.Map([]rune(line), func(x rune) T {
			return converter(string(x))
		})
	}
	return matrices.NewMatrixFromData(m)
}

func ReadDigitMatrix(filename string) matrices.IntMatrix[int] {
	return matrices.NewIntMatrixFromBase(ReadCharMatrix[int](filename))
}

// Split will split a string similar to strings.Split, but convert the result to the appriopriate type
func Split[T convert.Convertable](str string, sep string) []T {
	parts := strings.Split(str, sep)
	result := make([]T, len(parts))
	converter := convert.FuncFor[T]()
	for i, part := range parts {
		result[i] = converter(part)
	}
	return result
}

// SplitTrim will split a string similar to Split, but ignore any empty results and trim data
func SplitTrim[T convert.Convertable](str string, sep string) []T {
	parts := strings.Split(str, sep)
	result := []T{}
	converter := convert.FuncFor[T]()
	for _, part := range parts {
		if part != "" {
			result = append(result, converter(strings.TrimSpace(part)))
		}
	}
	return result
}
