package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func count(a int) int {
	n := a / 4
	switch a % 4 {
	case 0:
		// +0
	case 1:
		n++
	case 2:
		n += 2
	case 3:
		n += 2
	}
	return n
}

func run(in io.Reader, out io.Writer) error {
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return err
	}

	sum := 0
	for i := 0; i < n; i++ {
		var a int
		if _, err := fmt.Fscan(in, &a); err != nil {
			return err
		}
		sum += count(a)
	}

	fmt.Fprintln(out, sum)

	return nil
}

func main() {
	err := run(bufio.NewReader(os.Stdin), os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
