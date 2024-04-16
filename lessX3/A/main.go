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
	"sort"
	"unsafe"
)

func solve(words []string) [][]string {

	idxs := map[tKey]int{}
	var res [][]string

	for _, w := range words {
		key := genKey(w)

		if idx, ok := idxs[key]; ok {
			res[idx] = append(res[idx], w)
		} else {
			idxs[key] = len(res)
			res = append(res, []string{w})
		}
	}

	for _, part := range res {
		sort.Strings(part)
	}

	return res
}

// В качестве ключа используем отсортированную строку. Это хорошо работает на
// коротких словах. На словах больше 26 символов имеет смысл в качестве ключа
// использовать массив [26]int и считать буквы вместо сортировки.

type tKey string

func genKey(w string) tKey {
	buf := make([]byte, len(w))
	copy(buf, w)
	sort.Slice(buf, func(i, j int) bool {
		return buf[i] < buf[j]
	})
	return tKey(unsafe.String(unsafe.SliceData(buf), len(buf)))
}
