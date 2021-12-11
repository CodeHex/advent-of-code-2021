package main

import (
	"adventofcode2021/pkg/fileparser"
	"adventofcode2021/pkg/matrices"
	"fmt"
)

func main() {
	data := fileparser.ReadDigitMatrix("day11/input.txt")
	octopi := NewOctopi(data)
	octopi.RunSimulation()

	fmt.Printf("[Part 1] After 100 steps, total flashes: %d\n%s\n", octopi.FlashesAfter100Steps, octopi.Snapshot100Steps)
	fmt.Printf("[Part 2] Until all octopi flash, steps taken: %d\n%s", octopi.FirstSyncFlashStep, octopi.SnapshotFirstSync)
}

type Octopi struct {
	step                 int
	data                 matrices.IntMatrix[int]
	FlashesAfter100Steps int
	Snapshot100Steps     string
	FirstSyncFlashStep   int
	SnapshotFirstSync    string
}

func NewOctopi(data matrices.IntMatrix[int]) *Octopi {
	return &Octopi{data: data}
}

func (o *Octopi) RunSimulation() {
	// Keep progressing the octopi until we have at least reached 100 steps
	// and we have seen all the octopi flash at the same time
	syncFlashOccurred := false
	totalFlashes := 0
	for !(o.step > 100 && syncFlashOccurred) {
		o.step++
		stepFlashes, syncFlashed := o.progressStep()
		totalFlashes += stepFlashes

		if o.step == 100 {
			o.FlashesAfter100Steps = totalFlashes
			o.Snapshot100Steps = o.data.CompactString()
		}

		if syncFlashed && !syncFlashOccurred {
			syncFlashOccurred = true
			o.FirstSyncFlashStep = o.step
			o.SnapshotFirstSync = o.data.CompactString()
		}
	}
}

func (o *Octopi) progressStep() (flashes int, allFlashed bool) {
	// Increase all energy levels by 1
	o.data.IncrementAll()

	// Check for flashes on all octopi. If a flash is detected, also check neighbouring octopi as their
	// energy level will increase. Will also reset the octopi to 0 energy level if they do flash
	total := 0
	o.data.ForEach(func(x, y int, energy int) {
		total += o.checkFlash(x, y)
	})
	return total, o.data.Size == total
}

func (o *Octopi) checkFlash(x, y int) int {
	// Ignore point if already flashed or doesn't have enough energy
	if o.data.Get(x, y) < 10 {
		return 0
	}

	// Reset the octopus energy level as it has flashed
	o.data.Set(x, y, 0)

	// This point flashes, so increment all neighbouring octopi (if they haven't flashed)
	o.data.ForEachNeighbour(true, x, y, func(x1, y1 int) {
		o.incrementIfNotFlashed(x1, y1)
	})

	// Since we have incremented the neighbours, check if they have now flashed
	flashes := 1
	o.data.ForEachNeighbour(true, x, y, func(x1, y1 int) {
		flashes += o.checkFlash(x1, y1)
	})
	return flashes
}

func (o *Octopi) incrementIfNotFlashed(x, y int) {
	if o.data.Get(x, y) != 0 {
		o.data.Increment(x, y)
	}
}
