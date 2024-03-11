package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"
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
			args{strings.NewReader(`2
2015
1 January
8 January
Thursday
`)},
			`Monday Thursday`,
			false,
			false,
		},
		{
			"2",
			args{strings.NewReader(`3
2013
1 January
8 January
15 January
Tuesday
`)},
			`Monday Tuesday`,
			false,
			false,
		},
		{
			"3",
			args{strings.NewReader(`3
2013
6 February
13 February
20 February
Tuesday
`)},
			`Tuesday Wednesday`,
			false,
			false,
		},
		{
			"4+",
			args{strings.NewReader(`4
2000
6 February
13 February
20 February
1 April
` + time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Weekday().String())},
			`Monday Sunday`,
			false,
			false,
		},
		{
			"4++",
			args{strings.NewReader(`4
2000
6 February
13 February
29 February
1 April
` + time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Weekday().String())},
			`Monday Tuesday`,
			false,
			false,
		},
		{
			"5+",
			args{strings.NewReader(`4
2024
6 February
13 February
20 February
1 April
` + time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Weekday().String())},
			`Monday Tuesday`,
			false,
			false,
		},
		{
			"6+",
			args{strings.NewReader(`4
2100
6 February
13 February
20 February
1 April
` + time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC).Weekday().String())},
			`Friday Saturday`,
			false,
			false,
		},
		{
			"10",
			args{strings.NewReader(`3
2013
6 August
13 August
20 August
Tuesday
` + time.Date(2013, 1, 1, 0, 0, 0, 0, time.UTC).Weekday().String())},
			`Monday Tuesday`,
			false,
			false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		debugEnable = tt.debug
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
