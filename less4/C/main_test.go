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
			args{strings.NewReader(`5 2
1 3 5 7 9
2 4
1 3
`)},
			`1
2
`,
			false,
			true,
		},
		{
			"1.1",
			args{strings.NewReader(`1 1
1
1 1`)},
			`1`,
			false,
			true,
		},
		{
			"1.1",
			args{strings.NewReader(`1 1
1
2 1`)},
			`-1`,
			false,
			true,
		},
		{
			"1.2",
			args{strings.NewReader(`1 1
1
3 1`)},
			`-1`,
			false,
			true,
		},
		{
			"1.3",
			args{strings.NewReader(`5 3
1 3 5 7 9
2 9
1 9
1 10
`)},
			`-1
5
-1
`,
			false,
			true,
		},
		// {
		// {
		// 	"2",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	false,
		// 	true,
		// },
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
