package main

import (
	"adventofcode2021/pkg/fileparser"
	"adventofcode2021/pkg/tuples"
	"fmt"
)

func main() {
	instructions := fileparser.ReadPairs[string, int]("day02/input.txt", " ")

	pos, depth := calcBasicLoc(instructions)
	result := pos * depth
	fmt.Printf("[Part1] Horizontal pos: %d, Depth: %d, Result: %d (%d instructions)\n", pos, depth, result, len(instructions))

	pos, depth = calcAdvancedLoc(instructions)
	result = pos * depth
	fmt.Printf("[Part2] Horizontal pos: %d, Depth: %d, Result: %d (%d instructions)\n", pos, depth, result, len(instructions))
}

func calcBasicLoc(steps []tuples.Pair[string, int]) (pos, depth int) {
	pos, depth = 0, 0
	for _, pair := range steps {
		switch pair.Key {
		case "forward":
			pos = pos + pair.Value
		case "up":
			depth = depth - pair.Value
		case "down":
			depth = depth + pair.Value
		default:
			panic("unrecognized step")
		}
	}
	return pos, depth
}

func calcAdvancedLoc(steps []tuples.Pair[string, int]) (pos, depth int) {
	pos, depth = 0, 0
	aim := 0
	for _, pair := range steps {
		switch pair.Key {
		case "forward":
			pos = pos + pair.Value
			depth = depth + (aim * pair.Value)
		case "up":
			aim = aim - pair.Value
		case "down":
			aim = aim + pair.Value
		default:
			panic("unrecognized step")
		}
	}
	return pos, depth
}
