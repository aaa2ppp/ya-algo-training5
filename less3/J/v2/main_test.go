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
			args{strings.NewReader(`3 2`)},
			`3 3`,
			false,
			true,
		},
		// {
		// 	"2",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	false,
		// 	true,
		// },
		// {
		// 	"3",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	false,
		// 	true,
		// },
		{
			"19",
			args{strings.NewReader(`4 4`)},
			`7 8 8`,
			false,
			true,
		},
		{
			"50",
			args{strings.NewReader(`100 106`)},
			`3495 2889 3497 3414 3498 3395 2956 3422 3309 2337 3081 3268 2574 3407 3492 2683 2498 3496 3398 3465 3499 713 2495 3497 3448 3001 1180 3440 2936 2863 3194 2875 2868 3304 3489 3379 2700 3499 3465 2701 2276 3475 2291 2873 3500 3500 3478 2532 3501 3294 3138 2964 2584 2015 3491 3133 2135 3493 2136 2157 3091 3057 1948 2821 2824 3069 3498 3192 2155 3188 3502 2713 3238 2449 2044 3501 3194 3502 1855 3331 3503 3372 3503 3472 1827 3451 3284 3420 2309 3335 3504 2633 3008 3476 3411 3504 3293 1136 3505`,
			false,
			false,
		},
		{
			"100 200",
			args{strings.NewReader(`100 200`)},
			`5394 4883 1317 5430 6458 6447 2322 2897 5701 6419 5702 5261 4891 6150 6459 6166 6275 5232 6417 6266 6459 4894 5202 5793 3898 5790 6452 6030 4531 6154 6460 6363 4993 6337 4032 6458 5252 6114 6461 3910 6245 3408 4123 6235 6462 4261 6421 4279 6443 5391 5619 6461 4969 6244 4563 5990 6107 3788 3903 6094 4934 6264 6460 5794 6067 5498 6451 5425 5768 6462 5801 3466 6457 5530 6341 5944 6391 3367 6060 6036 5434 6172 6455 5922 4124 6463 4966 6435 6463 6339 6209 6464 6254 6464 5293 6465 5003 4242 4197`,
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

func Benchmark_solve_100_200(b *testing.B) {
	for i := 0; i < b.N; i++ {
		solve(100, 200) 
	}
}
