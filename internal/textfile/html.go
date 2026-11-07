package textfile

import (
	"bytes"
	"strings"

	nethtml "golang.org/x/net/html"
)

// block are the HTML elements that start a new line.
var block = map[string]bool{
	"p": true, "div": true, "br": true, "li": true, "h1": true, "h2": true, "h3": true, "h4": true, "h5": true,
	"h6": true, "tr": true, "section": true, "article": true, "blockquote": true, "pre": true, "dd": true, "dt": true,
	"header": true, "footer": true, "figcaption": true, "td": true, "th": true, "hr": true, "title": true,
}

// html returns an HTML (or XHTML) document's visible text.
func html(data []byte) string {
	z := nethtml.NewTokenizer(bytes.NewReader(data))
	var b strings.Builder
	skip := 0
	for {
		switch z.Next() {
		case nethtml.ErrorToken:
			return b.String()
		case nethtml.StartTagToken, nethtml.SelfClosingTagToken:
			name, _ := z.TagName()
			switch string(name) {
			case "script", "style", "head", "noscript", "svg":
				if t := z.Token(); t.Type == nethtml.StartTagToken {
					skip++
				}
			}
			if block[string(name)] {
				b.WriteString("\n")
			}
		case nethtml.EndTagToken:
			name, _ := z.TagName()
			switch string(name) {
			case "script", "style", "head", "noscript", "svg":
				if skip > 0 {
					skip--
				}
			}
			if block[string(name)] {
				b.WriteString("\n")
			}
		case nethtml.TextToken:
			if skip == 0 {
				b.WriteString(strings.ReplaceAll(string(z.Text()), "\n", " "))
			}
		}
	}
}
