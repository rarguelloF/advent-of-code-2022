package input

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
)

func getFilePath(name string) string {
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
	return path.Join("../../inputs/", fName)
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Printf("failed to close file: %v\n", err)
	}
}

type ProcessLineFunc func(line string) error

func ReadLines(name string, processLine ProcessLineFunc) error {
	fPath := getFilePath(name)

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

type ProcessCharFunc func(char rune) (bool, error)

func ReadChars(name string, processChar ProcessCharFunc) error {
	fPath := getFilePath(name)

	f, err := os.Open(fPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer closeFile(f)

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanRunes)

	for scanner.Scan() {
		charStr := scanner.Text()
		chars := []rune(charStr)
		if len(chars) != 1 {
			return fmt.Errorf("more than one character was given: %s", charStr)
		}

		stop, err := processChar(chars[0])
		if err != nil {
			return fmt.Errorf("failed to process char: %w", err)
		}
		if stop {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	return nil
}
