package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

func main() {

	file, err := os.Open("input")

	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)

	fmt.Print(numericOrSpelled(scanner))
}

// Part 1 - Must be a numeric character
func numericOnly(scanner *bufio.Scanner) int {
	total := 0
	for scanner.Scan() {
		first := 0
		second := 0
		lastSeen := ' '
		line := scanner.Text()

		for _, c := range line {
			if unicode.IsDigit(c) {
				if lastSeen == ' ' {
					first = int(c) - int('0')
				}
				lastSeen = c
			}
		}
		second = int(lastSeen) - int('0')
		total += 10*first + second
		fmt.Println(10*first + second)
	}

	return total
}

// Part 2 - Could be a numeric character OR a spelled-out number
func numericOrSpelled(scanner *bufio.Scanner) int {
	prefix_trie := new(Trie)
	insert(prefix_trie, 0, "one", 1)
	insert(prefix_trie, 0, "two", 2)
	insert(prefix_trie, 0, "three", 3)
	insert(prefix_trie, 0, "four", 4)
	insert(prefix_trie, 0, "five", 5)
	insert(prefix_trie, 0, "six", 6)
	insert(prefix_trie, 0, "seven", 7)
	insert(prefix_trie, 0, "eight", 8)
	insert(prefix_trie, 0, "nine", 9)

	total := 0
	for scanner.Scan() {
		first := -1
		second := -1
		lastSeen := -1
		line := scanner.Text()

		for i, c := range line {
			if unicode.IsDigit(c) {
				lastSeen = int(c) - int('0')
			} else {
				spelledValue := match(prefix_trie, i, line)
				if spelledValue > 0 {
					lastSeen = spelledValue
				}
			}

			if lastSeen >= 0 && first < 0 {
				first = lastSeen
			}
		}

		second = lastSeen
		total += 10*first + second
		fmt.Println(10*first + second)
	}

	return total
}

// Is a trie overkill? Yes. Am I doing it anyway? Also yes.
type Trie struct {
	value    int
	children map[string]*Trie
}

func insert(trie *Trie, index int, word string, value int) {
	if trie.children == nil {
		trie.children = make(map[string]*Trie)
	}

	if index == len(word) {
		trie.value = value
		return
	}

	charAt := string(word[index])
	child, ok := trie.children[charAt]
	if !ok {
		child = new(Trie)
		trie.children[charAt] = child
	}

	insert(child, index+1, word, value)
}

func match(trie *Trie, index int, text string) int {
	if trie.value > 0 {
		return trie.value
	}

	if index >= len(text) {
		return -1
	}

	charAt := string(text[index])
	child, ok := trie.children[charAt]

	if ok {
		return match(child, index+1, text)
	}

	return -1
}
