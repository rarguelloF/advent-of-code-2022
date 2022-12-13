package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/rarguelloF/advent-of-code-2022/input"
)

const inputName = "day11"

type Operation int

const (
	OperationSum Operation = iota
	OperationMul
)

func ParseOperation(s string) (Operation, error) {
	op, ok := map[string]Operation{
		"+": OperationSum,
		"*": OperationMul,
	}[s]

	if !ok {
		return -1, fmt.Errorf("unknown operation: %s", s)
	}

	return op, nil
}

type GetNewWorryLevelFunc func(item int64) int64

func ParseGetNewWorryLevelFunc(aStr, bStr, opStr string) (GetNewWorryLevelFunc, error) {
	var a, b int64

	aIsOld := aStr == "old"
	bIsOld := bStr == "old"

	if !aIsOld {
		n, err := strconv.ParseInt(aStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("a is not a number: %s", aStr)
		}
		a = n
	}

	if !bIsOld {
		n, err := strconv.ParseInt(bStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("b is not a number: %s", bStr)
		}
		b = n
	}

	op, err := ParseOperation(opStr)
	if err != nil {
		return nil, err
	}

	return func(item int64) int64 {
		first, second := a, b
		if aIsOld {
			first = item
		}
		if bIsOld {
			second = item
		}

		switch op {
		case OperationSum:
			return first + second

		case OperationMul:
			return first * second

		default:
			panic(fmt.Errorf("unimplemented operation: %d", op))
		}
	}, nil
}

type Monkey struct {
	ID               int
	Items            []int64
	GetNewWorryLevel GetNewWorryLevelFunc
	TestDivisible    int
	ThrowTrue        int
	ThrowFalse       int
	InspectedItems   int
}

func (m *Monkey) PopItem() (int64, bool) {
	if len(m.Items) == 0 {
		return 0, false
	}

	r := m.Items[0]
	m.Items = m.Items[1:]
	return r, true
}

type ReduceWorryFunc func(item int64) int64

func (m *Monkey) GetItemTarget(worryLevel int64, reduceWorry ReduceWorryFunc) (int64, int) {
	m.InspectedItems++
	newLevel := m.GetNewWorryLevel(worryLevel)
	newLevel = reduceWorry(newLevel)

	if newLevel%int64(m.TestDivisible) == 0 {
		return newLevel, m.ThrowTrue
	}
	return newLevel, m.ThrowFalse
}

type Monkeys map[int]*Monkey

func (ms Monkeys) GetCopy() Monkeys {
	cp := make(Monkeys, 0)
	for k, v := range ms {
		items := make([]int64, len(v.Items))
		copy(items, v.Items)

		cp[k] = &Monkey{
			ID:               v.ID,
			Items:            items,
			GetNewWorryLevel: v.GetNewWorryLevel,
			TestDivisible:    v.TestDivisible,
			ThrowTrue:        v.ThrowTrue,
			ThrowFalse:       v.ThrowFalse,
			InspectedItems:   v.InspectedItems,
		}
	}

	return cp
}

func (ms Monkeys) InspectedItems() []int {
	inspectedItems := make([]int, 0, len(ms))
	for i := 0; i < len(ms); i++ {
		inspectedItems = append(inspectedItems, ms[i].InspectedItems)
	}

	return inspectedItems
}

func GetMonkeyBusinessLevel(monkeys Monkeys, rounds int, reduceWorry ReduceWorryFunc) int {
	for i := 0; i < rounds; i++ {
		for id := 0; id < len(monkeys); id++ {
			m := monkeys[id]
			for len(m.Items) > 0 {
				item, ok := m.PopItem()
				if !ok {
					break
				}

				newLevel, target := m.GetItemTarget(item, reduceWorry)
				if target == m.ID {
					panic("monkey throwing to himself (infinite loop)")
				}

				monkeys[target].Items = append(monkeys[target].Items, newLevel)
			}
		}
	}

	inspectedItems := monkeys.InspectedItems()

	sort.Slice(inspectedItems, func(i, j int) bool {
		return inspectedItems[j] < inspectedItems[i]
	})

	return inspectedItems[0] * inspectedItems[1]
}

func PartOne(monkeys Monkeys) {
	const numRounds = 20

	reduceWorry := func(item int64) int64 {
		return int64(item / 3)
	}

	fmt.Printf("Part 1: %d\n", GetMonkeyBusinessLevel(monkeys, numRounds, reduceWorry))
}

