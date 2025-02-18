package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

func CountLines(file *os.File) (int, error) {
	lineCount := 0
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineCount++
	}

	file.Seek(0, io.SeekStart)
	return lineCount, scanner.Err()
}

func ReadInputToArray(filename string) ([]int, int, error) {
	file, err := os.Open("./input.txt")
	if err != nil {
		return []int{}, 0, err
	}

	lineCount, err := CountLines(file)
	if err != nil {
		return []int{}, 0, err
	}

	contents := make([]int, 2*lineCount)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	i := 0
	for scanner.Scan() && i < 2000 {
		word := scanner.Text()
		if word == " " || word == "\n" {
			continue
		}

		num, err := strconv.Atoi(word)
		if err != nil {
			return []int{}, 0, err
		}

		index := (i / 2) + (i%2)*lineCount
		contents[index] = num
		i++
	}

	return contents, lineCount, scanner.Err()
}

func SortLists(lists []int, offset int) []int {
	// Sort the left list.
	sort.Slice(lists[:offset], func(i, j int) bool {
		return lists[i] < lists[j]
	})
	// Sort the right list
	sort.Slice(lists[offset:], func(i, j int) bool {
		return lists[i+offset] < lists[j+offset]
	})

	return lists
}

func ComputeTotalDistance(lists []int, offset int) int {
	sum := 0
	for i := 0; i < offset; i++ {
		distance := lists[i] - lists[i+offset]
		if distance < 0 {
			distance = -distance
		}
		sum += distance
	}

	return sum
}

func ComputeSimilarity(lists []int, offset int) int {
	similarity := 0
	type tally struct {
		number int
		count  int
	}

	countMap := make(map[int]int)

	// Tally the number of times each item in the right list appears
	for i := offset; i < 2*offset; i++ {
		listID := lists[i]
		countMap[listID] = countMap[listID] + 1
	}

	currentID := lists[0] - 1 // The left list is sorted so this won't match any entry
	for i := 0; i < offset; i++ {
		if lists[i] == currentID {
			continue
		}

		currentID = lists[i]

		similarity += currentID * countMap[currentID]
	}

	return similarity
}

func main() {
	contents, offset, err := ReadInputToArray("input.txt")
	if err != nil {
		panic(err)
	}

	SortLists(contents, offset)

	totalDistance := ComputeTotalDistance(contents, offset)
	fmt.Printf("Total distance: %d\n", totalDistance)

	similarity := ComputeSimilarity(contents, offset)
	fmt.Printf("Similarity score: %d\n", similarity)
}
