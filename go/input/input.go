package input

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
)

type ProcessLineFunc func(line string) error

func ReadInput(name string, processLine ProcessLineFunc) error {
	useTest := false
	for _, arg := range os.Args[1:] {
		if arg == "-test" {
			useTest = true
		}
	}

	if useTest {
		name += "_test"
	}

	fName := fmt.Sprintf("%s.txt", name)
	fPath := path.Join("../../inputs/", fName)

	f, err := os.Open(fPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	defer closeFile(f)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if err := processLine(line); err != nil {
			return fmt.Errorf("failed to process line: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	return nil
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Printf("failed to close file: %v\n", err)
	}
}
