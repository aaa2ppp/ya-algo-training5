package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func run(in io.Reader, out io.Writer) error {

	var p, v, q, m int
	if _, err := fmt.Fscan(in, &p, &v, &q, &m); err != nil {
		return err
	}

	p1 := p - v
	p2 := p + v

	q1 := q - m
	q2 := q + m

	var res int
	if p2 >= q1 && q2 >= p1 {
		res = max(p2, q2) - min(p1, q1) + 1
	} else {
		res = (p2 - p1) + (q2 - q1) + 2
	}

	fmt.Fprintln(out, res)

	return nil
}

func main() {
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
