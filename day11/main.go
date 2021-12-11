package main

import (
	"adventofcode2021/pkg/convert"
	"adventofcode2021/pkg/fileparser"
	"adventofcode2021/pkg/matrices"
	"adventofcode2021/pkg/slices"
	"fmt"
	"strconv"
)

func main() {
	lines := fileparser.ReadLines("day11/input.txt")
	runeToStr := func(x rune) string { return string(x) }
	splitter := func(x string) []string { return slices.Map([]rune(x), runeToStr) }

	// Represents the heights of the seabed
	data := matrices.NewIntMatrixFromLines(lines, splitter, convert.FuncFor[int]())
	octopi := Octopi(data)

	step := 0
	totalFlashes := 0
	flashesAfter100Steps := 0
	firstAllFlashedOccurred := false
	stepWhenFirstAllFlashed := 0
	octopiView100Steps := ""
	octopiViewFirstAllFlashed := ""
	// Keep progressing the octopi until we have at least reached 100 steps
	// and we have seen all the octopi flash at the same time
	for !(step > 100 && firstAllFlashedOccurred) {
		step++
		stepFlashes, allFlashed := octopi.ProgressStep()
		totalFlashes += stepFlashes

		if step == 100 {
			flashesAfter100Steps = totalFlashes
			octopiView100Steps = octopi.String()
		}

		if allFlashed && !firstAllFlashedOccurred {
			firstAllFlashedOccurred = true
			stepWhenFirstAllFlashed = step
			octopiViewFirstAllFlashed = octopi.String()
		}
	}

	fmt.Println("[Part 1] After 100 steps")
	fmt.Println(octopiView100Steps)
	fmt.Printf("total flashes: %d\n", flashesAfter100Steps)
	fmt.Println()
	fmt.Println("[Part 2] Until all octopi flash")
	fmt.Println(octopiViewFirstAllFlashed)
	fmt.Printf("steps: %d\n", stepWhenFirstAllFlashed)
}

type Octopi matrices.IntMatrix[int]

func (o Octopi) m() matrices.IntMatrix[int] {
	return matrices.IntMatrix[int](o)
}

func (o Octopi) String() string {
	out := ""
	for _, line := range o {
		for _, val := range line {
			out += strconv.Itoa(val)
		}
		out += "\n"
	}
	return out
}

func (o Octopi) ProgressStep() (flashes int, allFlashed bool) {
	// Increase all energy levels by 1
	o.m().IncrementAll()

	// Check for flashes on all octopi. If a flash is detected, also check neighbouring octopi as their
	// energy level will increase. Will also reset the octopi to 0 energy level if they do flash
	total := 0
	o.m().ForEach(func(x, y int, energy int) {
		total += o.CheckFlash(x, y)
	})
	return total, o.m().NumOfElements() == total
}

func (o Octopi) CheckFlash(x, y int) int {
	// Ignore point if already flashed or doesn't have enough energy
	if o.m().Get(x, y) < 10 {
		return 0
	}

	// Reset the octopus energy level as it has flashed
	o.m().Set(x, y, 0)

	// This point flashes, so increment all neighbouring octopi (if they haven't flashed)
	o.m().ForEachNeighbour(true, x, y, func(x1, y1 int) {
		o.IncrementIfNotFlashed(x1, y1)
	})

	// Since we have incremented the neighbours, check if they have now flashed
	flashes := 1
	o.m().ForEachNeighbour(true, x, y, func(x1, y1 int) {
		flashes += o.CheckFlash(x1, y1)
	})
	return flashes
}

func (o Octopi) IncrementIfNotFlashed(x, y int) {
	if o.m().Get(x, y) != 0 {
		o.m().Increment(x, y)
	}
}
