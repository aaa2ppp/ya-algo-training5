package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func firstDay(n, k int) int {
	n *= 10

	for i := 0; i < 10; i++ {
		if n2 := n + i; n2%k == 0 {
			return n2
		}
	}

	return -1
}

func run(in io.Reader, out io.Writer) error {

	var n, k, d int
	if _, err := fmt.Fscan(in, &n, &k, &d); err != nil {
		return err
	}

	res := firstDay(n, k)
	if res == -1 {
		fmt.Fprintln(out, res)
		return nil
	}

	w := bufio.NewWriter(out)
	
	w.WriteString(strconv.Itoa(res))
	for i := 1; i < d; i++ {
		w.WriteByte('0')
	}

	w.WriteByte('\n')
	w.Flush()

	return nil
}

func main() {
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
