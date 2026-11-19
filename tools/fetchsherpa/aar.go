package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// writeAAR packs the libraries in dir (<abi>/lib*.so) into an Android
// library at path, as jni/<abi>/lib*.so. giode.toml lists it under
// android.libraries, so giode build android and aab pack the libraries
// with the app's own libgio.so, in the APK and in the Play bundle alike.
func writeAAR(path, dir string) error {
	found, err := filepath.Glob(filepath.Join(dir, "*", "*.so"))
	if err != nil {
		return err
	}
	sort.Strings(found)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := zip.NewWriter(f)
	for _, so := range found {
		name := "jni/" + filepath.Base(filepath.Dir(so)) + "/" + filepath.Base(so)
		if err := addFile(w, name, so); err != nil {
			return err
		}
	}
	if err := w.Close(); err != nil {
		return err
	}
	return f.Close()
}

func addFile(w *zip.Writer, name, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	dst, err := w.Create(name)
	if err != nil {
		return err
	}
	_, err = io.Copy(dst, in)
	return err
}
