// Package textfile reads documents as plain text to be spoken: text
// and Markdown, subtitles (SRT, WebVTT), HTML, Word (DOCX),
// OpenDocument (ODT), EPUB, PDF and RTF.
package textfile

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

// Read returns the text of the document in r. name (which may be
// empty, as on Android) picks the format by extension; without one the
// format is sniffed from the content.
func Read(name string, r io.Reader) (string, error) {
	data, err := io.ReadAll(io.LimitReader(r, maxSize))
	if err != nil {
		return "", err
	}
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(name)), ".")
	if ext == "" {
		ext = sniff(data)
	}
	var text string
	switch ext {
	case "docx":
		text, err = docx(data)
	case "odt":
		text, err = odt(data)
	case "epub":
		text, err = epub(data)
	case "pdf":
		text, err = pdf(data)
	case "rtf":
		text = rtf(decodeText(data))
	case "html", "htm", "xhtml":
		text = html(data)
	case "srt", "vtt":
		text = subtitles(decodeText(data))
	case "md", "markdown":
		text = markdown(decodeText(data))
	default:
		if !looksLikeText(data) {
			return "", fmt.Errorf("%s doesn't look like a text document", orDefault(name, "the file"))
		}
		text = decodeText(data)
	}
	if err != nil {
		return "", err
	}
	text = tidy(text)
	if text == "" {
		return "", errors.New("no text found in " + orDefault(name, "the file"))
	}
	return text, nil
}
