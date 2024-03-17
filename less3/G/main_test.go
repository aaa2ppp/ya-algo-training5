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
			"0",
			args{strings.NewReader(`0
`)},
			`4
0 0
0 1
1 0
1 1
`,
			false,
			true,
		},
		{
			"1",
			args{strings.NewReader(`2
0 1
1 0
`)},
			`2
0 0
1 1
`,
			false,
			true,
		},
		{
			"1.2",
			args{strings.NewReader(`2
0 0
0 1
`)},
			`2
1 0
1 1
`,
			false,
			true,
		},
		{
			"1.3",
			args{strings.NewReader(`2
0 0
1 0
`)},
			`2
0 1
1 1
`,
			false,
			true,
		},
		{
			"2",
			args{strings.NewReader(`3
0 2
2 0
2 2
`)},
			`1
0 0
`,
			false,
			true,
		},
		{
			"3",
			args{strings.NewReader(`4
-1 1
1 1
-1 -1
1 -1
`)},
			`0`,
			false,
			true,
		},
		// {
		// 	"4",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	false,
		// 	true,
		// },
		// TODO: Add test cases.
		{
			"6",
			args{strings.NewReader(`12
8 6
-9 6
-4 1
-5 3
6 4
7 -2
9 2
9 8
8 10
-7 -2
-5 -6
1 7
`)},
			``,
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
