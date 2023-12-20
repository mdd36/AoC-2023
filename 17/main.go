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

type State struct {
	row       int
	col       int
	direction Direction
	steps     int
}

type Costs = map[State]int

func (direction Direction) Turns() []Direction {
	switch direction {
	case Left, Right:
		return []Direction{Up, Down, direction}
	case Up, Down:
		return []Direction{Left, Right, direction}
	}
	return nil
}

type Step struct {
	estimatedTotalLoss int
	row                int
	col                int
	direction          Direction
	steps              int
	loss               int
}

func main() {
	file := lib.Must(os.Open("input"))
	scanner := bufio.NewScanner(file)
	maze := readInput(scanner)

	fmt.Println("Part 1 minimum heat loss:", part1(maze))
	fmt.Println("Part 2 minimum heat loss:", part2(maze))
}

func part1(maze [][]int) int {
	return astar(maze, 1, 3)
}

func part2(maze [][]int) int {
	return astar(maze, 4, 10)
}

func astar(maze [][]int, minStep, maxStep int) int {
	n, m := len(maze), len(maze[0])
	costs := make(Costs)
	pq := lib.NewHeap(func(s Step) int { return s.estimatedTotalLoss })
	pq.Push(Step{
		estimatedTotalLoss: n + m,
		direction:          Right,
	})
	pq.Push(Step{
		estimatedTotalLoss: n + m,
		direction:          Down,
	})

	for pq.Len() > 0 {
		node := pq.Pop()

		row, col, loss := node.row, node.col, node.loss
		direction, steps := node.direction, node.steps

		state := State{row, col, direction, steps}

		if cost, found := costs[state]; found && cost <= loss {
			continue
		}
		costs[state] = loss

		if row == n-1 && col == m-1 {
			return loss
		}

		for _, nextDirection := range direction.Turns() {
			if nextDirection == direction && steps == maxStep {
				continue
			}

			if steps < minStep && nextDirection != direction {
				continue
			}

			s := 1
			if nextDirection == direction {
				s = steps + 1
			}

			nextRow, nextCol := row, col
			switch nextDirection {
			case Up:
				nextRow -= 1
			case Right:
				nextCol += 1
			case Down:
				nextRow += 1
			case Left:
				nextCol -= 1
			}

			if nextRow < 0 || nextRow == n || nextCol < 0 || nextCol == m {
				continue
			}

			nextStep := Step{
				estimatedTotalLoss: loss + (n - nextRow) + (m - nextCol),
				row:                nextRow,
				col:                nextCol,
				direction:          nextDirection,
				steps:              s,
				loss:               loss + maze[nextRow][nextCol],
			}
			pq.Push(nextStep)
		}
	}

	return -1
}

func readInput(scanner *bufio.Scanner) [][]int {
	maze := make([][]int, 0)

	for scanner.Scan() {
		line := scanner.Text()
		row := make([]int, len(line))

		for i, coolingRate := range line {
			row[i] = int(coolingRate) - int('0')
		}
		maze = append(maze, row)
	}

	return maze
}
