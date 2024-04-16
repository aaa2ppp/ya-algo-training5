/*

# Дана строка (возможно, пустая), содержащая буквы A-Z:
# AAAABBBCCXYZDDDDEEEFFFAAAAAABBBBBBBBBBBBBBBBBBBBBBBBBBBB
# Нужно написать функцию RLE, которая вернет строку вида:
# A4B3C2XYZD4E3F3A6B28
# Еще надо выдавать ошибку, если на вход приходит недопустимая строка.

# empty string
# char not in ('A...Z')
# text 'Wrong string'

*/

package rle

import (
	"errors"
	"strconv"
	"strings"
)

var ErrWrongString = errors.New("wrong string")

func Encode(s string) (string, error) {

	if len(s) == 0 {
		return "", ErrWrongString
	}

	isValidChar := func(c byte) bool {
		return 'A' <= c && c <= 'Z'
	}

	var sb strings.Builder

	encodeChar := func(c byte, count int) {
		// it does nothing if count == 0
		switch {
		case count > 1:
			sb.WriteByte(c)
			sb.WriteString(strconv.Itoa(count))
		case count == 1:
			sb.WriteByte(c)
		}
	}

	var (
		prev  byte
		count int
	)

	for _, c := range []byte(s) {

		if !isValidChar(c) {
			return "", ErrWrongString
		}

		if c != prev {
			encodeChar(prev, count)
			prev = c
			count = 0
		}

		count++
	}

	encodeChar(prev, count)

	return sb.String(), nil
}
