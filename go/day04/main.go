package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/rarguelloF/advent-of-code-2022/input"
)

const inputName = "day04"

type Range [2]int

func (r Range) Size() int {
	return r[1] - r[0] + 1
}

func (r Range) ContainsAll(other Range) bool {
	rStart, rEnd := r[0], r[1]
	oStart, oEnd := other[0], other[1]

	return oStart >= rStart && oStart <= rEnd && oEnd >= rStart && oEnd <= rEnd
}

func (r Range) ContainsAny(other Range) bool {
	rStart, rEnd := r[0], r[1]
	oStart, oEnd := other[0], other[1]

	return (oStart < rStart && oEnd >= rStart) || (oStart >= rStart && oStart <= rEnd)
}

type ElfPair [2]Range

func (e ElfPair) GetBiggestRange() Range {
	if e[0].Size() > e[1].Size() {
		return e[0]
	}

	return e[1]
}

func (e ElfPair) FullyOverlap() bool {
	big, small := e[0], e[1]
	if small.Size() > big.Size() {
		big, small = small, big
	}

	return big.ContainsAll(small)
}

func (e ElfPair) Overlap() bool {
	return e[0].ContainsAny(e[1])
}

func PartOne(pairs []ElfPair) {
	sum := 0
	for _, p := range pairs {
		if p.FullyOverlap() {
			sum++
		}
	}

	fmt.Printf("Part 1: %d\n", sum)
}

func PartTwo(pairs []ElfPair) {
	sum := 0
	for _, p := range pairs {
		if p.Overlap() {
			sum++
		}
	}

	fmt.Printf("Part 2: %d\n", sum)
}

func readInput() ([]ElfPair, error) {
	pairs := make([]ElfPair, 0)

	processLine := func(line string) error {
		rangesStr := strings.Split(line, ",")
		if len(rangesStr) != 2 {
			return fmt.Errorf("unknown line format: %s", line)
		}

		pair := ElfPair{}

		for elfIdx, rangeStr := range rangesStr {
			numsStr := strings.Split(rangeStr, "-")
			if len(numsStr) != 2 {
				return fmt.Errorf("unknown line format: %s", line)
			}

			for i, ns := range numsStr {
				n, err := strconv.Atoi(ns)
				if err != nil {
					return fmt.Errorf("line contains non-numeric values: %s", line)
				}

				pair[elfIdx][i] = n
			}
		}

		pairs = append(pairs, pair)
		return nil
	}

	if err := input.ReadLines(inputName, processLine); err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	return pairs, nil
}

func main() {
	pairs, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	PartOne(pairs)
	PartTwo(pairs)
}
