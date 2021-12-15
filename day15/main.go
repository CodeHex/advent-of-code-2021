package main

import (
	"adventofcode2021/pkg/fileparser"
	"adventofcode2021/pkg/maps"
	"adventofcode2021/pkg/matrices"
	"fmt"
	"math"
)

func main() {
	data := fileparser.ReadDigitMatrix("day15/input.txt")

	solverPart1 := NewSolver(data, 1)
	fmt.Println("[Part 1] Shortest distance for map is", solverPart1.Solve())

	solverPart2 := NewSolver(data, 5)
	fmt.Println("[Part 2] Shortest distance for 5 by 5 map is", solverPart2.Solve())
}

type Pos struct{ x, y int }

type Node struct {
	pos      Pos
	risk     int
	prevNode *Node
	distance int
}

type PathSolver struct {
	source    *Node
	target    *Node
	nodeList  map[Pos]*Node
	unvisited map[Pos]*Node
	updated   map[Pos]*Node
	matrixRef matrices.Matrix[struct{}]
}

func NewSolver(data matrices.IntMatrix[int], repeat int) *PathSolver {
	solver := &PathSolver{
		nodeList:  make(map[Pos]*Node),
		unvisited: make(map[Pos]*Node),
		updated:   make(map[Pos]*Node),
	}
	shift := data.Rows
	data.ForEach(func(x, y, value int) {
		// Duplicate the tile based on the repeat factor
		for i := 0; i < repeat; i++ {
			for j := 0; j < repeat; j++ {
				// Ensure value is wrapped to the correct range
				newValue := value + i + j
				for newValue > 9 {
					newValue -= 9
				}
				newNode := &Node{
					pos:      Pos{x + (i * shift), y + (j * shift)},
					risk:     newValue,
					distance: math.MaxInt,
				}
				solver.nodeList[newNode.pos] = newNode
				solver.unvisited[newNode.pos] = newNode
			}
		}
	})
	solver.source = solver.nodeList[Pos{0, 0}]
	solver.source.distance = 0
	solver.updated[solver.source.pos] = solver.source
	solver.target = solver.nodeList[Pos{(data.Columns * repeat) - 1, (data.Rows * repeat) - 1}]

	// Create a blank matrix of the map to enable calculating nearest neighbours coord
	solver.matrixRef = matrices.NewMatrix[struct{}](data.Rows*repeat, data.Columns*repeat)
	return solver
}

func (s *PathSolver) Solve() int {
	for len(s.unvisited) > 0 {
		// Find the next node to progress
		nextPos, nextNode := maps.MinMappedValue(s.updated, func(n *Node) int { return n.distance })
		delete(s.unvisited, nextPos)
		delete(s.updated, nextPos)

		// If we are on the target node, we have explored everything we need and finished
		if nextPos == s.target.pos {
			break
		}

		s.matrixRef.ForEachNeighbour(false, nextPos.x, nextPos.y, func(x, y int) {
			lookAt := s.nodeList[Pos{x, y}]
			newDist := nextNode.distance + lookAt.risk
			if newDist < lookAt.distance {
				lookAt.distance = newDist
				lookAt.prevNode = nextNode
				s.updated[lookAt.pos] = lookAt
			}
		})
	}
	return s.target.distance
}
