package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const MaxDepthFactor int = 1000

type Direction int

const (
	Left Direction = iota
	Right
)

type Forks = map[string][2]string
type Graph map[string]*GraphNode
type GraphNode struct {
	label string
	left  *GraphNode
	right *GraphNode
}

func main() {
	file := must(os.Open("input"))
	scanner := bufio.NewScanner(file)
	directions, graph := readInput(scanner)

	fmt.Println("Part 1 shortest path:", part1(directions, graph))
	fmt.Println("Part 2 shortest path:", part2(directions, graph))
}

func part1(directions []Direction, graph Graph) int {
	sinkPredicate := func(label string) bool { return "ZZZ" == label }
	return distanceBetween("AAA", sinkPredicate, directions, graph)
}

func part2(directions []Direction, graph Graph) int {
	starts := make([]string, 0)
	for label := range graph {
		if strings.HasSuffix(label, "A") {
			starts = append(starts, label)
		}
	}

	numPaths := len(starts)
	distancesToSink := make([]int, numPaths)
	sinkPredicate := func(label string) bool { return strings.HasSuffix(label, "Z") }

	for i, start := range starts {
		distancesToSink[i] = distanceBetween(start, sinkPredicate, directions, graph)
	}

	return lcm(distancesToSink)
}

// ---------- Helpers ----------
type Predicate func(string) bool

func distanceBetween(startLabel string, sinkPredicate Predicate, directions []Direction, graph Graph) int {
	distance, n := 0, len(directions)
	node := graph[startLabel]

	for !sinkPredicate(node.label) && distance < (n*MaxDepthFactor) {
		node = node.next(directions[distance%n])
		distance++
	}

	return distance
}

func lcm(nums []int) int {
	lowestMultiple := 1

	for _, n := range nums {
		lowestMultiple = (lowestMultiple * n) / gcd(lowestMultiple, n)
	}

	return lowestMultiple
}

func gcd(a int, b int) int {
	for b > 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func readInput(scanner *bufio.Scanner) ([]Direction, Graph) {
	scanner.Scan()
	directionsStr := scanner.Text()
	directions := make([]Direction, len(directionsStr))
	for i, ch := range directionsStr {
		switch ch {
		case 'L':
			directions[i] = Left
		case 'R':
			directions[i] = Right
		}
	}

	scanner.Scan()

	pathForks := make(Forks)
	nodeRegex := regexp.MustCompile("[A-Z]{3}")
	for scanner.Scan() {
		line := scanner.Text()
		nodes := nodeRegex.FindAllString(line, -1)
		origin, left, right := nodes[0], nodes[1], nodes[2]
		pathForks[origin] = [2]string{left, right}
	}

	return directions, NewGraph(pathForks)
}

func NewGraph(forks Forks) Graph {
	graph := make(Graph)

	for nodeLabel, children := range forks {
		leftLabel, rightLabel := children[0], children[1]

		if _, ok := graph[leftLabel]; !ok {
			graph[leftLabel] = &GraphNode{label: leftLabel}
		}

		if _, ok := graph[rightLabel]; !ok {
			graph[rightLabel] = &GraphNode{label: rightLabel}
		}

		if _, ok := graph[nodeLabel]; !ok {
			graph[nodeLabel] = &GraphNode{label: nodeLabel}
		}

		graph[nodeLabel].left = graph[leftLabel]
		graph[nodeLabel].right = graph[rightLabel]
	}

	return graph
}

func (node *GraphNode) next(direction Direction) *GraphNode {
	if direction == Left {
		return node.left
	} else {
		return node.right
	}
}

func must[T any](val T, err any) T {
	if err != nil {
		panic(err)
	}

	return val
}
