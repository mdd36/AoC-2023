package main

import (
	"AoC_2023/lib"
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	file := lib.Must(os.Open("input"))
	scanner := bufio.NewScanner(file)
	sequences := readInput(scanner)
	constants := make([][]float64, len(sequences))

	for i, seq := range sequences {
		constants[i] = determineConstants(seq)
	}

	fmt.Println("Next element sums:", part1(constants))
	fmt.Println("Previous element sums:", part2(constants))
}

func part1(constants [][]float64) int {

	total := 0
	for _, sequenceConstants := range constants {
		total += evalDiscrete(sequenceConstants, len(sequenceConstants))
	}

	return total
}

func part2(constants [][]float64) int {
	total := 0
	for _, sequenceConstants := range constants {
		// We love when part1's impl answers part2 as well!
		total += evalDiscrete(sequenceConstants, -1)
	}
	return total
}

func determineConstants(sequence []int) []float64 {
	// Discrete calculus time! As described in the puzzle, we can
	// make a table of finite differences, and from the we can
	// then re-integrate them to get a general formula for the
	// function. We'll give back the array of constants for the
	// the formula ordered by increasing polynomial power.
	maxDegree := len(sequence)
	derivatives := sequence
	initialConditions := make([]int, maxDegree)
	for degree := 0; degree < maxDegree; degree++ {
		initialConditions[degree] = derivatives[0]
		derivatives = discreteDerivative(derivatives)
	}

	return discreteIntegral(initialConditions, maxDegree)
}

func discreteDerivative(sequence []int) []int {
	finiteDifferences := make([]int, len(sequence)-1)

	for i := 0; i < len(sequence)-1; i++ {
		finiteDifferences[i] = sequence[i+1] - sequence[i]
	}

	return finiteDifferences
}

func discreteIntegral(initialConditions []int, maxDegree int) []float64 {
	constants := make([]float64, maxDegree)
	for i := maxDegree - 1; i > -1; i-- {
		constants[i] = float64(initialConditions[i])
		for j := i + 1; j < maxDegree; j++ {
			// Polynomial integration: ∫a*x^n dx = (a/n)*x^(n+1) + C
			// Same holds in the discrete domain, so divide each constant
			// by its current power to get the new constant.
			constants[j] /= float64(j - i)
		}
	}

	return constants
}

// Find the n-th point of the discrete function described by the
// constants array
func evalDiscrete(constants []float64, n int) int {
	total := 0

	for i, constant := range constants {
		polynomialValue := 1.0
		for j := 0; j < i; j++ {
			polynomialValue *= float64(n - j)
		}
		total += int(math.Round(constant * polynomialValue))
	}

	return total
}

func readInput(scanner *bufio.Scanner) [][]int {
	sequences := make([][]int, 0)

	for scanner.Scan() {
		line := scanner.Text()
		seq := make([]int, 0)
		for _, ch := range strings.Fields(line) {
			seq = append(seq, lib.Must(strconv.Atoi(ch)))
		}
		sequences = append(sequences, seq)
	}

	return sequences
}
