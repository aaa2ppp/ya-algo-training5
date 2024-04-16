/*
++++++++++++++++++++++++++++++++++++


"""
Даны два отсортированных списка с интервалами присутствия
пользователей в онлайне в течение дня. Начало интервала строго меньше конца.
Нужно вычислить интервалы, когда оба пользователя были в онлайне.
"""

intersection(
    [(8, 12), (17, 22)],
    [(5, 11), (14, 18), (20, 23)]
) # [(8, 11), (17, 18), (20, 22)]

intersection(
    [(9, 15), (18, 21)],
    [(10, 14), (21, 22)]
) # [(10, 14)]

def intersection(user1, user2):
    # your code here
*/

package main

import (
	"reflect"
	"testing"
)

func Test_intersection(t *testing.T) {
	type args struct {
		user1 [][]int
		user2 [][]int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			"1",
			args{
				[][]int{{8, 12}, {17, 22}},
    			[][]int{{5, 11}, {14, 18}, {20, 23}},
			},
			[][]int{{8,11}, {17, 18}, {20,22}},
		},
		{
			"2",
			args{
				[][]int{{9, 15}, {18, 21}},
    			[][]int{{10, 14}, {21, 22}},
			},
			[][]int{{10, 14}},
		},
		{
			"3",
			args{
				[][]int{{9, 15}},
    			[][]int{{21, 22}},
			},
			nil,
		},
		{
			"4",
			args{
				[][]int{{9, 24}},
				[][]int{{21, 24}},
			},
			[][]int{{21, 24}},
		},
		{
			"5",
			args{
				[][]int{{0, 24}},
				[][]int{{0, 24}},
			},
			[][]int{{0, 24}},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := intersection(tt.args.user1, tt.args.user2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("intersection() = %v, want %v", got, tt.want)
			}
		})
	}
}
