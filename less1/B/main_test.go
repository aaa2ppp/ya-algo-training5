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
			args{strings.NewReader(`0:0
0:0
1
`)},
			`1`,
			false,
		},
		{
			"2",
			args{strings.NewReader(`0:2
0:3
1
`)},
			`5`,
			false,
		},
		{
			"3",
			args{strings.NewReader(`0:2
0:3
2
`)},
			`6`,
			false,
		},
		{
			"4+",
			args{strings.NewReader(`3:2
0:3
2
`)},
			`3`,
			false,
		},
		{
			"5+",
			args{strings.NewReader(`3:2
2:3
2
`)},
			`1`,
			false,
		},
		{
			"6+",
			args{strings.NewReader(`3:2
3:3
2
`)},
			`0`,
			false,
		},
		{
			"7+",
			args{strings.NewReader(`2:3
3:2
1
`)},
			`1`,
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
