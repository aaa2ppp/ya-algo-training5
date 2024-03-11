package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"unsafe"
)

type order struct {
	start int
	end   int
	cost  int
}

func (o order) duration() int {
	return o.end - o.start
}

type sumItem struct {
	t int // time
	v int // val
}

type sums []sumItem

func (s sums) get(start, end int) int {
	i1 := sort.Search(len(s), func(i int) bool {
		return s[i].t >= start
	}) - 1

	i2 := sort.Search(len(s), func(i int) bool {
		return s[i].t > end
	}) - 1

	return s[i2].v - s[i1].v
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

	orders := make([]order, 0, n)
	for i := 0; i < n; i++ {
		start, end, cost, err := scanThreeInt(sc)
		if err != nil {
			return err
		}

		orders = append(orders, order{
			start: start,
			end:   end,
			cost:  cost,
		})
	}

	sums1 := make(sums, 1, len(orders)+1)
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].start < orders[j].start
	})
	for i := range orders {
		sums1 = append(sums1, sumItem{
			t: orders[i].start,
			v: orders[i].cost + sums1[i].v,
		})
	}

	sums2 := make(sums, 1, len(orders)+1)
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].end < orders[j].end
	})
	for i := range orders {
		sums2 = append(sums2, sumItem{
			t: orders[i].end,
			v: orders[i].duration() + sums2[i].v,
		})
	}

	q, err := scanInt(sc)
	if err != nil {
		return err
	}

	for i := 0; i < q; i++ {
		start, end, typ, err := scanThreeInt(sc)
		if err != nil {
			return err
		}

		var res int
		switch typ {
		case 1:
			res = sums1.get(start, end)
		case 2:
			res = sums2.get(start, end)
		}

		writeInt(bw, res)
		bw.WriteByte(' ')
	}

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
