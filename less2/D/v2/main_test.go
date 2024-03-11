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
		debug bool
	}{
		{
			"1",
			args{strings.NewReader(`3
1 1
1 2
2 1`)},
			`8`,
			false,
			true,
		},
		{
			"2",
			args{strings.NewReader(`1
8 8
`)},
			`4`,
			false,
			true,
		},
		// {
		// 	"3",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	false,
		// 	true,
		// },
		// {
		// 	"4",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	false,
		// 	true,
		// },
		// TODO: Add test cases.
		{
			"10",
			args{strings.NewReader(`64
1 1
1 2
1 3
1 4
1 5
1 6
1 7
1 8
2 1
2 2
2 3
2 4
2 5
2 6
2 7
2 8
3 1
3 2
3 3
3 4
3 5
3 6
3 7
3 8
4 1
4 2
4 3
4 4
4 5
4 6
4 7
4 8
5 1
5 2
5 3
5 4
5 5
5 6
5 7
5 8
6 1
6 2
6 3
6 4
6 5
6 6
6 7
6 8
7 1
7 2
7 3
7 4
7 5
7 6
7 7
7 8
8 1
8 2
8 3
8 4
8 5
8 6
8 7
8 8
`)},
			`32`,
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
			if gotOut := out.String(); strings.TrimRight(gotOut, "\r\n") != strings.TrimRight(tt.wantOut, "\r\n") {
				t.Errorf("run() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
