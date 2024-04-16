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

import "unicode"

func isPalindrome(s string) bool {

	// Oops!..
	// В GO строки в utf-8. Т.е. символ переменного размера. Перебирать с конца проблематично.
	// Поэтому мы копируем его в слайс рун (мое почтенье, Роберт).
	ss := []rune(s)
	
	for l, r := 0, len(ss)-1; l < r; {

		// Ммм, а дигиты проверять надо? Если да, то '1' == '١'? И то и то еденица. unicode.IsDigit
		// считает оба знака цифрами, strconv.Atoi на второй ругается... 
		if !unicode.IsLetter(ss[l]) {
			l++
			continue
		}
		if !unicode.IsLetter(ss[r]) {
			r--
			continue
		}

		if unicode.ToLower(ss[l]) != unicode.ToLower(ss[r]) {
			return false
		}
		l++
		r--
	}

	return true
}