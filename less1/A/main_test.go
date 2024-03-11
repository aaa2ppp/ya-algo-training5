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
	}{
		{
			"1",
			args{strings.NewReader(`0 7
12 5`)},
			`25`,
			false,
		},
		{
			"2+",
			args{strings.NewReader(`0 2
4 1`)},
			`8`,
			false,
		},
		{
			"3+",
			args{strings.NewReader(`0 2
3 1`)},
			`7`,
			false,
		},
		{
			"4+",
			args{strings.NewReader(`1 1
-1 1`)},
			`5`,
			false,
		},
		
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
