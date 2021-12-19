package main

import (
	"adventofcode2021/pkg/fileparser"
	"fmt"
	"strconv"
)

func main() {
	nums := fileparser.ReadTypedLines("day18/input.txt", NewSnailPair)

	current := nums[0]
	for i := 1; i < len(nums); i++ {
		current = Add(current, nums[i])
		Reduce(current)
	}
	fmt.Println(Magnitude(current))

	nums = fileparser.ReadTypedLines("day18/input.txt", NewSnailPair)
	max := len(nums)
	maxVal := 0
	for i := 0; i < max; i++ {
		for j := 0; j < max; j++ {
			if i == j {
				continue
			}
			numsI := fileparser.ReadTypedLines("day18/input.txt", NewSnailPair)
			numsJ := fileparser.ReadTypedLines("day18/input.txt", NewSnailPair)
			num := Add(numsI[i], numsJ[j])
			Reduce(num)
			val := Magnitude(num)
			if val > maxVal {
				maxVal = val
			}

			numsI = fileparser.ReadTypedLines("day18/input.txt", NewSnailPair)
			numsJ = fileparser.ReadTypedLines("day18/input.txt", NewSnailPair)
			num = Add(numsJ[j], numsI[i])
			Reduce(num)
			val = Magnitude(num)
			if val > maxVal {
				maxVal = val
			}
		}
	}
	fmt.Println(maxVal)

}

type SnailPair struct {
	LeftVal    int
	RightVal   int
	LeftPair   *SnailPair
	RightPair  *SnailPair
	ParentPair *SnailPair
}

func NewSnailPair(line string) *SnailPair {
	var current *SnailPair
	stack := []*SnailPair{}
	pointer := "left"
	for _, c := range line {
		switch c {
		case '[':
			newPair := &SnailPair{}
			if current != nil {
				if pointer == "left" {
					current.LeftPair = newPair
				} else {
					current.RightPair = newPair
				}
				newPair.ParentPair = current
			}
			stack = append(stack, newPair)
			current = newPair
			pointer = "left"
		case ']':
			pointer = "left"
			if len(stack) == 1 {
				return stack[0]
			}
			stack = stack[0 : len(stack)-1]
			current = stack[len(stack)-1]
		case ',':
			pointer = "right"
		default:
			digit, err := strconv.Atoi(string(c))
			if err != nil {
				panic(err)
			}
			if pointer == "left" {
				current.LeftVal = digit
			} else {
				current.RightVal = digit
			}
		}
	}
	return nil
}

func (s *SnailPair) String() string {
	leftPart := ""
	if s.LeftPair == nil {
		leftPart = strconv.Itoa(s.LeftVal)
	} else {
		leftPart = s.LeftPair.String()
	}

	rightPart := ""

	if s.RightPair == nil {
		rightPart = strconv.Itoa(s.RightVal)
	} else {
		rightPart = s.RightPair.String()
	}
	return fmt.Sprintf("[%s,%s]", leftPart, rightPart)
}

func Reduce(s *SnailPair) {
	noop := false
	for !noop {
		if Explode(s) {
			continue
		}
		if Split(s) {
			continue
		}
		noop = true
	}
}

func Add(s1 *SnailPair, s2 *SnailPair) *SnailPair {
	new := &SnailPair{
		LeftPair:  s1,
		RightPair: s2,
	}
	s1.ParentPair = new
	s2.ParentPair = new
	return new
}

func Explode(s *SnailPair) bool {
	if s == nil {
		return false
	}

	boom := Explode(s.LeftPair)
	if boom {
		return true
	}
	boom = Explode(s.RightPair)
	if boom {
		return true
	}

	depth := 0
	depthCurrent := s.ParentPair
	for depthCurrent != nil {
		depth++
		depthCurrent = depthCurrent.ParentPair
	}
	if depth != 4 {
		return false
	}

	// Explode this pair
	if s.LeftPair != nil || s.RightPair != nil {
		panic("pairs detected greater than depth 4")
	}

	none := false
	prev := s
	current := s.ParentPair
	if current.LeftPair == nil {
		current.LeftVal += s.LeftVal
	} else {
		for current.LeftPair == prev {
			if current.ParentPair == nil {
				none = true
				break
			} else {
				prev = current
				current = current.ParentPair
			}
		}
		if !none {
			if current.LeftPair == nil {
				current.LeftVal += s.LeftVal
			} else {
				current = current.LeftPair
				for current.RightPair != nil {
					current = current.RightPair
				}
				current.RightVal += s.LeftVal
			}
		}
	}

	none = false
	prev = s
	current = s.ParentPair
	if current.RightPair == nil {
		current.RightVal += s.RightVal
	} else {
		for current.RightPair == prev {
			if current.ParentPair == nil {
				none = true
				break
			} else {
				prev = current
				current = current.ParentPair
			}
		}
		if !none {
			if current.RightPair == nil {
				current.RightVal += s.RightVal
			} else {
				current = current.RightPair
				for current.LeftPair != nil {
					current = current.LeftPair
				}
				current.LeftVal += s.RightVal
			}
		}
	}

	// Explode it
	parent := s.ParentPair
	if parent.LeftPair == s {
		parent.LeftVal = 0
		parent.LeftPair = nil
	}
	if parent.RightPair == s {
		parent.RightVal = 0
		parent.RightPair = nil
	}
	return true
}

func Split(s *SnailPair) bool {
	if s == nil {
		return false
	}

	plip := Split(s.LeftPair)
	if plip {
		return true
	}

	if s.LeftPair == nil && s.LeftVal >= 10 {
		newLeft := s.LeftVal / 2
		newRight := s.LeftVal - newLeft
		newPair := &SnailPair{
			LeftVal:    newLeft,
			RightVal:   newRight,
			ParentPair: s,
		}
		s.LeftVal = 0
		s.LeftPair = newPair
		return true
	}

	if s.RightPair == nil && s.RightVal >= 10 {
		newLeft := s.RightVal / 2
		newRight := s.RightVal - newLeft
		newPair := &SnailPair{
			LeftVal:    newLeft,
			RightVal:   newRight,
			ParentPair: s,
		}
		s.RightVal = 0
		s.RightPair = newPair
		return true
	}

	plip = Split(s.RightPair)
	if plip {
		return true
	}
	return false
}

func Magnitude(s *SnailPair) int {
	left := s.LeftVal
	if s.LeftPair != nil {
		left = Magnitude(s.LeftPair)
	}

	right := s.RightVal
	if s.RightPair != nil {
		right = Magnitude(s.RightPair)
	}

	return (3 * left) + (2 * right)
}
