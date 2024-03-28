package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"unsafe"
)

func checkRectangle(sum [][]int, i0, j0, i1, j1 int) bool {
	s := sum[i1][j1] - sum[i0][j1] - sum[i1][j0] + sum[i0][j0]
	return s == (i1-i0)*(j1-j0)
}

func checkPlus(sum [][]int, i, j, k int) bool {
	n := len(sum)
	m := len(sum[0])
	i0 := i - k
	j0 := j - k
	i1 := i + k*2
	j1 := j + k*2

	if i0 < 0 || j0 < 0 || i1 >= n || j1 >= m {
		return false
	}

	res := checkRectangle(sum, i, j0, i+k, j1) &&
		checkRectangle(sum, i0, j, i1, j+k)

	return res
}

func solve(sum [][]int) int {
	n := len(sum)
	m := len(sum[0])
	res := 1

	i0 := 1
	j0 := 1
	i1 := n - 1
	j1 := m - 1
	for i := i0; i < i1; i++ {
		for j := j0; j < j1; j++ {
			nn := min(n, m) / 3
			k := sort.Search(nn+1, func(k int) bool {
				return !checkPlus(sum, i, j, k)
			})
			if debugEnable && k > res {
				log.Printf("%d, %d, %d", i, j, k-1)
			}
			if k > res {
				i1 = n - k
				j0 = k
				j1 = m - k
			}
			res = max(res, k)
		}
	}

	return res - 1
}

func calcSum(matrix [][]byte) [][]int {
	n := len(matrix) + 1
	m := len(matrix[0]) + 1
	buf := make([]int, n*m)
	sum := make([][]int, n)
	for i, j := 0, 0; i < len(sum); i, j = i+1, j+m {
		sum[i] = buf[j : j+m]
	}

	for i, row := range matrix {
		for j, c := range row {
			sum[i+1][j+1] = sum[i][j+1] + sum[i+1][j] - sum[i][j]
			if c == '#' {
				sum[i+1][j+1]++
			}
		}
	}

	return sum
}

func readMatrix(br *bufio.Reader, n, m int) ([][]byte, error) {
	buf, err := io.ReadAll(br)
	if err != nil {
		return nil, err
	}
	matrix := bytes.Split(bytes.TrimSpace(buf), []byte("\n"))
	for i := range matrix {
		matrix[i] = matrix[i][:m]
	}
	return matrix, nil
}

func debugMatrix(matrix [][]byte) {
	for i := range matrix {
		log.Printf("%s", matrix[i])
	}
}

func debugSum(matrix [][]int) {
	for i := range matrix {
		log.Printf("%3d ", matrix[i])
	}
}

func run(in io.Reader, out io.Writer) error {
	br := bufio.NewReader(in)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	var n, m int
	if _, err := fmt.Fscan(br, &n, &m); err != nil {
		return err
	}
	br.ReadBytes('\n') // skip first line

	matrix, err := readMatrix(br, n, m)
	if err != nil {
		return err
	}

	sum := calcSum(matrix)

	if debugEnable {
		debugMatrix(matrix)
		debugSum(sum)
	}

	res := solve(sum)
	writeInt(bw, res, writeOpts{end: '\n'})

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
