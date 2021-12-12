package main

import (
	"adventofcode2021/pkg/fileparser"
	"adventofcode2021/pkg/sets"
	"adventofcode2021/pkg/slices"
	"adventofcode2021/pkg/tuples"
	"fmt"
	"sort"
	"strings"
	"unicode"
)

type validatorFunc func(w *Walker) func(c Cave) bool

func main() {
	segments := fileparser.ReadPairs[string, string]("day12/input.txt", "-")
	solver1 := NewSolver(segments)
	total := solver1.Solve(part1ValidateFunc)
	fmt.Printf("[Part 1] Total paths with single visit to small caves: %d\n", total)

	solver2 := NewSolver(segments)
	total = solver2.Solve(part2ValidateFunc)
	fmt.Printf("[Part 2] Total paths with double visit to a single small cave: %d\n", total)
}

type Cave struct {
	Name    string
	IsSmall bool
}

func NewCave(raw string) Cave {
	// Only check the first letter to see if its lowercase
	return Cave{Name: raw, IsSmall: unicode.IsLower([]rune(raw)[0])}
}

type Solver struct {
	completeWalkers []*Walker
	connectedCaves  map[Cave]sets.Set[Cave]
}

func NewSolver(segments []tuples.Pair[string, string]) *Solver {
	connected := make(map[Cave]sets.Set[Cave])
	for _, seg := range segments {
		cave1, cave2 := NewCave(seg.Key), NewCave(seg.Value)
		if connected[cave1] == nil {
			connected[cave1] = sets.NewEmptySet[Cave]()
		}
		if connected[cave2] == nil {
			connected[cave2] = sets.NewEmptySet[Cave]()
		}
		connected[cave1].Add(cave2)
		connected[cave2].Add(cave1)
	}
	return &Solver{connectedCaves: connected, completeWalkers: []*Walker{}}
}

func (s *Solver) PrintResult() {
	out := slices.Map(s.completeWalkers, func(w *Walker) string { return w.String() })
	sort.Strings(out)
	for _, val := range out {
		fmt.Println(val)
	}
	fmt.Printf("Number of paths: %d\n", len(s.completeWalkers))
}

type Walker struct {
	current          Cave
	path             []Cave
	caveVisitedCount map[Cave]int
	doubleVisited    bool
}

func NewStartingWalker() *Walker {
	w := &Walker{caveVisitedCount: make(map[Cave]int)}
	w.moveTo(NewCave("start"))
	return w
}

func (w *Walker) clone() *Walker {
	// Clone by retracing the path steps
	clone := &Walker{caveVisitedCount: make(map[Cave]int)}
	for _, step := range w.path {
		clone.moveTo(step)
	}
	return clone
}

func (w *Walker) moveTo(c Cave) {
	w.current = c
	w.path = append(w.path, c)
	w.caveVisitedCount[c]++
	if c.IsSmall && w.caveVisitedCount[c] == 2 {
		w.doubleVisited = true
	}
}

func (w *Walker) validMovePart1(c Cave) bool {
	return !c.IsSmall || w.caveVisitedCount[c] == 0
}

func part1ValidateFunc(w *Walker) func(c Cave) bool {
	return func(c Cave) bool { return w.validMovePart1(c) }
}

func (w *Walker) validMovePart2(c Cave) bool {
	switch {
	case c.Name == "start":
		return false
	case !c.IsSmall || w.caveVisitedCount[c] == 0:
		return true
	case c.IsSmall && w.caveVisitedCount[c] == 1 && !w.doubleVisited:
		return true
	default:
		return false
	}
}

func part2ValidateFunc(w *Walker) func(c Cave) bool {
	return func(c Cave) bool { return w.validMovePart2(c) }
}

func (w *Walker) String() string {
	pathNames := slices.Map(w.path, func(x Cave) string { return x.Name })
	return strings.Join(pathNames, ",")
}

func (s *Solver) Solve(validator validatorFunc) int {
	runningWalkers := []*Walker{NewStartingWalker()}
	// Keep solving until all running walkers have completed
	for len(runningWalkers) > 0 {

		// Progress each running walker
		nextWalkers := []*Walker{}
		for _, walker := range runningWalkers {
			// Check if the walker has finished and if so move it to the completed walkers
			if walker.current.Name == "end" {
				s.completeWalkers = append(s.completeWalkers, walker)
				continue
			}
			// Check if the walker has any valid moves, and if not kill it
			validMoves := s.connectedCaves[walker.current].Filter(validator(walker))
			if len(validMoves) == 0 {
				continue
			}

			// Generate walkers for each possible valid move. For the last move
			// reuse the original walker, otherwise, clone a new one
			for i, validMove := range validMoves.ToSlice() {
				lastIndex := len(validMoves) - 1
				var nextWalker *Walker
				if i == lastIndex {
					nextWalker = walker
				} else {
					nextWalker = walker.clone()
				}
				nextWalker.moveTo(validMove)
				nextWalkers = append(nextWalkers, nextWalker)
			}
		}
		runningWalkers = nextWalkers
	}
	return len(s.completeWalkers)
}
