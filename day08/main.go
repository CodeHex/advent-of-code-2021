package main

import (
	"adventofcode2021/pkg/fileparser"
	"adventofcode2021/pkg/maps"
	"adventofcode2021/pkg/slices"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func main() {
	entries := fileparser.ReadTypedLines("day08/input.txt", NewEntry)
	counts := slices.Map(entries, func(e *Entry) int { return e.CountUniqueOutput() })

	total := slices.Sum(counts)
	fmt.Printf("[Part 1] Total number of 1, 4, 7 and 8 in outputs : %d\n", total)

	totalCodes := slices.SumWeighted(entries, func(e *Entry) int { return e.outputCode })
	fmt.Printf("[Part 2] Total sum of output code : %d\n", totalCodes)
}

type Entry struct {
	signals       []string
	outputs       []string
	digitToSignal map[int]string
	signalToDigit map[string]int
	outputCode    int
}

func NewEntry(line string) *Entry {
	parts := strings.Split(line, "|")
	e := &Entry{
		signals:       fileparser.SplitTrim[string](parts[0], " "),
		outputs:       fileparser.SplitTrim[string](parts[1], " "),
		digitToSignal: make(map[int]string),
		signalToDigit: make(map[string]int),
	}
	e.signals = slices.Map(e.signals, orderString)
	e.outputs = slices.Map(e.outputs, orderString)
	e.solveMap()
	return e
}

func orderString(in string) string {
	copy := []rune(in)
	sort.Slice(copy, func(i, j int) bool {
		return copy[i] < copy[j]
	})
	return string(copy)
}

func (e *Entry) CountUniqueOutput() int {
	uniqueLengths := []int{2, 3, 4, 7}
	criteria := func(s string) bool {
		return slices.Contains(uniqueLengths, len(s))
	}
	return slices.CountIf(e.outputs, criteria)
}

func (e *Entry) recordMapping(digit int, signal string) {
	e.digitToSignal[digit] = signal
	e.signalToDigit[signal] = digit
}

// solveMap creates a map between the signal values and the numbers they represent.
func (e *Entry) solveMap() {
	// Step 1 : First identify 1, 4, 7, 8 by segment length
	e.recordMapping(1, slices.FirstOrDefault(e.signals, LengthCheckFunc(2))) // Digit 1 only has 2 segments
	e.recordMapping(4, slices.FirstOrDefault(e.signals, LengthCheckFunc(4))) // Digit 4 only has 4 segments
	e.recordMapping(7, slices.FirstOrDefault(e.signals, LengthCheckFunc(3))) // Digit 7 only has 3 segments
	e.recordMapping(8, slices.FirstOrDefault(e.signals, LengthCheckFunc(7))) // Digit 8 only has 8 segments

	// Step 2 : Removing segments in 1 from unknown numbers indentifies 3 and 6 from the remaining numbers
	// 0 - leaves 4 segments (2 segments removed) (not unique)
	// 2 - leaves 4 segments (1 segments removed) (not unique)
	// 3 - leaves 3 segments (2 segments removed)
	// 5 - leaves 4 segments (1 segments removed) (not unique)
	// 6 - leaves 5 segments (1 segments removed)
	// 9 - leaves 4 segments (2 segments removed) (not unique)
	e.recordMapping(3, slices.FirstOrDefault(e.signals, e.matchCheckFunc(1, 2, 3)))
	e.recordMapping(6, slices.FirstOrDefault(e.signals, e.matchCheckFunc(1, 1, 5)))

	// Step 3 : Removing segments in 4 from unknown numbers identifies remaining digits
	// 0 - leaves 3 segments (3 segments removed)
	// 2 - leaves 3 segments (2 segments removed)
	// 5 - leaves 2 segments (3 segments removed)
	// 9 - leaves 2 segments (4 segments removed)
	e.recordMapping(0, slices.FirstOrDefault(e.signals, e.matchCheckFunc(4, 3, 3)))
	e.recordMapping(2, slices.FirstOrDefault(e.signals, e.matchCheckFunc(4, 2, 3)))
	e.recordMapping(5, slices.FirstOrDefault(e.signals, e.matchCheckFunc(4, 3, 2)))
	e.recordMapping(9, slices.FirstOrDefault(e.signals, e.matchCheckFunc(4, 4, 2)))
	e.calculateOutputCode()
}

func (e *Entry) calculateOutputCode() {
	outputCodeStr := ""
	for _, outputEntry := range e.outputs {
		outputCodeStr += fmt.Sprintf("%d", e.signalToDigit[outputEntry])
	}
	e.outputCode, _ = strconv.Atoi(outputCodeStr)
}

// matchCheckFunc generates a function that will match a signal if
// - we haven't already matched the signal to a number AND
// - when we remove the segments in the provided digit
// -  - the number of expect segements removed match AND
// -  - the number of expect segements remaining match
func (e *Entry) matchCheckFunc(digit int, expRemoved int, expRemaing int) func(s string) bool {
	return func(s string) bool {
		// Don't match if already matched
		if maps.ContainsKey(e.signalToDigit, s) {
			return false
		}

		removed, remaining := slices.Divide([]rune(s), func(r rune) bool {
			return slices.Contains([]rune(e.digitToSignal[digit]), r)
		})
		return len(removed) == expRemoved && len(remaining) == expRemaing
	}
}

// LengthCheckFunc generates a function that will match a signal if the signal has
// the provided length
func LengthCheckFunc(length int) func(s string) bool {
	return func(d string) bool { return len(d) == length }
}
