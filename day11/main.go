package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type StoneCount map[int]int

func PopulateStones(file *os.File) []int {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	stones := make([]int, 0, 1)
	for scanner.Scan() {
		word := scanner.Text()
		num, err := strconv.Atoi(word)
		if err != nil {
			panic(err)
		}

		stones = append(stones, num)
	}

	return stones
}

func (stones *StoneCount) Blink(stone int, next *StoneCount) {
	nStones := (*stones)[stone]
	(*stones)[stone] = 0
	if stone == 0 {
		(*next)[1] += nStones
		return
	}

	stoneStr := strconv.Itoa(stone)
	nDigits := len(stoneStr)
	if nDigits%2 == 0 {
		left := stoneStr[:nDigits/2]
		right := stoneStr[nDigits/2:]

		num, _ := strconv.Atoi(left)
		(*next)[num] += nStones

		num, _ = strconv.Atoi(right)
		(*next)[num] += nStones
		return
	}

	(*next)[stone*2024] += nStones
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	stones := PopulateStones(file)

	nStones := 0
	for _, stone := range stones {
		stoneBufferA := make(StoneCount)
		stoneBufferB := make(StoneCount)
		var current, next *StoneCount

		stoneBufferA[stone] = 1
		toggle := true
		for range 75 {
			if toggle {
				current = &stoneBufferA
				next = &stoneBufferB
			} else {
				current = &stoneBufferB
				next = &stoneBufferA
			}

			for k, v := range *current {
				if v > 0 {
					current.Blink(k, next)
				}
			}

			toggle = !toggle
		}

		for _, v := range *next {
			nStones += v
		}
	}

	fmt.Printf("Number of stones: %d\n", nStones)
}
