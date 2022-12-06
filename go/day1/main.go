package main

import (
	"fmt"
	"log"
	"sort"
)

type Food struct {
	Calories int
}

type ElfInventory struct {
	FoodItems []*Food
}

func (e *ElfInventory) TotalCalories() (result int) {
	for _, f := range e.FoodItems {
		result += f.Calories
	}

	return result
}

func newElfInventory() *ElfInventory {
	return &ElfInventory{
		FoodItems: make([]*Food, 0),
	}
}

func findMaxCalories(inventories []*ElfInventory) (max int) {
	for _, inventory := range inventories {
		if total := inventory.TotalCalories(); total > max {
			max = total
		}
	}

	return max
}

func findTopNCalories(inventories []*ElfInventory, n int) []int {
	allTotalCalories := make([]int, 0, len(inventories))
	for _, inventory := range inventories {
		allTotalCalories = append(allTotalCalories, inventory.TotalCalories())
	}

	sort.Slice(allTotalCalories, func(i, j int) bool {
		return allTotalCalories[j] < allTotalCalories[i]
	})

	return allTotalCalories[:n]
}

func PartOne(inventories []*ElfInventory) {
	fmt.Printf("Part 1: %d\n", findMaxCalories(inventories))
}

func PartTwo(inventories []*ElfInventory) {
	sumTopThree := 0
	for _, t := range findTopNCalories(inventories, 3) {
		sumTopThree += t
	}

	fmt.Printf("Part 2: %d\n", sumTopThree)
}

func main() {
	inventories, err := readInput("../../inputs/day1.txt")
	if err != nil {
		log.Fatal(err)
	}

	PartOne(inventories)
	PartTwo(inventories)
}
