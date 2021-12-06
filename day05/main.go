package main

import (
	"adventofcode2021/pkg/fileparser"
	"adventofcode2021/pkg/maps"
	"adventofcode2021/pkg/slices"
	"fmt"
	"strings"
)

func main() {
	vents := fileparser.ReadTypedLines("day05/input.txt", NewVent)

	// Don't consider diagonal vents
	seabedFloorNoDiags := NewFloor(vents, true)
	fmt.Printf("[Part 1] No diagonal vents - Number of at least 2 overlaps: %d\n", seabedFloorNoDiags.Overlaps())

	// Consider diagonal vents
	seabedFloorAll := NewFloor(vents, false)
	fmt.Printf("[Part 2] All vents - Number of at least 2 overlaps: %d\n", seabedFloorAll.Overlaps())
}

type Vent struct {
	X1, Y1 int
	X2, Y2 int
}

func NewVent(data string) Vent {
	parts := strings.Split(data, "->")
	start := fileparser.SplitTrim[int](parts[0], ",")
	stop := fileparser.SplitTrim[int](parts[1], ",")
	return Vent{X1: start[0], Y1: start[1], X2: stop[0], Y2: stop[1]}
}

type Floor struct {
	ventCount [][]int
	stats     map[int]int
}

func NewFloor(vents []Vent, ignoreDiag bool) Floor {
	X1s := slices.Map(vents, func(v Vent) int { return v.X1 })
	X2s := slices.Map(vents, func(v Vent) int { return v.X2 })
	Y1s := slices.Map(vents, func(v Vent) int { return v.Y1 })
	Y2s := slices.Map(vents, func(v Vent) int { return v.Y2 })
	Xs := append(X1s, X2s...)
	Ys := append(Y1s, Y2s...)
	maxX := slices.Max(Xs)
	maxY := slices.Max(Ys)

	// Initialize floor
	floor := Floor{
		ventCount: slices.InitGrid[int](maxX+1, maxY+1),
		stats:     make(map[int]int),
	}
	// Set all counts to 0
	floor.stats[0] = (maxX + 1) * (maxY + 1)

	// Create vent walkers for each vent and walk them until they are finished
	for _, vent := range vents {
		if ignoreDiag && vent.X1 != vent.X2 && vent.Y1 != vent.Y2 {
			continue
		}

		walker := NewVentWalker(vent)
		for !walker.Finished {
			// Decrease the previous vent count and increase the next one
			currentVentCount := floor.ventCount[walker.CurrentY][walker.CurrentX]
			floor.stats[currentVentCount]--
			floor.ventCount[walker.CurrentY][walker.CurrentX]++
			floor.stats[currentVentCount+1]++
			walker.NextStep()
		}
	}
	return floor
}

func (f Floor) PrintFloor() {
	for _, line := range f.ventCount {
		for _, val := range line {
			fmt.Printf("%2d", val)
		}
		fmt.Println()
	}
}

func (f Floor) Overlaps() int {
	selectFunc := func(key int) bool { return key >= 2 }
	return maps.SumValuesFor(f.stats, selectFunc)
}

type VentWalker struct {
	vent               Vent
	CurrentX, CurrentY int
	stepX, stepY       int
	Finished           bool
}

func NewVentWalker(v Vent) *VentWalker {
	walker := &VentWalker{vent: v, CurrentX: v.X1, CurrentY: v.Y1}
	switch {
	case v.X2 > v.X1:
		walker.stepX = 1
	case v.X1 > v.X2:
		walker.stepX = -1
	}
	switch {
	case v.Y2 > v.Y1:
		walker.stepY = 1
	case v.Y1 > v.Y2:
		walker.stepY = -1
	}
	return walker
}

func (w *VentWalker) NextStep() {
	if w.CurrentX == w.vent.X2 && w.CurrentY == w.vent.Y2 {
		w.Finished = true
	} else {
		w.CurrentX += w.stepX
		w.CurrentY += w.stepY
	}
}
