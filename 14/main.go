package main

import (
	"AoC_2023/lib"
	"bufio"
	"fmt"
	"os"
)

type Rock int

const (
	Empty Rock = iota
	Square
	Round
)

func main() {
	file := lib.Must(os.Open("input"))
	scanner := bufio.NewScanner(file)
	rocks := readInput(scanner)

	fmt.Println("Part 1 load:", part1(rocks))
	fmt.Println("Part 2 load:", part2(rocks))
}

// Rather than altering the [][]Rock, we can just track where the next
// round rock should stop based on the other rocks in the same column
// closer to the north edge.
func part1(rocks [][]Rock) int {
	nextStop := make([]int, len(rocks[0]))
	size, load := len(rocks), 0

	for row, r := range rocks {
		for col, rock := range r {
			switch rock {
			case Square:
				nextStop[col] = row + 1
			case Round:
				load += (size - nextStop[col])
				nextStop[col]++
			}
		}
	}

	return load
}

func part2(rocks [][]Rock) int {
	verticalBlock, horizontalBlock := make([]int, len(rocks)), make([]int, len(rocks[0]))
	rockArrangements := make(map[int]int)
	currentHash := hash(&rocks)
	i := 0

	for cycleFound := false; !cycleFound; _, cycleFound = rockArrangements[currentHash] {
		rockArrangements[currentHash] = i
		i++

		cycle(&rocks, &verticalBlock, &horizontalBlock)
		currentHash = hash(&rocks)
	}

	// Cycle length = total # of arrangements seen - the index of the entrypoint to the cycle
	cycleLen := len(rockArrangements) - rockArrangements[currentHash]

	// How far through the cycle is the 1,000,000,000th rock arrangement?
	cycleIndex := (1_000_000_000 - rockArrangements[currentHash]) % cycleLen

	for i := 0; i < cycleIndex; i++ {
		cycle(&rocks, &verticalBlock, &horizontalBlock)
	}

	// Okay, looks like altering the [][]Rock was possibly the way to go for part 1
	load := 0
	for i := 0; i < len(rocks); i++ {
		for j := 0; j < len(rocks[0]); j++ {
			if rocks[i][j] == Round {
				load += len(rocks) - i
			}
		}
	}

	return load
}

func cycle(r *[][]Rock, verticalBlock, horizontalBlock *[]int) {
	vb, hb, rocks := *verticalBlock, *horizontalBlock, *r
	verticalSize, horizontalSize := len(rocks), len(rocks[0])

	// Tilt north
	lib.Fill(verticalBlock, 0)
	for i := 0; i < verticalSize; i++ {
		for j := 0; j < horizontalSize; j++ {
			switch rocks[i][j] {
			case Square:
				vb[j] = i + 1
			case Round:
				rocks[i][j] = Empty
				rocks[vb[j]][j] = Round
				vb[j]++
			}
		}
	}

	// Tilt west
	lib.Fill(horizontalBlock, 0)
	for j := 0; j < horizontalSize; j++ {
		for i := 0; i < verticalSize; i++ {
			switch rocks[i][j] {
			case Square:
				hb[i] = j + 1
			case Round:
				rocks[i][j] = Empty
				rocks[i][hb[i]] = Round
				hb[i]++
			}
		}
	}

	// Tilt south
	lib.Fill(verticalBlock, verticalSize-1)
	for i := verticalSize - 1; i > -1; i-- {
		for j := 0; j < horizontalSize; j++ {
			switch rocks[i][j] {
			case Square:
				vb[j] = i - 1
			case Round:
				rocks[i][j] = Empty
				rocks[vb[j]][j] = Round
				vb[j]--
			}
		}
	}

	// Tilt east
	lib.Fill(horizontalBlock, horizontalSize-1)
	for j := horizontalSize - 1; j > -1; j-- {
		for i := 0; i < verticalSize; i++ {
			switch rocks[i][j] {
			case Square:
				hb[i] = j - 1
			case Round:
				rocks[i][j] = Empty
				rocks[i][hb[i]] = Round
				hb[i]--
			}
		}
	}
}

func hash(r *[][]Rock) int {
	rocks := *r
	hash := 2166136261
	prime := 16777619

	for i := 0; i < len(rocks); i++ {
		for j := 0; j < len(rocks[0]); j++ {
			hash = (hash ^ int(rocks[i][j]) ^ (i * j)) * prime
		}
	}

	return hash
}

func readInput(scanner *bufio.Scanner) [][]Rock {
	rocks := make([][]Rock, 0)

	for scanner.Scan() {
		line := scanner.Text()
		row := make([]Rock, len(line))

		for i, ch := range line {
			switch ch {
			case '.':
				row[i] = Empty
			case '#':
				row[i] = Square
			case 'O':
				row[i] = Round
			}
		}
		rocks = append(rocks, row)
	}

	return rocks
}
