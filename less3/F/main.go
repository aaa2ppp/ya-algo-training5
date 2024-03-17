package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

type hashValue uint64

const initHashValue hashValue = 5381

func djb2Hash(text string) hashValue {
	hash := initHashValue
	for _, c := range []byte(text) {
		hash = djb2HashAdd(hash, c)
	}
	return hash
}

func djb2HashAdd(hash hashValue, c byte) hashValue {
	return (hash << 5) + hash + hashValue(c) // (hash * 33) + c
}

// simple hash table
type tDict [][]string

func (d tDict) add(w string) {
	h := djb2Hash(w)
	idx := h % hashValue(len(d))
	d[idx] = append(d[idx], w)
}

func (d tDict) findShort(w string) string {
	h := initHashValue

	for i, c := range []byte(w) {
		h = djb2HashAdd(h, c)
		idx := h % hashValue(len(d))

		if list := d[idx]; list != nil {
			short := w[:i+1]

			for _, it := range list {
				if it == short {
					// NOTE: here we need to return the value stored in the dictionary.
					//  The input value may not be safe!
					return it 
				}
			}
		}
	}

	return ""
}

func run(in io.Reader, out io.Writer) error {
	br := bufio.NewReader(in)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	buf, err := io.ReadAll(br)
	if err != nil {
		return err
	}

	p := bytes.IndexByte(buf, '\n')

	// scan first line

	dict := make(tDict, 1024)
	sc := bufio.NewScanner(bytes.NewReader(buf[:p]))
	sc.Split(bufio.ScanWords)

	for sc.Scan() {
		dict.add(sc.Text())
	}

	// scan second line

	sc = bufio.NewScanner(bytes.NewReader(buf[p+1:]))
	sc.Split(bufio.ScanWords)

	for i := 0; sc.Scan(); i++ {

		if i > 0 {
			bw.WriteByte(' ')
		}

		b := sc.Bytes()
		w := dict.findShort(unsafeString(b))

		if w != "" {
			bw.WriteString(w)
		} else {
			bw.Write(b)
		}
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

// ----------------------------------------------------------------------------
