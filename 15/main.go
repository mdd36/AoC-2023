package main

import (
	"AoC_2023/lib"
	"bufio"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const HASH_PRIME int = 17
const HASH_MODULO int = 256

type Lense struct {
	label       string
	focalLength int
}

func main() {
	file := lib.Must(os.Open("input"))
	scanner := bufio.NewScanner(file)
	steps := readInput(scanner)

	fmt.Println("Part 1 hash:", part1(steps))
	fmt.Println("Part 2 focusing power:", part2(steps))
}

func part1(steps []string) int {
	total := 0

	for _, s := range steps {
		total += hash(s)
	}

	return total
}

func part2(steps []string) int {
	focusingPower := 0

	boxes := make([][]Lense, 256)
	for i := 0; i < 256; i++ {
		boxes[i] = make([]Lense, 0)
	}

	pattern := regexp.MustCompile("([a-z]+)(=|-)([0-9]*)")
	for _, s := range steps {
		groups := pattern.FindStringSubmatch(s)
		label, operation := groups[1], groups[2] // groups[0] is the whole string
		boxNo := hash(label)

		if operation == "-" {
			remove(&boxes[boxNo], label)
		} else {
			focalLength := lib.Must(strconv.Atoi(groups[3]))
			lense := Lense{label, focalLength}
			insert(&boxes[boxNo], lense)
		}
	}

	for i, box := range boxes {
		for j, lense := range box {
			focusingPower += (i + 1) * (j + 1) * lense.focalLength
		}
	}

	return focusingPower
}

func insert(box *[]Lense, lense Lense) {
	b := *box
	if index := slices.IndexFunc(b, hasLabel(lense.label)); index > -1 {
		b[index] = lense
	} else {
		*box = append(*box, lense)
	}
}

func remove(box *[]Lense, label string) {
	b := *box
	if index := slices.IndexFunc(b, hasLabel(label)); index > -1 {
		for i := index; i < len(b)-1; i++ {
			b[i] = b[i+1]
		}
		lib.PopSlice(box)
	}
}

type LensePredicate func(Lense) bool

func hasLabel(label string) LensePredicate {
	return func(l Lense) bool {
		return l.label == label
	}
}

func hash(s string) int {
	hash := 0

	for _, ch := range s {
		hash += int(ch)
		hash *= HASH_PRIME
		hash = hash % HASH_MODULO
	}

	return hash
}

func readInput(scanner *bufio.Scanner) []string {
	scanner.Scan()
	line := scanner.Text()

	return strings.Split(line, ",")
}
