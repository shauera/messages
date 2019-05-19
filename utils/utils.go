package utils

import (
	"regexp"
	"strings"
	"unsafe"
)

// https://en.wikipedia.org/wiki/Palindrome
// A palindrome is a word, number, phrase, or other sequence of characters
// which reads the same backward as forward, such as madam or racecar or the
// number 10801. Sentence-length palindromes may be written when allowances
// are made for adjustments to capital letters, punctuation, and word dividers,
// such as "A man, a plan, a canal, Panama!", "Was it a car or a cat I saw?"
// or "No 'x' in Nixon".

// IsPalindrome - returns true if the given string is a palindrome.
// Note that this function is considerring English characters only.
// Specifically it will ignore any character that is not in {a-z,A-Z,0-9}
// Case is also ignored
func IsPalindrome(str string) bool {
	str = prepare(str)
	for i := 0; i < len(str)/2; i++ {
		if str[i] != str[len(str)-i-1] {
			return false
		}
	}
	return true
}

func prepare(str string) string {
	reg, _ := regexp.Compile("[^A-Za-z0-9]+")
	clean := reg.ReplaceAllString(str, "")
	return strings.ToLower(strings.Trim(clean, ""))
}

//IsNilValue - returns true if the value passed in is nil
func IsNilValue(value interface{}) bool {
	return (*[2]uintptr)(unsafe.Pointer(&value))[1] == 0
}
