package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("input")

	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	races := readInput(scanner)

	fmt.Println("Ways to win divided races:", part1(races))
	fmt.Println("Ways to win one big race:", part2(races))
}

type Race struct {
	time int
	dist int
}

func part1(races []Race) int {
	waysProd := 1
	for _, race := range races {
		waysProd *= ways(race)
	}
	return waysProd
}

func part2(races []Race) int {
	totalTime := 0
	totalDist := 0

	for _, race := range races {
		t, d := race.time, race.dist
		tNumDigits := int(math.Floor(math.Log10(float64(t)) + 1))
		dNumDigits := int(math.Floor(math.Log10(float64(d)) + 1))
		totalTime = (int(math.Pow10(tNumDigits)) * totalTime) + t
		totalDist = (int(math.Pow10(dNumDigits)) * totalDist) + d
	}
	return ways(Race{time: totalTime, dist: totalDist})
}

// ------- Helpers -------
func ways(race Race) int {
	// Simple quadratic. We can find the zeros of t(T-t)-d=0,
	// where t is the time we hold the button, T is the maximum
	// race time, and d is the best distance. Round the lower root
	// up, the upper root down, and find the distance between those
	// integers
	T, d := (float64)(race.time), (float64)(race.dist)
	lZero := math.Floor(((T - math.Sqrt(math.Pow(T, 2)-4*d)) / 2) + 1)
	rZero := math.Ceil(((T + math.Sqrt(math.Pow(T, 2)-4*d)) / 2) - 1)
	return int(rZero-lZero) + 1
}

func readInput(scanner *bufio.Scanner) []Race {
	scanner.Scan()
	times := strings.Fields(strings.Split(scanner.Text(), ":")[1])
	scanner.Scan()
	distances := strings.Fields(strings.Split(scanner.Text(), ":")[1])

	races := make([]Race, len(times))
	for i, t := range times {
		time := must(strconv.Atoi(t))
		dist := must(strconv.Atoi(distances[i]))
		races[i] = Race{time, dist}
	}

	return races
}

func must[T any](val T, err any) T {
	if err != nil {
		panic(err)
	}
	return val
}
