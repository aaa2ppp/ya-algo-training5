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
	x, y int
}

type segment struct {
	p1, p2 point
}

func findTwoPnt(p1, p2 point) (p3, p4 point) {
	dx := p2.x - p1.x
	dy := p2.y - p1.y
	p3 = point{p1.x - dy, p1.y + dx}
	p4 = point{p2.x - dy, p2.y + dx}
	return p3, p4
}

func checkSegment(p1, p2 point, ps map[point]struct{}, res []point) []point {
	tmp := make([]point, 2)

	check := func(p1, p2 point) {
		p3, p4 := findTwoPnt(p1, p2)
		tmp = tmp[:0]

		if _, ok := ps[p3]; !ok {
			tmp = append(tmp, p3)
		}

		if _, ok := ps[p4]; !ok {
			tmp = append(tmp, p4)
		}

		if len(tmp) < len(res) {
			res = res[:len(tmp)]
			copy(res, tmp)
		}
	}

	check(p1, p2)
	check(p2, p1)

	return res
}

func solution(points []point) (res []point) {
	n := len(points)

	switch n {
	case 0:
		return []point{
			{0, 0},
			{0, 1},
			{1, 0},
			{1, 1},
		}
	case 1:
		p := points[0]
		return []point{
			{p.x, p.y + 1},
			{p.x + 1, p.y},
			{p.x + 1, p.y + 1},
		}
	default:
		// в худшем случае нужно добавить две точки
		p3, p4 := findTwoPnt(points[0], points[1])
		res = []point{p3, p4}
	}

	pointSet := make(map[point]struct{}, n)
	dists := make(map[int]int, n*(n-1)/2)

	for i, p1 := range points {
		pointSet[p1] = struct{}{}

		for _, p2 := range points[i+1:] {
			dx := (p1.x - p2.x)
			dy := (p1.y - p2.y)
			d2 := dx*dx + dy*dy
			dists[d2]++
		}
	}

	for i, p1 := range points {

		for _, p2 := range points[i+1:] {
			dx := (p1.x - p2.x)
			dy := (p1.y - p2.y)
			d2 := dx*dx + dy*dy
			n1 := dists[d2]
			n2 := dists[d2*2]
			if len(res) == 1 && n1 >= 4 && n2 >= 2 ||
				n1 >= 2 && n2 >= 1 {
				res = checkSegment(p1, p2, pointSet, res)
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
		points = append(points, point{x, y})
	}

	res := solution(points)

	writeInt(bw, len(res))
	bw.WriteByte('\n')

	for _, p := range res {
		writeInts(bw, []int{p.x, p.y}, " ")
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
