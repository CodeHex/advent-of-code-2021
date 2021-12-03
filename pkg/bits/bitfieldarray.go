package bits

import (
	"adventofcode2021/pkg/slices"
	"bytes"
	"fmt"
)

type BitFieldArray []BitField

// MostCommon create a BitField where each flag in each postion represets the most common flag for that position
// across all entries in the array
func (b BitFieldArray) MostCommon() BitField {
	if len(b) == 0 {
		panic("No entries in array")
	}
	counts := make([]int, b[0].Length)

	// Count how many 1s are in each position
	for pos := 0; pos < len(counts); pos++ {
		for _, reading := range b {
			if reading.Get(pos) {
				counts[pos]++
			}
		}
	}

	// Create a bitfield where each 1 indicates that position had half or more
	var resultBytes bytes.Buffer
	for _, count := range counts {
		if 2*count >= len(b) {
			resultBytes.WriteString("1")
		} else {
			resultBytes.WriteString("0")
		}
	}

	result, err := NewBitField(resultBytes.String())
	if err != nil {
		panic(err)
	}
	return result
}

// FilterByPos reduces the bit field entries to only ones where the flag in the provided position matches
// either most common or least common flag in that position
func (b BitFieldArray) FilterByPos(pos int, useCommon bool) BitFieldArray {
	criteria := b.MostCommon()
	if !useCommon {
		criteria = criteria.Invert()
	}

	return slices.Filter(b, func(field BitField) bool {
		return field.Get(pos) == criteria.Get(pos)
	})
}

// ReduceToRating iterates through each position in the bit field and reduces the entries by position until
// there is only one result left
func (b BitFieldArray) ReduceToRating(useCommon bool) BitField {
	possibleResults := b
	for pos := 0; pos < b[0].Length; pos++ {
		possibleResults = possibleResults.FilterByPos(pos, useCommon)
		if val, ok := slices.IsSingle(possibleResults); ok {
			return val
		}
	}

	panic(fmt.Sprintf("reduction failed to find 1 result %v", possibleResults))
}
