package service

const (
	MaxSummaryLength = 100
)

// TruncateByWords 按单词截断字符串
func TruncateByWords(s string, maxWords int) string {
	runes := []rune(s)
	if len(runes) <= maxWords {
		return s
	}
	return string(runes[:maxWords]) + "..."
}

// isSeparator 判断是否为分隔符
// func isSeparator(r rune) bool {
// 	// Letters and digits are not separators
// 	if unicode.IsLetter(r) || unicode.IsDigit(r) {
// 		return false
// 	}
// 	// 空白字符视为分隔符
// 	if unicode.IsSpace(r) {
// 		return true
// 	}
// 	// 标点符号等特殊字符也视为分隔符
// 	if unicode.IsPunct(r) || unicode.IsSymbol(r) {
// 		return true
// 	}
// 	// 其他字符根据需要自行决定
// 	return true
// }
