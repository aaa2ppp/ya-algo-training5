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
			args{strings.NewReader(`6 3 1 1 1`)},
			`YES
1.0000000000`,
			false,
			false,
		},
		{
			"2",
			args{strings.NewReader(`12 8 10 5 20
`)},
			`YES
0.3000000000`,
			false,
			false,
		},
		{
			"3",
			args{strings.NewReader(`5 0 0 1 2`)},
			`YES
2.0000000000`,
			false,
			false,
		},
		{
			"4",
			args{strings.NewReader(`10 7 -3 1 4`)},
			`YES
0.8571428571`,
			false,
			false,
		},
		{
			"5+",
			args{strings.NewReader(`10 8 0 2 0`)},
			`YES
0`,
			false,
			false,
		},
		{
			"6+",
			args{strings.NewReader(`10 3 0 7 0`)},
			`YES
0`,
			false,
			false,
		},
		{
			"7+",
			args{strings.NewReader(`10 4 0 4 0`)},
			`YES
0`,
			false,
			false,
		},
		{
			"8+",
			args{strings.NewReader(`10 5 0 5 0`)},
			`YES
0`,
			false,
			false,
		},
		{
			"9+",
			args{strings.NewReader(`10 4 0 3 0`)},
			`NO`,
			false,
			false,
		},
		{
			"12",
			args{strings.NewReader(`82 42 -354891707 42 -354891707`)},
			`YES
0`,
			false,
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

func Test_intersection(t *testing.T) {
	type args struct {
		f1 F
		f2 F
	}
	tests := []struct {
		name   string
		args   args
		wantT  float64
		wantL  float64
		wantOk bool
	}{
		{
			"1",
			args{
				f1: F{t0: 0, l0: 0, v: 1},
				f2: F{t0: 0, l0: 1, v: -1},
			},
			0.5,
			0.5,
			true,
		},
		{
			"2",
			args{
				f1: F{t0: 0, l0: 0, v: 1},
				f2: F{l0: 0.5},
			},
			0.5,
			0.5,
			true,
		},
		{
			"3",
			args{
				f1: F{t0: 0, l0: 3, v: -1},
				f2: F{t0: 0, l0: 1, v: 1},
			},
			1, 2,
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, gotL, gotOk := intersection(tt.args.f1, tt.args.f2)
			if gotOk != tt.wantOk {
				t.Errorf("intersection() gotOk = %v, wantOk %v", gotOk, tt.wantOk)
			}
			if gotT != tt.wantT || gotL != tt.wantL {
				t.Errorf("intersection() got = (%v, %v), want (%v, %v)", gotT, gotL, tt.wantT, tt.wantL)
			}
		})
	}
}
