package service

import (
	"unicode"
	"unicode/utf8"
)

// TruncateByWords 用于按单词截断字符串, maxWords 为保留的单词数.
// 如果字符串中的单词数小于 maxWords, 则返回原字符串.
// 如果字符串中的单词数大于等于 maxWords, 则返回截断后的字符串, 并在末尾添加 "...".
// 该函数会忽略连续的空白字符, 例如 "  hello  world  " 中的空白字符会被忽略.
// @param s string 要截断的字符串
// @param maxWords int 保留的单词数
// @return string 截断后的字符串
func TruncateByWords(s string, maxWords int) string {
	processedWords := 0
	wordStarted := false
	for i := 0; i < len(s); {
		r, width := utf8.DecodeRuneInString(s[i:])
		if !isSeparator(r) {
			i += width
			wordStarted = true
			continue
		}

		if !wordStarted {
			i += width
			continue
		}

		wordStarted = false
		processedWords++
		if processedWords == maxWords {
			const ending = "..."
			if (i + len(ending)) >= len(s) {
				// Source string ending is shorter than "..."
				return s
			}

			return s[:i] + ending
		}

		i += width
	}

	// Source string contains less words count than maxWords.
	return s
}

// isSeparator 用于判断 r 是否为分隔符.
// 该函数会忽略字母和数字, 例如 "hello123" 中的字母和数字会被忽略.
// 该函数会忽略连续的空白字符, 例如 "  hello  world  " 中的空白字符会被忽略.
// @param r rune 要判断的字符
// @return bool 是否为分隔符
func isSeparator(r rune) bool {
	// ASCII alphanumerics and underscore are not separators
	if r <= 0x7F {
		switch {
		case '0' <= r && r <= '9':
			return false
		case 'a' <= r && r <= 'z':
			return false
		case 'A' <= r && r <= 'Z':
			return false
		case r == '_':
			return false
		}
		return true
	}
	// Letters and digits are not separators
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return false
	}
	// Otherwise, all we can do for now is treat spaces as separators.
	return unicode.IsSpace(r)
}
