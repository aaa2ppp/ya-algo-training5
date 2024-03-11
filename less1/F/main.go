package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
)

func run(in io.Reader, out io.Writer) error {

	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)

	sc.Scan()
	n, err := strconv.Atoi(sc.Text())
	if err != nil {
		return err
	}

	oddCountIsEven := true
	firstOddPos := -1

	for i := 0; i < n; i++ {

		sc.Scan()
		a, err := strconv.Atoi(sc.Text())
		if err != nil {
			return err
		}

		if a&1 == 1 {
			oddCountIsEven = !oddCountIsEven 
			if firstOddPos == -1 {
				firstOddPos = i
			}
		}
	}

	w := bufio.NewWriter(out)

	i := 1
	if oddCountIsEven {
		for n := firstOddPos+1; i < n; i++ {
			w.WriteByte('+')
		}
		w.WriteByte('x')
		i++
	}

	for ; i < n; i++ {
		w.WriteByte('+')
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
