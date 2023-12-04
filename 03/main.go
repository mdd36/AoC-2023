package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("input")

	if err != nil {
		panic(err)
	}

	grid := createGrid(bufio.NewScanner(file))

	fmt.Println(gearRatio(&grid))
}

// Part 2
func gearRatio(gridPtr *[][]rune) int {
	total := 0
	grid := *gridPtr

	for i, row := range grid {
		for j, symbol := range row {
			if symbol != '*' {
				// Not a possible gear
				continue
			}

			neighboringNumbers := make([]int, 8)
			for i, _ := range neighboringNumbers {
				neighboringNumbers[i] = -1
			}

			if i > 0 {
				lastRow := &grid[i-1]
				neighboringNumbers[0] = readNumber(lastRow, j-1)
				neighboringNumbers[1] = readNumber(lastRow, j)
				neighboringNumbers[2] = readNumber(lastRow, j+1)
			}

			if i+1 < len(grid) {
				nextRow := &grid[i+1]
				neighboringNumbers[3] = readNumber(nextRow, j-1)
				neighboringNumbers[4] = readNumber(nextRow, j)
				neighboringNumbers[5] = readNumber(nextRow, j+1)
			}

			neighboringNumbers[6] = readNumber(&row, j-1)
			neighboringNumbers[7] = readNumber(&row, j+1)

			ratio := 1
			found := 0
			for _, val := range neighboringNumbers {
				if val > -1 {
					found += 1
					ratio *= val
				}
			}

			if found == 2 {
				total += ratio
			}
		}
	}

	return total
}

// Part 1
func sumTouchingSymbol(gridPtr *[][]rune) int {
	total := 0
	grid := *gridPtr
	for i, row := range grid {
		for j, symbol := range row {
			if symbol == '.' || parseInt(symbol) > -1 {
				// Is a number or a period, so it can't force us to include its
				// neighbors
				continue
			}

			if i > 0 {
				lastRow := &grid[i-1]
				total += max(0, readNumber(lastRow, j-1))
				total += max(0, readNumber(lastRow, j))
				total += max(0, readNumber(lastRow, j+1))
			}

			if i+1 < len(grid) {
				nextRow := &grid[i+1]
				total += max(0, readNumber(nextRow, j-1))
				total += max(0, readNumber(nextRow, j))
				total += max(0, readNumber(nextRow, j+1))
			}

			total += max(0, readNumber(&row, j-1))
			total += max(0, readNumber(&row, j+1))
		}
	}

	return total
}

func readNumber(gridPtr *[]rune, start int) int {
	grid := *gridPtr

	if start < 0 || start >= len(grid) || parseInt(grid[start]) < 0 {
		return -1
	}

	left := start
	for left > -1 && parseInt(grid[left]) > -1 {
		left -= 1
	}

	left += 1
	total := 0
	for left < len(grid) && parseInt(grid[left]) > -1 {
		total = total*10 + parseInt(grid[left])
		grid[left] = '.' // Avoid double counts!
		left += 1
	}

	return total
}

func parseInt(r rune) int {
	val := int(r) - int('0')
	if val < 0 || val > 9 {
		return -1
	}
	return val
}

func createGrid(scanner *bufio.Scanner) [][]rune {
	grid := make([][]rune, 0, 10)

	for scanner.Scan() {
		line := scanner.Text()
		lineArr := make([]rune, len(line))
		for i, ch := range line {
			lineArr[i] = ch
		}
		grid = append(grid, lineArr)
	}

	return grid
}
