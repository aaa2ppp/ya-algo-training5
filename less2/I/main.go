package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func run(in io.Reader, out io.Writer) error {

	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)

	n, err := scanInt(sc)
	if err != nil {
		return err
	}

	rows := make([]int, n)
	cols := make([]int, n)

	for i := 0; i < n; i++ {
		row, col, err := scanTwoInt(sc)
		if err != nil {
			return err
		}

		rows[row-1]++
		cols[col-1]++
	}

	if debugEnable {
		log.Println("rows:", rows, "cols:", cols)
	}

	steps := 0
	k := 0

	for i := range rows {
		steps += k

		if rows[i] == 1 {
			continue
		}

		if rows[i] == 0 {
			k++
			continue
		}

		m := min(k, rows[i]-1)
		k -= m
		rows[i] -= m

		if rows[i] > 1 {
			m := rows[i] - 1
			rows[i+1] += m
			rows[i] = 1
			steps += m
		}
	}

	if debugEnable {
		log.Println("rows:", rows, "steps:", steps, "k:", k)
	}

	for i, j := 0, len(cols)-1; i < j; {
		if cols[i] < cols[j] {
			steps += cols[i]
			cols[i+1] += cols[i]
			cols[i] = 0
			i++
		} else {
			steps += cols[j]
			cols[j-1] += cols[j] 
			cols[j] = 0
			j--
		}
	}

	if debugEnable {
		log.Println("cols:", cols, "steps:", steps)
	}

	fmt.Fprintln(out, steps)

	return nil
}

func scanInt(sc *bufio.Scanner) (int, error) {
	sc.Scan()
	return strconv.Atoi(sc.Text())
}

func scanTwoInt(sc *bufio.Scanner) (v1, v2 int, err error) {
	v1, err = scanInt(sc)
	if err == nil {
		v2, err = scanInt(sc)
	}
	return v1, v2, err
}

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
