package main

import (
	"adventofcode2021/pkg/fileparser"
	"adventofcode2021/pkg/slices"
	"fmt"
)

func main() {
	input := fileparser.ReadLines("day04/input.txt")

	numbersToCall := fileparser.Split[int](input[0], ",")

	// Get the data for each board, starting at the beginning of the first board,
	// reading the 5 lines that represent the board and then moving to the beginning
	// of the next board. Each board is stored in 6 lines (5 lines of numbers and 1 blank line)
	boards := []*Board{}
	for line := 2; line < len(input); line += 6 {
		nextBoard := NewBoard(input[line : line+5])
		boards = append(boards, nextBoard)
	}

	// Run through each round checking each called number against each board
	// After every round, only progress boards that have not won and are still playing.
	boardsInPlay := boards
	completedBoards := []*Board{}
	completedBoardsThisRound := []*Board{}
	isCompletedFunc := func(b *Board) bool { return b.Won }

	for _, calledNumber := range numbersToCall {
		for _, board := range boardsInPlay {
			board.MarkNumber(calledNumber)
		}
		completedBoardsThisRound, boardsInPlay = slices.Divide(boardsInPlay, isCompletedFunc)
		completedBoards = append(completedBoards, completedBoardsThisRound...)
	}

	// Print results
	if len(boardsInPlay) == 0 {
		fmt.Printf("all boards completed (%d boards)\n\n", len(completedBoards))
	} else {
		fmt.Printf("%d boards completed (out of %d boards)\n\n", len(completedBoards), len(boards))
	}

	fmt.Println("first board to win is ")
	completedBoards[0].Print()

	fmt.Println()
	fmt.Println("last board to win is")
	completedBoards[len(completedBoards)-1].Print()
}
