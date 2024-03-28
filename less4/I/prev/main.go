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
	x, y float64
}

type gamer struct {
	x, y, v float64
}

const eps = 0.001

// возвращает значение true, если точка достижима хотя бы для одного игрока за время t
func checkPoint(gamers []gamer, p point, t float64) bool {

	for _, g := range gamers {
		r := g.v * t //
		dx := g.x - p.x
		dy := g.y - p.y

		if r*r-(dx*dx+dy*dy) > 0 {
			return true
		}
	}

	return false
}

// возвращает значение true, если все точки достижимы хотя бы для одного игрока за время t
func checkPoints(gamers []gamer, points []point, t float64) bool {

	for _, p := range points {
		if !checkPoint(gamers, p, t) {
			return false
		}
	}

	return true
}

func appendFronRect(points []point, d float64, x1, y1, x2, y2, step float64) []point {
	d2 := d * d

	for x := x1; x-x2 < eps; x += step {

		for y := y1; y-y2 < eps; y += step {

			if y < -eps || (x*x+y*y)-d2 > eps {
				// пропускаем, если точка за пределами броска
				continue
			}

			points = append(points, point{x, y})
		}
	}

	return points
}

func unavailablePoints(gamers []gamer, points []point, t float64) []point {
	var res []point

	for _, p := range points {
		if !checkPoint(gamers, p, t) {
			res = append(res, p)
		}
	}

	return res
}

func searchTimes(gamers []gamer, points []point, t1, t2, teps float64) (tt1, tt2 float64) {

	for t2-t1 > teps {
		t := (t2 + t1) / 2

		if checkPoints(gamers, points, t) {
			t2 = t
		} else {
			t1 = t
		}
	}

	return t1, t2
}

func solve(d int, gamers []gamer) (t, x, y float64) {
	dd := float64(d)
	x1, y1 := -dd, 0.0
	x2, y2 := dd, dd
	t1, t2 := 0.0, 2000.0

	const scale = 10 // увеличиваем разрешение на каждой итерации
	// почему 10? 100x200 не много не мало?

	step := dd / scale
	points := appendFronRect(nil, dd, x1, y1, x2, y2, step)

	for {
		t1, t2 = searchTimes(gamers, points, t1, t2*2, step/1000) // максимальная скорость 1000
		if debugEnable {
			log.Println("times:", t1, t2)
		}

		unaval := unavailablePoints(gamers, points, t1)
		if debugEnable {
			log.Println("unaval:", unaval)
		}

		if step < eps && t2-t1 < eps { // проверка времени не имеет смысла дельта в 1000 раз меньше шага
			points = unaval
			break
		}

		t1, t2 = max(0, t1-step*2), min(2000, t2+step*2) // минимальная скорость 1
		if debugEnable {
			log.Println("new times:", t1, t2)
		}

		step = step * 2 / scale
		if debugEnable {
			log.Println("new step:", step)
		}

		points = points[:0]
		for _, p := range unaval {
			x1, y1 := max(-dd, p.x-step), max(0.0, p.y-step)
			x2, y2 := min(dd, p.x+step), min(dd, p.y+step)
			points = appendFronRect(points, dd, x1, y1, x2, y2, step)
			if debugEnable {
				log.Println("rect:", x1, y1, x2, y2, points)
			}
		}
	}

	return (t1 + t2) / 2, points[0].x, points[0].y
}

func run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	d, n, err := scanTwoInt(sc)
	if err != nil {
		return err
	}

	gamers := make([]gamer, n)
	for i := 0; i < n; i++ {
		x, y, v, err := scanThreeInt(sc)
		if err != nil {
			return err
		}
		gamers[i] = gamer{float64(x), float64(y), float64(v)}
	}

	t, x, y := solve(d, gamers)

	fmt.Fprintf(bw, "%g\n", t)
	fmt.Fprintf(bw, "%g %g\n", x, y)

	return nil
}

// ----------------------------------------------------------------------------

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

type writeOpts struct {
	sep byte
	end byte
}

func writeInt[I Int](bw *bufio.Writer, v I, opts writeOpts) error {
	var buf [32]byte

	_, err := bw.Write(strconv.AppendInt(buf[:0], int64(v), 10))

	if err == nil && opts.end != 0 {
		bw.WriteByte(opts.end)
	}

	return err
}

func writeInts[I Int](bw *bufio.Writer, a []I, opts writeOpts) error {
	var err error

	if len(a) != 0 {
		var buf [32]byte

		if opts.sep == 0 {
			opts.sep = ' '
		}

		_, err = bw.Write(strconv.AppendInt(buf[:0], int64(a[0]), 10))

		for i := 1; err == nil && i < len(a); i++ {
			err = bw.WriteByte(opts.sep)
			if err == nil {
				_, err = bw.Write(strconv.AppendInt(buf[:0], int64(a[i]), 10))
			}
		}
	}

	if err == nil && opts.end != 0 {
		err = bw.WriteByte(opts.end)
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
