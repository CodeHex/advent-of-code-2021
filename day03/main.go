package main

import (
	"adventofcode2021/pkg/bits"
	"adventofcode2021/pkg/fileparser"
	"fmt"
)

func main() {
	var readings bits.BitFieldArray = fileparser.ReadSingles[bits.BitField]("day03/input.txt")

	gammaRate := readings.MostCommon()
	epsilionRate := gammaRate.Invert()
	powerConsumption := gammaRate.Value * epsilionRate.Value

	fmt.Printf("gamma rate is %d [%s], epsilion rate is %d [%s], power consumption is %d\n",
		gammaRate.Value, gammaRate,
		epsilionRate.Value, epsilionRate,
		powerConsumption)

	oxyGenRating := readings.ReduceToRating(true)
	co2ScubberRating := readings.ReduceToRating(false)
	lifeSupportRating := oxyGenRating.Value * co2ScubberRating.Value

	fmt.Printf("oxygen generator rating is %d [%s], CO2 scrubber rating is %d [%s], life support rating is %d\n",
		oxyGenRating.Value, oxyGenRating,
		co2ScubberRating.Value, co2ScubberRating,
		lifeSupportRating)
}
