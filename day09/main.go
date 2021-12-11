package main

import (
	"adventofcode2021/pkg/fileparser"
	"adventofcode2021/pkg/matrices"
	"fmt"
	"sort"
)

const maxHeight = 9

func main() {
	// Represents the heights of the seabed
	seabed := fileparser.ReadDigitMatrix("day09/input.txt")

	// Represents if we have already mapped this point when mapping basins
	mapped := matrices.NewMatrix[bool](seabed.Rows, seabed.Columns)

	lowPoints := 0
	riskLevel := 0
	basinSizes := []int{}
	seabed.ForEach(func(pointX, pointY int, height int) {
		isLowPoint := true
		// Check all neighbours, this point will still be a low point
		// if the point is lower than all of its neighbours
		seabed.ForEachNeighbour(false, pointX, pointY, func(x, y int) {
			// If we still think its a low point, check the next neighbour that its lower
			if isLowPoint {
				neighbourHeight := seabed.Get(x, y)
				isLowPoint = height < neighbourHeight
			}
		})

		if isLowPoint {
			lowPoints++
			riskLevel += (1 + height)
		}

		// Map the basin containing this point is in if it hasn't already been mapped
		if !mapped.Get(pointX, pointY) && height != maxHeight {
			basinSizes = append(basinSizes, mapBasin(seabed, mapped, pointX, pointY))
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

func mapBasin(seabed matrices.IntMatrix[int], mapped matrices.Matrix[bool], pointX, pointY int) int {
	// Ignore point if its out of bounds, already mapped or at the maximum height
	if mapped.Get(pointX, pointY) || seabed.Get(pointX, pointY) == maxHeight {
		return 0
	}
	mapped.Set(pointX, pointY, true)

	// Iterative determine the size by considering this point as 1 and
	// adding the size of neighbouring points
	size := 1
	seabed.ForEachNeighbour(false, pointX, pointY, func(x, y int) {
		size += mapBasin(seabed, mapped, x, y)
	})
	return size
}
