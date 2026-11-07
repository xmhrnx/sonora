package textfile

import (
	"strconv"
	"strings"
)

// subtitles keeps the spoken lines of SRT and WebVTT files.
func subtitles(s string) string {
	var out []string
	skip := false // inside a WebVTT NOTE or STYLE block
	for _, l := range strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n") {
		t := strings.TrimSpace(l)
		switch {
		case t == "":
			skip = false
			out = append(out, "")
			continue
		case skip:
			continue
		case strings.HasPrefix(t, "WEBVTT"):
			continue
		case strings.HasPrefix(t, "NOTE") || t == "STYLE" || t == "REGION":
			skip = true
			continue
		case strings.Contains(t, "-->"):
			continue
		}
		if _, err := strconv.Atoi(t); err == nil {
			continue // cue number
		}
		out = append(out, htmlTag.ReplaceAllString(t, ""))
	}
	// Cues are separated by blank lines; join each cue's lines.
	return joinCues(out)
}

func joinCues(lines []string) string {
	var paras []string
	var cur []string
	for _, l := range append(lines, "") {
		if l == "" {
			if len(cur) > 0 {
				paras = append(paras, strings.Join(cur, " "))
				cur = nil
			}
			continue
		}
		cur = append(cur, l)
	}
	return strings.Join(paras, "\n")
}
