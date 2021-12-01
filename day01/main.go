package main

import (
	"adventofcode2021/tools/fileparser"
	"constraints"
	"fmt"
)

func main() {
	measurements := fileparser.ReadInts("day01/input.txt")
	calculateDiffCounts(measurements, 1)
	calculateDiffCounts(measurements, 3)
}

func calculateDiffCounts[T constraints.Ordered](measurements []T, windowSize int) {
	incCount := 0
	decCount := 0

	for i := range measurements {
		// Ignore the first entry and if the final entry has no full window
		if i == 0 || i+windowSize-1 >= len(measurements) {
			continue
		}
		var previousWindowValue T
		var currentWindowValue T
		for j := 0; j < windowSize; j++ {
			previousWindowValue += measurements[i-1+j]
			currentWindowValue += measurements[i+j]
		}

		if currentWindowValue > previousWindowValue {
			incCount++
		} else {
			decCount++
		}
	}
	fmt.Printf("[Window size %d] %d increased and %d decreased (%d measurements)\n", windowSize, incCount, decCount, len(measurements))
}
