package main

import (
	"adventofcode2021/pkg/fileparser"
	"fmt"
	"strconv"
)

func main() {
	ops := fileparser.ReadTypedLines("day24/input.txt", NewOp)

	// Inspecting the input, the operations are split into blocks with the following properties
	// Each block starts an input to w (wiping out the previous value of w)
	// Previous values of x and y are irrelevant during the block as they are wiped before they are used
	// Only z's value is used in future blocks

	highestCode := SearchHighest(ops)
	fmt.Println("[Part 1] Highest code is ", highestCode)

	lowestCode := SearchLowest(ops)
	fmt.Println("[Part 2] Lowest code is ", lowestCode)
}

func SplitOps(ops []Op) [][]Op {
	result := [][]Op{}
	current := []Op{}
	for _, op := range ops {
		if op.RequiresInput {
			if len(current) > 0 {
				result = append(result, current)
			}
			current = []Op{}
		}
		current = append(current, op)
	}
	result = append(result, current)
	return result
}

func SearchLowest(ops []Op) int64 {
	return Search(ops, false)
}

func SearchHighest(ops []Op) int64 {
	return Search(ops, true)
}

func Search(ops []Op, isHigher bool) int64 {
	isLowerFunc := func(stored, attempt int64) bool {
		return attempt < stored || stored == 0
	}
	isHigherFunc := func(stored, attempt int64) bool {
		return attempt > stored
	}
	compareFunc := isLowerFunc
	if isHigher {
		compareFunc = isHigherFunc
	}

	opsBlocks := SplitOps(ops)

	digits := make(map[int64]int64)
	digits[0] = 0

	for b, block := range opsBlocks {
		newDigits := make(map[int64]int64)

		// Try each digit for the input to the block
		for digit := int64(1); digit <= 9; digit++ {

			// Try each previous value of z from the previous block
			for zVal, maxVal := range digits {
				endState := ApplyOps(State{z: zVal}, block, []int64{digit})

				// Calculate what the digits would be using this digit
				nextVal := (maxVal * 10) + int64(digit)
				if compareFunc(newDigits[endState.z], nextVal) {
					newDigits[endState.z] = nextVal
				}
			}
		}
		digits = newDigits
		fmt.Printf("Block %d of %d processed (%d entries)\n", b+1, len(opsBlocks), len(newDigits))
	}
	return digits[0]
}

type State struct {
	x, y, z, w int64
}

type Op struct {
	RequiresInput bool
	InputFunc     func(int64, State) State
	Func          func(State) State
	Label         string
}

func (s State) String() string {
	return fmt.Sprintf("(x:%d, y:%d, z:%d, w:%d)", s.x, s.y, s.z, s.w)
}

func (o Op) String() string {
	return o.Label
}

func ApplyOps(s State, ops []Op, input []int64) State {
	result := s
	inputIndex := 0
	for _, o := range ops {
		if o.RequiresInput {
			if inputIndex >= len(input) {
				panic("not enough inputs")
			}
			result = o.InputFunc(input[inputIndex], result)
			inputIndex++
		} else {
			result = o.Func(result)
		}
	}
	return result
}

func NewOp(line string) Op {
	parts := fileparser.SplitTrim[string](line, " ")
	var op Op
	switch parts[0] {
	case "inp":

		op = Op{RequiresInput: true, InputFunc: func(i int64, s State) State {
			return s.Set(parts[1], i)
		}}
	case "add":
		op = Op{Func: func(s State) State {
			return s.Set(parts[1], s.Get(parts[1])+s.GetOrConvert(parts[2]))
		}}
	case "mul":
		op = Op{Func: func(s State) State {
			return s.Set(parts[1], s.Get(parts[1])*s.GetOrConvert(parts[2]))
		}}
	case "div":
		op = Op{Func: func(s State) State {
			return s.Set(parts[1], s.Get(parts[1])/s.GetOrConvert(parts[2]))
		}}
	case "mod":
		op = Op{Func: func(s State) State {
			return s.Set(parts[1], s.Get(parts[1])%s.GetOrConvert(parts[2]))
		}}
	case "eql":
		op = Op{Func: func(s State) State {
			result := int64(0)
			if s.Get(parts[1]) == s.GetOrConvert(parts[2]) {
				result = 1
			}
			return s.Set(parts[1], result)
		}}
	default:
		panic("unrecognized op")
	}
	op.Label = line
	return op
}

func (s State) Get(d string) int64 {
	switch d {
	case "x":
		return s.x
	case "y":
		return s.y
	case "z":
		return s.z
	case "w":
		return s.w
	default:
		panic("unrecognized state dimension")
	}
}

func (s State) Set(d string, val int64) State {
	result := s
	switch d {
	case "x":
		result.x = val
	case "y":
		result.y = val
	case "z":
		result.z = val
	case "w":
		result.w = val
	default:
		panic("unrecognized state dimension")
	}
	return result
}

func (s State) GetOrConvert(d string) int64 {
	switch d {
	case "x", "y", "z", "w":
		return s.Get(d)
	default:
		val, err := strconv.Atoi(d)
		if err != nil {
			panic(err)
		}
		return int64(val)
	}
}
