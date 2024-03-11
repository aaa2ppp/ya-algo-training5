package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strconv"
)

type rectangle struct {
	x1, y1, x2, y2 int
}

func (r rectangle) width() int {
	return abs(r.x1-r.x2) + 1
}

func (r rectangle) height() int {
	return abs(r.y1-r.y2) + 1
}

func (r rectangle) area() int {
	return r.height() * r.width()
}

func findCorner(matrix [][]byte) (int, int, int, int, bool) {
	n := len(matrix)
	m := len(matrix[0])

	var ok bool
	var x1, y1, x2, y2 = n, m, -1, -1

	for x := 0; x < n; x++ {
		for y := 0; y < m; y++ {
			c := matrix[x][y]
			if c == '#' {
				ok = true
				x1 = min(x1, x)
				y1 = min(y1, y)
				x2 = max(x2, x)
				y2 = max(y2, y)
			}
		}
	}

	if !ok {
		return 0, 0, 0, 0, false
	}

	switch {
	case matrix[x1][y1] == '#':
		return x1, y1, 1, 1, true

	case matrix[x1][y2] == '#':
		return x1, y2, 1, -1, true

	case matrix[x2][y1] == '#':
		return x2, y1, -1, 1, true

	case matrix[x2][y2] == '#':
		return x2, y2, -1, -1, true
	}

	panic("*** logical error ***")
}

func findRect(matrix [][]byte) (rectangle, bool) {
	n := len(matrix)
	m := len(matrix[0])

	x1, y1, dx, dy, ok := findCorner(matrix)
	if !ok {
		return rectangle{}, false
	}

	x2, y2 := x1, y1
	for f1, f2 := true, true; f1 || f2; {

		f1 = f1 && x2+dx < n && checkX(matrix, x2+dx, y1, y2)
		if f1 {
			x2 += dx
		}

		f2 = f2 && y2+dy < m && checkY(matrix, y2+dy, x1, x2)
		if f2 {
			y2 += dy
		}
	}

	return rectangle{x1, y1, x2, y2}, true
}

func checkX(matrix [][]byte, x, y1, y2 int) bool {
	m := len(matrix[0])

	if y1 > y2 {
		y1, y2 = y2, y1
	}

	for y := y1; y < m && y <= y2; y++ {
		if matrix[x][y] != '#' {
			return false
		}
	}

	return true
}

func checkY(matrix [][]byte, y, x1, x2 int) bool {
	n := len(matrix)

	if x1 > x2 {
		x1, x2 = x2, x1
	}

	for x := x1; x < n && x <= x2; x++ {
		if matrix[x][y] != '#' {
			return false
		}
	}

	return true
}

func fillRect(matrix [][]byte, r rectangle, c byte) {
	n := len(matrix)
	m := len(matrix[0])

	x1 := min(r.x1, r.x2)
	y1 := min(r.y1, r.y2)
	x2 := max(r.x1, r.x2)
	y2 := max(r.y1, r.y2)

	for x := x1; x < n && x <= x2; x++ {
		for y := y1; y < m && y <= y2; y++ {
			matrix[x][y] = c
		}
	}
}

func countPoints(matrix [][]byte) int {
	var n int

	for x := range matrix {
		for y := range matrix[x] {
			if matrix[x][y] == '#' {
				n++
			}
		}
	}

	return n
}

func solution(matrix [][]byte) bool {

	n := countPoints(matrix)
	if n < 2 {
		return false
	}

	r1, ok := findRect(matrix)
	if !ok {
		// при n > 0 такое невозможно
		panic("r1 not found")
	}

	fillRect(matrix, r1, 'b') // красим сначала в 'b', чтобы соотвествовать примерам

	if r1.area() == n {
		r2 := r1
		if r1.width() > 1 {
			r2.x2 = r2.x1
			fillRect(matrix, r2, 'a')
		} else {
			r2.y2 = r2.y1
			fillRect(matrix, r2, 'a')
		}
		return true
	}

	r2, ok := findRect(matrix)
	if !ok {
		// при n > 1 такое невозможно
		panic("r2 not found")
	}

	fillRect(matrix, r2, 'a')

	return r1.area()+r2.area() == n
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func run(in io.Reader, out io.Writer) error {

	r := bufio.NewReader(in)
	buf, err := r.ReadBytes('\n')
	if err != nil {
		return err
	}

	sc := bufio.NewScanner(bytes.NewReader(buf))
	sc.Split(bufio.ScanWords)
	n, m, err := scanTwoInt(sc)
	if err != nil {
		return err
	}

	buf, err = io.ReadAll(r)
	if err != nil {
		return err
	}

	matrix := bytes.Split(buf, []byte("\n"))
	matrix = matrix[:n]
	for i := range matrix {
		matrix[i] = matrix[i][:m]
	}

	w := bufio.NewWriter(out)
	defer w.Flush()

	if solution(matrix) {
		w.WriteString("YES\n")
		w.Write(buf)
	} else {
		w.WriteString("NO\n")
		// w.Write(buf)
	}

	return nil
}

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

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
