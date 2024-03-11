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
			args{strings.NewReader(`10
11
15
`)},
			`4`,
			false,
			false,
		},
		{
			"2",
			args{strings.NewReader(`1
2
1`)},
			`-1`,
			false,
			false,
		},
		{
			"3",
			args{strings.NewReader(`1
1
1`)},
			`1`,
			false,
			false,
		},
		{
			"4",
			args{strings.NewReader(`25
200
10`)},
			`13`,
			false,
			false,
		},
		{
			"9",
			args{strings.NewReader(`250 500 187`)},
			`4`,
			false,
			false,
		},
		{
			"13",
			args{strings.NewReader(`250 500 218`)},
			`6`,
			false,
			false,
		},
		{
			"20",
			args{strings.NewReader(`250 500 230`)},
			`8`,
			false,
			false,
		},
		{
			"123",
			args{strings.NewReader(`7 10 8`)},
			`4`,
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
