package main

import (
	"adventofcode2021/pkg/fileparser"
	"fmt"
	"strconv"
)

func main() {
	nums := fileparser.ReadTypedLines("day18/input.txt", NewSnailPair)

	sum := Sum(nums)
	fmt.Println("[Part 1] Magnitude of sum of numbers is:", sum.Magnitude())

	max := len(nums)
	maxVal := 0
	// Consider every permutation of each number
	for i := 0; i < max; i++ {
		for j := 0; j < max; j++ {
			if i == j {
				continue
			}
			// Check addition both ways
			testij := nums[i].Add(nums[j])
			if val := testij.Magnitude(); val > maxVal {
				maxVal = val
			}

			testji := nums[j].Add(nums[i])
			if val := testji.Magnitude(); val > maxVal {
				maxVal = val
			}
		}
	}
	fmt.Println("[Part 2] Maximum mangitude of adding 2 numbers: ", maxVal)
}

type SnailPair struct {
	LeftVal    int
	RightVal   int
	LeftPair   *SnailPair
	RightPair  *SnailPair
	ParentPair *SnailPair
}

func (s *SnailPair) Clone() *SnailPair {
	pairString := s.String()
	return NewSnailPair(pairString)
}

func NewSnailPair(line string) *SnailPair {
	// Create a binary tree representing the snail pair
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

func Sum(pairs []*SnailPair) *SnailPair {
	var current *SnailPair
	for _, s := range pairs {
		if current == nil {
			current = s
			continue
		}
		current = current.Add(s)
	}
	return current
}

func (s *SnailPair) Add(a *SnailPair) *SnailPair {
	result := &SnailPair{
		LeftPair:  s.Clone(),
		RightPair: a.Clone(),
	}
	result.LeftPair.ParentPair = result
	result.RightPair.ParentPair = result
	result.reduce()
	return result
}

func (s *SnailPair) Magnitude() int {
	left := s.LeftVal
	if s.LeftPair != nil {
		left = s.LeftPair.Magnitude()
	}

	right := s.RightVal
	if s.RightPair != nil {
		right = s.RightPair.Magnitude()
	}
	return (3 * left) + (2 * right)
}

func (s *SnailPair) String() string {
	leftPart := strconv.Itoa(s.LeftVal)
	if s.LeftPair != nil {
		leftPart = s.LeftPair.String()
	}

	rightPart := strconv.Itoa(s.RightVal)
	if s.RightPair != nil {
		rightPart = s.RightPair.String()
	}
	return fmt.Sprintf("[%s,%s]", leftPart, rightPart)
}

func (s *SnailPair) Depth() int {
	current := s
	depth := 0
	for current.ParentPair != nil {
		depth++
		current = current.ParentPair
	}
	return depth
}

func (s *SnailPair) reduce() {
	nothingToDo := false
	for !nothingToDo {
		if s.explode() {
			continue
		}
		if s.split() {
			continue
		}
		nothingToDo = true
	}
}

func (s *SnailPair) explode() bool {
	if s == nil {
		return false
	}

	// Try to explode left most pairs before right most
	// If we have performed an explosion, exit immediately
	if boom := s.LeftPair.explode(); boom {
		return true
	}
	if boom := s.RightPair.explode(); boom {
		return true
	}

	// Only explode pairs nested 4 times or more, i.e. depth >= 4
	if s.Depth() < 4 {
		return false
	}

	// Propogate the left hand side value to the next node with a value
	atRoot := false
	prev := s
	current := s.ParentPair
	// Move to parent, if the immediate left is a value, add it to that
	if current.LeftPair == nil {
		current.LeftVal += s.LeftVal
	} else {
		// Keep moving up the tree until we are
		// at the root of the tree or come from a right branch
		for current.LeftPair == prev {
			if current.ParentPair == nil {
				atRoot = true
				break
			} else {
				prev = current
				current = current.ParentPair
			}
		}

		// If we reached the root, we can't move any further left
		if !atRoot {
			// If the left has a value, use it, otherwise move down the
			// right branch as far as possible
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

	// Propogate the right hand side value to the next node with a value
	// same logic as above but across the opposite side
	atRoot = false
	prev = s
	current = s.ParentPair
	if current.RightPair == nil {
		current.RightVal += s.RightVal
	} else {
		for current.RightPair == prev {
			if current.ParentPair == nil {
				atRoot = true
				break
			} else {
				prev = current
				current = current.ParentPair
			}
		}
		if !atRoot {
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

	// Explode the actually pair, by moving to the parent node
	// and setting the corresponding value to 0 depending on
	// if we came from the left or right side
	parent := s.ParentPair
	if parent.LeftPair == s {
		parent.LeftVal = 0
		parent.LeftPair = nil
	} else {
		parent.RightVal = 0
		parent.RightPair = nil
	}
	return true
}

func (s *SnailPair) split() bool {
	if s == nil {
		return false
	}

	// Attept to split the left pairs first
	if rip := s.LeftPair.split(); rip {
		return true
	}

	if s.LeftVal >= 10 {
		newLeft := s.LeftVal / 2
		newRight := s.LeftVal - newLeft
		newPair := &SnailPair{LeftVal: newLeft, RightVal: newRight, ParentPair: s}
		s.LeftVal = 0
		s.LeftPair = newPair
		return true
	}

	// Now start the right pairs
	if s.RightVal >= 10 {
		newLeft := s.RightVal / 2
		newRight := s.RightVal - newLeft
		newPair := &SnailPair{LeftVal: newLeft, RightVal: newRight, ParentPair: s}
		s.RightVal = 0
		s.RightPair = newPair
		return true
	}

	if rip := s.RightPair.split(); rip {
		return true
	}
	return false
}
