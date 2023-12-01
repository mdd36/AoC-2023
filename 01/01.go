package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

func main() {

	file, err := os.Open("01")

	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)

	fmt.Print(numeric_or_spelled(scanner))
}

// Part 1 - Must be a numeric character
func numeric_only(scanner *bufio.Scanner) int {
	total := 0
	for scanner.Scan() {
		first := 0
		second := 0
		last_seen := ' '
		line := scanner.Text()

		for _, c := range line {
			if unicode.IsDigit(c) {
				if last_seen == ' ' {
					first = int(c) - int('0')
				}
				last_seen = c
			}
		}
		second = int(last_seen) - int('0')
		total += 10*first + second
		fmt.Println(10*first + second)
	}

	return total
}

// Part 2 - Could be a numeric character OR a spelled-out number
func numeric_or_spelled(scanner *bufio.Scanner) int {
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
		last_seen := -1
		line := scanner.Text()

		for i, c := range line {
			if unicode.IsDigit(c) {
				last_seen = int(c) - int('0')
			} else {
				spelled_value := match(prefix_trie, i, line)
				if spelled_value > 0 {
					last_seen = spelled_value
				}
			}

			if last_seen >= 0 && first < 0 {
				first = last_seen
			}
		}

		second = last_seen
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

	char_at := string(word[index])
	child, ok := trie.children[char_at]
	if !ok {
		child = new(Trie)
		trie.children[char_at] = child
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

	char_at := string(text[index])
	child, ok := trie.children[char_at]

	if ok {
		return match(child, index+1, text)
	}

	return -1
}
