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

type BitSet uint64

func (s *BitSet) Add(i int) {
	*s |= 1 << i
}

func (s BitSet) Includes(i int) bool {
	return s>>i&1 == 1
}

func (s BitSet) Intersection(other BitSet) BitSet {
	return s & other
}

func segmentsToSet(segments [][]int) BitSet {
	var s BitSet
	for _, seg := range segments {
		for i := seg[0]; i < seg[1]; i++ {
			s.Add(i)
		}
	}
	return s
}

func setToSegments(s BitSet) [][]int {
	var segments [][]int
	for i := 0; i < 24; {
		var l, r int
		for i < 24 && !s.Includes(i) {
			i++
		}
		if i == 24 {
			break
		}
		l = i
		for i < 24 && s.Includes(i) {
			i++
		}
		r = i
		segments = append(segments, []int{l, r})
	}
	return segments
}

func intersection(user1, user2 [][]int) [][]int {
	s1 := segmentsToSet(user1)
	s2 := segmentsToSet(user2)
	return setToSegments(s1.Intersection(s2))
}
