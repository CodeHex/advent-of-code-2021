package main

import (
	"adventofcode2021/pkg/fileparser"
	"fmt"
)

type Board struct {
	numbers [5][5]int
	marked  [5][5]bool

	Won           bool
	winningNumber int
	unmarkedTotal int
	score         int
}

func NewBoard(input []string) *Board {
	var result [5][5]int
	for y, line := range input {
		values := fileparser.SplitTrim[int](line, " ")
		for x, val := range values {
			result[y][x] = val
		}
	}
	return &Board{numbers: result}
}

// MarkNumber checks a called number against the board. If the number is present,
// the number is marked and we check if this has led to Bingo
func (b *Board) MarkNumber(num int) {
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			if b.numbers[y][x] == num {
				b.marked[y][x] = true
				b.bingoCheck(num, x, y)
				return
			}
		}
	}

}

// bingoCheck will check if the called number triggered Bingo. If so, the
// board is placed into a winning state and scores calculated
func (b *Board) bingoCheck(calledNumber int, xCalled, yCalled int) {
	// Check if all numbers in the column are now marked
	bingo := true
	for y := 0; y < 5; y++ {
		if !b.marked[y][xCalled] {
			bingo = false
			break
		}
	}

	// If we haven't found bingo yet, check all numbers in the row
	if !bingo {
		bingo = true
		for x := 0; x < 5; x++ {
			if !b.marked[yCalled][x] {
				bingo = false
				break
			}
		}
	}

	// If we have found bingo, mark the board as won
	if bingo {
		b.completeBoard(calledNumber)
	}
}

// calcUnmarkedTotal scans through the board and sums all numbers that
// are not marked
func (b *Board) calcUnmarkedTotal() int {
	total := 0
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			if !b.marked[y][x] {
				total += b.numbers[y][x]
			}
		}
	}
	return total
}

func (b *Board) completeBoard(called int) {
	b.Won = true
	b.winningNumber = called
	b.unmarkedTotal = b.calcUnmarkedTotal()
	b.score = b.winningNumber * b.unmarkedTotal
}

func (b *Board) Print() {
	colorReset := "\033[0m"
	colorBold := "\033[1m"
	for y, dataline := range b.numbers {
		for x, val := range dataline {
			if b.marked[y][x] {
				fmt.Printf("%s%2d%s ", colorBold, val, colorReset)
			} else {
				fmt.Printf("%2d ", val)
			}
		}
		fmt.Println()
	}

	if !b.Won {
		fmt.Println("board is still in play")
	} else {
		fmt.Printf("winning called number %d, unmarked total %d, score: %d\n", b.winningNumber, b.unmarkedTotal, b.score)
	}
}
