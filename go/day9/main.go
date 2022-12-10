package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/rarguelloF/advent-of-code-2022/input"
)

const inputName = "day9"

type Direction int

const (
	DirectionUp Direction = iota
	DirectionDown
	DirectionLeft
	DirectionRight
)

func ParseDirection(s string) (Direction, error) {
	dirMap := map[string]Direction{
		"R": DirectionRight,
		"L": DirectionLeft,
		"U": DirectionUp,
		"D": DirectionDown,
	}

	dir, ok := dirMap[s]
	if !ok {
		return 0, fmt.Errorf("unknown direction: %s", s)
	}

	return dir, nil
}

type Movement struct {
	Steps     int
	Direction Direction
}

type Position struct {
	X int
	Y int
}

func (p *Position) String() string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}

func (p *Position) Move(direction Direction) {
	switch direction {
	case DirectionUp:
		p.Y++

	case DirectionDown:
		p.Y--

	case DirectionRight:
		p.X++

	case DirectionLeft:
		p.X--
	}
}

func (p *Position) Follow(other *Position) {
	// they are in the same column, so need a vertical movement
	if p.X == other.X {
		if other.Y > p.Y {
			p.Y++
		} else {
			p.Y--
		}

		return
	}

	// they are in the same row, so need an horizontal movement
	if p.Y == other.Y {
		if other.X > p.X {
			p.X++
		} else {
			p.X--
		}

		return
	}

	// they are in different row and column, so need a diagonal movement
	if other.X > p.X {
		p.X++
	} else {
		p.X--
	}

	if other.Y > p.Y {
		p.Y++
	} else {
		p.Y--
	}

	return
}

func (p *Position) Distance(other *Position) float64 {
	return math.Sqrt(
		math.Pow(float64(other.X-p.X), 2) +
			math.Pow(float64(other.Y-p.Y), 2))
}

func (p *Position) IsAdjacent(other *Position) bool {
	return p.Distance(other) < 2.0
}

func PartOne(movements []*Movement) {
	tailPositions := make(map[string]int, 0)

	headPos := &Position{X: 0, Y: 0}
	tailPos := &Position{X: 0, Y: 0}

	tailPositions[tailPos.String()]++

	for _, m := range movements {
		for s := 0; s < m.Steps; s++ {
			headPos.Move(m.Direction)

			if !tailPos.IsAdjacent(headPos) {
				tailPos.Follow(headPos)
				tailPositions[tailPos.String()]++
			}
		}
	}

	fmt.Printf("Part 1: %d\n", len(tailPositions))
}

func PartTwo(movements []*Movement) {
	const numberOfKnots = 10

	knots := make([]*Position, 0, numberOfKnots)
	for i := 0; i < numberOfKnots; i++ {
		knots = append(knots, &Position{X: 0, Y: 0})
	}

	tailPositions := make(map[string]int, 0)
	tailPositions[knots[len(knots)-1].String()]++

	for _, m := range movements {
		for s := 0; s < m.Steps; s++ {
			head := knots[0]
			head.Move(m.Direction)

			for i := 1; i < len(knots); i++ {
				cur := knots[i]
				prev := knots[i-1]

				if !cur.IsAdjacent(prev) {
					cur.Follow(prev)
					if i == len(knots)-1 {
						tailPositions[cur.String()]++
					}
				}
			}
		}
	}

	fmt.Printf("Part 2: %d\n", len(tailPositions))
}

func readInput() ([]*Movement, error) {
	movements := make([]*Movement, 0)

	processLine := func(line string) error {
		values := strings.Split(line, " ")
		if len(values) != 2 {
			return fmt.Errorf("unexpected line format: %s", line)
		}

		dirStr, stepsStr := values[0], values[1]

		dir, err := ParseDirection(dirStr)
		if err != nil {
			return err
		}

		steps, err := strconv.Atoi(stepsStr)
		if err != nil {
			return fmt.Errorf("steps is not a number: %s", line)
		}

		movements = append(movements, &Movement{
			Direction: dir,
			Steps:     steps,
		})

		return nil
	}

	if err := input.ReadLines(inputName, processLine); err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	return movements, nil
}

func main() {
	movements, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	PartOne(movements)
	PartTwo(movements)
}
