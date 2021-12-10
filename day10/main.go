package main

import (
	"adventofcode2021/pkg/fileparser"
	"adventofcode2021/pkg/slices"
	"fmt"
	"sort"
)

func main() {
	navResults := fileparser.ReadTypedLines("day10/input.txt", NewNavResult)
	resultCount := len(navResults)
	incomplete := slices.Filter(navResults, IsIncompleteFunc())
	corrupted := slices.Filter(navResults, IsCorruptFunc())
	fmt.Printf("Number of lines %d, incomplete: %d, corrupted %d\n", resultCount, len(incomplete), len(corrupted))

	syntaxScore := slices.SumWeighted(corrupted, SyntaxScoreFunc())
	fmt.Printf("[Part 1] Syntax error score %d\n", syntaxScore)

	autoCompleteScores := slices.Map(incomplete, AutocompleteScoreFunc())
	sort.Ints(autoCompleteScores)
	middleIndex := ((len(autoCompleteScores) + 1) / 2) - 1
	fmt.Printf("[Part 2] Autocomplete score %d\n", autoCompleteScores[middleIndex])
}

type NavResult struct {
	IsIncomplete bool
	IsCorrupt    bool
	completeSeq  string // If incomplete, represents the sequence that would complete it
	corruptChar  rune   // If corrupt, represents the corrupt char that doesn't match the last opening char
}

func IsIncompleteFunc() func(NavResult) bool {
	return func(x NavResult) bool {
		return x.IsIncomplete
	}
}

func IsCorruptFunc() func(NavResult) bool {
	return func(x NavResult) bool {
		return x.IsCorrupt
	}
}
func NewNavResult(line string) NavResult {
	openChars := []rune{'(', '[', '{', '<'}
	closeChars := []rune{')', ']', '}', '>'}
	matchChars := map[rune]rune{'(': ')', '[': ']', '{': '}', '<': '>'}
	result := NavResult{}
	currentChunks := []rune{}

	for _, c := range line {
		switch {
		case slices.Contains(openChars, c):
			// Start another chunk
			currentChunks = append(currentChunks, c)
		case slices.Contains(closeChars, c):
			currentChunk := slices.Last(currentChunks)
			closingForCurrentChunk := matchChars[currentChunk]
			if c == closingForCurrentChunk {
				// Close chunk (by removing it from the end) if it matches the current chunk
				currentChunks = slices.TrimEnd(currentChunks, 1)
			} else {
				// If it doesn't match, the line must be corrupt
				result.IsCorrupt = true
				result.corruptChar = c
				return result
			}
		}
	}

	// Calculate the chars required to close the remaining open chunks by
	// swapping them with their corresponding closing chars and reversing
	// the order
	completeChars := slices.Map(currentChunks, func(x rune) rune { return matchChars[x] })
	completeChars = slices.Reverse(completeChars)

	// If there are any completing characters, then we must have not completed the line
	result.IsIncomplete = len(completeChars) != 0
	result.completeSeq = string(completeChars)
	return result
}

func (n NavResult) SyntaxScore() int {
	if !n.IsCorrupt {
		panic("expected result to be corrupt to provide syntax score")
	}
	syntaxScoreTable := map[rune]int{
		')': 3,
		']': 57,
		'}': 1197,
		'>': 25137,
	}
	return syntaxScoreTable[n.corruptChar]
}

func SyntaxScoreFunc() func(NavResult) int {
	return func(x NavResult) int {
		return x.SyntaxScore()
	}
}

func (n NavResult) AutocompleteScore() int {
	if !n.IsIncomplete {
		panic("expected result to be incomplete to provide auto correct score")
	}
	autocompleteScoreTable := map[rune]int{
		')': 1,
		']': 2,
		'}': 3,
		'>': 4,
	}
	score := 0
	for _, c := range n.completeSeq {
		score = score * 5
		score += autocompleteScoreTable[c]
	}
	return score
}

func AutocompleteScoreFunc() func(NavResult) int {
	return func(x NavResult) int {
		return x.AutocompleteScore()
	}
}
