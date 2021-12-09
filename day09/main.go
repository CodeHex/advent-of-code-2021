package main

import (
	"adventofcode2021/pkg/convert"
	"adventofcode2021/pkg/fileparser"
	"adventofcode2021/pkg/matrices"
	"adventofcode2021/pkg/slices"
	"fmt"
	"sort"
)

const maxHeight = 9

func main() {
	lines := fileparser.ReadLines("day09/input.txt")
	runeToStr := func(x rune) string { return string(x) }
	splitter := func(x string) []string { return slices.Map([]rune(x), runeToStr) }

	// Represents the heights of the seabed
	seabed := matrices.NewMatrixFromLines(lines, splitter, convert.FuncFor[int]())

	// Represents if we have already mapped this point when mapping basins
	rows, columns := seabed.Dimensions()
	mapped := matrices.NewMatrix[bool](rows, columns)

	lowPoints := 0
	riskLevel := 0
	basinSizes := []int{}
	seabed.ForEach(func(x, y int, height int) {
		// Detect low point
		switch {
		case x != columns-1 && seabed.Get(x+1, y) <= height: // point to the right is lower (if there is one)
		case x != 0 && seabed.Get(x-1, y) <= height: // point to the left is lower (if there is one)
		case y != rows-1 && seabed.Get(x, y+1) <= height: // point downwards is lower (if there is one)
		case y != 0 && seabed.Get(x, y-1) <= height: // point upwards is lower (if there is one)
		default:
			lowPoints++
			riskLevel += (1 + height) // Found point
		}

		// Map the basin containing this point is in if it hasn't already been mapped
		if !mapped.Get(x, y) && height != maxHeight {
			basinSizes = append(basinSizes, mapBasin(seabed, mapped, x, y))
		}

	})

	fmt.Printf("[Part 1] Detected %d low points, total risk level is %d\n", lowPoints, riskLevel)

	sort.Ints(basinSizes)
	maxBasin1 := basinSizes[len(basinSizes)-1]
	maxBasin2 := basinSizes[len(basinSizes)-2]
	maxBasin3 := basinSizes[len(basinSizes)-3]
	outputBasinSize := maxBasin1 * maxBasin2 * maxBasin3

	fmt.Printf("[Part 2] Detected %d basins, largest 3 basins sizes multiplied is %d\n", len(basinSizes), outputBasinSize)
}

func mapBasin(seabed matrices.Matrix[int], mapped matrices.Matrix[bool], x, y int) int {
	// Mark this point as mapped
	rows, columns := seabed.Dimensions()
	outOfBounds := x < 0 || x > columns-1 || y < 0 || y > rows-1

	// Ignore point if its out of bounds, already mapped or at the maximum height
	if outOfBounds || mapped.Get(x, y) || seabed.Get(x, y) == maxHeight {
		return 0
	}
	mapped.Set(x, y, true)

	// Iterative determine the size by considering this point as 1 and
	// adding the size of neighbouring points
	size := 1
	size += mapBasin(seabed, mapped, x+1, y)
	size += mapBasin(seabed, mapped, x-1, y)
	size += mapBasin(seabed, mapped, x, y-1)
	size += mapBasin(seabed, mapped, x, y+1)
	return size
}
