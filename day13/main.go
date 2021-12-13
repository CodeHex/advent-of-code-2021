package main

import (
	"adventofcode2021/pkg/convert"
	"adventofcode2021/pkg/fileparser"
	"adventofcode2021/pkg/sets"
	"adventofcode2021/pkg/slices"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	data := fileparser.ReadLines("day13/input.txt")
	dots, folds := ParseInstructionsLines(data)

	count := 0
	for i, fold := range folds {
		count++
		fold.Apply(dots)
		if count == 1 {
			fmt.Printf("[Part 1] Dots after %s (fold %d): %d\n", folds[i].Name, count, len(dots))
		}
	}
	fmt.Printf("[Part 2] Final paper image\n\n")
	PrintDots(dots)
}

func ParseInstructionsLines(data []string) (sets.Set[Coord], []Fold) {
	isDotsLine := func(t string) bool { return !strings.HasPrefix(t, "fold") }
	dotsData, foldsData := slices.Divide(data, isDotsLine)
	dotsData = slices.TrimEnd(dotsData, 1) // Get rid of the last blank line
	dotsSet := sets.NewSetFromSlice(slices.Map(dotsData, NewCoord))
	folds := slices.Map(foldsData, NewFold)
	return dotsSet, folds
}

type Coord struct{ X, Y int }

func NewCoord(line string) Coord {
	parts := strings.Split(line, ",")
	return Coord{
		X: convert.Apply[int](parts[0]),
		Y: convert.Apply[int](parts[1]),
	}
}

type Reflector func(pos Coord) Coord

type Fold struct {
	Name    string
	axis    string
	lineVal int
	reflect Reflector
}

// ReflectXFunc generates a function that will reflect a coord around the defined x line
func ReflectXFunc(xr int) Reflector {
	return func(pos Coord) Coord { return Coord{X: 2*xr - pos.X, Y: pos.Y} }
}

// ReflectYFunc generates a function that will reflect a coord around the defined y line
func ReflectYFunc(yr int) Reflector {
	return func(pos Coord) Coord { return Coord{X: pos.X, Y: 2*yr - pos.Y} }
}

func NewFold(line string) Fold {
	foldStr := strings.Split(line, " ")[2]
	refLine := strings.Split(foldStr, "=")
	refAxis := refLine[0]
	lineNumber, _ := strconv.Atoi(refLine[1])
	var reflect Reflector
	switch refAxis {
	case "x":
		reflect = ReflectXFunc(lineNumber)
	case "y":
		reflect = ReflectYFunc(lineNumber)
	default:
		panic("unrecognised reflect line")
	}
	return Fold{Name: foldStr, axis: refAxis, lineVal: lineNumber, reflect: reflect}
}

func (f Fold) ShouldReflect(coord Coord) bool {
	// Only reflect coords that are above or the left of the reflect line
	switch f.axis {
	case "x":
		return coord.X > f.lineVal
	case "y":
		return coord.Y > f.lineVal
	default:
		panic("unrecognised reflect line")
	}
}

func (f Fold) Apply(dots sets.Set[Coord]) {
	for _, coord := range dots.ToSlice() {
		if f.ShouldReflect(coord) {
			dots.Remove(coord)
			dots.Add(f.reflect(coord))
		}
	}
}

func PrintDots(dots sets.Set[Coord]) {
	xSelectFunc := func(c Coord) int { return c.X }
	ySelectFunc := func(c Coord) int { return c.Y }
	maxX := slices.Max(slices.Map(dots.ToSlice(), xSelectFunc))
	maxY := slices.Max(slices.Map(dots.ToSlice(), ySelectFunc))
	for j := 0; j <= maxY; j++ {
		for i := 0; i <= maxX; i++ {
			if dots.IsMember(Coord{i, j}) {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
