package main

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"os"
	"strconv"
)

type Grid struct {
	contents []int
	width    int
	height   int
}

type Vector struct {
	x int
	y int
}

func (vec Vector) Add(other Vector) Vector {
	return Vector{vec.x + other.x, vec.y + other.y}
}

func (grid *Grid) InBounds(point Vector) bool {
	return point.x >= 0 && point.x < grid.width && point.y >= 0 && point.y < grid.height
}

func PopulateMapFromReader(r *bufio.Reader) (*Grid, []Vector) {
	grid := Grid{contents: make([]int, 0, 128), width: 0, height: 0}
	zeroes := make([]Vector, 0, 20)

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
			chCount = 0
			grid.height++
			continue
		} else if ch == '0' {
			zeroes = append(zeroes, Vector{chCount, grid.height})
		}

		num, err := strconv.Atoi(string(ch))
		if err != nil {
			panic(err)
		}
		grid.contents = append(grid.contents, num)

		// Once the grid's width is locked in, we don't need to keep track of the character count.
		if grid.height == 0 {
			grid.width++
		}
		chCount++
	}

	return &grid, zeroes
}

func (grid *Grid) Get(point Vector) int {
	return grid.contents[point.y*grid.width+point.x]
}

/**
 * Yields all adjacent cells that can be pathed to from the given point.
 */
func (mapData *Grid) Neighbours(point Vector) iter.Seq[Vector] {
	return func(yield func(Vector) bool) {
		checkDirn := []Vector{{0, -1}, {-1, 0}, {1, 0}, {0, 1}}
		for _, dir := range checkDirn {
			neighbour := point.Add(dir)
			if !mapData.InBounds(neighbour) {
				continue
			}

			canPath := (mapData.Get(neighbour) - mapData.Get(point)) == 1
			if canPath && !yield(neighbour) {
				return
			}
		}
	}
}

func FindPaths(zeroLoc Vector, mapData *Grid) (int, int) {
	frontier := make([]Vector, 0, 4) // Contains the cells that are queued for checking.
	ninesSet := make(map[Vector]bool)
	nines := 0

	frontier = append(frontier, zeroLoc)
	for len(frontier) != 0 {
		// Pop the first cell to check from the queue.
		current := frontier[0]
		frontier = frontier[1:]
		for next := range mapData.Neighbours(current) {
			frontier = append(frontier, next)
			if mapData.Get(next) == 9 {
				ninesSet[next] = true
				nines++
			}
		}
	}

	return len(ninesSet), nines
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	hikingMap, zeroes := PopulateMapFromReader(reader)

	totalScore := 0
	totalRating := 0
	for _, z := range zeroes {
		score, rating := FindPaths(z, hikingMap)
		totalScore += score
		totalRating += rating
	}

	fmt.Printf("Total Score : %4d\nTotal Rating: %4d\n", totalScore, totalRating)
}
