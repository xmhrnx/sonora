package textfile

import (
	"regexp"
	"strings"
)

var (
	mdImage    = regexp.MustCompile(`!\[[^\]]*\]\([^)]*\)`)
	mdLink     = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
	mdHeading  = regexp.MustCompile(`^\s{0,3}#{1,6}\s+`)
	mdBullet   = regexp.MustCompile(`^\s*(?:[-*+]|\d+[.)])\s+`)
	mdQuote    = regexp.MustCompile(`^\s*>\s?`)
	mdRule     = regexp.MustCompile(`^\s*(?:[-*_]\s*){3,}$`)
	mdEmphasis = regexp.MustCompile("(\\*\\*|__|\\*|_|~~|`)")
	htmlTag    = regexp.MustCompile(`<[^>]+>`)
)

// markdown keeps what a reader would say: no markup, links as their
// text, code blocks dropped.
func markdown(s string) string {
	var out []string
	fence := false
	for _, l := range strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n") {
		if strings.HasPrefix(strings.TrimSpace(l), "```") || strings.HasPrefix(strings.TrimSpace(l), "~~~") {
			fence = !fence
			continue
		}
		if fence || mdRule.MatchString(l) {
			continue
		}
		heading := mdHeading.MatchString(l)
		l = mdHeading.ReplaceAllString(l, "")
		l = mdBullet.ReplaceAllString(l, "")
		l = mdQuote.ReplaceAllString(l, "")
		l = mdImage.ReplaceAllString(l, "")
		l = mdLink.ReplaceAllString(l, "$1")
		l = htmlTag.ReplaceAllString(l, "")
		l = mdEmphasis.ReplaceAllString(l, "")
		if heading && l != "" {
			// A heading is its own sentence.
			if !strings.ContainsAny(l[len(l)-1:], ".!?:") {
				l += "."
			}
			out = append(out, "", l, "")
			continue
		}
		out = append(out, l)
	}
	return strings.Join(out, "\n")
}
