package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

func findMaximum(f []int, i, j int) int {
	maximum := 0
	for {
		maximum = max(maximum, f[i])
		if i == j {
			break
		}
		i++
		if i == len(f) {
			i = 0
		}
	}
	return maximum
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

	f := make([]int, n)
	if err := scanInts(sc, f); err != nil {
		return err
	}

	if debugEnable {
		log.Println("f:", f)
	}

	a, b, k, err := scanThreeInt(sc)
	if err != nil {
		return err
	}

	// XXX возможно ошибся +-1 -- опыт - критерий истины
	var res int
	if (b-a-1)/k >= len(f)-1 {
		if debugEnable {
			log.Println("full round")
		}
		res = findMaximum(f, 0, len(f)-1)
	} else {
		m := len(f) * k
		i := (a - 1) % m / k
		j := (b - 1) % m / k
		if debugEnable {
			log.Println("m:", m, "i:", i, "j:", j)
		}
		res = findMaximum(f, i, j)
		res = max(res, findMaximum(f, (len(f)-j)%len(f), (len(f)-i)%len(f)))
	}

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

func writeInts[I Int](bw *bufio.Writer, a []I, delm string) error {
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
