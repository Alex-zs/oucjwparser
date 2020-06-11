package util

const (
	Digit = iota	// 数字
	Capital			// 大写字母
	Lowercase		// 小写字母
	Others			// 其他字符
)

// 返回 字符类型
func CharType(c rune) int {
	switch {
	case '0' <= c &&  c <= '9' :
		return Digit
	case 'A' <= c && c <= 'Z':
		return Capital
	case 'a' <= c && c <= 'z':
		return Lowercase
	default:
		return Others
	}
}



