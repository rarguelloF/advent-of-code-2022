package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/rarguelloF/advent-of-code-2022/input"
)

const inputName = "day05"

type Stack []string

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Pop() (string, bool) {
	if s.IsEmpty() {
		return "", false
	}

	lastIdx := len(*s) - 1
	val := (*s)[lastIdx]
	*s = (*s)[:lastIdx]
	return val, true
}

func (s *Stack) PopN(n int) ([]string, bool) {
	if s.IsEmpty() {
		return []string{""}, false
	}

	if n > len(*s) {
		n = len(*s)
	}

	vals := (*s)[len(*s)-n : len(*s)]
	*s = (*s)[:len(*s)-n]
	return vals, true
}

func (s *Stack) Peek() string {
	if s.IsEmpty() {
		return ""
	}

	lastIdx := len(*s) - 1
	return (*s)[lastIdx]
}

func (s *Stack) Push(val string) {
	*s = append(*s, val)
}

func (s *Stack) PushN(vals []string) {
	*s = append(*s, vals...)
}

func (s *Stack) PushLeft(val string) {
	*s = append([]string{val}, *s...)
}

type CrateStacks map[int]*Stack

func (c CrateStacks) GetCopy() CrateStacks {
	cp := make(CrateStacks, 0)
	for k, v := range c {
		cpVal := make(Stack, len(*v))
		copy(cpVal, *v)
		cp[k] = &cpVal
	}

	return cp
}

type Instruction struct {
	NumItems int
	From     int
	To       int
}

func PartOne(stacks CrateStacks, instructions []*Instruction) {
	// create a copy to not modify the original
	stacks = stacks.GetCopy()

	topCrates := ""

	for _, ins := range instructions {
		for i := 0; i < ins.NumItems; i++ {
			item, ok := stacks[ins.From].Pop()
			if ok {
				stacks[ins.To].Push(item)
			}
		}
	}

	numStacks := len(stacks)
	for i := 1; i <= numStacks; i++ {
		if val := stacks[i].Peek(); val != "" {
			topCrates += val
		}
	}

	fmt.Printf("Part 1: %s\n", topCrates)
}

func PartTwo(stacks CrateStacks, instructions []*Instruction) {
	// create a copy to not modify the original
	stacks = stacks.GetCopy()

	topCrates := ""

	for _, ins := range instructions {
		items, ok := stacks[ins.From].PopN(ins.NumItems)
		if ok {
			stacks[ins.To].PushN(items)
		}
	}

	numStacks := len(stacks)
	for i := 1; i <= numStacks; i++ {
		if val := stacks[i].Peek(); val != "" {
			topCrates += val
		}
	}

	fmt.Printf("Part 2: %s\n", topCrates)
}

func readInput() (CrateStacks, []*Instruction, error) {
	stacks := make(CrateStacks, 0)
	instructions := make([]*Instruction, 0)

	processStack := func(line string) error {
		stackID := 1
		for i := 0; i < len(line); i += 4 {
			if line[i] == '[' && line[i+2] == ']' {
				val := string(line[i+1])

				s, ok := stacks[stackID]
				if ok {
					s.PushLeft(val)
				} else {
					stacks[stackID] = &Stack{val}
				}
			}

			stackID++
		}

		return nil
	}

	processInstruction := func(line string) error {
		re := regexp.MustCompile(`move (?P<n1>\d+) from (?P<n2>\d+) to (?P<n3>\d+)`)
		match := re.FindStringSubmatch(line)

		if len(match) != 4 {
			return fmt.Errorf("wrong instruction format: %s", line)
		}

		nonNumberErr := fmt.Errorf("instruction contains non-numeric value: %s", line)

		numItems, err := strconv.Atoi(match[1])
		if err != nil {
			return nonNumberErr
		}

		from, err := strconv.Atoi(match[2])
		if err != nil {
			return nonNumberErr
		}

		to, err := strconv.Atoi(match[3])
		if err != nil {
			return nonNumberErr
		}

		instructions = append(instructions, &Instruction{
			NumItems: numItems,
			From:     from,
			To:       to,
		})

		return nil
	}

	processingInstructions := false

	processLine := func(line string) error {
		if line == "" {
			// now processing instructions
			processingInstructions = true
			return nil
		}

		if processingInstructions {
			return processInstruction(line)
		}

		return processStack(line)
	}

	if err := input.ReadLines(inputName, processLine); err != nil {
		return nil, nil, fmt.Errorf("failed to read input: %w", err)
	}

	return stacks, instructions, nil
}

func main() {
	stacks, instructions, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	PartOne(stacks, instructions)
	PartTwo(stacks, instructions)
}
