/*
"""
Написать функцию, которая определяет, является ли переданная строка палиндромом
(читается слева-направо и справа-налево одинаково).

Примеры палиндромов:
- Казак
- А роза упала на лапу Азора
- Do geese see God?
- Madam, I’m Adam

Ограничение по памяти O(1).
"""
*/

package main

import "testing"

func Test_isPalindrome(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"1. Казак",
			args{"Казак"},
			true,
		},
		{
			"2. А роза упала на лапу Азора",
			args{"А роза упала на лапу Азора"},
			true,
		},
		{
			"3. Do geese see God?",
			args{"Do geese see God?"},
			true,
		},
		{
			"4. Madam, I’m Adam",
			args{"Madam, I’m Adam"},
			true,
		},
		{
			"5. empty",
			args{""},
			true,
		},
		{
			"6. no any letters",
			args{"! # ;,,,,"},
			true,
		},
		// {
		// 	"7. no any letters",
		// 	args{"12345 5 4 3,2   1,"},
		// 	true,
		// },
		{
			"4. Madam, you are not Adam",
			args{"Madam, you are not Adam"},
			false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPalindrome(tt.args.s); got != tt.want {
				t.Errorf("isPalindrome() = %v, want %v", got, tt.want)
			}
		})
	}
}
