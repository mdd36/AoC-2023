package main

import (
	"AoC_2023/lib"
	"bufio"
	"fmt"
	"os"
)

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

type OpticalElement int

const (
	None OpticalElement = iota
	LeftSlantMirror
	RightSlantMirror
	VerticalSplitter
	HorizontalSplitter
)

func (direction Direction) clockwise() Direction {
	return Direction((direction + 1) % 4)
}

func (direction Direction) anticlockwise() Direction {
	return Direction((direction - 1 + 4) % 4)
}

func (direction Direction) isHorizontal() bool {
	return direction == Left || direction == Right
}

func (element OpticalElement) Next(in Direction) []Direction {
	switch element {
	case None:
		return []Direction{in}
	case LeftSlantMirror:
		if in.isHorizontal() {
			return []Direction{in.anticlockwise()}
		} else {
			return []Direction{in.clockwise()}
		}
	case RightSlantMirror:
		if in.isHorizontal() {
			return []Direction{in.clockwise()}
		} else {
			return []Direction{in.anticlockwise()}
		}
	case HorizontalSplitter:
		if in.isHorizontal() {
			return []Direction{in}
		} else {
			return []Direction{Left, Right}
		}
	case VerticalSplitter:
		if in.isHorizontal() {
			return []Direction{Up, Down}
		} else {
			return []Direction{in}
		}
	}
	return []Direction{}
}

type Location struct {
	row             int
	col             int
	travelDirection Direction
}

func main() {
	file := lib.Must(os.Open("input"))
	scanner := bufio.NewScanner(file)
	elements := readInput(scanner)

	fmt.Println("Part 1 energized squares", part1(elements))
	fmt.Println("Part 2 max energized squares", part2(elements))
}

func part1(elements [][]OpticalElement) int {
	start := Location{row: 0, col: 0, travelDirection: Right}
	return energizeSquares(elements, start)
}

func part2(elements [][]OpticalElement) int {
	n, m := len(elements), len(elements[0])
	max_energized := 0

	// TODO There should be away to reuse the number of
	// squares energized by entering a particular square from a
	// particular direction across runs, which would hugely reduce
	// duplicated work. But since brute force is already pretty fast
	// (under a second), I'm going to call this good for now.

	for i := 0; i < n; i++ {
		left_start := Location{row: i, col: 0, travelDirection: Right}
		right_start := Location{row: i, col: m - 1, travelDirection: Left}

		max_energized = max(
			max_energized,
			energizeSquares(elements, left_start),
			energizeSquares(elements, right_start),
		)
	}

	for j := 0; j < m; j++ {
		top_start := Location{row: 0, col: j, travelDirection: Down}
		bottom_start := Location{row: n - 1, col: j, travelDirection: Up}
		max_energized = max(
			max_energized,
			energizeSquares(elements, top_start),
			energizeSquares(elements, bottom_start),
		)
	}

	return max_energized
}

func energizeSquares(elements [][]OpticalElement, start Location) int {
	n, m := len(elements), len(elements[0])
	seen := lib.NewSet[Location]()
	stack := lib.NewStack[Location]()
	stack.Append(start)

	energyGrid := make([][]int, n)
	for i := 0; i < n; i++ {
		energyGrid[i] = make([]int, m)
	}

	for len(stack) > 0 {
		location := stack.Pop()
		row, col, direction := location.row, location.col, location.travelDirection
		if seen.Contains(location) ||
			row < 0 || row >= n ||
			col < 0 || col >= m {
			continue
		}

		seen.Add(location)
		element := elements[row][col]
		energyGrid[row][col] = 1

		for _, d := range element.Next(direction) {
			switch d {
			case Left:
				stack.Append(Location{row: row, col: col - 1, travelDirection: Left})
			case Right:
				stack.Append(Location{row: row, col: col + 1, travelDirection: Right})
			case Up:
				stack.Append(Location{row: row + 1, col: col, travelDirection: Up})
			case Down:
				stack.Append(Location{row: row - 1, col: col, travelDirection: Down})
			}
		}
	}

	energized := 0
	for _, r := range energyGrid {
		for _, v := range r {
			energized += v
		}
	}
	return energized
}

func readInput(scanner *bufio.Scanner) [][]OpticalElement {
	elements := make([][]OpticalElement, 0)

	for scanner.Scan() {
		line := scanner.Text()
		row := make([]OpticalElement, len(line))
		for i, ch := range line {
			switch ch {
			case '.':
				row[i] = None
			case '/':
				row[i] = RightSlantMirror
			case '\\':
				row[i] = LeftSlantMirror
			case '|':
				row[i] = VerticalSplitter
			case '-':
				row[i] = HorizontalSplitter
			}
		}
		elements = append(elements, row)
	}
	return elements
}
