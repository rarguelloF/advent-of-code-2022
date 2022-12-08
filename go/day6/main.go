package main

import (
	"fmt"
	"log"

	"github.com/rarguelloF/advent-of-code-2022/input"
)

const inputName = "day6"

func isUniqueChars(foundChars map[rune]int) bool {
	for _, count := range foundChars {
		if count > 1 {
			return false
		}
	}
	return true
}

func findFirstUniqueSequenceStart(seqLength int) (int, error) {
	idx := -1
	curIdx := 1
	chars := make([]rune, 0, seqLength)
	foundChars := make(map[rune]int, 0)

	processChar := func(char rune) (bool, error) {
		if len(chars) >= seqLength {
			first := chars[0]
			if foundChars[first] == 1 {
				delete(foundChars, first)
			} else {
				foundChars[first]--
			}
			chars = chars[1:]
		}

		chars = append(chars, char)
		foundChars[char]++
		if len(chars) >= seqLength && isUniqueChars(foundChars) {
			idx = curIdx
			return true, nil
		}

		curIdx++
		return false, nil
	}

	if err := input.ReadChars(inputName, processChar); err != nil {
		return -1, fmt.Errorf("failed to read input: %w", err)
	}

	return idx, nil
}

func PartOne() {
	startIdx, err := findFirstUniqueSequenceStart(4)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Part 1: %d\n", startIdx)
}

func PartTwo() {
	startIdx, err := findFirstUniqueSequenceStart(14)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Part 2: %d\n", startIdx)
}

func main() {
	PartOne()
	PartTwo()
}
