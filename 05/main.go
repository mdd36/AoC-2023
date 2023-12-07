package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

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

	outputChannel := make(chan int)

	for _, seedRange := range seedRanges {
		go processRange(seedRange, almanac, outputChannel)
	}

	lowest := math.MaxInt
	for i := 0; i < len(seedRanges); i++ {
		select {
		case l := <-outputChannel:
			lowest = min(lowest, l)
		}
	}
	return lowest
}

// ------- Helpers -------
func processRange(seedRange [2]int, almanac Almanac, channel chan int) {
	lowest := math.MaxInt

	for seed := seedRange[0]; seed <= seedRange[1]; seed++ {
		x := findPlotFor(seed, almanac)
		lowest = min(lowest, x)
	}

	channel <- lowest
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
	soil := getNext(seed, almanac.seedToSoil)
	fertilizer := getNext(soil, almanac.soilToFertilizer)
	water := getNext(fertilizer, almanac.fertilizerToWater)
	light := getNext(water, almanac.waterToLight)
	temperature := getNext(light, almanac.lightToTemperature)
	humidity := getNext(temperature, almanac.temperatureToHumidity)
	return getNext(humidity, almanac.humidityToLocation)
}

func getNext(target int, intervals []Interval) int {
	extractor := sourceRangeExtractor
	interval := binarySearch(intervals, target, extractor)

	if interval != nil {
		offset := target - interval.sourceStart
		return interval.destStart + offset
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

		if start <= target && target < end {
			return &midRange
		}

		if end <= target {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}

	return nil
}
