package main

import (
	"adventofcode2021/pkg/bits"
	"adventofcode2021/pkg/fileparser"
	"adventofcode2021/pkg/matrices"
	"fmt"
)

func main() {
	lines := fileparser.ReadLines("day20/input.txt")

	alg := make(map[string]string)
	for i := 0; i < 512; i++ {
		b := bits.NewBitFieldForVal(uint64(i), 9)
		out := ""
		for i := 0; i < 9; i++ {
			if b.Get(i) {
				out += "#"
			} else {
				out += "."
			}
		}
		alg[out] = string(lines[0][i])
	}

	field := fileparser.ReadCharMatrixFromLines[string](lines[2:])

	flip := alg["........."] == "#"
	def := "."
	for i := 0; i < 50; i++ {
		field = field.Expand(2)
		field.ForEach(func(x, y int, value string) {
			if value == "" {
				field.Set(x, y, def)
			}
		})
		field = Enhance(field, alg, def)
		if flip {

			if def == "." {
				def = "#"
			} else {
				def = "."
			}
		}
	}

	fmt.Println(CountPixels(field))
}

func Enhance(f matrices.Matrix[string], alg map[string]string, d string) matrices.Matrix[string] {
	result := matrices.NewMatrix[string](f.Rows, f.Columns)
	f.ForEach(func(x, y int, value string) {
		out := Pixel(d, f, x-1, y-1) + Pixel(d, f, x, y-1) + Pixel(d, f, x+1, y-1) +
			Pixel(d, f, x-1, y) + Pixel(d, f, x, y) + Pixel(d, f, x+1, y) +
			Pixel(d, f, x-1, y+1) + Pixel(d, f, x, y+1) + Pixel(d, f, x+1, y+1)
		result.Set(x, y, alg[out])
	})
	return result
}

func Pixel(d string, field matrices.Matrix[string], x, y int) string {
	if x < 0 || y < 0 || x > field.Columns-1 || y > field.Rows-1 {
		return d
	}
	return field.Get(x, y)
}

func CountPixels(field matrices.Matrix[string]) int {
	count := 0
	field.ForEach(func(x, y int, value string) {
		if value == "#" {
			count++
		}
	})
	return count
}

func PrintField(field matrices.Matrix[string]) {
	for j := 0; j < field.Rows; j++ {
		for i := 0; i < field.Columns; i++ {
			fmt.Printf(field.Get(i, j))
		}
		fmt.Println()
	}
	fmt.Println()
}
