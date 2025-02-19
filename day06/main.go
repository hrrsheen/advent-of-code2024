package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
)

type Grid struct {
	contents []rune
	width    int
	height   int
}

type Vec struct {
	x int
	y int
}

type Guard struct {
	pos Vec
	dir Vec
}

func ToIndex(pos Vec, grid *Grid) int {
	return pos.y*grid.width + pos.x
}

func InBounds(pos Vec, grid *Grid) bool {
	return pos.x >= 0 && pos.x < grid.width && pos.y >= 0 && pos.y < grid.height
}

func AddVec(a Vec, b Vec) Vec {
	return Vec{a.x + b.x, a.y + b.y}
}

func Step(guard *Guard, grid *Grid) bool {
	nextPos := AddVec(guard.pos, guard.dir)
	if !InBounds(nextPos, grid) {
		return false
	}

	if grid.contents[ToIndex(nextPos, grid)] == '#' {
		// Rotate the direction 90 degrees.
		guard.dir = Vec{-guard.dir.y, guard.dir.x}
		return Step(guard, grid)
	}

	guard.pos = nextPos
	return true
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

	var nonGuardTiles = []rune{'.', '#', '\n'}
	var guard Guard
	// Search for the guard to find the initial direction.
	for i, c := range grid.contents {
		if !slices.Contains(nonGuardTiles, c) {
			guard.pos = Vec{i % grid.width, i / grid.width}
			switch c {
			case '^':
				guard.dir = Vec{0, -1}
			case '>':
				guard.dir = Vec{1, 0}
			case 'V':
				guard.dir = Vec{0, 1}
			case '<':
				guard.dir = Vec{-1, 0}
			}
		}
	}

	coveredTiles := make(map[Vec]bool)
	coveredTiles[guard.pos] = true
	for Step(&guard, &grid) {
		coveredTiles[guard.pos] = true
	}

	fmt.Printf("Tiles covered: %d\n", len(coveredTiles))
}
