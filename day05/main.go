package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func ValidateUpdate(pages []string, rules *map[string]map[string]bool) bool {
	updated := make(map[string]bool)

	for _, page := range pages {
		updated[page] = true
		// If a page in the corresponding rules already has an entry in "updates", then it clearly has been updated
		// before the current page. Thus, the update as a whole is invalid.
		for rule := range (*rules)[page] {
			if updated[rule] {
				return false
			}
		}
	}

	return true
}

func FixUpdate(pages []string, rules *map[string]map[string]bool) []string {
	CompareFn := func(a string, b string) int {
		if (*rules)[a][b] {
			return -1
		} else if (*rules)[b][a] {
			return 1
		}
		return 0
	}

	slices.SortFunc(pages, CompareFn)

	return pages
}

func main() {
	bytes, err := os.ReadFile("./input.txt")
	if err != nil {
		panic(err)
	}
	contents := string(bytes)

	matchRules, _ := regexp.Compile(`(\d{2})\|(\d{2})`)
	ruleStrings := matchRules.FindAllStringSubmatch(contents, -1)

	rules := make(map[string]map[string]bool)
	for _, match := range ruleStrings {
		before := match[1]
		after := match[2]
		if rules[before] == nil {
			rules[before] = make(map[string]bool)
		}
		rules[before][after] = true
	}

	matchUpdates, _ := regexp.Compile(`((?:,?\d{2}){2,})`)
	updateStrings := matchUpdates.FindAllString(contents, -1)

	middleTotal := 0
	for _, PagesStr := range updateStrings {
		pages := strings.Split(PagesStr, ",")
		if !ValidateUpdate(pages, &rules) {
			FixUpdate(pages, &rules)
			fmt.Printf("Fixed: %t - %v\n", ValidateUpdate(pages, &rules), pages)

			middleValue, _ := strconv.Atoi(pages[len(pages)/2])
			middleTotal += middleValue
		}
	}

	fmt.Printf("Middle value total: %d\n", middleTotal)
}
