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
	"testing"
)

func TestEncode(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
        {
            "AAAABBBCCXYZDDDDEEEFFFAAAAAABBBBBBBBBBBBBBBBBBBBBBBBBBBB",
            args{"AAAABBBCCXYZDDDDEEEFFFAAAAAABBBBBBBBBBBBBBBBBBBBBBBBBBBB"},
            "A4B3C2XYZD4E3F3A6B28",
            false,
        },
        {
            "A",
            args{"A"},
            "A",
            false,
        },
        {
            "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
            args{"ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
            "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
            false,
        },
        {
            "empty string",
            args{""},
            "",
            true,
        },
        {
            "AAAAaBBBB",
            args{"AAAAaBBBB"},
            "",
            true,
        },
        {
            "AAAA1BBBB",
            args{"AAAA1BBBB"},
            "",
            true,
        },
        {
            "AAAA\x00BBBB",
            args{"AAAA\x00BBBB"},
            "",
            true,
        },
        {
            "AAAA\nBBBB",
            args{"AAAA\nBBBB"},
            "",
            true,
        },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
