package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type HandType int

const (
	HighCard HandType = iota + 1
	OnePair
	TwoPair
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

type Hand struct {
	bet      int
	handType HandType
	cards    [5]int
}

func main() {
	file := must(os.Open("input"))

	scanner := bufio.NewScanner(file)
	hands := readInput(scanner)

	fmt.Println("Part 1 total winnings:", part1(hands))
	fmt.Println("Part 2 total winnings:", part2(hands))
}

func part1(hands []Hand) int {
	total := 0
	quicksort(hands, 0, len(hands)-1)

	for i, hand := range hands {
		total += (i + 1) * hand.bet
	}

	return total
}

func part2(hands []Hand) int {
	jokerifiedHand := make([]Hand, len(hands))

	for i, hand := range hands {
		jokerifiedHand[i] = hand.rescoreWithJoker()
	}

	quicksort(jokerifiedHand, 0, len(hands)-1)
	total := 0
	for i, hand := range jokerifiedHand {
		total += (i + 1) * hand.bet
	}
	return total
}

// -------- Helpers --------
func quicksort(hands []Hand, left int, right int) {
	if left >= len(hands) || left >= right {
		return
	}

	pivot := hands[left]
	swapDex := left
	swap(hands, left, right)

	for i := left; i < right; i++ {
		if hands[i].compareTo(&pivot) < 0 {
			swap(hands, i, swapDex)
			swapDex++
		}
	}
	swap(hands, swapDex, right)
	quicksort(hands, left, swapDex-1)
	quicksort(hands, swapDex+1, right)
}

func swap(hands []Hand, i int, j int) {
	tmp := hands[i]
	hands[i] = hands[j]
	hands[j] = tmp
}

func NewHand(cardsStr string, betStr string) Hand {
	cards := [5]int{}

	for i, c := range cardsStr {
		switch c {
		case 'A':
			cards[i] = 14
		case 'K':
			cards[i] = 13
		case 'Q':
			cards[i] = 12
		case 'J':
			cards[i] = 11
		case 'T':
			cards[i] = 10
		default:
			cards[i] = int(c) - int('0')
		}
	}

	cardCounts := make(map[int]int)
	for _, card := range cards {
		cardCounts[card] += 1
	}

	groupsOfSize := [6]int{}
	for _, count := range cardCounts {
		groupsOfSize[count] += 1
	}

	handType := scoreHandGroups(groupsOfSize)
	bet := must(strconv.Atoi(betStr))

	return Hand{bet, handType, cards}
}

func (hand *Hand) compareTo(other *Hand) int {
	delta := (int)(hand.handType) - (int)(other.handType)

	if delta != 0 {
		return delta
	}

	for i := 0; i < 5; i++ {
		delta = hand.cards[i] - other.cards[i]
		if delta != 0 {
			return delta
		}
	}

	return 0
}

func (hand *Hand) rescoreWithJoker() Hand {
	newCardValues := [5]int{}
	for i := 0; i < 5; i++ {
		if hand.cards[i] == 11 {
			newCardValues[i] = 1
		} else {
			newCardValues[i] = hand.cards[i]
		}
	}

	cardCounts := make(map[int]int)
	for _, card := range newCardValues {
		cardCounts[card] += 1
	}

	groupsOfSize := [6]int{}
	for card, count := range cardCounts {
		if card == 1 {
			continue
		}
		groupsOfSize[count] += 1
	}

	numJokers := cardCounts[1]
	if numJokers > 0 {
		for i := 5; i > -1; i-- {
			if groupsOfSize[i] > 0 || i == 0 {
				groupsOfSize[i+numJokers] = 1
				groupsOfSize[i] = max(0, groupsOfSize[i]-1)
				break
			}
		}
	}

	cards := newCardValues
	handType := scoreHandGroups(groupsOfSize)
	bet := hand.bet

	return Hand{bet, handType, cards}
}

func scoreHandGroups(groupsOfSize [6]int) HandType {
	handType := HighCard
	if groupsOfSize[5] == 1 {
		handType = FiveOfAKind
	} else if groupsOfSize[4] == 1 {
		handType = FourOfAKind
	} else if groupsOfSize[3] == 1 && groupsOfSize[2] == 1 {
		handType = FullHouse
	} else if groupsOfSize[3] == 1 {
		handType = ThreeOfAKind
	} else if groupsOfSize[2] == 2 {
		handType = TwoPair
	} else if groupsOfSize[2] == 1 {
		handType = OnePair
	}

	return handType
}

func readInput(scanner *bufio.Scanner) []Hand {
	hands := make([]Hand, 0)

	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		cards, bet := line[0], line[1]
		hands = append(hands, NewHand(cards, bet))
	}
	return hands
}

func must[T any](val T, err any) T {
	if err != nil {
		panic(err)
	}
	return val
}
