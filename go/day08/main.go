package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/rarguelloF/advent-of-code-2022/input"
)

const inputName = "day08"

type Direction int

const (
	DirectionTop Direction = iota
	DirectionBottom
	DirectionLeft
	DirectionRight
)

type Trees [][]int

func (t Trees) IsVisible(row, col int) bool {
	// outer trees are always visible
	if t.isOuterTree(row, col) {
		return true
	}

	return t.isVisibleFromBottom(row, col) ||
		t.isVisibleFromTop(row, col) ||
		t.isVisibleFromLeft(row, col) ||
		t.isVisibleFromRight(row, col)
}

func (t Trees) VisibilityScore(row, col int) int {
	treesBottom, _ := t.countVisibleTrees(row, col, DirectionBottom)
	treesTop, _ := t.countVisibleTrees(row, col, DirectionTop)
	treesLeft, _ := t.countVisibleTrees(row, col, DirectionLeft)
	treesRight, _ := t.countVisibleTrees(row, col, DirectionRight)

	return treesBottom * treesTop * treesLeft * treesRight
}

func (t Trees) countVisibleTrees(row, col int, direction Direction) (numTrees int, reachedEnd bool) {
	if t.isOuterTree(row, col) {
		return 0, true
	}

	height := t[row][col]
	numRows, numCols := len(t), len(t[0])

	var loopRows bool
	var loopStart, loopEnd, loopIncrease int

	switch direction {
	case DirectionTop:
		loopRows = true
		loopStart = row - 1
		loopEnd = 0
		loopIncrease = -1

	case DirectionBottom:
		loopRows = true
		loopStart = row + 1
		loopEnd = numRows - 1
		loopIncrease = 1

	case DirectionLeft:
		loopStart = col - 1
		loopEnd = 0
		loopIncrease = -1

	case DirectionRight:
		loopStart = col + 1
		loopEnd = numCols - 1
		loopIncrease = 1

	default:
		return 0, false
	}

	checkLoop := func(i int) bool {
		if loopIncrease < 0 {
			return i >= loopEnd
		}
		return i <= loopEnd
	}

	for i := loopStart; checkLoop(i); i += loopIncrease {
		numTrees++
		r, c := row, col
		if loopRows {
			r = i
		} else {
			c = i
		}

		otherHeight := t[r][c]
		if otherHeight >= height {
			return numTrees, false
		}
	}

	return numTrees, true
}

func (t Trees) isVisibleFromBottom(row, col int) bool {
	_, ok := t.countVisibleTrees(row, col, DirectionBottom)
	return ok
}

func (t Trees) isVisibleFromTop(row, col int) bool {
	_, ok := t.countVisibleTrees(row, col, DirectionTop)
	return ok
}

func (t Trees) isVisibleFromLeft(row, col int) bool {
	_, ok := t.countVisibleTrees(row, col, DirectionLeft)
	return ok
}

func (t Trees) isVisibleFromRight(row, col int) bool {
	_, ok := t.countVisibleTrees(row, col, DirectionRight)
	return ok
}

func (t Trees) isOuterTree(row, col int) bool {
	numCols := len(t[0])
	return row <= 0 || row >= len(t)-1 || col == 0 || col >= numCols-1
}

func PartOne(trees Trees) {
	totalVisibleTrees := 0

	for row := 0; row < len(trees); row++ {
		for col := 0; col < len(trees[row]); col++ {
			if trees.IsVisible(row, col) {
				totalVisibleTrees++
			}
		}
	}

	fmt.Printf("Part 1: %d\n", totalVisibleTrees)
}

func PartTwo(trees Trees) {
	maxVisibility := 0

	for row := 0; row < len(trees); row++ {
		for col := 0; col < len(trees[row]); col++ {
			s := trees.VisibilityScore(row, col)
			if s > maxVisibility {
				maxVisibility = s
			}
		}
	}

	fmt.Printf("Part 2: %d\n", maxVisibility)
}

func readInput() (Trees, error) {
	trees := make(Trees, 0)
	numCols := -1

	processLine := func(line string) error {
		if len(line) == 0 {
			return errors.New("found empty line")
		}

		cols := len(line)
		if numCols == -1 {
			numCols = cols
		} else if cols != numCols {
			return errors.New("number of columns should be the same for every row")
		}

		treeRow := make([]int, 0, cols)
		for _, h := range line {
			height := int(h - '0')
			treeRow = append(treeRow, height)
		}

		trees = append(trees, treeRow)
		return nil
	}

	if err := input.ReadLines(inputName, processLine); err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	return trees, nil
}

func main() {
	trees, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	PartOne(trees)
	PartTwo(trees)
}
