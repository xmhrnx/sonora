package textfile

import (
	"bytes"
	"unicode/utf16"
	"unicode/utf8"
)

// Extensions are the files worth offering in a file picker.
var Extensions = []string{"txt", "md", "markdown", "srt", "vtt", "html", "htm", "xhtml", "docx", "odt", "epub", "pdf", "rtf"}

// maxSize bounds how much of a file is read.
const maxSize = 64 << 20

// decodeText turns bytes into a string, honoring UTF-16 byte-order
// marks and falling back to Latin-1 for invalid UTF-8.
func decodeText(data []byte) string {
	switch {
	case bytes.HasPrefix(data, []byte{0xEF, 0xBB, 0xBF}):
		data = data[3:]
	case bytes.HasPrefix(data, []byte{0xFF, 0xFE}), bytes.HasPrefix(data, []byte{0xFE, 0xFF}):
		big := data[0] == 0xFE
		data = data[2:]
		u := make([]uint16, len(data)/2)
		for i := range u {
			if big {
				u[i] = uint16(data[2*i])<<8 | uint16(data[2*i+1])
			} else {
				u[i] = uint16(data[2*i+1])<<8 | uint16(data[2*i])
			}
		}
		return string(utf16.Decode(u))
	}
	if utf8.Valid(data) {
		return string(data)
	}
	r := make([]rune, len(data))
	for i, b := range data {
		r[i] = rune(b)
	}
	return string(r)
}

// looksLikeText rejects binary files: text has no NUL bytes outside
// UTF-16.
func looksLikeText(data []byte) bool {
	if bytes.HasPrefix(data, []byte{0xFF, 0xFE}) || bytes.HasPrefix(data, []byte{0xFE, 0xFF}) {
		return true
	}
	return !bytes.Contains(data[:min(len(data), 8192)], []byte{0})
}
