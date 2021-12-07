package main

import (
	"adventofcode2021/pkg/fileparser"
	"adventofcode2021/pkg/maps"
	"adventofcode2021/pkg/slices"
	"constraints"
	"fmt"
)

func main() {
	crabPositions := fileparser.ReadCSVLine[int]("day07/input.txt")
	min, max := slices.MinMax(crabPositions)
	fuelCalc := NewFuelCalculator()

	basicAttempts := make(map[int]int)
	advAttempts := make(map[int]int)

	// Loop through each possible alignment position calculating the total fuel used
	for alignAttempt := min; alignAttempt <= max; alignAttempt++ {
		basicAttempts[alignAttempt] = slices.SumWeighted(crabPositions, fuelCalc.BasicFuelCostFunc(alignAttempt))
		advAttempts[alignAttempt] = slices.SumWeighted(crabPositions, fuelCalc.AdvancedFuelCostFunc(alignAttempt))
	}

	pos1, fuel1 := maps.MinValue(basicAttempts)
	pos2, fuel2 := maps.MinValue(advAttempts)
	fmt.Printf("[Part 1] Least fuel is at position %d, using %d fuel\n", pos1, fuel1)
	fmt.Printf("[Part 2] Least fuel is at position %d, using %d fuel\n", pos2, fuel2)
}

type FuelCalculator struct {
	resultsCache       map[int]int
	maxResultsCacheVal int
}

func NewFuelCalculator() *FuelCalculator {
	return &FuelCalculator{resultsCache: make(map[int]int)}
}

func (f *FuelCalculator) BasicFuelCost(start, end int) int {
	return Distance(start, end)
}

func (f *FuelCalculator) BasicFuelCostFunc(end int) func(start int) int {
	return func(start int) int { return f.BasicFuelCost(start, end) }
}

func (f *FuelCalculator) AdvancedFuelCost(start, end int) int {
	dist := Distance(start, end)

	// Check if we have already calculated the fuel cost for this distance.
	// If not calculate all values from the biggest value in the cache to the final distance
	// since it uses these values as intermediate steps for the final fuel cost
	if dist > f.maxResultsCacheVal {
		for i := f.maxResultsCacheVal; i <= dist; i++ {
			f.resultsCache[i] = f.resultsCache[i-1] + i
		}
		f.maxResultsCacheVal = dist
	}
	return f.resultsCache[dist]
}

func (f *FuelCalculator) AdvancedFuelCostFunc(end int) func(start int) int {
	return func(start int) int { return f.AdvancedFuelCost(start, end) }
}

func Distance[T constraints.Integer](start, end T) T {
	dist := end - start
	if dist < 0 {
		dist = dist * (interface{})(-1).(T)
	}
	return dist
}
