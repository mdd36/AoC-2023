package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("input")

	if err != nil {
		panic(err)
	}

	// fmt.Println(countPossible(bufio.NewScanner(file)))
	fmt.Println(minRequired(bufio.NewScanner(file)))
}

// Part 1
func countPossible(scanner *bufio.Scanner) int {
	numPossible := 0

	// See my comment below for why this can't be a constant -- dang it Go!
	totalInBag := make(map[string]int)
	totalInBag["green"] = 13
	totalInBag["red"] = 12
	totalInBag["blue"] = 14

	for scanner.Scan() {
		line := scanner.Text()

		colonSplit := strings.Split(line, ":")
		gameInformation, gameActions := colonSplit[0], colonSplit[1]

		gameId, err := strconv.Atoi(strings.Split(gameInformation, " ")[1])

		if err != nil {
			panic(err)
		}

		counts := countMaxSeen(gameActions)
		possible := true

		possible = possible && counts["green"] <= totalInBag["green"]
		possible = possible && counts["red"] <= totalInBag["red"]
		possible = possible && counts["blue"] <= totalInBag["blue"]

		if possible {
			numPossible += gameId
		}
	}

	return numPossible
}

// Part 2
func minRequired(scanner *bufio.Scanner) int {
	total := 0

	for scanner.Scan() {
		line := scanner.Text()
		counts := countMaxSeen(strings.Split(line, ":")[1])
		total += counts["green"] * counts["red"] * counts["blue"]
	}

	return total
}

func countMaxSeen(round string) map[string]int {

	// If this was Rust I'd make a constant array of colors and iterate over it,
	// but since Go doesn't include mutability in its type system it can't do a
	// compile time check to enforce that values in constant arrays are never changed
	// and hence it just doesn't support constant arrays. Big thumbs down from me.
	colorCounts := make(map[string]int)
	colorCounts["red"] = 0
	colorCounts["blue"] = 0
	colorCounts["green"] = 0

	for _, round := range strings.Split(round, ";") {
		for _, cubeVariant := range strings.Split(round, ",") {
			spaceSplit := strings.Split(strings.Trim(cubeVariant, " "), " ")
			numTakenStr, color := spaceSplit[0], spaceSplit[1]
			numTakenInt, err := strconv.Atoi(numTakenStr)
			if err != nil {
				panic(err)
			}

			// Ignoring the potential missing check since we know that all the colors
			// are in the map -- we put them there before the loops! If we dynamically
			// inserted, out map might be missing a color if we happened to never draw
			// it from the bag
			existingCount, _ := colorCounts[color]
			colorCounts[color] = max(existingCount, numTakenInt)
		}
	}

	return colorCounts
}
