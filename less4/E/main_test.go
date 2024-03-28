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
			"1(1)",
			args{strings.NewReader(`1`)},
			`1/1`,
			false,
			true,
		},
		{
			"2(6)",
			args{strings.NewReader(`6`)},
			`3/1`,
			false,
			true,
		},
		{
			"3(2)",
			args{strings.NewReader(`2`)},
			`2/1`,
			false,
			true,
		},
		{
			"(8)",
			args{strings.NewReader(`8`)},
			`3/2`,
			false,
			true,
		},
		{
			"(7)",
			args{strings.NewReader(`7`)},
			`4/1`,
			false,
			true,
		},
		{
			"(12)",
			args{strings.NewReader(`12`)},
			`2/4`,
			false,
			true,
		},
		{
			"(13)",
			args{strings.NewReader(`13`)},
			`3/3`,
			false,
			true,
		},
		{
			"(3495349085345)",
			args{strings.NewReader(`3495349085345`)},
			`915317/1728677`,
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
