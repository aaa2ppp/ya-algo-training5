package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func run(in io.Reader, out io.Writer) error {
	var s1, s2 string
	var f int
	if _, err := fmt.Fscan(in, &s1, &s2, &f); err != nil {
		return err
	}

	prev1, prev2, err := parseMatchScore(s1)
	if err != nil {
		return err
	}

	cur1, cur2, err := parseMatchScore(s2)
	if err != nil {
		return err
	}

	res := cur2 + prev2 - (cur1 + prev1) + 1

	if f == 1 && (cur1+res-1) > prev2 {
		res--
	}

	if f == 2 && prev1 > cur2 {
		res--
	}

	if res < 0 {
		res = 0
	}

	fmt.Fprintln(out, res)

	return nil
}

func parseMatchScore(s string) (int, int, error) {
	a := strings.Split(s, ":")

	g1, err := strconv.Atoi(a[0])
	if err != nil {
		return 0, 0, err
	}

	g2, err := strconv.Atoi(a[1])
	if err != nil {
		return 0, 0, err
	}

	return g1, g2, nil
}

func main() {
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
