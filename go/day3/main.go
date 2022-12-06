package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/rarguelloF/advent-of-code-2022/input"
)

type Rucksack struct {
	AllItems     map[rune]int
	Compartments [2]map[rune]int
}

func (r *Rucksack) GetRepeatedItem() (rune, error) {
	first, second := r.Compartments[0], r.Compartments[1]

	for item := range first {
		if _, ok := second[item]; ok {
			return item, nil
		}
	}

	return 0, errors.New("not found")
}

func itemTypeToPriority(itemType rune) int {
	asciiPos := int(itemType)
	if asciiPos >= 97 {
		// a to z
		// a is 97 in ascii -> 1 priority
		return asciiPos - 96
	}

	// A to Z
	// A is position 65 in ascii -> 27 priority
	return asciiPos - 38
}

func findCommonItem(group [3]*Rucksack) (rune, error) {
	g1, g2, g3 := group[0], group[1], group[2]
	for item := range g1.AllItems {
		if _, ok := g2.AllItems[item]; ok {
			if _, ok := g3.AllItems[item]; ok {
				return item, nil
			}
		}
	}

	return 0, errors.New("not found")
}

func readInput(inputName string) ([]*Rucksack, error) {
	rucksacks := make([]*Rucksack, 0)

	processLine := func(line string) error {
		numElems := len(line)
		if numElems%2 != 0 {
			return fmt.Errorf("odd number of elements in rucksack (%d)", numElems)
		}

		r := &Rucksack{
			AllItems: make(map[rune]int, 0),
			Compartments: [2]map[rune]int{
				{},
				{},
			},
		}

		for idx, item := range line {
			if idx >= numElems/2 {
				r.Compartments[1][item]++
			} else {
				r.Compartments[0][item]++
			}

			r.AllItems[item]++
		}

		rucksacks = append(rucksacks, r)
		return nil
	}

	if err := input.ReadInput(inputName, processLine); err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	return rucksacks, nil
}

func PartOne(rucksacks []*Rucksack) {
	sum := 0
	for _, r := range rucksacks {
		item, err := r.GetRepeatedItem()
		if err != nil {
			log.Fatal(err)
		}

		sum += itemTypeToPriority(item)
	}

	fmt.Printf("Part 1: %d\n", sum)
}

func PartTwo(rucksacks []*Rucksack) {
	if len(rucksacks)%3 != 0 {
		log.Fatal("the number of rucksacks is not divisible by 3")
	}

	groups := make([][3]*Rucksack, 0, len(rucksacks)/3)

	for i := 0; i < len(rucksacks); i += 3 {
		g := [3]*Rucksack{
			rucksacks[i],
			rucksacks[i+1],
			rucksacks[i+2],
		}

		groups = append(groups, g)
	}

	sum := 0
	for _, g := range groups {
		item, err := findCommonItem(g)
		if err != nil {
			log.Fatal(err)
		}

		sum += itemTypeToPriority(item)
	}

	fmt.Printf("Part 2: %d\n", sum)
}

func main() {
	rucksacks, err := readInput("day3")
	if err != nil {
		log.Fatal(err)
	}

	PartOne(rucksacks)
	PartTwo(rucksacks)
}
