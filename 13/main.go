package main

import (
	"AoC_2023/lib"
	"bufio"
	"fmt"
	"os"
)

type Terrain int

const (
	Empty Terrain = iota
	Ash
	Rock
)

type Landscape [][]Terrain

type Direction int

const (
	Row Direction = iota
	Column
)

func main() {
	file := lib.Must(os.Open("input"))
	scanner := bufio.NewScanner(file)
	landscapes := readInput(scanner)

	fmt.Println("Part 1 mirror summaries:", part1(landscapes))
	fmt.Println("Part 2 mirror summaries:", part2(landscapes))
}

func part1(landscapes []Landscape) int {
	summary := 0

	for _, landscape := range landscapes {
		n, m := len(landscape), len(landscape[0])
		rowPalindromes, colPalindromes := make([][]int, n), make([][]int, m)

		for i := range landscape {
			rowPalindromes[i] = palindromeLengths(landscape, i, Row)
		}

		for j := range landscape[0] {
			colPalindromes[j] = palindromeLengths(landscape, j, Column)
		}

		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				// Even indices are between two rows or columns, which is where all the mirrors must
				// be. Odd indices are palindromes centered on a character from the input, so we
				// can just ignore them.
				rowPalindromes[0][2*j+1] = 0
				colPalindromes[0][2*i+1] = 0

				rowPalindromes[0][2*j] = min(rowPalindromes[0][2*j], rowPalindromes[i][2*j])
				colPalindromes[0][2*i] = min(colPalindromes[0][2*i], colPalindromes[j][2*i])
			}
		}

		for i, l := range rowPalindromes[0] {
			if l > 1 && (i+l == (2*m) || i-l == 0) {
				summary += i / 2
				break
			}
		}

		for j, l := range colPalindromes[0] {
			if l > 1 && (j+l == (2*n) || j-l == 0) {
				summary += 100 * j / 2
				break
			}
		}
	}

	return summary
}

func part2(landscapes []Landscape) int {
	summary := 0

	for _, landscape := range landscapes {
		n, m := len(landscape), len(landscape[0])
		rowPalindromes, colPalindromes := make([][]int, n), make([][]int, m)

		for i := range landscape {
			rowPalindromes[i] = palindromeLengths(landscape, i, Row)
		}

		for j := range landscape[0] {
			colPalindromes[j] = palindromeLengths(landscape, j, Column)
		}

		maxRowLengths, maxColLengths := make([]int, m), make([]int, n)
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				maxRowLengths[j] = max(maxRowLengths[j], rowPalindromes[i][2*j])
				maxColLengths[i] = max(maxColLengths[i], colPalindromes[j][2*i])

				if i > 0 {
					rowPalindromes[0][2*j+1] += (rowPalindromes[i][2*j] * (1 << i))
				} else {
					rowPalindromes[0][2*j+1] = rowPalindromes[0][2*j]
				}

				if j > 0 {
					colPalindromes[0][2*i+1] += (colPalindromes[j][2*i] * (1 << j))
				} else {
					colPalindromes[0][2*i+1] = colPalindromes[0][2*i]
				}
			}
		}

		// TODO
		// Broad idea:
		// We've encoded row/col index information into the summed palindrome lengths, so if
		// one was off we can try flipping the character at (palindromeCenter + palindromeRadius + 1),
		// and re-calculating the palindromes for that row/col. If it's at least as long as the others,
		// then we've found our smudge! We know they must be equal length since the palindrome *must*
		// touch an edge of the puzzle and hence is clamped, so the fixed row/col can't have a longer
		// palindrome than the others and the others must have the same length.

		// Just a little sick of this puzzle, so I'm parking it for now.
	}

	return summary
}

func isPowerOfTwo(n int) bool {
	return n&(n-1) == 0
}

func flip(landscape *Landscape, r, c int) {
	l := *landscape
	if l[r][c] == '#' {
		l[r][c] = '.'
	} else {
		l[r][c] = '#'
	}
}

func palindromeLengths(landscape Landscape, index int, direction Direction) []int {
	interspersed := intersperse(landscape, index, direction)
	center, radius, n := 0, 0, len(interspersed)
	palindromeRadii := make([]int, len(interspersed))

	for center < n {

		for center < n &&
			center+radius+1 < n &&
			center-radius-1 > -1 &&
			interspersed[center-radius-1] == interspersed[center+radius+1] {
			radius += 1
		}

		palindromeRadii[center] = radius

		oldCenter, oldRadius := center, radius
		center, radius = center+1, 0

		for center <= oldCenter+oldRadius {
			mirroredCenter := oldCenter - (center - oldCenter)
			maxMirroredRadius := oldCenter + oldRadius - center

			if palindromeRadii[mirroredCenter] < maxMirroredRadius {
				palindromeRadii[center] = palindromeRadii[mirroredCenter]
				center += 1
			} else if palindromeRadii[mirroredCenter] > maxMirroredRadius {
				palindromeRadii[center] = maxMirroredRadius
				center += 1
			} else {
				radius = maxMirroredRadius
				break
			}
		}
	}

	return palindromeRadii
}

func intersperse(landscape Landscape, index int, direction Direction) []Terrain {
	n := len(landscape)
	if direction == Row {
		n = len(landscape[0])
	}

	interspersed := make([]Terrain, 2*n+1)
	interspersed[0] = Empty
	for i := 0; i < n; i++ {
		var terrain Terrain
		if direction == Row {
			terrain = landscape[index][i]
		} else {
			terrain = landscape[i][index]
		}
		interspersed[2*i+1] = terrain
		interspersed[2*i+2] = Empty
	}

	return interspersed
}

func readInput(scanner *bufio.Scanner) []Landscape {
	landscapes := make([]Landscape, 0)
	area := make(Landscape, 0)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			landscapes = append(landscapes, area)
			area = make(Landscape, 0)
			continue
		}

		row := make([]Terrain, len(line))
		for i, ch := range line {
			switch ch {
			case '.':
				row[i] = Ash
			case '#':
				row[i] = Rock
			}
		}
		area = append(area, row)
	}

	landscapes = append(landscapes, area)
	return landscapes
}
