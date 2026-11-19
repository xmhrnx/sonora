package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteAAR(t *testing.T) {
	dir := t.TempDir()
	libs := filepath.Join(dir, "android")
	for _, p := range []string{"arm64-v8a/libonnxruntime.so", "arm64-v8a/libsherpa-onnx-c-api.so", "x86_64/libonnxruntime.so"} {
		os.MkdirAll(filepath.Dir(filepath.Join(libs, p)), 0o755)
		os.WriteFile(filepath.Join(libs, p), []byte(p), 0o644)
	}
	aar := libs + ".aar"
	if err := writeAAR(aar, libs); err != nil {
		t.Fatal(err)
	}
	r, err := zip.OpenReader(aar)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	var got []string
	for _, e := range r.File {
		got = append(got, e.Name)
	}
	want := []string{"jni/arm64-v8a/libonnxruntime.so", "jni/arm64-v8a/libsherpa-onnx-c-api.so", "jni/x86_64/libonnxruntime.so"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}
