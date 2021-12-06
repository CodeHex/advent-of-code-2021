package main

import (
	"adventofcode2021/pkg/fileparser"
	"adventofcode2021/pkg/maps"
	"fmt"
)

func main() {
	input := fileparser.ReadLines("day06/input.txt")
	startingFish := fileparser.Split[int](input[0], ",")
	stats := NewFishStats(startingFish)

	// Progress for 80 days
	for i := 0; i < 80; i++ {
		stats = ProgressDay(stats)
	}
	fmt.Println(stats)
	fmt.Printf("Total fish after 80 days is %d\n", maps.SumValues(stats))

	// Progress for up to 256 days
	for i := 0; i < 256-80; i++ {
		stats = ProgressDay(stats)
	}
	fmt.Println(stats)
	fmt.Printf("Total fish after 256 days is %d\n", maps.SumValues(stats))
}

func NewFishStats(startingFish []int) map[int]int {
	stats := make(map[int]int)
	for _, fish := range startingFish {
		stats[fish]++
	}
	return stats
}

func ProgressDay(stats map[int]int) map[int]int {
	results := make(map[int]int)
	for fish, count := range stats {
		// For the fish that has reached the end of their timer
		// reset them to 6, but also spawn new ones of 8
		if fish == 0 {
			results[6] += count
			results[8] += count
		} else {
			results[fish-1] += count
		}
	}
	return results
}
