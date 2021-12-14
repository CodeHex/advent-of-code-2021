package main

import (
	"adventofcode2021/pkg/fileparser"
	"adventofcode2021/pkg/maps"
	"fmt"
)

func main() {
	lines := fileparser.ReadLines("day14/input.txt")
	template := lines[0]
	mappingData := fileparser.ReadPairsFromStrings[string, string](lines[2:], " -> ")

	mapper := make(map[string]string)
	for _, t := range mappingData {
		mapper[t.Key] = t.Value
	}

	// Keep a counter of each pair ( e.g. ABCDE gets stored as AB, BC, CD and DE)
	counters := make(map[string]int)
	for i := 0; i < len(template)-1; i++ {
		counters[template[i:i+2]]++
	}

	// Remember the first element
	firstLetter := string(template[0])

	for i := 1; i <= 40; i++ {
		ProgressStep(counters, mapper)
		if i == 10 {
			stats := Stats(firstLetter, counters)
			maxLetter, maxVal := maps.MaxValue(stats)
			minLetter, minVal := maps.MinValue(stats)
			fmt.Printf("[Part 1] After 10 steps, most common letter (%s) minus least common (%s): %d\n", maxLetter, minLetter, maxVal-minVal)
		}

		if i == 40 {
			stats := Stats(firstLetter, counters)
			maxLetter, maxVal := maps.MaxValue(stats)
			minLetter, minVal := maps.MinValue(stats)
			fmt.Printf("[Part 2] After 40 steps, most common letter (%s) minus least common (%s): %d\n", maxLetter, minLetter, maxVal-minVal)
		}
	}
}

func ProgressStep(counters map[string]int, mapper map[string]string) {
	// Map each counter to its 2 new pairs e.g if the mapper indicates AC -> B
	// then reset the counter for AC to 0 and remember to add the same count to AB and BC
	additions := make(map[string]int)
	for oldPair, count := range counters {
		additions[string(oldPair[0])+mapper[oldPair]] += count
		additions[mapper[oldPair]+string(oldPair[1])] += count
		counters[oldPair] = 0
	}

	for newPair, newCount := range additions {
		counters[newPair] += newCount
	}
}

func Stats(firstLetter string, counters map[string]int) map[string]int {
	// Use the count against the last value of each pair, adding the first letter. e.g.
	// for ABCBC, the pairs are
	// AB, BC, CB, BC  which translates to
	// AB:1, BC:2, CB:1
	// Looking at the last letter of each pair
	// B:1, C:2, B:1
	// B:2, C:2
	// We need to add a single count of the first letter
	// as this isn't considered when looking at the last value in the pair
	// A:1, B:2, C:2
	result := make(map[string]int)
	for pair, count := range counters {
		result[string(pair[1])] += count
	}
	result[firstLetter]++
	return result
}
