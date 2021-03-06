package main

import (
	"adventofcode2021/pkg/fileparser"
	"adventofcode2021/pkg/slices"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	lines := fileparser.ReadLines("day19/input.txt")

	scanners := ParseScanners(lines)

	unalignedScanners := scanners[1:]          // Scanners where we don't know how to orientate to origin
	alignedScanners := []*Scanner{scanners[0]} // Scanners where we have know how to orientate to origin

	type attempt struct{ unaligned, aligned int }
	tries := make(map[attempt]struct{})
	// Loop tthough all unaligned scanners and attempt to align with
	// already aligned scanners
	for len(unalignedScanners) > 0 {
		fmt.Printf("Aligning %d scanner(s)...\n", len(unalignedScanners))
		for _, unaligned := range unalignedScanners {
			for _, aligned := range alignedScanners {
				// ignore this attempt as we've tried previously
				thisAttempt := attempt{unaligned.label, aligned.label}
				if _, ok := tries[thisAttempt]; ok {
					continue
				}
				tries[thisAttempt] = struct{}{}
				if ok := AttemptToAlign(unaligned, aligned); ok {
					break
				}
			}
		}
		// Reset lists based on what has been aligned
		alignedScanners, unalignedScanners = slices.Divide(scanners, IsAligned)
	}

	// Loop through all scanners, transform the beacons and add them to the list of found beacons
	finalBeacons := make(map[Coord]struct{})
	for _, scanner := range alignedScanners {
		for _, beacon := range scanner.beacons {
			finalBeacons[AddCoords(scanner.rotateToOrigin(beacon), *scanner.transformToOrigin)] = struct{}{}
		}
	}
	fmt.Println("[Part 1] Total number of beacons:", len(finalBeacons))

	// Loop through combinations of 2 scanners and calculate distance between them
	// based on how they are transformed to origin
	maxDist := 0
	for _, s1 := range alignedScanners {
		for _, s2 := range alignedScanners {
			dist := Distance(s1, s2)
			if dist > maxDist {
				maxDist = dist
			}
		}
	}
	fmt.Println("[Part 2] Max distance between scanners:", maxDist)
}

func IsAligned(s *Scanner) bool {
	return s.transformToOrigin != nil
}

func ParseScanners(lines []string) []*Scanner {
	markers := []int{}
	for i, line := range lines {
		if strings.HasPrefix(line, "--- scanner") {
			markers = append(markers, i)
		}
	}
	scanners := []*Scanner{}
	for i := range markers {
		start := markers[i]
		end := len(lines)
		if i < len(markers)-1 {
			end = markers[i+1]
		}
		scanners = append(scanners, NewScanner(lines[start:end]))
	}
	return scanners
}

type Coord struct{ x, y, z int }

type Rotation func(Coord) Coord

type Scanner struct {
	label             int
	beacons           []Coord
	transformToOrigin *Coord
	rotateToOrigin    Rotation
}

func NewScanner(data []string) *Scanner {
	name := data[0][4:]
	name = name[:len(name)-4]
	label, err := strconv.Atoi(name[8:])
	if err != nil {
		panic(err)
	}

	coords := []Coord{}
	for _, coordStr := range data[1:] {
		if coordStr != "" {
			parts := fileparser.SplitTrim[int](coordStr, ",")
			coords = append(coords, Coord{parts[0], parts[1], parts[2]})
		}
	}
	s := &Scanner{label: label, beacons: coords}
	// Define the origin scanner as aligned at 0,0,0 (no rotatio required)
	if label == 0 {
		s.transformToOrigin = &Coord{0, 0, 0}
		s.rotateToOrigin = func(c Coord) Coord { return c }
	}
	return s
}

func (s *Scanner) String() string {
	out := fmt.Sprintf("scanner %d (beacons:%d)\n\n", s.label, len(s.beacons))
	for _, b := range s.beacons {
		out += fmt.Sprintf("%5d%5d%5d\n", b.x, b.y, b.z)
	}
	return out
}

