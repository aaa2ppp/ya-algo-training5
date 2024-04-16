/*
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

func intersection(user1, user2 [][]int) [][]int {
	var res [][]int

	for len(user1) > 0 && len(user2) > 0 {
		s1 := user1[0]
		s2 := user2[0]

		if s1[1] > s2[0] && s2[1] > s1[0] {
			res = append(res, []int{max(s1[0], s2[0]), min(s1[1], s2[1])})
		}

		if s1[1] < s2[1] {
			user1 = user1[1:] // это не приводит к копированию
		} else {
			user2 = user2[1:]
		}
	}

	return res
}
