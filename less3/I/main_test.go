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
			args{strings.NewReader(`"Juventus" - "Milan" 3:1
Inzaghi 45'
Del Piero 67'
Del Piero 90'
Shevchenko 34'
Total goals for "Juventus"
Total goals by Pagliuca
Mean goals per game by Inzaghi
"Juventus" - "Lazio" 0:0
Mean goals per game by Inzaghi
Mean goals per game by Shevchenko
Score opens by Inzaghi
`)},
			`3
0
1
0.5
1
0
`,
			false,
			true,
		},
		{
			"1.1",
			args{strings.NewReader(`"Juventus 2" - "Milan" 3:1
Inzaghi 45'
Del Piero 67'
Del Piero 90'
Shevchenko 34'
Total goals for "Juventus 2"
Total goals by Pagliuca
Mean goals per game by Inzaghi
"Juventus 2" - "Lazio" 0:0
Mean goals per game by Inzaghi
Mean goals per game by Shevchenko
Score opens by Inzaghi
Goals on last 1 minutes by Del Piero
Goals on first 50 minutes by Inzaghi
Mean goals per game for "Juventus 2"
Score opens by "Juventus 2"
Score opens by "Milan"
"Juventus 2" - "Milan" 0:1
Shevchenko 34'
Goals on minute 34 by Shevchenko
Total goals for "Milan"
`)},
			`3
0
1
0.5
1
0
1
1
1.5
0
1
2
2
`,
			false,
			true,
		},
		{
			"2",
			args{strings.NewReader(`Total goals by Arshavin
`)},
			`0`,
			false,
			true,
		},
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
