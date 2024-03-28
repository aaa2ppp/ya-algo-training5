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

func solve(n int, a, d []int) []int {
	var (
		prev     = make([]int, n)
		next     = make([]int, n)
		toPass   = make(map[int]struct{}, n)
		toRemove = make([]int, 0, n)
		res      = make([]int, 0, n)
	)

	for i := 0; i < n; i++ {
		prev[i] = i - 1
		next[i] = i + 1
		toPass[i] = struct{}{}
	}

	for i := 0; i < n; i++ {
		toRemove = toRemove[:0]

		if debugEnable {
			log.Println("toPass:", toPass)
		}

		for k := range toPass {
			v := d[k]
			if prev[k] >= 0 {
				v -= a[prev[k]]
			}
			if next[k] < n {
				v -= a[next[k]]
			}
			if v < 0 {
				toRemove = append(toRemove, k)
			}			
		}

		if debugEnable {
			log.Println("toRemove:", toRemove)
		}

		if len(toRemove) == 0 {
			break
		}

		res = append(res, len(toRemove))

		for k := range toPass {
			delete(toPass, k)
		}

		for _, k := range toRemove {
			delete(toPass, k)

			if prev[k] >= 0 {
				next[prev[k]] = next[k]
				toPass[prev[k]] = struct{}{}
			}
			if next[k] < n {
				prev[next[k]] = prev[k]
				toPass[next[k]] = struct{}{}
			}
		}
	}

	return res[:n]
}

func task(sc *bufio.Scanner, bw *bufio.Writer) error {
	n, err := scanInt(sc)
	if err != nil {
		return err
	}

	a := make([]int, n)
	if err := scanInts(sc, a); err != nil {
		return err
	}

	b := make([]int, n)
	if err := scanInts(sc, b); err != nil {
		return err
	}

	res := solve(n, a, b)
	writeInts(bw, res, writeOpts{end: '\n'})

	return nil
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

	for i := 1; i <= t; i++ {
		if err := task(sc, bw); err != nil {
			return fmt.Errorf("task %d: %w", i, err)
		}
	}

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
