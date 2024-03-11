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
			args{strings.NewReader(`120 10 8
start (image layout=embedded width=12 height=5)
(image layout=surrounded width=25 height=58)
and word is 
(image layout=floating dx=18 dy=-15 width=25 height=20)
here new 
(image layout=embedded width=20 height=22)
another
(image layout=embedded width=40 height=19)
longword

new paragraph
(image layout=surrounded width=5 height=30)
(image layout=floating width=20 height=35 dx=50 dy=-16)
`)},
			`48 0
60 0
74 -5
32 20
0 52
104 81
100 65
`,
			false,
			true,
		},
		{
			"2",
			args{strings.NewReader(`1000 2 3
`)},
			``,
			false,
			false,
		},
		{
			"3",
			args{strings.NewReader(`100 2 3
(image dx=10 dy=11 height=100 width=20 layout=floating)
`)},
			`10 11
`,
			false,
			false,
		},
		{
			"4",
			args{strings.NewReader(`100 2 3
(image dx=-10 dy=11 height=100 width=20 layout=floating)
(image dx=0 dy=11 height=100 width=20 layout=floating)
(image dx=0 dy=11 height=100 width=20 layout=floating)
`)},
			`0 11
20 22
40 33`,
			false,
			false,
		},
		{
			"13",
			args{strings.NewReader(`20 2 1
(image layout=surrounded width=1 height=8) 0123 (image layout=surrounded width=1 height=5) 123 (image layout=surrounded width=1 height=100) 1234 (image layout=surrounded width=2 height=6)
321 (image layout=surrounded width=1 height=3) (image layout=embedded width=4 height=1)
ab
(image layout=embedded width=1 height=3) (image layout=embedded width=3 height=1)
(image layout=embedded width=5 height=1)
ab
(image layout=embedded width=2 height=1)
(image layout=embedded width=6 height=1)
(image layout=surrounded width=1 height=20)
ab
(image layout=surrounded width=1 height=20)
abcde
ab bc cd x
(image layout=surrounded width=1 height=30)
yzah
abc bc a cde
(image layout=surrounded width=2 height=1)
abc
(image layout=surrounded width=4 height=20)
bcd (image layout=surrounded width=2 height=20)
abcd
(image layout=embedded width=4 height=1)
(image layout=embedded width=5 height=1)
(image layout=embedded width=10 height=1)
(image layout=embedded width=11 height=1)
`)},
			`0 0
5 0
9 0
14 0
19 0
1 2
10 2
16 2
1 5
10 5
1 7
7 7
12 7
14 9
18 11
3 13
18 13
10 29
0 33
10 39
0 101`,
			false,
			false,
		},
		{
			"14",
			args{strings.NewReader(`3 2 1
(image layout=surrounded width=2 height=1000) (image layout=surrounded width=2 height=1000) (image layout=surrounded width=2 height=1000) (image layout=surrounded width=2 height=1000) (image layout=surrounded width=2 height=1000) (image layout=surrounded width=2 height=1000) (image layout=surrounded width=2 height=1000) (image layout=surrounded width=2 height=1000) (image layout=surrounded width=2 height=1000) (image layout=surrounded width=2 height=1000) (image layout=surrounded width=2 height=1000) (image layout=surrounded width=2 height=1000) (image layout=surrounded width=2 height=1000) (image layout=surrounded width=2 height=1000) (image layout=surrounded width=2 height=1000) (image layout=surrounded width=2 height=1000) (image layout=surrounded width=2 height=1000) (image layout=surrounded width=2 height=1000) (image layout=surrounded width=2 height=1000) (image layout=surrounded width=2 height=1000) ab abc (image layout=surrounded width=2 height=1000)
`)},
			`0 0
0 1000
0 2000
0 3000
0 4000
0 5000
0 6000
0 7000
0 8000
0 9000
0 10000
0 11000
0 12000
0 13000
0 14000
0 15000
0 16000
0 17000
0 18000
0 19000
0 20004`,
			false,
			false,
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