// Defines all possible rotation mappings of a cube (24 possible mapping)
var rotations = []Rotation{
	func(c Coord) Coord { return Coord{c.x, c.y, c.z} },
	func(c Coord) Coord { return Coord{c.x, -1 * c.z, c.y} },
	func(c Coord) Coord { return Coord{c.x, -1 * c.y, -1 * c.z} },
	func(c Coord) Coord { return Coord{c.x, c.z, -1 * c.y} },

	func(c Coord) Coord { return Coord{-1 * c.y, c.x, c.z} },
	func(c Coord) Coord { return Coord{c.z, c.x, c.y} },
	func(c Coord) Coord { return Coord{c.y, c.x, -1 * c.z} },
	func(c Coord) Coord { return Coord{-1 * c.z, c.x, -1 * c.y} },

	func(c Coord) Coord { return Coord{-1 * c.x, -1 * c.y, c.z} },
	func(c Coord) Coord { return Coord{-1 * c.x, -1 * c.z, -1 * c.y} },
	func(c Coord) Coord { return Coord{-1 * c.x, c.y, -1 * c.z} },
	func(c Coord) Coord { return Coord{-1 * c.x, c.z, c.y} },

	func(c Coord) Coord { return Coord{c.y, -1 * c.x, c.z} },
	func(c Coord) Coord { return Coord{c.z, -1 * c.x, -1 * c.y} },
	func(c Coord) Coord { return Coord{-1 * c.y, -1 * c.x, -1 * c.z} },
	func(c Coord) Coord { return Coord{-1 * c.z, -1 * c.x, c.y} },

	func(c Coord) Coord { return Coord{-1 * c.z, c.y, c.x} },
	func(c Coord) Coord { return Coord{c.y, c.z, c.x} },
	func(c Coord) Coord { return Coord{c.z, -1 * c.y, c.x} },
	func(c Coord) Coord { return Coord{-1 * c.y, -1 * c.z, c.x} },

	func(c Coord) Coord { return Coord{-1 * c.z, -1 * c.y, -1 * c.x} },
	func(c Coord) Coord { return Coord{-1 * c.y, c.z, -1 * c.x} },
	func(c Coord) Coord { return Coord{c.z, c.y, -1 * c.x} },
	func(c Coord) Coord { return Coord{c.y, -1 * c.z, -1 * c.x} },
}

func AttemptToAlign(unaligned *Scanner, aligned *Scanner) bool {
	// Create a list of beacons aligned to origin, from the alligned scaller
	mapToOriginFunc := func(b Coord) Coord { return aligned.rotateToOrigin(b) }
	alignedBeacons := slices.Map(aligned.beacons, mapToOriginFunc)

	// Loop through all rotations
	for _, r := range rotations {

		// Create a list of beacons from the unaligned scanner, aligned to this test rotation
		testRotationFunc := func(b Coord) Coord { return r(b) }
		rotatedBeaconsUnaligned := slices.Map(unaligned.beacons, testRotationFunc)

		// Compare the sets of beacons to see if they match
		shift, ok := CompareBeacons(alignedBeacons, rotatedBeaconsUnaligned)
		if !ok {
			continue
		}

		// If so, this rotation should apply to the unaligned beacon
		unaligned.rotateToOrigin = r
		transform := AddCoords(shift, *aligned.transformToOrigin)
		unaligned.transformToOrigin = &transform
		return true
	}
	return false
}

func CompareBeacons(beacons1 []Coord, beacons2 []Coord) (Coord, bool) {
	// Loop through each combination of the two sets of beacons.
	// For each combo, create a test shift vector that maps b2 beacons to b1
	for _, b1 := range beacons1 {
		for _, b2 := range beacons2 {
			matched := 0
			testShift := SubtractCoords(b1, b2)

			// Loop through the second set of beacons and apply the test shift
			// counting how many beacons line up exactly with beacons in the first set
			for _, test := range beacons2 {
				shifted := AddCoords(test, testShift)
				if slices.Contains(beacons1, shifted) {
					matched++
				}
			}
			// If we have lined up more than 12, these becons match. Return the shift
			// vector required to map the second set of beacons to the first set (and
			// hence is the shift from the two scanners)
			if matched >= 12 {
				// BINGO
				return testShift, true
			}
		}
	}
	return Coord{0, 0, 0}, false
}

func AddCoords(c1 Coord, c2 Coord) Coord {
	return Coord{c1.x + c2.x, c1.y + c2.y, c1.z + c2.z}
}

func SubtractCoords(c1 Coord, c2 Coord) Coord {
	return Coord{c1.x - c2.x, c1.y - c2.y, c1.z - c2.z}
}

func Distance(s1 *Scanner, s2 *Scanner) int {
	dx := Mod(s1.transformToOrigin.x - s2.transformToOrigin.x)
	dy := Mod(s1.transformToOrigin.y - s2.transformToOrigin.y)
	dz := Mod(s1.transformToOrigin.z - s2.transformToOrigin.z)
	return dx + dy + dz
}

func Mod(x int) int {
	if x > 0 {
		return x
	}
	return x * -1
}
