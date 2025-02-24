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

type HikingMap struct {
	grid   Grid
	zeroes []Vector
}

func (vec Vector) Add(other Vector) Vector {
	return Vector{vec.x + other.x, vec.y + other.y}
}

func (grid *Grid) InBounds(point Vector) bool {
	return point.x >= 0 && point.x < grid.width && point.y >= 0 && point.y < grid.height
}

func NewMap() *HikingMap {
	var topMap HikingMap
	topMap.grid = Grid{contents: make([]int, 0, 128), width: 0, height: 0}
	topMap.zeroes = make([]Vector, 0, 20)

	return &topMap
}

func PopulateMapFromReader(r *bufio.Reader) *HikingMap {
	hikingMap := NewMap()

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
			hikingMap.grid.height++
			continue
		} else if ch == '0' {
			hikingMap.zeroes = append(hikingMap.zeroes, Vector{chCount, hikingMap.grid.height})
		}

		num, err := strconv.Atoi(string(ch))
		if err != nil {
			panic(err)
		}
		hikingMap.grid.contents = append(hikingMap.grid.contents, num)

		// Once the grid's width is locked in, we don't need to keep track of the character count.
		if hikingMap.grid.height == 0 {
			hikingMap.grid.width++
		}
		chCount++
	}

	return hikingMap
}

func (grid *Grid) Get(point Vector) int {
	return grid.contents[point.y*grid.width+point.x]
}

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
	frontier := make([]Vector, 0, 4)
	ninesSet := make(map[Vector]bool)
	nines := 0

	frontier = append(frontier, zeroLoc)
	for len(frontier) != 0 {
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
	hikingMap := PopulateMapFromReader(reader)

	totalScore := 0
	totalRating := 0
	for _, z := range hikingMap.zeroes {
		score, rating := FindPaths(z, &hikingMap.grid)
		totalScore += score
		totalRating += rating
	}

	fmt.Printf("Total Score: %d\nTotal Rating: %d\n", totalScore, totalRating)
}
