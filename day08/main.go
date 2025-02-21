package main

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"os"
)

type vector struct {
	x int
	y int
}

/**
 * Functions to query the grid.
 */

func ToCoords(index int, dimensions vector) vector {
	return vector{index % dimensions.x, index / dimensions.x}
}

func InBounds(pos vector, dimensions vector) bool {
	return pos.x >= 0 && pos.x < dimensions.x && pos.y >= 0 && pos.y < dimensions.y
}

func PopulateTowersFromReader(r *bufio.Reader) (*map[rune][]int, vector) {
	antennas := make(map[rune][]int)
	dimensions := vector{0, 0}

	index := 0
	for {
		ch, _, err := r.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		if ch == '\n' {
			dimensions.y++
			continue
		} else if ch != '.' {
			if antennas[ch] == nil {
				antennas[ch] = make([]int, 0, 1)
			}
			antennas[ch] = append(antennas[ch], index)
		}

		// Once the grid's height has been incremented, we know the grid width.
		if dimensions.y == 0 {
			dimensions.x++
		}
		index++
	}

	return &antennas, dimensions
}

/**
 * Functions to modify vectortors.
 */

func (vec vector) Add(other vector) vector {
	return vector{vec.x + other.x, vec.y + other.y}
}

func (vec vector) Sub(other vector) vector {
	return vector{vec.x - other.x, vec.y - other.y}
}

/**
 * Problem solution.
 */

func Abs(val int) int {
	if val < 0 {
		return -val
	}

	return val
}

func GCD(a int, b int) int {
	if a == 0 {
		return b
	} else if b == 0 {
		return a
	}

	larger := max(a, b)
	smaller := min(a, b)
	r := larger % smaller
	return GCD(smaller, r)
}

func PopulateAntinodes(antennaA int, antennaB int, dimensions vector, cachedAntinodes *map[vector]bool) {
	p0 := ToCoords(antennaA, dimensions)
	p1 := ToCoords(antennaB, dimensions)
	dir := p1.Sub(p0)
	gcd := GCD(Abs(dir.x), Abs(dir.y))
	dir.x /= gcd
	dir.y /= gcd

	antinode := p0
	for InBounds(antinode, dimensions) {
		(*cachedAntinodes)[antinode] = true
		antinode = antinode.Add(dir)
	}

	antinode = p0.Sub(dir)
	for InBounds(antinode, dimensions) {
		(*cachedAntinodes)[antinode] = true
		antinode = antinode.Sub(dir)
	}
}

func AllPairs[T any](sequence []T) iter.Seq2[T, T] {
	return func(yield func(T, T) bool) {
		if len(sequence) <= 1 {
			return
		}
		for i, a := range sequence[:len(sequence)-1] {
			for _, b := range sequence[i+1:] {
				if !yield(a, b) {
					return
				}
			}
		}
	}
}

func FindAllAntinodes(frequencyMapping *map[rune][]int, dimensions vector) int {
	antinodeLocations := make(map[vector]bool)

	for _, antennas := range *frequencyMapping {
		if len(antennas) == 1 {
			continue
		}
		for antennaA, antennaB := range AllPairs(antennas) {
			PopulateAntinodes(antennaA, antennaB, dimensions, &antinodeLocations)
		}
	}

	return len(antinodeLocations)
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	freqToLocation, dimensions := PopulateTowersFromReader(reader)
	numAntinodes := FindAllAntinodes(freqToLocation, dimensions)

	fmt.Printf("Number of antinodes: %d\n", numAntinodes)
}
