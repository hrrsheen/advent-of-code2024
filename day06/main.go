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

/**
 * Functions to query the grid.
 */

func ToIndex(pos Vec, grid *Grid) int {
	return pos.y*grid.width + pos.x
}

func InBounds(pos Vec, grid *Grid) bool {
	return pos.x >= 0 && pos.x < grid.width && pos.y >= 0 && pos.y < grid.height
}

/**
 * Functions to manipulate vectors.
 */

func AddVec(a Vec, b Vec) Vec {
	return Vec{a.x + b.x, a.y + b.y}
}

func EqualVec(a Vec, b Vec) bool {
	return a.x == b.x && a.y == b.y
}

/**
 * Functions for the problem.
 */

func AvoidObstruction(guard *Guard, grid *Grid) bool {
	stepAhead := AddVec(guard.pos, guard.dir)
	if !InBounds(stepAhead, grid) {
		return false
	}

	if grid.contents[ToIndex(stepAhead, grid)] == '#' {
		// Rotate the guard 90 degrees CW if an obstruction is directly ahead.
		guard.dir = Vec{-guard.dir.y, guard.dir.x}
		return true
	}

	return false
}

func Step(guard *Guard, grid *Grid) bool {
	for AvoidObstruction(guard, grid) {
	}

	guard.pos = AddVec(guard.pos, guard.dir)
	return InBounds(guard.pos, grid)
}

func TestLoop(newObstruction Vec, guard *Guard, grid *Grid) bool {
	if !InBounds(newObstruction, grid) {
		return false
	}

	virtualGuard := *guard           // The simulated guard
	pathHistory := make(map[Vec]Vec) // Records the position the guard was last facing at each position.

	obstructionIndex := ToIndex(newObstruction, grid)
	looping := false
	// Place the obstruction for the similation.
	grid.contents[obstructionIndex] = '#'
	for Step(&virtualGuard, grid) {
		if EqualVec(pathHistory[virtualGuard.pos], virtualGuard.dir) {
			// Loops are achieved when the guard reaches any given point and faces the same direction as the
			// last time it was at that point.
			looping = true
			break
		}
		pathHistory[virtualGuard.pos] = virtualGuard.dir
	}

	// Remove the obstruction now that the simulation is over
	grid.contents[obstructionIndex] = '.'

	return looping
}

func WalkPatrol(guard *Guard, grid *Grid) (int, int) {
	coveredTiles := make(map[Vec]bool) // A set of all tiles that the guard has visited.

	loopsFormed := 0
	for InBounds(guard.pos, grid) {
		coveredTiles[guard.pos] = true
		// Turn untill the path ahead of the guard is empty.
		for AvoidObstruction(guard, grid) {
		}

		nextPosition := AddVec(guard.dir, guard.pos)
		// Place an obstruction in front of the guard and simulate the new path
		// to determine if a loop forms.
		// We don't do this for tiles that the guard has already patroled though. They might notice!
		if !coveredTiles[nextPosition] && TestLoop(nextPosition, guard, grid) {
			loopsFormed++
		}

		// Move the guard to the next position.
		guard.pos = nextPosition
	}

	return loopsFormed, len(coveredTiles)
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

	loopsFound, tilesCovered := WalkPatrol(&guard, &grid)

	fmt.Printf("Tiles covered: %d\n", tilesCovered)
	fmt.Printf("  Loops found: %d\n", loopsFound)
}
