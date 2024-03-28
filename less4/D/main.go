package main

import (
	"bufio"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"unsafe"
)

const infinity = math.MaxInt

func countRows(pageWidth int, wordLens []int) int {
	count := 0
	row := pageWidth

	for _, word := range wordLens {
		if word > pageWidth {
			count = infinity
			break
		}
		row += word + 1
		if (row - 1) > pageWidth {
			count++
			row = word + 1
		}
	}

	if debugEnable {
		log.Println(pageWidth, wordLens, count)
	}

	return count
}

func solve(w int, a, b []int) int {

	i := sort.Search(w, func(i int) bool {
		if debugEnable {
			log.Println("i:", i)
		}
		return countRows(i, a) <= countRows(w-i, b)
	})

	v := max(countRows(i, a), countRows(w-i, b))
	if i > 0 {
		i--
		v = min(v, max(countRows(i, a), countRows(w-i, b)))
	}

	return v
}

func run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	w, n, m, err := scanThreeInt(sc)
	if err != nil {
		return err
	}

	a := make([]int, n)
	if err := scanInts(sc, a); err != nil {
		return err
	}

	b := make([]int, m)
	if err := scanInts(sc, b); err != nil {
		return err
	}

	res := solve(w, a, b)

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
