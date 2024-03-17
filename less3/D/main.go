package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

func run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, k,err := scanTwoInt(sc)
	if err != nil {
		return err
	}
	
	q := newQueue(k)
	set := make(map[int]struct{}, k)

	for i := 0; i < n; i++ {
		v, err := scanInt(sc)
		if err != nil {
			return err
		}

		if _, ok := set[v]; ok {
			bw.WriteString("YES\n")
			return nil
		}

		if q.Len() == k {
			delete(set, q.Pop())
		}

		q.Push(v)
		set[v] = struct{}{}
	}

	bw.WriteString("NO\n")

	return nil
}

// simple round queue
type queue struct {
	items []int
	first int
	size int
}

func newQueue(n int) *queue {
	return &queue{
		items: make([]int, n),
	}
}

func (q queue) Len() int {
	return q.size
}

func (q *queue) Push(v int) {
	if q.size == len(q.items) {
		panic("queue is full")
	}
	i := q.first + q.size
	if n:=len(q.items); i >= n {
		i -= n
	}
	q.items[i] = v
	q.size++
}

func (q *queue) Pop() int {
	if q.size == 0 {
		panic("queue is empty")
	}
	v := q.items[q.first]
	q.first++
	if q.first == len(q.items) {
		q.first = 0
	}
	q.size--
	return v
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
