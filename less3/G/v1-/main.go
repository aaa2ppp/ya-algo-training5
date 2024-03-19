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

type point struct {
	x, y int32
}

type checker struct {
	pnts map[point]struct{}
	x    map[int32][]point
	y    map[int32][]point
	xy1  map[int32][]point
	xy2  map[int32][]point
}

func newChecker(n int) checker {
	return checker{
		pnts: make(map[point]struct{}, n),
		x:    make(map[int32][]point, n),
		y:    make(map[int32][]point, n),
		xy1:  make(map[int32][]point, n),
		xy2:  make(map[int32][]point, n),
	}
}

func (w *checker) addPnt(p point) {
	w.pnts[p] = struct{}{}
	w.x[p.y] = append(w.x[p.y], p)
	w.y[p.x] = append(w.y[p.x], p)
	w.xy1[p.x-p.y] = append(w.xy1[p.x-p.y], p)
	w.xy2[p.x+p.y] = append(w.xy2[p.x+p.y], p)
}

func (w *checker) check() []point {
	res := []point{{0, 0}, {0, 1}, {1, 0}, {1, 1}}

	for p := range w.pnts {
		res = w.checkPnt(p, res)
		if len(res) == 0 {
			break
		}
	}

	return res
}

func (w *checker) checkPnt(p point, res []point) []point {

	if len(res) > 3 {
		res = res[:3]
		res[0] = point{p.x + 1, p.y}
		res[1] = point{p.x, p.y + 1}
		res[2] = point{p.x + 1, p.y + 1}
	}

	res, w.x[p.y] = w.checkPntAxis(p, w.x[p.y], w.checkX, res)
	if len(res) > 0 {
		res, w.y[p.x] = w.checkPntAxis(p, w.y[p.x], w.checkY, res)
	}
	if len(res) > 0 {
		res, w.xy1[p.x-p.y] = w.checkPntAxis(p, w.xy1[p.x-p.y], w.checkXY, res)
	}
	if len(res) > 0 {
		res, w.xy2[p.x+p.y] = w.checkPntAxis(p, w.xy2[p.x+p.y], w.checkXY, res)
	}
	if len(res) > 0 {
		delete(w.pnts, p)
	}

	return res
}

func (w *checker) checkPntAxis(
	p point,
	axis []point,
	check func(point, point, []point) []point,
	res []point,
) ([]point, []point) {

	delete := func(i int) []point {
		n := len(axis) - 1
		axis[i] = axis[n]
		return axis[:n]
	}

	for i := len(axis) - 1; i >= 0; i-- {
		if axis[i] == p {
			axis = delete(i)
			continue
		}
		res = check(p, axis[i], res)
		if len(res) == 0 {
			break
		}
	}

	return res, axis
}

func (w *checker) checkX(p1, p2 point, res []point) []point {

	if p1.y != p2.y {
		panic(fmt.Sprintf("checkX: p1.y must be equal p2.y, got %v, %v", p1, p2))
	}

	tmp := make([]point, 2)

	d := p2.x - p1.x
	for _, y := range []int32{p1.y + d, p1.y - d} {
		tmp = tmp[:0]

		for _, x := range []int32{p1.x, p2.x} {
			p := point{x, y}

			if _, ok := w.pnts[p]; !ok {
				tmp = append(tmp, p)
			}
		}

		if len(tmp) < len(res) {
			res = res[:len(tmp)]
			copy(res, tmp)
		}

		if len(res) == 0 {
			break
		}
	}

	return res
}

func (w *checker) checkY(p1, p2 point, res []point) []point {

	if p1.x != p2.x {
		panic(fmt.Sprintf("checkY: p1.x must be equal p2.x, got %v, %v", p1, p2))
	}

	tmp := make([]point, 2)

	d := p2.y - p1.y
	for _, x := range []int32{p1.x + d, p1.x - d} {
		tmp = tmp[:0]

		for _, y := range []int32{p1.y, p2.y} {
			p := point{x, y}

			if _, ok := w.pnts[p]; !ok {
				tmp = append(tmp, p)
			}
		}

		if len(tmp) < len(res) {
			res = res[:len(tmp)]
			copy(res, tmp)
		}

		if len(res) == 0 {
			break
		}
	}

	return res
}

func abs(a int32) int32 {
	if a < 0 {
		return -a
	}
	return a
}

func (w *checker) checkXY(p1, p2 point, res []point) []point {

	if abs(p1.x-p2.x) != abs(p1.y-p2.y) {
		panic(fmt.Sprintf("checkXY: points must be on the same diagonal, got %v, %v", p1, p2))
	}

	tmp := make([]point, 0, 2)

	for _, p := range []point{{p1.x, p2.y}, {p2.x, p1.y}} {
		if _, ok := w.pnts[p]; !ok {
			tmp = append(tmp, p)
		}
	}

	if len(tmp) < len(res) {
		res = res[:len(tmp)]
		copy(res, tmp)
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

	w := newChecker(n)

	for i := 0; i < n; i++ {
		x, y, err := scanTwoInt(sc)
		if err != nil {
			return err
		}
		w.addPnt(point{int32(x), int32(y)})
	}

	res := w.check()

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
