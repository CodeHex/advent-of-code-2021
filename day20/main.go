package main

import (
	"adventofcode2021/pkg/bits"
	"adventofcode2021/pkg/fileparser"
	"adventofcode2021/pkg/matrices"
	"fmt"
)

func main() {
	lines := fileparser.ReadLines("day20/sample.txt")

	enhancer := NewEnhancer(lines)
	for i := 1; i <= 50; i++ {
		enhancer.Enhance()
		if i == 2 {
			enhancer.PrintField()
			fmt.Println("[Part 1] Number of pixels after 2 enhances:", enhancer.CountPixels())
		}
	}
	fmt.Println("[Part 2] Number of pixels after 50 enhances:", enhancer.CountPixels())
}

type Enhancer struct {
	alg                     map[string]string
	darkMapPixel            string
	flipDarkMapAfterEnhance bool
	trench                  matrices.Matrix[string]
}

func NewEnhancer(lines []string) *Enhancer {
	alg := newAlg(lines[0])
	return &Enhancer{
		alg:                     alg,
		darkMapPixel:            ".",                     // Tile to use when extending the map
		flipDarkMapAfterEnhance: alg["........."] == "#", // Indicates that after an enhance, a pixel woth all dark pixels will become bright
		trench:                  fileparser.ReadCharMatrixFromLines[string](lines[2:]),
	}
}

func newAlg(data string) map[string]string {
	if len(data) != 512 {
		panic("unexpected algorithm length")
	}
	alg := make(map[string]string)
	for i := 0; i < 512; i++ {
		// Encode the number into a 9 bit value
		b := bits.NewBitFieldForVal(uint64(i), 9)
		key := ""
		for i := 0; i < b.Length; i++ {
			pixelVal := "."
			if b.Get(i) {
				pixelVal = "#"
			}
			key += pixelVal
		}
		alg[key] = string(data[i])
	}
	return alg
}

func (e *Enhancer) Enhance() {
	// Expand matrix by one row/column on all sides
	e.trench = e.trench.Expand(1, 1, 1, 1, e.darkMapPixel)
	result := matrices.NewMatrix[string](e.trench.Rows, e.trench.Columns)

	// Calculate algorithm key based on all neighbours
	e.trench.ForEach(func(x, y int, value string) {
		key := e.SurroundingPixelsKey(x, y)
		result.Set(x, y, e.alg[key])
	})
	e.trench = result

	// Flip the pixel to use for out of bounds if the algorithm key flips them
	if e.flipDarkMapAfterEnhance {
		if e.darkMapPixel == "#" {
			e.darkMapPixel = "."
		} else {
			e.darkMapPixel = "#"
		}
	}
}

func (e *Enhancer) SurroundingPixelsKey(x, y int) string {
	return e.Pixel(x-1, y-1) +
		e.Pixel(x, y-1) +
		e.Pixel(x+1, y-1) +
		e.Pixel(x-1, y) +
		e.Pixel(x, y) +
		e.Pixel(x+1, y) +
		e.Pixel(x-1, y+1) +
		e.Pixel(x, y+1) +
		e.Pixel(x+1, y+1)
}

func (e *Enhancer) Pixel(x, y int) string {
	if x < 0 || y < 0 || x > e.trench.Columns-1 || y > e.trench.Rows-1 {
		return e.darkMapPixel
	}
	return e.trench.Get(x, y)
}

func (e *Enhancer) CountPixels() int {
	count := 0
	e.trench.ForEach(func(x, y int, value string) {
		if value == "#" {
			count++
		}
	})
	return count
}

func (e *Enhancer) PrintField() {
	for j := 0; j < e.trench.Rows; j++ {
		for i := 0; i < e.trench.Columns; i++ {
			fmt.Printf(e.trench.Get(i, j))
		}
		fmt.Println()
	}
	fmt.Println()
}
