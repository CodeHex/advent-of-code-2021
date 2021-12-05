package fileparser

import (
	"adventofcode2021/pkg/convert"
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

func ReadPairs[T, U convert.Convertable](filename string) []tuples.Pair[T, U] {
	strParts := ReadSingles[string](filename)
	result := make([]tuples.Pair[T, U], len(strParts))

	convertKey := convert.FuncFor[T]()
	convertValue := convert.FuncFor[U]()

	for i, part := range strParts {
		vals := strings.Split(part, " ")
		if len(vals) != 2 {
			panic(fmt.Sprintf("expecting 2 parts, '%s'", part))
		}
		result[i] = tuples.Pair[T, U]{
			Key:   convertKey(vals[0]),
			Value: convertValue(vals[1]),
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
