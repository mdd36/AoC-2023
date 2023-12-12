package main

import (
	"AoC_2023/lib"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Spring int

const (
	Damaged Spring = iota
	Operational
	Unknown
)

type SpringRow struct {
	springs []Spring
	runLens []int
}

type Triplet struct {
	a, b, c int
}

type Cache map[Triplet]int

func main() {
	file := lib.Must(os.Open("input"))
	scanner := bufio.NewScanner(file)
	springRows := readInput(scanner)

	fmt.Println("Total ways to arrange each row:", part1(springRows))
	fmt.Println("Total ways to arrange unfolded rows:", part2(springRows))
}

func part1(rows []SpringRow) int {
	total := 0

	for _, row := range rows {
		total += countWays(row)
	}

	return total
}

func part2(rows []SpringRow) int {
	total := 0

	// There's definitely a mathematical way to combine the answer from part 1,
	// but I'm kinda sick of this problem and brute force still takes under a second
	// with memoization so why bother.
	for _, row := range rows {
		total += countWays(unfold(row))
	}

	return total
}

func countWays(springRow SpringRow) int {
	cache := make(Cache)

	// I don't really feel like spending the time to handle a 3D bottom-up
	// problem, so I'm gonna be lazy and do top-down instead
	var recurse func(int, int, int) int
	recurse = func(s, g, r int) int {
		t := Triplet{a: s, b: g, c: r}

		if v, ok := cache[t]; ok {
			return v
		}

		sp := Operational
		if s < len(springRow.springs) {
			sp = springRow.springs[s]
		}

		if g == len(springRow.runLens) && r > 0 {
			return 0
		}

		gs := 0
		if g < len(springRow.runLens) {
			gs = springRow.runLens[g]
		}

		if r > 0 && r == gs {
			if sp == Damaged {
				return 0
			}
			v := recurse(s+1, g+1, 0)
			cache[t] = v
			return v
		}

		if s == len(springRow.springs)+1 {
			if g == len(springRow.runLens) {
				return 1
			}

			return 0
		}

		ways := 0
		if sp == Damaged || sp == Unknown {
			ways += recurse(s+1, g, r+1)
		}
		if sp == Operational || sp == Unknown {
			if r > 0 && r == gs && g < len(springRow.runLens) {
				ways += recurse(s+1, g+1, 0)
			} else if r == 0 {
				ways += recurse(s+1, g, 0)
			}
		}
		cache[t] = ways
		return ways
	}

	return recurse(0, 0, 0)
}

func unfold(row SpringRow) SpringRow {
	springs, runLens := make([]Spring, 0), make([]int, 0)

	for i := 0; i < 5; i++ {
		springs = append(springs, row.springs...)
		springs = append(springs, Unknown)
		runLens = append(runLens, row.runLens...)
	}

	lib.PopSlice(&springs) // Remove the last ? from the input
	return SpringRow{springs, runLens}
}

func readInput(scanner *bufio.Scanner) []SpringRow {
	rows := make([]SpringRow, 0)

	for scanner.Scan() {
		line := scanner.Text()
		spl := strings.Fields(line)

		springs := make([]Spring, 0)
		for _, ch := range spl[0] {
			switch ch {
			case '#':
				springs = append(springs, Damaged)
			case '.':
				springs = append(springs, Operational)
			case '?':
				springs = append(springs, Unknown)
			}
		}

		runLens := make([]int, 0)
		for _, num := range strings.Split(spl[1], ",") {
			runLens = append(runLens, lib.Must(strconv.Atoi(num)))
		}

		rows = append(rows, SpringRow{springs, runLens})
	}

	return rows
}
