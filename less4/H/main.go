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

type party struct {
	id int
	v int
	p int
}

func checkParty(pp []*party, idx int, cost int) int {
	votes := pp[idx].v + cost

	for i := range pp {
		if i == idx {
			continue
		}

		v := pp[i].v - votes + 1
		if debugEnable {
			log.Println("  :", pp[i].id, v)
		}
		if v <= 0 {
			break
		}

		cost -= v
		if cost < 0 {
			break
		}
	}

	return cost
}

func solve(parties []party) (int, int){

	pp := make([]*party, len(parties))
	for i := range parties {
		pp[i] = &parties[i]
	}

	sort.Slice(pp, func(i, j int) bool {
		return pp[i].v > pp[j].v
	})

	totalVotes := 0
	for i := range parties {
		totalVotes += parties[i].v
	}

	minCost := infinity
	minBribe := infinity
	winIdx := -1

	for i, p := range pp {
		if p.p == -1 || p.p >= minBribe {
			continue
		}

		votesCost := sort.Search(totalVotes, func(c int) bool {
			if debugEnable {
				log.Println("check:", pp[i].id, c)
			}
			v := checkParty(pp, i, c)
			if debugEnable {
				log.Println("  v:", v)
			}
			return v >= 0
		})

		if debugEnable {
			log.Println("found:", pp[i].id, votesCost, pp[i].p+votesCost)
		}

		cost := pp[i].p + votesCost
		if cost < minCost {
			minCost = cost
			winIdx = i
		}

		minBribe = min(minBribe, pp[i].p) 
	}

	votesCost := minCost - pp[winIdx].p
	if debugEnable {
		log.Printf("winID:%d minCost:%d votesCost:%d", pp[winIdx].id, minCost, votesCost)
	}

	pp[winIdx].v += votesCost
	for i, p := range pp {
		if i == winIdx {
			continue
		}
		v := p.v - pp[winIdx].v + 1
		if v <= 0 {
			break
		}
		p.v -= v
		votesCost -= v		
	}

	for i, p := range pp {
		if i == winIdx {
			continue
		}
		v := min(p.v, votesCost)
		p.v -= v
		votesCost -= v
		if votesCost == 0 {
			break
		}
	}

	return minCost, pp[winIdx].id
}

func scanParties(sc *bufio.Scanner, parties []party) error {
	for i := range parties {
		v, p, err := scanTwoInt(sc)
		if err != nil {
			return err
		}
		parties[i] = party{id:i+1, v:v, p:p}
	}
	return nil
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

	parties := make([]party, n)
	if err := scanParties(sc, parties); err != nil {
		return err
	}

	minCost, winID := solve(parties)

	writeInts(bw, []int{minCost, winID}, writeOpts{sep:'\n', end:'\n'})
	for i := range parties {
		writeInt(bw, parties[i].v, writeOpts{end: ' '})
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
