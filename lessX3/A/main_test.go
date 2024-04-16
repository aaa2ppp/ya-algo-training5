// Группировка анаграмм
// Дан массив строк, необходимо сгруппировать анаграммы.
// Слово X является анаграммой слова Y если одно может быть получено из другого перестановкой букв.
// В итоговом массиве каждый массив анаграмм должен быть отсортирован в лексикографическом порядке.
// Все слова в исходном массиве состоят только из строчных латинских букв

// Sample Input
// ["eat", "tea", "tan", "ate", "nat", "bat"]

// Sample Output
// [
//   ["ate", "eat","tea"],
//   ["nat","tan"],
//   ["bat"]
// ]

package main

import (
	"reflect"
	"sort"
	"testing"
)

func Test_solve(t *testing.T) {
	type args struct {
		words []string
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			"1",
			args{[]string{"eat", "tea", "tan", "ate", "nat", "bat"}},
			[][]string{
				{"ate", "eat", "tea"},
				{"nat", "tan"},
				{"bat"},
			},
		},
		{
			"2",
			args{nil},
			nil,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			
			if got := solve(tt.args.words); !reflect.DeepEqual(normResult(got), normResult(tt.want)) {
				t.Errorf("solve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func normResult(res [][]string) [][]string {
	if len(res) == 0 {
		return nil
	}
	sort.Slice(res, func(i, j int) bool{
		return len(res[i]) > len(res[j]) || len(res[i]) == len(res[j]) && res[i][0] < res[j][0]
	})
	return res
}
