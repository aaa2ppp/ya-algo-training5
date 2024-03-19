package main

import (
	"bytes"
	"io"
	"log"
	"os"
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
			args{strings.NewReader(`5
0 0 1 2
1 0 0 2
2 0 2 2
4 0 3 2
4 0 5 2
9 -1 10 1
10 1 9 3
8 1 10 1
8 1 9 -1
8 1 9 3
`)},
			`3`,
			false,
			true,
		},
		{
			"2",
			args{strings.NewReader(`1
3 4 7 9
-1 3 3 8
`)},
			`0`,
			false,
			true,
		},
		{
			"3",
			args{strings.NewReader(`1
-4 5 2 -3
-12 4 -2 4
`)},
			`1`,
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

func Benchmark_Run(b *testing.B) {
	fileName, ok := os.LookupEnv("TESTFILE")
	if !ok {
		log.Fatal("env TESTFILE required")
	} 
	buf, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	out := io.Discard

	for i := 0; i < b.N; i++ {
		in := bytes.NewReader(buf)
		if err := run(in, out); err != nil {
			log.Fatal(err)
		}
	}
}
