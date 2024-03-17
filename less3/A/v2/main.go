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

func run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, err := scanInt(sc)
	if err != nil {
		return err
	}

	dict := map[string]int{}

	for i := 0; i < n; i++ {
		k, err := scanInt(sc)
		if err != nil {
			return err
		}
		for j := 0; j < k; j++ {
			sc.Scan()

			// чтобы лишний раз не копировать из буфера, сначала проверяем наличие слова в словаре
			b := sc.Bytes()
			w := unsafeString(b)
			if v, ok := dict[w]; ok {
				// !!! это не работеет !!!, т.к. мапа imho переписывет ключ всегда. верно для 1.21.8
				// TODO: нужно проверить в исходнике мапы
				dict[w] = v + 1
			} else {
				// здесь нужно скопировать иначе при следующем чтении строка испортится
				dict[string(b)] = 1
			}
		}
	}

	var res []string
	for w, v := range dict {
		if v == n {
			res = append(res, w)
		}
	}

	sort.Strings(res)

	writeInt(bw, len(res))
	bw.WriteByte('\n')

	for _, w := range res {
		bw.WriteString(w)
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
