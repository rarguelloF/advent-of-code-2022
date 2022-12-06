package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Printf("failed to close file: %v\n", err)
	}
}

func readInput(path string) ([]*Strategy, error) {
	f, err := os.Open(path)
	if err != nil {

	}
	defer closeFile(f)

	strategies := make([]*Strategy, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		values := strings.Split(line, " ")

		if len(values) != 2 {
			return nil, fmt.Errorf("input line contains unknown format: %s", line)
		}

		s, err := NewStrategy(values[0], values[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse strategy: %w", err)
		}

		strategies = append(strategies, s)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return strategies, nil
}
