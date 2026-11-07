package pyenv

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

func extractZip(data []byte, dest func(string) string) error {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		err = writeEntry(dest(f.Name), f.Mode(), rc)
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// safeJoin joins an archive entry name under root, refusing entries
// that would escape it.
func safeJoin(root, name string) string {
	p := filepath.Join(root, filepath.FromSlash(name))
	if !strings.HasPrefix(p, filepath.Clean(root)+string(os.PathSeparator)) {
		return ""
	}
	return p
}
