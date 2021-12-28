package main

import (
	"adventofcode2021/pkg/fileparser"
	"adventofcode2021/pkg/matrices"
	"fmt"
)

type Pos struct{ x, y int }

func main() {
	seabed := fileparser.ReadCharMatrix[string]("day25/input.txt")
	count := 0
	for {
		count++
		movedRight := MoveCucumbersRight(seabed)
		movedDown := MoveCucumbersDown(seabed)
		if !movedRight && !movedDown {
			break
		}
	}
	fmt.Println("[Part 1] Steps until cucumbers stop moving:", count)
}

func MoveCucumbers(icon string, seabed matrices.Matrix[string], move func(x Pos) Pos) bool {
	moved := false
	moves := make(map[Pos]Pos)
	seabed.ForEach(func(x, y int, value string) {
		if value != icon {
			return
		}

		newPos := move(Pos{x, y})
		if seabed.Get(newPos.x, newPos.y) == "." {
			moves[Pos{x, y}] = newPos
			moved = true
		}
	})

	for start, end := range moves {
		seabed.Set(start.x, start.y, ".")
		seabed.Set(end.x, end.y, icon)
	}
	return moved
}

func MoveCucumbersRight(seabed matrices.Matrix[string]) bool {
	moveFunc := func(p Pos) Pos {
		newX := p.x + 1
		if newX >= seabed.Columns {
			newX = 0
		}
		return Pos{newX, p.y}
	}
	return MoveCucumbers(">", seabed, moveFunc)
}

func MoveCucumbersDown(seabed matrices.Matrix[string]) bool {
	moveFunc := func(p Pos) Pos {
		newY := p.y + 1
		if newY >= seabed.Rows {
			newY = 0
		}
		return Pos{p.x, newY}
	}
	return MoveCucumbers("v", seabed, moveFunc)
}
