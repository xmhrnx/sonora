package textfile

import (
	"bytes"
	"path/filepath"
	"regexp"
	"strings"
)

// sniff guesses the format of an unnamed file.
func sniff(data []byte) string {
	switch {
	case bytes.HasPrefix(data, []byte("%PDF")):
		return "pdf"
	case bytes.HasPrefix(data, []byte(`{\rtf`)):
		return "rtf"
	case bytes.HasPrefix(data, []byte("PK\x03\x04")):
		return zipKind(data)
	case bytes.HasPrefix(data, []byte("WEBVTT")):
		return "vtt"
	}
	head := strings.ToLower(string(data[:min(len(data), 512)]))
	if strings.Contains(head, "<html") || strings.Contains(head, "<!doctype html") {
		return "html"
	}
	if srtTiming.MatchString(head) {
		return "srt"
	}
	return "txt"
}

var (
	blankLines = regexp.MustCompile(`\n{3,}`)
	lineSpaces = regexp.MustCompile(`[ \t\f\v\x{00a0}]+`)
	srtTiming  = regexp.MustCompile(`\d{1,2}:\d{2}:\d{2}[,.]\d{3}\s*-->`)
)

// tidy normalizes line endings and whitespace, keeping paragraphs.
func tidy(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimSpace(lineSpaces.ReplaceAllString(l, " "))
	}
	s = strings.Join(lines, "\n")
	return strings.TrimSpace(blankLines.ReplaceAllString(s, "\n\n"))
}

func orDefault(s, def string) string {
	if s == "" {
		return def
	}
	return filepath.Base(s)
}
