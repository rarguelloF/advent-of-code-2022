package main

import (
	"fmt"
	"log"
)

type Play struct {
	ID      PlayID
	Points  int
	WinsTo  PlayID
	LosesTo PlayID
}

type PlayID int

const (
	PlayIDRock PlayID = iota + 1
	PlayIDPaper
	PlayIDScissor
)

var plays = map[PlayID]Play{
	PlayIDRock: {
		ID:      PlayIDRock,
		Points:  1,
		WinsTo:  PlayIDScissor,
		LosesTo: PlayIDPaper,
	},

	PlayIDPaper: {
		ID:      PlayIDPaper,
		Points:  2,
		WinsTo:  PlayIDRock,
		LosesTo: PlayIDScissor,
	},

	PlayIDScissor: {
		ID:      PlayIDScissor,
		Points:  3,
		WinsTo:  PlayIDPaper,
		LosesTo: PlayIDRock,
	},
}

type Strategy struct {
	OponentPlay Play
	OwnPlayV1   Play
	OwnPlayV2   Play
}

func NewStrategy(oponentPlay, ownPlay string) (*Strategy, error) {
	s := &Strategy{}

	switch oponentPlay {
	case "A":
		s.OponentPlay = plays[PlayIDRock]

	case "B":
		s.OponentPlay = plays[PlayIDPaper]

	case "C":
		s.OponentPlay = plays[PlayIDScissor]

	default:
		return nil, fmt.Errorf("unknown play for oponent: %s", oponentPlay)
	}

	switch ownPlay {
	case "X":
		s.OwnPlayV1 = plays[PlayIDRock]
		s.OwnPlayV2 = plays[s.OponentPlay.WinsTo]

	case "Y":
		s.OwnPlayV1 = plays[PlayIDPaper]
		s.OwnPlayV2 = plays[s.OponentPlay.ID]

	case "Z":
		s.OwnPlayV1 = plays[PlayIDScissor]
		s.OwnPlayV2 = plays[s.OponentPlay.LosesTo]

	default:
		return nil, fmt.Errorf("unknown play for own player: %s", ownPlay)
	}

	return s, nil
}

func (s *Strategy) PointsV1() int {
	return s.Points(s.OwnPlayV1)
}

func (s *Strategy) PointsV2() int {
	return s.Points(s.OwnPlayV2)
}

func (s *Strategy) Points(ownPlay Play) (result int) {
	result += ownPlay.Points

	if ownPlay.ID == s.OponentPlay.ID {
		return result + 3
	}

	if ownPlay.WinsTo == s.OponentPlay.ID {
		return result + 6
	}

	return result
}

func PartOne(strategies []*Strategy) {
	sum := 0
	for _, s := range strategies {
		sum += s.PointsV1()
	}

	fmt.Printf("Part 1: %d\n", sum)
}

func PartTwo(strategies []*Strategy) {
	sum := 0
	for _, s := range strategies {
		sum += s.PointsV2()
	}

	fmt.Printf("Part 2: %d\n", sum)
}

func main() {
	strategies, err := readInput("../../inputs/day2.txt")
	if err != nil {
		log.Fatal(err)
	}

	PartOne(strategies)
	PartTwo(strategies)
}
