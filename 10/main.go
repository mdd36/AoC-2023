package main

import (
	"AoC_2023/lib"
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
)

type Connection int

const (
	Up Connection = iota
	Right
	Down
	Left
)

type Maze [][][]Connection

type Coordinate struct {
	row int
	col int
}

func main() {
	file := lib.Must(os.Open("input"))
	scanner := bufio.NewScanner(file)
	start, maze := readInput(scanner)

	fmt.Println("Part 1 max distance:", part1(start, maze))
	fmt.Println("Part 2 enclosed area:", part2(start, maze))
}

func part1(start Coordinate, maze Maze) int {
	pointsOnLoop := dfs(start, maze)
	return int(math.Round(float64(len(pointsOnLoop)) / 2.0))
}

func part2(start Coordinate, maze Maze) int {
	// An implementation of the point-in-polygon algorithm
	// Needed a refresher, so I'm assuming future me will need it also:
	// https://en.wikipedia.org/wiki/Point_in_polygon
	total := 0
	inside := false
	pointsOnLoop := dfs(start, maze)

	for row := range maze {
		for col := range maze[row] {
			point := Coordinate{row, col}
			connections := maze[row][col]
			isVertical := slices.Contains(connections, Up)
			if pointsOnLoop.Contains(point) && isVertical {
				inside = !inside
			}

			if !pointsOnLoop.Contains(point) && inside {
				total++
			}
		}
	}

	return total
}

func dfs(start Coordinate, maze Maze) lib.Set[Coordinate] {
	n, m := len(maze), len(maze[0])
	stack := lib.NewStack[Coordinate]()
	seen := lib.NewSet[Coordinate]()

	stack.Append(start)

	for len(stack) > 0 {
		cell := stack.Pop()
		row, col := cell.row, cell.col

		if seen.Contains(cell) || row < 0 || row >= n || col < 0 || col >= m {
			continue
		}

		seen.Add(cell)

		for _, connection := range maze[cell.row][cell.col] {
			switch connection {
			case Up:
				stack.Append(Coordinate{row: row - 1, col: col})
			case Right:
				stack.Append(Coordinate{row: row, col: col + 1})
			case Down:
				stack.Append(Coordinate{row: row + 1, col: col})
			case Left:
				stack.Append(Coordinate{row: row, col: col - 1})
			}
		}
	}

	return seen
}

func readInput(scanner *bufio.Scanner) (Coordinate, Maze) {
	maze := make(Maze, 0)
	start := Coordinate{}

	for rowNum := 0; scanner.Scan(); rowNum++ {
		line := scanner.Text()
		row := make([][]Connection, len(line))
		for colNum, ch := range line {
			row[colNum] = connectionsOf(ch)

			if ch == 'S' {
				start.row = rowNum
				start.col = colNum
			}
		}
		maze = append(maze, row)
	}

	for row := 0; row < len(maze); row++ {
		for col := 0; col < len(maze[row]); col++ {
			validConnections := make([]Connection, 0)

			if row-1 > -1 &&
				slices.Contains(maze[row][col], Up) &&
				slices.Contains(maze[row-1][col], Down) {
				validConnections = append(validConnections, Up)
			}

			if row+1 < len(maze) &&
				slices.Contains(maze[row][col], Down) &&
				slices.Contains(maze[row+1][col], Up) {
				validConnections = append(validConnections, Down)
			}

			if col-1 > -1 &&
				slices.Contains(maze[row][col], Left) &&
				slices.Contains(maze[row][col-1], Right) {
				validConnections = append(validConnections, Left)
			}

			if col+1 < len(maze[row]) &&
				slices.Contains(maze[row][col], Right) &&
				slices.Contains(maze[row][col+1], Left) {
				validConnections = append(validConnections, Right)
			}

			maze[row][col] = validConnections
		}
	}

	return start, maze
}

func connectionsOf(pipe rune) []Connection {
	switch pipe {
	case '|':
		return []Connection{Up, Down}
	case '-':
		return []Connection{Left, Right}
	case 'L':
		return []Connection{Up, Right}
	case 'J':
		return []Connection{Up, Left}
	case '7':
		return []Connection{Left, Down}
	case 'F':
		return []Connection{Right, Down}
	case 'S':
		return []Connection{Up, Right, Down, Left}
	default:
		return []Connection{}
	}
}
