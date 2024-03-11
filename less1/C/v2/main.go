package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

func count(a int) int {
	return a/4 + [4]int{0, 1, 2, 2}[a%4]
}

func run(in io.Reader, out io.Writer) error {
	scan := newScanner(in)

	n, err := scan.Int()
	if err != nil {
		return err
	}

	sum := 0
	for i := 0; i < n; i++ {
		a, err := scan.Int()
		if err != nil {
			return err
		}
		sum += count(a)
	}

	fmt.Fprintln(out, sum)

	return nil
}

func main() {
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

type scanner struct {
	*bufio.Scanner
}

func newScanner(in io.Reader) scanner {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	return scanner{sc}
}

func orEOF(err error) error {
	if err == nil {
		return io.EOF
	}
	return err
}

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func (sc scanner) Int() (int, error) {
	if sc.Scan() {
		return strconv.Atoi(unsafeString(sc.Bytes()))
	}
	return 0, orEOF(sc.Err())
}

func (sc scanner) TwoInt() (v1, v2 int, err error) {
	if sc.Scan() {
		v1, err = strconv.Atoi(unsafeString(sc.Bytes()))
	}
	if err == nil && sc.Scan() {
		v2, err = strconv.Atoi(unsafeString(sc.Bytes()))
		return
	}
	if err == nil {
		err = orEOF(sc.Err())
	}
	return
}

func (sc scanner) ThreeInt() (v1, v2, v3 int, err error) {
	if sc.Scan() {
		v1, err = strconv.Atoi(unsafeString(sc.Bytes()))
	}
	if err == nil && sc.Scan() {
		v2, err = strconv.Atoi(unsafeString(sc.Bytes()))
	}
	if err == nil && sc.Scan() {
		v3, err = strconv.Atoi(unsafeString(sc.Bytes()))
		return
	}
	if err == nil {
		err = orEOF(sc.Err())
	}
	return
}
