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

func solution(a []int, buf []int) []int {
	minimum := a[0]
	size := 0

	for _, v := range a {
		size++
		minimum = min(minimum, v)

		if minimum < size {
			buf = append(buf, size-1)
			minimum = v
			size = 1
		}
	}

	buf = append(buf, size)

	return buf
}

func task(sc *bufio.Scanner, bw *bufio.Writer, buf []int) ([]int, error) {
	n, err := scanInt(sc)
	if err != nil {
		return buf, err
	}

	a := make([]int, n)
	if err := scanInts(sc, a); err != nil {
		return buf, err
	}

	buf = solution(a, buf)

	writeInt(bw, len(buf))
	bw.WriteByte('\n')
	writeInts(bw, buf, " ")
	bw.WriteByte('\n')

	return buf, nil
}

func run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)

	bw := bufio.NewWriter(out)
	defer bw.Flush()

	t, err := scanInt(sc)
	if err != nil {
		return err
	}

	buf := make([]int, 1000) // жадный я, в детстве было мало памяти
	for i := 0; i < t; i++ {
		buf, err = task(sc, bw, buf[:0])
		if err != nil {
			return fmt.Errorf("task%d: %w", i+1, err)
		}
	}

	return nil
}

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func scanInt(sc *bufio.Scanner) (int, error) {
	sc.Scan()
	return strconv.Atoi(unsafeString(sc.Bytes()))
}

func scanTwoInt(sc *bufio.Scanner) (v1, v2 int, err error) {
	v1, err = scanInt(sc)
	if err == nil {
		v2, err = scanInt(sc)
	}
	return v1, v2, err
}

func scanInts(sc *bufio.Scanner, a []int) error {
	for i := range a {
		v, err := scanInt(sc)
		if err != nil {
			return err
		}
		a[i] = v
	}
	return nil
}

func writeInt(bw *bufio.Writer, v int) error {
	var buf [32]byte
	_, err := bw.Write(strconv.AppendInt(buf[:0], int64(v), 10))
	return err
}

func writeInts(bw *bufio.Writer, a []int, delm string) error {
	if len(a) == 0 {
		return nil
	}

	var buf [32]byte

	_, err := bw.Write(strconv.AppendInt(buf[:0], int64(a[0]), 10))
	for i := 1; err == nil && i < len(a); i++ {
		_, err = bw.WriteString(delm)
		if err == nil {
			_, err = bw.Write(strconv.AppendInt(buf[:0], int64(a[i]), 10))
		}
	}

	return err
}

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
