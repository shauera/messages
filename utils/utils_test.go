package utils

import (
	"testing"
)

// https://en.wikipedia.org/wiki/Palindrome
// A palindrome is a word, number, phrase, or other sequence of characters
// which reads the same backward as forward, such as madam or racecar or the
// number 10801. Sentence-length palindromes may be written when allowances
// are made for adjustments to capital letters, punctuation, and word dividers,
// such as "A man, a plan, a canal, Panama!", "Was it a car or a cat I saw?"
// or "No 'x' in Nixon".

// This function is limited to English characters only.
// Also and will ignore any character that is not in {a-z,A-Z,0-9}
func TestIsPlaindrome(t *testing.T) {
	testCases := []struct {
		inputStr string
		expected bool
	} {
		{"", true}, //debatebale
		{" ", true}, //debatebale
		{"      ", true}, //debatebale
		{"1", true},
		{"12", false},
		{"121", true},
		{"1221", true},
		{"A man, a plan, a canal, Panama!", true},
		{"Was it a car or a cat I saw?", true},
		{"No 'x' in Nixon", true},
		{"madam ", true},
		{" madam", true},
		{" 296", false},
		{" 13331", true},
		{"13331 ", true},
		{"123.321", true},
		{"123!!321", true},
		{"What is 34.5 this 1 SIht 543 ", false},
		{"What is 34.5 @#=this 1 SIht 543 Si!! t A H w", true},
	}

	for _, testCase := range testCases {
		if IsPalindrome(testCase.inputStr) != testCase.expected {
			var toBeOrNotToBe string //;-)
			if testCase.expected == false {
				toBeOrNotToBe = " not"
			}
			t.Errorf("expected '%s'%s to be a palindrome", testCase.inputStr, toBeOrNotToBe)
		}
	}
}