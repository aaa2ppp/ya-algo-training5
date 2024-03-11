package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func mark(desk [][]byte) {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			switch desk[i][j] {
			case 'R':
				markR(desk, i, j)
			case 'B':
				markB(desk, i, j)
			}
		}
	}
}

func isFigure(c byte) bool {
	return c == 'B' || c == 'R'
}

func markR(desk [][]byte, i, j int) {
	for i := i + 1; i < 8; i++ {
		if isFigure(desk[i][j]) {
			break
		}
		desk[i][j] = '+'
	}
	for i := i - 1; i >= 0; i-- {
		if isFigure(desk[i][j]) {
			break
		}
		desk[i][j] = '+'
	}
	for j := j + 1; j < 8; j++ {
		if isFigure(desk[i][j]) {
			break
		}
		desk[i][j] = '+'
	}
	for j := j - 1; j >= 0; j-- {
		if isFigure(desk[i][j]) {
			break
		}
		desk[i][j] = '+'
	}
}

func markB(desk [][]byte, i, j int) {
	for i, j := i+1, j+1; i < 8 && j < 8; i, j = i+1, j+1 {
		if isFigure(desk[i][j]) {
			break
		}
		desk[i][j] = '+'
	}
	for i, j := i-1, j+1; i >= 0 && j < 8; i, j = i-1, j+1 {
		if isFigure(desk[i][j]) {
			break
		}
		desk[i][j] = '+'
	}
	for i, j := i+1, j-1; i < 8 && j >= 0; i, j = i+1, j-1 {
		if isFigure(desk[i][j]) {
			break
		}
		desk[i][j] = '+'
	}
	for i, j := i-1, j-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
		if isFigure(desk[i][j]) {
			break
		}
		desk[i][j] = '+'
	}
}

func run(in io.Reader, out io.Writer) error {

	sc := bufio.NewScanner(in)
	buf := make([]byte, 64)
	desk := make([][]byte, 8)

	for i, j := 0, 0; i < 8; i, j = i+1, j+8 {
		if !sc.Scan() {
			return fmt.Errorf("сликом мало строк %d", i)
		}
		if len(sc.Bytes()) < 8 {
			return fmt.Errorf("сликом короткая строка %d", i)
		}
		row := buf[j : j+8]
		copy(row, sc.Bytes())
		desk[i] = row
	}

	mark(desk)

	count := 0
	for _, c := range buf {
		if c == '*' {
			count++
		}
	}

	fmt.Fprintln(out, count)

	return nil
}

func main() {
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
