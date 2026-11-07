package pyenv

import (
	"path/filepath"
	"strings"
)

// Extract unpacks a .zip or .tar.gz archive. dest maps each archive
// entry to its output path; "" skips the entry.
func Extract(data []byte, ext string, dest func(name string) string) error {
	if ext == ".zip" {
		return extractZip(data, dest)
	}
	return extractTarGz(data, dest)
}

// StripTo returns a dest function for Extract that drops the archive's
// top-level folder and places everything under root.
func StripTo(root string) func(string) string {
	return func(name string) string {
		_, rest, ok := strings.Cut(filepath.ToSlash(name), "/")
		if !ok || rest == "" {
			return ""
		}
		return safeJoin(root, rest)
	}
}
