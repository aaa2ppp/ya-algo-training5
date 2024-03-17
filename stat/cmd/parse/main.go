package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"html"
	"log"
	"os"
	"strings"
	"unicode"
	"unsafe"
)

func main() {
	for _, fileName := range os.Args[1:] {

		buf, err := os.ReadFile(fileName)
		if err != nil {
			log.Fatal(err)
		}

		rows, err := parseDoc(buf)
		if err != nil {
			log.Fatal(err)
		}

		w := csv.NewWriter(os.Stdout)
		w.WriteAll(rows)
	}
}

func parseDoc(buf []byte) ([][]string, error) {

	p := bytes.Index(buf, []byte(`<table class="table table_role_standings`))
	if p == -1 {
		return nil, fmt.Errorf("can't find table")
	}
	buf = buf[skipTag(buf, p):]

	p = bytes.Index(buf, []byte(`<tbody`))
	if p == -1 {
		return nil, fmt.Errorf("can't find <tbody")
	}
	buf = buf[skipTag(buf, p):]

	p = bytes.Index(buf, []byte(`</tbody>`))
	if p == -1 {
		return nil, fmt.Errorf("can't find </tbody>")
	}

	return parseTableBody(buf[:p])
}

func parseTableBody(buf []byte) ([][]string, error) {
	var rows [][]string

	for i, p := 1, bytes.Index(buf, []byte(`<tr`)); p != -1; i, p = i+1, bytes.Index(buf, []byte(`<tr`)) {
		buf = buf[skipTag(buf, p):]

		p := bytes.Index(buf, []byte(`</tr>`))
		if p == -1 {
			return nil, fmt.Errorf("row%d: can't find </tr>", i)
		}

		row, err := parseRow(buf[:p])
		if err != nil {
			return nil, fmt.Errorf("row%d: %v", i, err)
		}

		rows = append(rows, row)
		buf = buf[skipTag(buf, p):]
	}

	return rows, nil
}

func parseRow(buf []byte) ([]string, error) {
	var cells []string

	for i, p := 1, bytes.Index(buf, []byte(`<td`)); p != -1; i, p = i+1, bytes.Index(buf, []byte(`<td`)) {
		buf = buf[skipTag(buf, p):]

		p := bytes.Index(buf, []byte(`</td>`))
		if p == -1 {
			return nil, fmt.Errorf("cell%d: can't find </td>", i)
		}

		cell, err := parseCell(buf[:p])
		if err != nil {
			return nil, fmt.Errorf("cell%d: %v", i, err)
		}

		cells = append(cells, cell)
		buf = buf[skipTag(buf, p):]
	}

	return cells, nil
}

func parseCell(buf []byte) (string, error) {
	var sb strings.Builder
	var prev rune

	for i := 0; i < len(buf); {

		if buf[i] == '<' {

			// insert space between blocks
			if len(buf) - i >= 4 && unsafeString(buf[i:i+4]) == "<div" && prev != ' ' {
				sb.WriteByte(' ')
				prev = ' '
			}

			i = skipTag(buf, i)
			continue
		}

		text := buf[i:]

		if j := bytes.IndexByte(text, '<'); j == -1 {
			i += len(text)
		} else {
			text = text[:j]
			i += j
		}

		// log.Printf("text:\n%s", text)

		for _, c := range unsafeString(text) { // to range by rune

			if unicode.IsSpace(c) {
				c = ' '
			}

			if c == ' ' && prev == ' ' {
				continue
			}

			sb.WriteRune(c)
			prev = c
		}
	}

	return strings.TrimSpace(html.UnescapeString(sb.String())), nil
}

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func skipTag(buf []byte, i int) int {
	for i < len(buf) && buf[i] != '>' {
		i++
	}
	return i+1
}
