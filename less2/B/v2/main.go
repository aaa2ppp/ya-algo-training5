package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

type item struct {
	val int
	min int
}

// stack with miminum
type stack []item

func (s stack) empty() bool {
	return len(s) == 0
}

func (s stack) min() int {
	if len(s) == 0 {
		panic("stack empty")
	}
	return s[len(s)-1].min
}

func (s *stack) pop() int {
	old := *s
	if old.empty() {
		panic("stack empty")
	}
	v := old[len(old)-1].val
	*s = old[:len(old)-1]
	return v
}

func (s *stack) push(v int) {
	if s.empty() {
		*s = append(*s, item{v, v})
	} else {
		*s = append(*s, item{v, min(v, s.min())})
	}
}

// queue with minimum
type queue struct {
	in  stack
	out stack
}

func newQueueSize(n int) *queue {
	return &queue{
		in:  make(stack, 0, n),
		out: make(stack, 0, n),
	}
}

func (q *queue) len() int {
	return len(q.in) + len(q.out)
}

func (q *queue) empty() bool {
	return q.len() == 0
}

func (q *queue) min() int {
	if q.empty() {
		panic("queue is empty")
	}
	if len(q.out) == 0 {
		return q.in.min()
	}
	if len(q.in) == 0 {
		return q.out.min()
	}
	return min(q.in.min(), q.out.min())
}

func (q *queue) pop() int {
	if q.empty() {
		panic("queue is empty")
	}
	if len(q.out) == 0 {
		for len(q.in) > 0 {
			q.out.push(q.in.pop())
		}
	}
	return q.out.pop()
}

func (q *queue) push(v int) {
	q.in.push(v)
}

func run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, k, err := scanTwoInt(sc)
	if err != nil {
		return err
	}

	q := newQueueSize(k)
	maximum := 0

	for i := 0; i < n; i++ {
		price, err := scanInt(sc)
		if err != nil {
			return err
		}

		if q.len() > 0 {
			maximum = max(maximum, price-q.min())
		}

		q.push(price)
		if q.len() > k {
			q.pop()
		}
	}

	writeInt(bw, maximum)
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

type Int interface {
	~int | ~int64 | ~int32 | ~int16 | ~int8
}

func writeInt[I Int](bw *bufio.Writer, v I) error {
	var buf [32]byte
	_, err := bw.Write(strconv.AppendInt(buf[:0], int64(v), 10))
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
