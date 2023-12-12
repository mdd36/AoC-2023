package main

import (
	"AoC_2023/lib"
	"bufio"
	"fmt"
	"os"
)

type Sector int

const (
	Empty Sector = iota
	Occupied
)

type Coordinate struct {
	x int
	y int
}

type StarChart struct {
	image    [][]Sector
	galaxies []Coordinate
}

func main() {
	file := lib.Must(os.Open("input"))
	scanner := bufio.NewScanner(file)
	starChart := readInput(scanner)

	fmt.Println("Pairwise distance:", part1(starChart))
	fmt.Println("Pairwise distance with expansion:", part2(starChart))
}

func part1(starChart StarChart) int {
	total := 0

	for i, originGalaxy := range starChart.galaxies {
		for _, destinationGalaxy := range starChart.galaxies[i+1:] {
			total += expandedManhattanDistance(originGalaxy, destinationGalaxy, starChart, 1)
		}
	}

	return total
}

func part2(starChart StarChart) int {
	total := 0
	expansionFactor := 1000000 - 1 // Minus 1 to account for the single existing empty row/col
	for i, originGalaxy := range starChart.galaxies {
		for _, destinationGalaxy := range starChart.galaxies[i+1:] {
			total += expandedManhattanDistance(originGalaxy, destinationGalaxy, starChart, expansionFactor)
		}
	}

	return total
}

func expandedManhattanDistance(origin, dest Coordinate, starChart StarChart, expansionFactor int) int {
	dist := 0

	// Just the manhattan distance, plus a some extra constant for every row or column
	// that expanded
	top, bottom := min(origin.y, dest.y), max(origin.y, dest.y)
	for j := top; j < bottom; j++ {
		if starChart.image[j][0] == Empty {
			dist += expansionFactor
		}
		dist++
	}

	left, right := min(origin.x, dest.x), max(origin.x, dest.x)
	for j := left; j < right; j++ {
		if starChart.image[0][j] == Empty {
			dist += expansionFactor
		}
		dist++
	}

	return dist
}

func readInput(scanner *bufio.Scanner) StarChart {
	image := make([][]Sector, 0)
	stars := make([]Coordinate, 0)
	starNumber := Sector(1)

	for row := 0; scanner.Scan(); row++ {
		line := scanner.Text()

		// Insert the row that tracks empty cols
		if row == 0 {
			image = append(image, make([]Sector, len(line)+1))
		}

		// len(line)+1 to add a column to track empty rows
		image = append(image, make([]Sector, len(line)+1))

		for col, ch := range line {
			if ch == '#' {
				image[0][col+1] = Occupied
				image[row+1][0] = Occupied
				image[row+1][col+1] = starNumber
				stars = append(stars, Coordinate{x: col + 1, y: row + 1})
				starNumber++
			}
		}
	}

	return StarChart{image, stars}
}
