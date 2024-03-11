package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

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

func run(in io.Reader, out io.Writer) error {

	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)

	n, err := scanInt(sc)
	if err != nil {
		return err
	}

	// |Xi|, |Yi| ≤ 10^9
	var x1, y1, x2, y2 int = 1e9, 1e9, -1e9, -1e9

	for i := 0; i < n; i++ {

		x, y, err := scanTwoInt(sc)
		if err != nil {
			return err
		}

		x1 = min(x1, x)
		y1 = min(y1, y)

		x2 = max(x2, x)
		y2 = max(y2, y)
	}

	fmt.Fprintln(out, x1, y1, x2, y2)

	return nil
}

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
