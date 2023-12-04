package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("input")

	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	winningNumbers := make([]map[string]bool, 0)
	ticketNumbers := make([][]string, 0)

	for scanner.Scan() {
		line := scanner.Text()
		ticketWinners, numbersOnTicket := parseTicket(line)
		winningNumbers = append(winningNumbers, ticketWinners)
		ticketNumbers = append(ticketNumbers, numbersOnTicket)
	}

	fmt.Println("Total score:", countWinnings(&winningNumbers, &ticketNumbers))
	fmt.Println("Total cards:", totalCards(&winningNumbers, &ticketNumbers))
}

// Part 1
func countWinnings(winningNumbersArr *[]map[string]bool, ticketNumbersArr *[][]string) int {
	total := 0

	for i, ticketNums := range *ticketNumbersArr {
		winningNums := (*winningNumbersArr)[i]
		countInWinning := 0
		for _, num := range ticketNums {
			if winningNums[num] {
				countInWinning++
			}
		}
		if countInWinning > 0 {
			toAdd := int(math.Pow(2, float64(countInWinning-1)))
			total += toAdd
		}
	}

	return total
}

// Part 2
func totalCards(winningNumbersArr *[]map[string]bool, ticketNumbersArr *[][]string) int {
	numCards := make([]int, len(*ticketNumbersArr))
	for i := range numCards {
		numCards[i] = 1
	}

	for i, ticketNumbers := range *ticketNumbersArr {
		winningNumbers := (*winningNumbersArr)[i]
		numFound := 0

		for _, n := range ticketNumbers {
			if winningNumbers[n] {
				numFound++
			}
		}

		for j := i + 1; j < len(numCards) && j <= i+numFound; j++ {
			numCards[j] += numCards[i] // One extra j for each copy of i we have
		}
	}

	sum := 0
	for _, n := range numCards {
		sum += n
	}

	return sum
}

func parseTicket(line string) (map[string]bool, []string) {
	numbers := strings.Split(strings.Split(line, ":")[1], "|")
	winningNumbersSet := make(map[string]bool) // Why is there no set type?? (╯°□°)╯︵ ┻━┻

	for _, numStr := range strings.Fields(numbers[0]) {
		winningNumbersSet[strings.Trim(numStr, " ")] = true
	}

	return winningNumbersSet, strings.Fields(numbers[1])
}
