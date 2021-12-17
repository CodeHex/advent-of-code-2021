package main

import (
	"adventofcode2021/pkg/maps"
	"fmt"
)

type Box struct {
	left, right, top, bottom int
}

type Velocity struct {
	vX, vY int
}

func main() {
	//target := Box{left: 20, right: 30, top: -5, bottom: -10}
	target := Box{left: 269, right: 292, top: -44, bottom: -68}

	// Consider y only.
	// When the projectile as on the opposite side of the trajectory at y=0,
	// the velocity will be -Vy. The step after that, y = -Vy - 1. The maximum, this
	// can be and still hit the target box would be at the bottom of the box i.e.
	// T_bottom = - Vy_max - 1
	// Vy_max = -T_bottom - 1
	maxVY := -target.bottom - 1

	// After each step, distance travel reduces by 1 so max height is sum of arithmetic sequence
	maxHeight := (maxVY * (maxVY + 1) / 2)

	// The minimum Vy can be is if we shoot the projectile straight down and hits the
	// bottom of the target in the next step
	minVY := target.bottom

	// Create a map of steps to possible Y velocities where at that step
	// the projectile would be in the target box
	possibleVY := map[int][]int{}
	for testVY := minVY; testVY <= maxVY; testVY++ {
		// Track the end position on every step
		endPosY := 0
		step := 0
		currentVY := testVY
		for {
			endPosY += currentVY
			step++
			currentVY--
			// If after the step to the Y coord is in the target, its a possible velocity
			if endPosY <= target.top && endPosY >= target.bottom {
				possibleVY[step] = append(possibleVY[step], testVY)
			}

			// Don't bother calculating if we're already past the bottom
			if endPosY < target.bottom {
				break
			}
		}
	}

	// Find the maximum number of steps where a velocity would result in hitting the target (y only)
	maxStepY := maps.MaxKey(possibleVY)

	// The max Vx can be is when reaching the right most of the target after 1 step
	maxVX := target.right

	// Now consider X only, try test velocities
	possibleVX := map[int][]int{}
	for testVX := 1; testVX <= maxVX; testVX++ {
		endPosX := 0
		step := 0
		currentVX := testVX
		for {
			endPosX += currentVX
			step++
			currentVX--

			// If after the step to the X coord is in the target, its a possible velocity
			if endPosX >= target.left && endPosX <= target.right {
				possibleVX[step] = append(possibleVX[step], testVX)
			}
			// Projectile has stopped, if we are in the zone, add this velocity to all further
			// steps
			if currentVX == 0 {
				if endPosX >= target.left && endPosX <= target.right {
					for s := maxStepY; s > step; s-- {
						possibleVX[s] = append(possibleVX[s], testVX)
					}
				}
				break
			}
		}
	}

	// Create a map of all possible velocities that the target on both the
	// x and y axis on the same step
	type velocity struct{ vX, vY int }
	totalCombos := make(map[Velocity]struct{})
	for t := 1; t <= maxStepY; t++ {
		for _, vx := range possibleVX[t] {
			for _, vy := range possibleVY[t] {
				totalCombos[Velocity{vx, vy}] = struct{}{}
			}
		}
	}

	fmt.Println("[Part 1] Max velocity for y:", maxVY, "height:", maxHeight)
	fmt.Println("[Part 2] Total initial velocities", len(totalCombos))
}
