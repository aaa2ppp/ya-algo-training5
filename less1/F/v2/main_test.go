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
			args{strings.NewReader(`3
5 7 2
`)},
			`x+`,
			false,
		},
		{
			"2",
			args{strings.NewReader(`2
4 -5
`)},
			`+`,
			false,
		},
		{
			"2+",
			args{strings.NewReader(`5
4 -5 4 -3 1
`)},
			`++++`,
			false,
		},
		{
			"6",
			args{strings.NewReader(`6
-76959846 -779700294 380306679 -340361999 58979764 -392237502`)},
			`++x++`,
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