func PartTwo(monkeys Monkeys) {
	const numRounds = 10_000

	// Trick to reduce worry level to not overflow int64.
	// Good explanation from a reddit user:
	// (Link to the comment: https://www.reddit.com/r/adventofcode/comments/zifqmh/comment/izv7hpx)

	/*
		When you check if a number is divisible by 23, you check what the remainder is, right?
		For instance, 123 mod 23 = 8, so 123 is not divisible by 23.

		But what if you thought 123 was WAY too big a number and you wanted to reduce it first?
		But you STILL want to do that "div by 23" test, with the same result.
		Well, you could first mod it by 23 (...wait for it...). That gives 8. Now instead of 123
		you have 8. Then you do the divisibility test with this reduced number: 8 mod 23 = 8.
		That's the same, so it worked :)

		Now watch what happens if you mod by a multiple, ANY multiple of 23 first to reduce your big number.
		For instance, you could mod by 46 = 2 x 23: 123 mod 46 = 31.
		What does your div test give? 31 mod 23 = 8. The same! And this is always true, for any multiple.
		You can think of "mod 23" as "take the 23-ness out of your number" while the remainder
		remains the same.

		So that 2 could be the div test number of another monkey.
		And another one has 5, etc. etc. If you want to take the 23-ness,
		the 2-ness and the 5-ness out of a number, you mod by 2 x 5 x 23.
		As we saw above, the remainder will remain the same for our div tests with 2, 5, and 23.
		Because you can see this factor as a multiple of 23, or a multiple of 5, or a multiple of 2.
	*/
	divisor := 1
	for _, m := range monkeys {
		divisor *= m.TestDivisible
	}

	reduceWorry := func(item int64) int64 {
		return item % int64(divisor)
	}

	fmt.Printf("Part 2: %d\n", GetMonkeyBusinessLevel(monkeys, numRounds, reduceWorry))
}

func readInput() (Monkeys, error) {
	regexLines := [6]string{
		`^Monkey (?P<id>.*):$`,
		`^  Starting items: (?P<items>.*)$`,
		`^  Operation: new = (?P<a>.*) (?P<op>.*) (?P<b>.*)$`,
		`^  Test: divisible by (?P<num>.*)$`,
		`^    If true: throw to monkey (?P<id>.*)$`,
		`^    If false: throw to monkey (?P<id>.*)$`,
	}

	monkeys := make(Monkeys, 0)

	var curMonkey *Monkey
	curLine := 0

	processLine := func(line string) error {
		regexNotMatchErr := errors.New("regex didn't match")
		expectedNumberErr := errors.New("expected a number")

		if len(line) == 0 {
			curLine = 0
			return nil
		}

		re := regexp.MustCompile(regexLines[curLine])
		match := re.FindStringSubmatch(line)

		switch curLine {
		case 0:
			if len(match) != 2 {
				return regexNotMatchErr
			}

			id, err := strconv.Atoi(match[1])
			if err != nil {
				return expectedNumberErr
			}

			curMonkey = &Monkey{ID: id}

		case 1:
			if len(match) != 2 {
				return regexNotMatchErr
			}

			itemsStr := strings.Split(match[1], ", ")
			curMonkey.Items = make([]int64, 0, len(itemsStr))

			for _, s := range itemsStr {
				item, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return expectedNumberErr
				}

				curMonkey.Items = append(curMonkey.Items, item)
			}

		case 2:
			if len(match) != 4 {
				return regexNotMatchErr
			}

			f, err := ParseGetNewWorryLevelFunc(match[1], match[3], match[2])
			if err != nil {
				return err
			}

			curMonkey.GetNewWorryLevel = f

		case 3:
			if len(match) != 2 {
				return regexNotMatchErr
			}

			num, err := strconv.Atoi(match[1])
			if err != nil {
				return expectedNumberErr
			}

			curMonkey.TestDivisible = num

		case 4:
			if len(match) != 2 {
				return regexNotMatchErr
			}

			id, err := strconv.Atoi(match[1])
			if err != nil {
				return expectedNumberErr
			}

			curMonkey.ThrowTrue = id

		case 5:
			if len(match) != 2 {
				return regexNotMatchErr
			}

			id, err := strconv.Atoi(match[1])
			if err != nil {
				return expectedNumberErr
			}

			curMonkey.ThrowFalse = id

			monkeys[curMonkey.ID] = curMonkey
			curMonkey = nil
		}

		curLine++
		return nil
	}

	if err := input.ReadLines(inputName, processLine); err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	return monkeys, nil
}

func main() {
	monkeys, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	PartOne(monkeys.GetCopy())
	PartTwo(monkeys.GetCopy())
}
