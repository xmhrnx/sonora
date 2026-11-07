package textfile

import (
	"strconv"
	"strings"
)

// rtf strips RTF control words and groups that hold no text.
func rtf(s string) string {
	var b strings.Builder
	// Groups whose content is not text (font tables, pictures…).
	ignore := map[string]bool{"fonttbl": true, "colortbl": true, "stylesheet": true, "info": true, "pict": true,
		"header": true, "footer": true, "listtable": true, "listoverridetable": true, "rsidtbl": true, "generator": true,
		"xmlnstbl": true, "themedata": true, "colorschememapping": true, "datastore": true, "latentstyles": true}
	type group struct{ skip bool }
	stack := []group{{}}
	skipNext := 0 // characters to skip after \uN
	for i := 0; i < len(s); i++ {
		c := s[i]
		top := &stack[len(stack)-1]
		switch c {
		case '{':
			stack = append(stack, group{skip: top.skip})
			if strings.HasPrefix(s[i:], `{\*`) {
				stack[len(stack)-1].skip = true
			}
		case '}':
			if len(stack) > 1 {
				stack = stack[:len(stack)-1]
			}
		case '\\':
			if i+1 >= len(s) {
				break
			}
			n := s[i+1]
			switch {
			case n == '\\' || n == '{' || n == '}':
				if !top.skip {
					b.WriteByte(n)
				}
				i++
			case n == '\'' && i+3 < len(s):
				if v, err := strconv.ParseUint(s[i+2:i+4], 16, 8); err == nil && !top.skip && skipNext == 0 {
					b.WriteRune(rune(v)) // Windows-1252, close enough to Latin-1
				}
				if skipNext > 0 {
					skipNext--
				}
				i += 3
			case n == '~':
				if !top.skip {
					b.WriteByte(' ')
				}
				i++
			case (n >= 'a' && n <= 'z') || (n >= 'A' && n <= 'Z'):
				j := i + 1
				for j < len(s) && ((s[j] >= 'a' && s[j] <= 'z') || (s[j] >= 'A' && s[j] <= 'Z')) {
					j++
				}
				word := s[i+1 : j]
				k := j
				if k < len(s) && (s[k] == '-' || (s[k] >= '0' && s[k] <= '9')) {
					k++
					for k < len(s) && s[k] >= '0' && s[k] <= '9' {
						k++
					}
				}
				arg := s[j:k]
				if k < len(s) && s[k] == ' ' {
					k++
				}
				i = k - 1
				switch {
				case ignore[word]:
					top.skip = true
				case top.skip:
				case word == "par" || word == "line" || word == "sect" || word == "page":
					b.WriteByte('\n')
				case word == "tab" || word == "cell":
					b.WriteByte(' ')
				case word == "u":
					if v, err := strconv.Atoi(arg); err == nil {
						if v < 0 {
							v += 65536
						}
						b.WriteRune(rune(v))
						skipNext = 1
					}
				}
			default:
				i++
			}
		case '\r', '\n':
		default:
			if skipNext > 0 {
				skipNext--
				continue
			}
			if !top.skip {
				b.WriteByte(c)
			}
		}
	}
	return b.String()
}
