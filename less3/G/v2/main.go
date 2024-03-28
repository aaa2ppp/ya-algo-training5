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
	x, y int32
}

func checkSegment(p1, p2 point, ps map[point]struct{}, res []point) []point {
	dx := p2.x - p1.x
	dy := p2.y - p1.y

	dy2 := dx + dy
	if dy2%2 != 0 {
		return res
	}

	dy2 /= 2
	dx2 := (dy - dx) / 2

	p3 := point{p1.x - dx2, p1.y + dy2}
	p4 := point{p2.x + dx2, p2.y - dy2}

	var f, n int

	if _, ok := ps[p3]; !ok {
		f |= 1
		n++
	}

	if _, ok := ps[p4]; !ok {
		f |= 2
		n++
	}

	if n < len(res) {
		res = res[:0]
		if f&1 != 0 {
			res = append(res, p3)
		}
		if f&2 != 0 {
			res = append(res, p4)
		}
	}

	return res
}

func solution(points []point) (res []point) {

	p := points[0] // нам гарантируют как минимум одну точку
	res = []point{
		{p.x, p.y + 1},
		{p.x + 1, p.y},
		{p.x + 1, p.y + 1},
	}

	pointSet := make(map[point]struct{}, len(points))
	for _, p := range points {
		pointSet[p] = struct{}{}
	}

	for i, p1 := range points {
		for _, p2 := range points[i+1:] {
			res = checkSegment(p1, p2, pointSet, res)
			if len(res) == 0 {
				return nil
			}
		}
	}

	return res
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

	points := make([]point, 0, n)

	for i := 0; i < n; i++ {
		x, y, err := scanTwoInt(sc)
		if err != nil {
			return err
		}
		points = append(points, point{int32(x), int32(y)})
	}

	res := solution(points)

	writeInt(bw, len(res))
	bw.WriteByte('\n')

	for _, p := range res {
		writeInts(bw, []int32{p.x, p.y}, " ")
		bw.WriteByte('\n')
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
