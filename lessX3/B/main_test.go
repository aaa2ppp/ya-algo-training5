/*
# Для заданной строки найти длину наибольшей подстроки без повторяющихся символов.

# abcabcbbddee -> 3 (abc)
# bbbbb -> 1 (b)
# pwwkew -> 3 (wke)

# abcab -> 3 (abc)
*/

package main

import "testing"

func Test_solve(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"1",
			args{"abcabcbbddee"},
			3,
		},
		{
			"2",
			args{"bbbbb"},
			1,
		},
		{
			"3",
			args{"pwwkew"},
			3,
		},
		{
			"4",
			args{"abcab"},
			3,
		},
		{
			"5",
			args{""},
			0,
		},
		{
			"6",
			args{"Шла Маша по шоссе"},
			5,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := solve(tt.args.s); got != tt.want {
				t.Errorf("solve() = %v, want %v", got, tt.want)
			}
		})
	}
}
