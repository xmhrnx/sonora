package core

import "strings"

// FileName suggests a file name for an export titled title.
func FileName(title, ext string) string { return clipName(title, ext) }

// clipName suggests a file name for an exported clip.
func clipName(title, ext string) string {
	clean := strings.Map(func(r rune) rune {
		if strings.ContainsRune(`<>:"/\|?*`, r) {
			return '-'
		}
		return r
	}, strings.TrimSpace(title))
	if clean == "" {
		clean = "speech"
	}
	if ext == "" {
		ext = ".wav"
	}
	if len([]rune(clean)) > 60 {
		clean = string([]rune(clean)[:60])
	}
	return "Sonora - " + clean + ext
}
