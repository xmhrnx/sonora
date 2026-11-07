package onnx

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// JoinText joins transcript pieces with spaces, except between pieces
// of languages written without them (Chinese, Japanese…).
func JoinText(parts []string) string {
	var b strings.Builder
	for i, p := range parts {
		if i > 0 {
			prev, _ := utf8.DecodeLastRuneInString(parts[i-1])
			next, _ := utf8.DecodeRuneInString(p)
			if !(unspaced(prev) && unspaced(next)) {
				b.WriteByte(' ')
			}
		}
		b.WriteString(p)
	}
	return b.String()
}

func unspaced(r rune) bool {
	return unicode.In(r, unicode.Han, unicode.Hiragana, unicode.Katakana, unicode.Thai) ||
		strings.ContainsRune("。，！？、；：」』）", r)
}
