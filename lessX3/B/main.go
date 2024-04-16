/*
# Для заданной строки найти длину наибольшей подстроки без повторяющихся символов.

# abcabcbbddee -> 3 (abc)
# bbbbb -> 1 (b)
# pwwkew -> 3 (wke)

# abcab
*/

package main

func solve(s string) int {
	runes := []rune(s)
	maximum := 0
	freq := map[rune]int{}
	count := 0 // количество разных символов в окне число которых > 1

	for l, r := 0, 0; r < len(runes); r++ {

		if freq[runes[r]] == 1 {
			count++
		}

		freq[runes[r]]++

		if count == 0 {
			maximum = max(maximum, r-l+1)

		} else {
			freq[runes[l]]--

			if freq[runes[l]] == 1 {
				count--
			}

			l++
		}
	}

	return maximum
}
