package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Grid struct {
	contents []rune
	width    int
	height   int
}

func (grid *Grid) ToIndex(x int, y int) int {
	return y*grid.width + x
}

func (grid *Grid) InBounds(x int, y int) bool {
	return x >= 0 && x < grid.width && y >= 0 && y < grid.height
}

/**
 * Recursively searches for the remainder of the word "XMAS" by looking in the given direction.
 */
func SearchNext(x int, y int, xdir int, ydir int, prev rune, grid *Grid) int {
	if !grid.InBounds(x+xdir, y+ydir) {
		return 0
	}

	index := grid.ToIndex(x+xdir, y+ydir)
	character := grid.contents[index]
	if (character == 'M' && prev == 'X') || (character == 'A' && prev == 'M') {
		return SearchNext(x+xdir, y+ydir, xdir, ydir, character, grid)
	}

	if character == 'S' && prev == 'A' {
		return 1
	}

	return 0
}

func Search(x int, y int, grid *Grid) int {
	character := grid.contents[grid.ToIndex(x, y)]

	count := 0
	if character == 'X' {
		// An "XMAS" may appear in any of the 8 directions, outwardly from the 'X'.
		count += SearchNext(x, y, 1, 0, 'X', grid)
		count += SearchNext(x, y, 0, 1, 'X', grid)
		count += SearchNext(x, y, 1, 1, 'X', grid)
		count += SearchNext(x, y, 0, -1, 'X', grid)
		count += SearchNext(x, y, -1, 0, 'X', grid)
		count += SearchNext(x, y, 1, -1, 'X', grid)
		count += SearchNext(x, y, -1, 1, 'X', grid)
		count += SearchNext(x, y, -1, -1, 'X', grid)
	}

	return count
}

/**
 * Assuming that the point at (x, y) contains an 'A'. Check the diagonally-adjacent cells for
 * instances of the word "MAS". Returns the number of times it occurs
 */
func CheckMAS(x int, y int, grid *Grid) int {
	if grid.contents[grid.ToIndex(x, y)] != 'A' {
		return 0
	}

	above := grid.contents[grid.ToIndex(x-1, y-1)]
	below := grid.contents[grid.ToIndex(x+1, y+1)]

	if !((above == 'M' && below == 'S') || (above == 'S' && below == 'M')) {
		return 0
	}

	above = grid.contents[grid.ToIndex(x+1, y-1)]
	below = grid.contents[grid.ToIndex(x-1, y+1)]

	if (above == 'M' && below == 'S') || (above == 'S' && below == 'M') {
		return 1
	}

	return 0
}

func SearchGrid(grid *Grid) int {
	// By changing these variables and the SearchFunc, I can easily switch between
	// the words I'm searching for.
	var (
		xStart int = 1
		xEnd   int = grid.width - 1
		yStart int = 1
		yEnd   int = grid.height - 1
	)

	SearchFunc := CheckMAS

	numFound := 0
	for y := yStart; y < yEnd; y++ {
		for x := xStart; x < xEnd; x++ {
			numFound += SearchFunc(x, y, grid)
		}
	}

	return numFound
}

func PopulateGridFromReader(r *bufio.Reader) Grid {
	grid := Grid{contents: make([]rune, 0, 128)}

	width := 0
	height := 0
	chCount := 0

	for {
		ch, _, err := r.ReadRune()
		ch = rune(ch)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		if ch == '\n' {
			width = chCount
			height++
			continue
		}

		grid.contents = append(grid.contents, ch)
		// Once the grid's width is locked in, we don't need to keep track of the character count.
		if width == 0 {
			chCount++
		}
	}

	grid.width = width
	grid.height = height

	return grid
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	grid := PopulateGridFromReader(reader)

	count := SearchGrid(&grid)

	fmt.Printf("XMAS found: %d\n", count)
}
