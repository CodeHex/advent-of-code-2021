package fileparser

import (
	"os"
	"strconv"
	"strings"
)

// Readlines reads a file contents into a slice of strings seperated by newline characters
// Panics if there are any errors with reading the file
func ReadLines(filename string) []string {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	dataString := string(data)
	dataString = strings.TrimSpace(dataString)
	return strings.Split(dataString, "\n")
}

// ReadInts reads a files contents into a slice if ints, seperated by newline characters
// Panics if there are any errors reading the file or parsing the numbers
func ReadInts(filename string) []int {
	strParts := ReadLines(filename)
	intParts := make([]int, len(strParts))
	var err error
	for i, strPart := range strParts {
		intParts[i], err = strconv.Atoi(strPart)
		if err != nil {
			panic(err)
		}
	}
	return intParts
}
