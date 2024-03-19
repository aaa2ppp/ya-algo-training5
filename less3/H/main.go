package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

type point struct {
	x, y int16
}

type segment struct {
	p point
	d point
}

func makeSegment(x1, y1, x2, y2 int) segment {
	if x1 > x2 || x1 == x2 && y1 > y2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}
	dx := int16(x2 - x1)
	dy := int16(y2 - y1)
	return segment{point{int16(x1), int16(y1)}, point{dx, dy}}
}

func solve(a, b []segment) int {
	n := len(a)

	aDict := make(map[point][]point, n)
	for _, s := range a {
		aDict[s.d] = append(aDict[s.d], s.p)
	}

	offsetCount := make(map[point]int16, n*n)

	for _, s := range b {
		for _, p2 := range aDict[s.d] {
			offset := point{p2.x - s.p.x, p2.y - s.p.y}
			offsetCount[offset]++
		}
	}

	maximum := 0
	for _, count := range offsetCount {
		maximum = max(maximum, int(count))
	}

	return n - maximum
}

func scanSegments(sc *bufio.Scanner, ss []segment) error {
	for i := range ss {
		x1, y1, x2, y2, err := scanFourInt(sc)
		if err != nil {
			return err
		}
		ss[i] = makeSegment(x1, y1, x2, y2)
	}
	return nil
}

func run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, err := scanInt(sc)
	if err != nil {
		return err
	}

	a := make([]segment, n)
	if err := scanSegments(sc, a); err != nil {
		return err
	}

	b := make([]segment, n)
	if err := scanSegments(sc, b); err != nil {
		return err
	}

	res := solve(a, b)

	writeInt(bw, res)
	bw.WriteByte('\n')

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

func scanThreeInt(sc *bufio.Scanner) (v1, v2, v3 int, err error) {
	v1, err = scanInt(sc)
	if err == nil {
		v2, err = scanInt(sc)
	}
	if err == nil {
		v3, err = scanInt(sc)
	}
	return v1, v2, v3, err
}

func scanFourInt(sc *bufio.Scanner) (v1, v2, v3, v4 int, err error) {
	v1, err = scanInt(sc)
	if err == nil {
		v2, err = scanInt(sc)
	}
	if err == nil {
		v3, err = scanInt(sc)
	}
	if err == nil {
		v4, err = scanInt(sc)
	}
	return v1, v2, v3, v4, err
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

type Int interface {
	~int | ~int64 | ~int32 | ~int16 | ~int8
}

func writeInt[I Int](bw *bufio.Writer, v I) error {
	var buf [32]byte
	_, err := bw.Write(strconv.AppendInt(buf[:0], int64(v), 10))
	return err
}

func writeInts[I Int](bw *bufio.Writer, a []I, sep string) error {
	if len(a) == 0 {
		return nil
	}

	var buf [32]byte

	_, err := bw.Write(strconv.AppendInt(buf[:0], int64(a[0]), 10))
	for i := 1; err == nil && i < len(a); i++ {
		_, err = bw.WriteString(sep)
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
