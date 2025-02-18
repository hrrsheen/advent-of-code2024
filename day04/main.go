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

func NewGrid(width int, height int) Grid {
	grid := Grid{width: width, height: height}
	grid.contents = make([]rune, width*height)

	return grid
}

func NewKernel(width int, height int, filter []rune) Grid {
	grid := NewGrid(width, height)
	grid.contents = filter

	return grid
}

func ToIndex(x int, y int, grid *Grid) int {
	return y*grid.width + x
}

func ToCoords(index int, grid *Grid) (int, int) {
	return index % grid.width, index / grid.width
}

func InBounds(x int, y int, grid *Grid) bool {
	return x >= 0 && x < grid.width && y >= 0 && y < grid.height
}

func SearchNext(x int, y int, xdir int, ydir int, prev rune, grid *Grid) int {
	if !InBounds(x+xdir, y+ydir, grid) {
		return 0
	}

	index := ToIndex(x+xdir, y+ydir, grid)
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
	character := grid.contents[ToIndex(x, y, grid)]

	count := 0
	if character == 'X' {
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

func SearchGrid(grid *Grid) int {
	numFound := 0
	for i := 0; i < grid.width*grid.height; i++ {
		x, y := ToCoords(i, grid)
		numFound += Search(x, y, grid)
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
