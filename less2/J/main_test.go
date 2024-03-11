package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func Test_run(t *testing.T) {
	type args struct {
		in io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
		debug   bool
	}{
		{
			"1",
			args{strings.NewReader(`2 1
#
.
`)},
			`NO`,
			false,
			true,
		},
		{
			"2",
			args{strings.NewReader(`2 2
..
##
`)},
			`YES
..
ab
`,
			false,
			true,
		},
		{
			"3",
			args{strings.NewReader(`1 3
###
`)},
			`YES
abb
`,
			false,
			true,
		},
		{
			"4",
			args{strings.NewReader(`1 5
####.
`)},
			`YES
abbb.
`,
			false,
			true,
		},
		{
			"23",
			args{strings.NewReader(`2 4
.#..
###.
`)},
			`YES
.b..
aaa.
`,
			false,
			true,
		},
		{
			"1+",
			args{strings.NewReader(`5 7
.......
.##....
.##....
....##.
....##.
`)},
			`YES
.......
.bb....
.bb....
....aa.
....aa.
`,
			false,
			true,
		},
		{
			"2+",
			args{strings.NewReader(`5 7
.......
.##....
.####..
...##..
.......
`)},
			`YES
.......
.bb....
.bbaa..
...aa..
.......
`,
			false,
			true,
		},
		{
			"3+",
			args{strings.NewReader(`5 7
.......
.##....
.###...
..##...
.......
`)},
			`NO`,
			false,
			true,
		},
		{
			"4+",
			args{strings.NewReader(`5 7
.......
.###...
.##....
.###...
.......
`)},
			`NO`,
			false,
			true,
		},
		{
			"5+",
			args{strings.NewReader(`5 7
.......
.##....
.##....
....#..
.......
`)},
			`YES
.......
.bb....
.bb....
....a..
.......
`,
			false,
			true,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugEnable = tt.debug
			out := &bytes.Buffer{}
			if err := run(tt.args.in, out); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut := out.String(); strings.TrimRight(gotOut, "\r\n") != strings.TrimRight(tt.wantOut, "\r\n") &&
				strings.TrimRight(invert(gotOut), "\r\n") != strings.TrimRight(tt.wantOut, "\r\n") {
				t.Errorf("run() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func invert(s string) string {
	var sb strings.Builder
	sb.Grow(len(s))
	for _, c := range []byte(s) {
		switch c {
		case 'a':
			sb.WriteByte('b')
		case 'b':
			sb.WriteByte('a')
		default:
			sb.WriteByte(c)
		}
	}
	return sb.String()
}
