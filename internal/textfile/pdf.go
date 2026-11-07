package textfile

import (
	"bytes"
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"

	pdflib "github.com/ledongthuc/pdf"
)

// pdf extracts the text layer page by page. Scanned PDFs have none.
func pdf(data []byte) (text string, err error) {
	defer func() {
		// The PDF parser panics on some malformed files.
		if recover() != nil {
			text, err = "", errors.New("this PDF can't be read")
		}
	}()
	r, err := pdflib.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", errors.New("this PDF can't be read: " + err.Error())
	}
	var b strings.Builder
	for i := 1; i <= r.NumPage(); i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}
		t, err := p.GetPlainText(nil)
		if err != nil {
			continue
		}
		b.WriteString(t)
		b.WriteString("\n\n")
	}
	if strings.TrimSpace(b.String()) == "" {
		return "", errors.New("this PDF has no text layer (it may be scanned images)")
	}
	return joinPDFLines(b.String()), nil
}

// joinPDFLines undoes the line breaks a PDF puts inside sentences: a
// line that starts in lower case or with punctuation continues the one
// before, and lone bullets are dropped.
func joinPDFLines(s string) string {
	var out []string
	for _, l := range strings.Split(s, "\n") {
		t := strings.TrimSpace(l)
		if t == "•" || t == "·" || t == "◦" || t == "▪" {
			continue
		}
		if n := len(out); n > 0 && t != "" && out[n-1] != "" {
			r, _ := utf8.DecodeRuneInString(t)
			switch {
			case strings.ContainsRune(".,;:!?)]’”", r):
				out[n-1] += t
				continue
			case unicode.IsLower(r):
				out[n-1] += " " + t
				continue
			}
		}
		out = append(out, t)
	}
	return strings.Join(out, "\n")
}
