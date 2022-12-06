package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Printf("failed to close file: %v\n", err)
	}
}

func readInput(path string) ([]*ElfInventory, error) {
	f, err := os.Open(path)
	if err != nil {

	}
	defer closeFile(f)

	data := make([]*ElfInventory, 0)
	cur := newElfInventory()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		val := scanner.Text()
		if val == "" {
			data = append(data, cur)
			cur = newElfInventory()
			continue
		}

		calories, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("found non-intenger value in input: %s", val)
		}

		cur.FoodItems = append(cur.FoodItems, &Food{
			Calories: calories,
		})
	}

	data = append(data, cur)

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}
