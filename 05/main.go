package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const FORWARDS int = 1
const BACKWARDS int = -1

type SearchExtractor[T any] func(T) (int, int)
type KeyExtractor[T any] func(T) int

type Interval struct {
	sourceStart int
	destStart   int
	length      int
}

func sourceRangeExtractor(interval Interval) (int, int) {
	return interval.sourceStart, interval.sourceStart + interval.length
}

func destRangeExtractor(interval Interval) (int, int) {
	return interval.destStart, interval.destStart + interval.length
}

type Almanac struct {
	seeds                 []int
	seedToSoil            []Interval
	soilToFertilizer      []Interval
	fertilizerToWater     []Interval
	waterToLight          []Interval
	lightToTemperature    []Interval
	temperatureToHumidity []Interval
	humidityToLocation    []Interval
}

func main() {
	file, err := os.Open("input")

	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	almanac := buildAlmanac(scanner)

	sortAlmanac(almanac, func(interval Interval) int { return interval.sourceStart })
	fmt.Println("Closest plot to plant is", closestSeedPlot(almanac))
	sortAlmanac(almanac, func(interval Interval) int { return interval.destStart })
	fmt.Println("Closest plot for seed ranges is", closestSeedPlotWithRange(almanac))
}

// Part 1
func closestSeedPlot(almanac Almanac) int {
	closest := math.MaxInt
	for _, seed := range almanac.seeds {
		closest = min(closest, findPlotFor(seed, almanac))
	}

	return closest
}

// Part 2
func closestSeedPlotWithRange(almanac Almanac) int {
	seedRanges := make([][2]int, 0)

	for i := 0; i < len(almanac.seeds); i += 2 {
		start, length := almanac.seeds[i], almanac.seeds[i+1]
		seedRanges = append(seedRanges, [2]int{start, start + length})
	}

	quicksort(seedRanges, 0, len(seedRanges)-1, func(i [2]int) int { return i[0] })

	foundSeeds := make([]int, 0)
	for location := 0; location <= 51399229; location++ {
		humidity := getNext(location, almanac.humidityToLocation, BACKWARDS)
		temperature := getNext(humidity, almanac.temperatureToHumidity, BACKWARDS)
		light := getNext(temperature, almanac.lightToTemperature, BACKWARDS)
		water := getNext(light, almanac.waterToLight, BACKWARDS)
		fertilizer := getNext(water, almanac.fertilizerToWater, BACKWARDS)
		soil := getNext(fertilizer, almanac.soilToFertilizer, BACKWARDS)
		seed := getNext(soil, almanac.seedToSoil, BACKWARDS)
		foundSeeds = append(foundSeeds, seed)
		if haveSeed(seed, seedRanges) {
			return location
		}
	}

	return -1
}

// ------- Helpers -------
func haveSeed(seed int, seedRanges [][2]int) bool {
	return nil != binarySearch(seedRanges, seed, func(i [2]int) (int, int) { return i[0], i[1] })
}

func sortAlmanac(almanac Almanac, extractor KeyExtractor[Interval]) {
	quicksort(almanac.seedToSoil, 0, len(almanac.seedToSoil)-1, extractor)
	quicksort(almanac.soilToFertilizer, 0, len(almanac.soilToFertilizer)-1, extractor)
	quicksort(almanac.fertilizerToWater, 0, len(almanac.fertilizerToWater)-1, extractor)
	quicksort(almanac.waterToLight, 0, len(almanac.waterToLight)-1, extractor)
	quicksort(almanac.lightToTemperature, 0, len(almanac.lightToTemperature)-1, extractor)
	quicksort(almanac.temperatureToHumidity, 0, len(almanac.temperatureToHumidity)-1, extractor)
	quicksort(almanac.humidityToLocation, 0, len(almanac.humidityToLocation)-1, extractor)
}

func findPlotFor(seed int, almanac Almanac) int {
	soil := getNext(seed, almanac.seedToSoil, FORWARDS)
	fertilizer := getNext(soil, almanac.soilToFertilizer, FORWARDS)
	water := getNext(fertilizer, almanac.fertilizerToWater, FORWARDS)
	light := getNext(water, almanac.waterToLight, FORWARDS)
	temperature := getNext(light, almanac.lightToTemperature, FORWARDS)
	humidity := getNext(temperature, almanac.temperatureToHumidity, FORWARDS)
	return getNext(humidity, almanac.humidityToLocation, FORWARDS)
}

func getNext(target int, intervals []Interval, direction int) int {
	extractor := sourceRangeExtractor
	if direction == BACKWARDS {
		extractor = destRangeExtractor
	}
	interval := binarySearch(intervals, target, extractor)

	if interval != nil {
		if direction == FORWARDS {
			offset := target - interval.sourceStart
			return interval.destStart + offset
		} else {
			offset := target - interval.destStart
			return interval.sourceStart + offset
		}
	} else {
		return target
	}
}

func buildAlmanac(scanner *bufio.Scanner) Almanac {
	scanner.Scan()
	seedsStrs := strings.Fields(strings.Split(scanner.Text(), ":")[1])
	seeds := make([]int, len(seedsStrs))
	for i, seed := range seedsStrs {
		seedInt, err := strconv.Atoi(seed)
		if err != nil {
			panic(err)
		}
		seeds[i] = seedInt
	}
	scanner.Scan() // Skip the blank line

	seedToSoil := createIntervals(scanner)
	soilToFertilizer := createIntervals(scanner)
	fertilizerToWater := createIntervals(scanner)
	waterToLight := createIntervals(scanner)
	lightToTemperature := createIntervals(scanner)
	temperatureToHumidity := createIntervals(scanner)
	humidityToLocation := createIntervals(scanner)

	return Almanac{
		seeds,
		seedToSoil,
		soilToFertilizer,
		fertilizerToWater,
		waterToLight,
		lightToTemperature,
		temperatureToHumidity,
		humidityToLocation,
	}
}

func createIntervals(scanner *bufio.Scanner) []Interval {
	scanner.Scan() // skip name line
	ret := make([]Interval, 0)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		spl := strings.Fields(line)
		sourceStart, err := strconv.Atoi(spl[1])
		if err != nil {
			panic(err)
		}

		destStart, err := strconv.Atoi(spl[0])
		if err != nil {
			panic(err)
		}

		length, err := strconv.Atoi(spl[2])
		if err != nil {
			panic(err)
		}

		ret = append(ret, Interval{sourceStart, destStart, length})
	}

	return ret
}

func quicksort[T any](rangeArr []T, start int, end int, key KeyExtractor[T]) {
	if start >= len(rangeArr) || start >= end {
		return
	}

	pivot := rangeArr[start]
	swapIndex := start
	swap(rangeArr, start, end)

	for i := start; i < end; i++ {
		if key(rangeArr[i]) < key(pivot) {
			swap(rangeArr, i, swapIndex)
			swapIndex++
		}
	}

	swap(rangeArr, swapIndex, end)

	quicksort(rangeArr, start, swapIndex-1, key)
	quicksort(rangeArr, swapIndex+1, end, key)
}

func swap[T any](arr []T, i int, j int) {
	tmp := arr[i]
	arr[i] = arr[j]
	arr[j] = tmp
}

func binarySearch[T any](rangeArr []T, target int, extractor SearchExtractor[T]) *T {
	n := len(rangeArr)
	l, r := 0, n-1

	for l <= r {
		mid := l + (r-l)/2
		midRange := rangeArr[mid]
		start, end := extractor(midRange)

		if start <= target && target <= end {
			return &midRange
		}

		if end < target {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}

	return nil
}
