package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/rarguelloF/advent-of-code-2022/input"
)

const inputName = "day10"

type Operation interface {
	Run(register, cycle int) (int, int)
}

type OpNoop struct{}

func (*OpNoop) Run(register, cycle int) (int, int) {
	return register, cycle + 1
}

type OpAddx struct {
	Value int
}

func (o *OpAddx) Run(register, cycle int) (int, int) {
	return register + o.Value, cycle + 2
}

func GetSignalStrength(register, cycle int) int {
	return register * cycle
}

type CRT struct {
	Values []bool
}

func NewCRT() *CRT {
	return &CRT{
		Values: make([]bool, 240),
	}
}

func (c *CRT) DrawPixel(register, cycle int) {
	curRow := int(cycle / 40)
	register = register + (40 * curRow)

	val := register == cycle || register-1 == cycle || register+1 == cycle
	c.Values[cycle] = val
}

func (c *CRT) String() string {
	s := ""

	for i, v := range c.Values {
		if i > 0 && i%40 == 0 {
			s += "\n"
		}
		if v {
			s += "#"
		} else {
			s += " "
		}
	}

	return s
}

func PartOne(ops []Operation) {
	sumSignalStrengths := 0
	observeCycles := []int{20, 60, 100, 140, 180, 220}

	register, cycle := 1, 1

	for _, o := range ops {
		curObserveCycle := observeCycles[0]
		prevRegister := register

		register, cycle = o.Run(register, cycle)

		if cycle == curObserveCycle {
			sumSignalStrengths += GetSignalStrength(register, curObserveCycle)
			observeCycles = observeCycles[1:]
		} else if cycle > curObserveCycle {
			sumSignalStrengths += GetSignalStrength(prevRegister, curObserveCycle)
			observeCycles = observeCycles[1:]
		}

		if len(observeCycles) == 0 {
			break
		}
	}

	fmt.Printf("Part 1: %d\n", sumSignalStrengths)
}

func PartTwo(ops []Operation) {
	crt := NewCRT()

	register := 1
	op, ops := ops[0], ops[1:]
	updatedRegister, cycleApplyOp := op.Run(register, 0)

	for cycle := 0; cycle < 240; cycle++ {
		if cycle == cycleApplyOp {
			register = updatedRegister
			op, ops = ops[0], ops[1:]
			updatedRegister, cycleApplyOp = op.Run(register, cycle)
		}

		crt.DrawPixel(register, cycle)
	}

	fmt.Printf("Part 2: \n%s\n", crt.String())
}

func readInput() ([]Operation, error) {
	ops := make([]Operation, 0)

	processLine := func(line string) error {
		if len(line) == 0 {
			return errors.New("empty line")
		}

		opAndArgs := strings.Split(line, " ")

		switch opAndArgs[0] {
		case "noop":
			ops = append(ops, &OpNoop{})
			return nil

		case "addx":
			addVal, err := strconv.Atoi(opAndArgs[1])
			if err != nil {
				return fmt.Errorf("addx needs to be followed by a number: %s", line)
			}
			ops = append(ops, &OpAddx{Value: addVal})
			return nil

		default:
			return fmt.Errorf("unknown operation: %s", line)
		}
	}

	if err := input.ReadLines(inputName, processLine); err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	return ops, nil
}

func main() {
	ops, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	PartOne(ops)
	PartTwo(ops)
}
