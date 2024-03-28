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

type item struct {
	v   int
	min int
	max int
}

func (it *item) updateMinMax(val int) {
	it.min = min(it.min, val)
	it.max = max(it.max, val)
}

func solve(w, h int, rows []item) int {

	// XXX это порно!

	if debugEnable {
		log.Println("-------------")
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].v < rows[j].v
	})

	rows2 := make([]item, len(rows)+1)
	rows2[len(rows2)-1] = item{
		v:   rows[len(rows)-1].v,
		min: w,
		max: 1,
	}

	for i := len(rows) - 1; i >= 1; i-- {
		rows2[i] = item{
			v:   rows[i-1].v,
			min: min(rows[i].min, rows2[i+1].min),
			max: max(rows[i].max, rows2[i+1].max),
		}
	}

	rows2[0] = item{
		v:   0,
		min: min(rows[0].min, rows2[1].min),
		max: max(rows[0].max, rows2[1].max),
	}

	if debugEnable {
		for i := range rows {
			log.Println(rows[i], rows2[i])
		}
	}

	topMin, topMax := w, 1

	verW := func(k int) int {
		curMin := min(topMin, rows2[k].min)
		curMax := max(topMax, rows2[k].max)
		return curMax - curMin + 1
	}

	horW := func(k, i int) int {
		return rows2[k].v - rows[i].v + 1
	}

	wid := min(w, h)

	for i := 0; i < len(rows); i++ {
		if debugEnable {
			log.Println("wid:", wid)
		}

		if i > 0 {
			topMin = min(topMin, rows[i-1].min)
			topMax = max(topMax, rows[i-1].max)
		}

		if debugEnable {
			log.Println("i:", i, "tmin,tmax:", topMin, topMax)
		}

		vw := topMax - topMin + 1
		hw := rows[len(rows)-1].v - rows[i].v + 1

		if debugEnable {
			log.Println("+i:", i, "vw,hw:", vw, hw)
		}

		wid = min(wid, max(vw, hw))
		if vw > hw {
			continue
		}

		k := sort.Search(len(rows)-i, func(k int) bool {
			k += i
			vw := verW(k)
			hw := horW(k, i)

			if debugEnable {
				log.Println("i,k:", i, k, "vw,hw:", vw, hw)
			}

			return vw <= hw
		})
		k += i

		{
			vw := verW(k)
			hw := horW(k, i)

			if debugEnable {
				log.Println("=i,k:", i, k, "vw,hw:", vw, hw)
			}

			wid = min(wid, max(verW(k), horW(k, i)))
		}

		if k > i {
			k--
			vw := verW(k)
			hw := horW(k, i)

			if debugEnable {
				log.Println("=i,k:", i, k, "vw,hw:", vw, hw)
			}

			wid = min(wid, max(verW(k), horW(k, i)))
		}
	}

	if debugEnable {
		log.Println("wid:", wid)
	}

	return wid
}

func run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	w, h, n, err := scanThreeInt(sc)
	if err != nil {
		return err
	}

	rows := make([]item, 0, n)
	rowMap := make(map[int]*item, n)

	cols := make([]item, 0, n)
	colMap := make(map[int]*item, n)

	for i := 0; i < n; i++ {
		x, y, err := scanTwoInt(sc)
		if err != nil {
			return err
		}

		if it := rowMap[y]; it == nil {
			n := len(rows)
			rows = append(rows, item{v: y, min: x, max: x})
			rowMap[y] = &rows[n]
		} else {
			it.updateMinMax(x)
		}

		if it := colMap[x]; it == nil {
			n := len(cols)
			cols = append(cols, item{v: x, min: y, max: y})
			colMap[x] = &cols[n]
		} else {
			it.updateMinMax(y)
		}
	}

	if len(rows) > len(cols) {
		rows, cols = cols, rows
		w, h = h, w
	}

	wid := solve(w, h, rows)

	// FOR TEST ONLY: если поменять местами сроки с колонками, решение не должно меняться
	// wid2 := solve(h, w, cols)
	// if wid != wid2 {
	// 	log.Fatalln("wid != wid2", wid, wid2)
	// }

	writeInt(bw, wid, writeOpts{end: '\n'})

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
