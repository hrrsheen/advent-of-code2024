package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func Abs(val int) int {
	if val < 0 {
		return -val
	}

	return val
}

func IsSafe(levels []int) (bool, int) {
	const maxDiff = 3

	prevLevel := levels[0]
	prevDiff := levels[1] - levels[0]

	for i, level := range levels[1:] {
		diff := level - prevLevel

		if Abs(diff) > maxDiff || Abs(diff) < 1 {
			return false, i
		}

		if diff*prevDiff < 0 {
			return false, i
		}

		prevLevel = level
		prevDiff = diff
	}

	return true, -1
}

func IsSafeWithDampening(levels []int) bool {
	isSafe, where := IsSafe(levels)

	if !isSafe {
		if where > 0 {
			where--
		}

		// Attempt removal of an element.
		// Build a slice with the level removed.
		levelsWithRemoval := append([]int{}, levels[:where]...)
		levelsWithRemoval = append(levelsWithRemoval, levels[(where+1):]...)

		for ; where < len(levels); where++ {
			// Test whether the new levels slice is safe.
			safeWithRemoval, _ := IsSafe(levelsWithRemoval)
			if safeWithRemoval {
				return true
			}

			if where < len(levelsWithRemoval) {
				levelsWithRemoval[where] = levels[where]
			}
		}
	}

	return isSafe
}

func CountSafeReports(r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	safeCount := 0
	unsafeCount := 0
	levels := make([]int, 0, 8)
	for scanner.Scan() {
		lineScanner := bufio.NewScanner(strings.NewReader(scanner.Text()))
		lineScanner.Split(bufio.ScanWords)
		for lineScanner.Scan() {
			levelText := lineScanner.Text()
			level, err := strconv.Atoi(levelText)
			if err != nil {
				return 0, err
			}

			levels = append(levels, level)
		}

		if safeCount == 5 && unsafeCount == 0 { // DEBUGGING point
			fmt.Printf("lol\n")
		}

		safe := IsSafeWithDampening(levels)
		if safe {
			safeCount++
		} else {
			unsafeCount++
		}
		fmt.Printf("%t: %v\n", safe, levels)
		levels = levels[:0]
	}

	return safeCount, scanner.Err()
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	safeCount, err := CountSafeReports(file)

	fmt.Printf("Safe count: %d\n", safeCount)
}
