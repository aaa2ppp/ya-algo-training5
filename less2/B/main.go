package main

import (
	"bufio"
	"container/heap"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

type Item struct {
	price int // The priority of the item in the queue.
	index int // The index of the item in the heap.
}

type queue struct {
	items []Item
	first int
	size  int
}

func (q *queue) Len() int {
	return q.size
}

func (q *queue) Push(v Item) *Item {
	if q.size == len(q.items) {
		panic("queue is full: " + strconv.Itoa(q.size))
	}
	i := q.first + q.size
	if i >= len(q.items) {
		i -= len(q.items)
	}
	q.items[i] = v
	q.size++
	return &q.items[i] // xxx
}

// func (q *queue) Peek() Item {...}

func (q *queue) Pop() int {
	if q.size == 0 {
		panic("queue is empty")
	}
	index := q.items[q.first].index
	q.first++
	q.size--
	if q.first == len(q.items) {
		q.first = 0
	}
	return index
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

	a := queue{items: make([]Item, k+1)} // yes k+1
	q := make(PriorityQueue, 0, k+1)
	maximum := 0

	for i := 0; i < n; i++ {
		curPrice, err := scanInt(sc)
		if err != nil {
			return err
		}

		if q.Len() > k {
			heap.Remove(&q, a.Pop())
		}

		if len(q) > 0 {
			maximum = max(maximum, curPrice-q[0].price)
		}

		it := a.Push(Item{price: curPrice})
		heap.Push(&q, it)
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

// ----------------------------------------------------------------------------
// https://pkg.go.dev/container/heap@go1.22.1#example-package-PriorityQueue

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].price < pq[j].price
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
