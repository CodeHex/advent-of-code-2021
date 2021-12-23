package main

import (
	"adventofcode2021/pkg/fileparser"
	"adventofcode2021/pkg/sets"
	"adventofcode2021/pkg/slices"
	"fmt"
	"sort"
	"strings"
)

func main() {
	instructions := fileparser.ReadTypedLines("day22/input.txt", NewRebootStep)
	smallSteps := slices.Filter(instructions, SmallStep)

	fmt.Println("[Part 1] Total cubes on for initialization:", RunSteps(smallSteps).SumWeighted(SizeFunc))
	fmt.Println("[Part 2] Total cubes on for all instructions:", RunSteps(instructions).SumWeighted(SizeFunc))
}

type Box struct {
	minX, maxX int
	minY, maxY int
	minZ, maxZ int
}

type RebootStep struct {
	box   Box
	state bool
}

func NewRebootStep(line string) RebootStep {
	result := RebootStep{}
	parts := fileparser.SplitTrim[string](line, " ")
	if parts[0] == "on" {
		result.state = true
	}

	// Split into coords
	parts = fileparser.SplitTrim[string](parts[1], ",")
	coordsX := strings.TrimPrefix(parts[0], "x=")
	coordsY := strings.TrimPrefix(parts[1], "y=")
	coordsZ := strings.TrimPrefix(parts[2], "z=")

	xParts := fileparser.SplitTrim[int](coordsX, "..")
	yParts := fileparser.SplitTrim[int](coordsY, "..")
	zParts := fileparser.SplitTrim[int](coordsZ, "..")

	result.box.minX, result.box.maxX = xParts[0], xParts[1]+1
	result.box.minY, result.box.maxY = yParts[0], yParts[1]+1
	result.box.minZ, result.box.maxZ = zParts[0], zParts[1]+1
	return result
}

func SmallStep(step RebootStep) bool {
	if step.box.minX < -50 || step.box.maxX > 51 ||
		step.box.minY < -50 || step.box.maxY > 51 ||
		step.box.minZ < -50 || step.box.maxZ > 51 {
		return false
	}
	return true
}

func (b Box) Size() int {
	return (b.maxX - b.minX) * (b.maxY - b.minY) * (b.maxZ - b.minZ)
}

func SizeFunc(b Box) int {
	return b.Size()
}

// Decomposes box 1 into component boxes with the intersection of b1 and b2 removed
// Boolean indicates if there is no intersection
func Decompose(b1 Box, b2 Box) []Box {
	intersect := Intersect(b1, b2)
	if intersect == nil {
		return []Box{b1}
	}
	x := []int{b1.minX, b1.maxX, b2.minX, b2.maxX}
	y := []int{b1.minY, b1.maxY, b2.minY, b2.maxY}
	z := []int{b1.minZ, b1.maxZ, b2.minZ, b2.maxZ}
	sort.Ints(x)
	sort.Ints(y)
	sort.Ints(z)

	inBox1 := []Box{}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			for k := 0; k < 3; k++ {
				b := Box{x[i], x[i+1], y[j], y[j+1], z[k], z[k+1]}
				if Inside(b1, b) && b != *intersect {
					inBox1 = append(inBox1, b)
				}
			}
		}
	}
	return inBox1
}

func Inside(outer, inner Box) bool {
	return outer.minX <= inner.minX && outer.maxX >= inner.maxX &&
		outer.minY <= inner.minY && outer.maxY >= inner.maxY &&
		outer.minZ <= inner.minZ && outer.maxZ >= inner.maxZ
}

func RunSteps(steps []RebootStep) sets.Set[Box] {
	onLightBoxes := sets.NewEmptySet[Box]() // Defines a list of non intersecting boxes that are switched on

	for _, step := range steps {
		// If any of the light boxes intersect with the step box, decompose it
		// into non overlapping cuboids, with the intersect removed
		nextLightBoxes := sets.NewEmptySet[Box]()
		for onBox := range onLightBoxes {
			onBoxWithStepRemoved := Decompose(onBox, step.box)
			nextLightBoxes.AddSlice(onBoxWithStepRemoved)
		}

		// This should leave the box in the step undefined in the list of boxes. Add
		// it if we're turning them on
		if step.state {
			nextLightBoxes.Add(step.box)
		}
		onLightBoxes = nextLightBoxes
	}
	return onLightBoxes
}

func Intersect(b1 Box, b2 Box) *Box {
	b := &Box{
		minX: Max(b1.minX, b2.minX),
		maxX: Min(b1.maxX, b2.maxX),
		minY: Max(b1.minY, b2.minY),
		maxY: Min(b1.maxY, b2.maxY),
		minZ: Max(b1.minZ, b2.minZ),
		maxZ: Min(b1.maxZ, b2.maxZ),
	}

	if b.maxX >= b.minX && b.maxY >= b.minY && b.maxZ >= b.minZ {
		return b
	}
	return nil
}

func Min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
