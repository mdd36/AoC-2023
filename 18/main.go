package main

import (
	"AoC_2023/lib"
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Direction int

const (
	Right Direction = iota
	Down
	Left
	Up
)

type Color uint

type Edge struct {
	direction Direction
	length    int
	color     Color
}

func (e Edge) fixEdge() Edge {
	distance, direction := e.color&0xfffffff0, e.color&0x0000000f
	return Edge{
		direction: Direction(direction),
		length:    int(distance) >> 4,
	}
}

func main() {
	file := lib.Must(os.Open("input"))
	scanner := bufio.NewScanner(file)
	edges := readInput(scanner)

	fmt.Println("Part 1 enclosed area:", part1(edges))
	fmt.Println("Part 2 enclosed area:", part2(edges))
}

func part1(edges []Edge) int {
	return greens(edges)
}

func part2(errantEdges []Edge) int {
	fixedEdges := make([]Edge, len(errantEdges))
	for i, e := range errantEdges {
		fixedEdges[i] = e.fixEdge()
	}

	return greens(fixedEdges)
}

func greens(edges []Edge) int {
	area := 1
	height := 0

	for _, e := range edges {
		if e.direction == Up {
			area += e.length
			height += e.length
		} else if e.direction == Down {
			height -= e.length
		} else if e.direction == Right {
			area += (height + 1) * e.length
		} else {
			area -= height * e.length
		}
	}

	return area
}

func readInput(scanner *bufio.Scanner) []Edge {
	edges := make([]Edge, 0)
	hexPattern := regexp.MustCompile("[0-9a-f]+")

	for scanner.Scan() {
		line := scanner.Text()
		spl := strings.Fields(line)
		var direction Direction
		switch spl[0] {
		case "U":
			direction = Up
		case "R":
			direction = Right
		case "D":
			direction = Down
		case "L":
			direction = Left
		}

		length := lib.Must(strconv.Atoi(spl[1]))

		colorHex := hexPattern.FindString(spl[2])
		color := Color(lib.Must(strconv.ParseInt(colorHex, 16, 32)))

		edge := Edge{direction, length, color}
		edges = append(edges, edge)
	}

	return edges
}
