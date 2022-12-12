package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/rarguelloF/advent-of-code-2022/input"
)

const inputName = "day07"

type File struct {
	Name string
	Size int64
}

type Dir struct {
	Name   string
	Parent *Dir
	Files  map[string]*File
	Dirs   map[string]*Dir
}

func (d *Dir) Size() (r int64) {
	for _, f := range d.Files {
		r += f.Size
	}

	for _, sd := range d.Dirs {
		r += sd.Size()
	}

	return r
}

func newDir(name string, parent *Dir) *Dir {
	return &Dir{
		Name:   name,
		Parent: parent,
		Files:  make(map[string]*File, 0),
		Dirs:   make(map[string]*Dir, 0),
	}
}

func PartOne(rootDir *Dir) {
	const maxSize = 100_000

	sumSizes := int64(0)

	dirs := []*Dir{rootDir}

	for len(dirs) > 0 {
		d := dirs[0]
		dirs = dirs[1:]

		if s := d.Size(); s <= maxSize {
			sumSizes += s
		}

		for _, sd := range d.Dirs {
			dirs = append(dirs, sd)
		}
	}

	fmt.Printf("Part 1: %d\n", sumSizes)
}

func PartTwo(rootDir *Dir) {
	const (
		availableDisk = 70_000_000
		minUnused     = 30_000_000
	)

	totalSize := rootDir.Size()
	curFreeSize := availableDisk - totalSize
	chosenDeleteSize := int64(math.MaxInt64)

	dirs := []*Dir{rootDir}

	for len(dirs) > 0 {
		d := dirs[0]
		dirs = dirs[1:]

		s := d.Size()

		if s+curFreeSize >= minUnused && s < chosenDeleteSize {
			chosenDeleteSize = s
		}

		for _, sd := range d.Dirs {
			dirs = append(dirs, sd)
		}
	}

	fmt.Printf("Part 2: %d\n", chosenDeleteSize)
}

const (
	CommandChangeDir = "cd"
	CommandList      = "ls"
)

func readInput() (*Dir, error) {
	rootDir := newDir("/", nil)

	curDir := rootDir
	listingDir := false

	processLine := func(line string) error {
		if len(line) == 0 {
			return errors.New("found empty line")
		}

		if line[0] == '$' {
			listingDir = false
			cmdAndArgs := strings.Split(line[2:], " ")

			cmd := cmdAndArgs[0]
			switch cmd {
			case CommandList:
				listingDir = true

			case CommandChangeDir:
				targetDir := cmdAndArgs[1]

				switch targetDir {
				case "/":
					curDir = rootDir

				case "..":
					if curDir.Parent == nil {
						return fmt.Errorf("cannot go to previous dir since current dir does not have a parent: %+v", curDir)
					}

					curDir = curDir.Parent

				default:
					dir, ok := curDir.Dirs[targetDir]
					if !ok {
						return fmt.Errorf("dir not found or did not run ls before: %s", targetDir)
					}

					curDir = dir
				}

			default:
				return fmt.Errorf("unknown command: %s", cmd)
			}

			return nil
		}

		if listingDir {
			data := strings.Split(line, " ")

			if data[0] == "dir" {
				dirName := data[1]

				if _, ok := curDir.Dirs[dirName]; !ok {
					curDir.Dirs[dirName] = newDir(dirName, curDir)
				}
			} else {
				sizeStr, fileName := data[0], data[1]
				size, err := strconv.ParseInt(sizeStr, 10, 64)
				if err != nil {
					return fmt.Errorf("expected a number as first parameter: %s", line)
				}

				curDir.Files[fileName] = &File{
					Name: fileName,
					Size: size,
				}
			}

			return nil
		}

		return fmt.Errorf("unexpected line: %s", line)
	}

	if err := input.ReadLines(inputName, processLine); err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	return rootDir, nil
}

func main() {
	rootDir, err := readInput()
	if err != nil {
		log.Fatal(err)
	}

	PartOne(rootDir)
	PartTwo(rootDir)
}
