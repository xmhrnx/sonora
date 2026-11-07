package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// addLibs copies the APK at in to out without its signature, adding
// every <abi>/*.so under libs for the ABIs the APK was built for. The
// libraries are stored uncompressed, as Android loads them in place.
func addLibs(in, out, libs string) (int, error) {
	r, err := zip.OpenReader(in)
	if err != nil {
		return 0, err
	}
	defer r.Close()
	f, err := os.Create(out)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	w := zip.NewWriter(f)
	abis := map[string]bool{}
	have := map[string]bool{}
	for _, e := range r.File {
		if strings.HasPrefix(e.Name, "META-INF/") {
			continue // the old signature
		}
		if parts := strings.Split(e.Name, "/"); len(parts) == 3 && parts[0] == "lib" {
			abis[parts[1]] = true
		}
		have[e.Name] = true
		if err := w.Copy(e); err != nil {
			return 0, err
		}
	}
	if len(abis) == 0 {
		return 0, fmt.Errorf("%s has no native libraries; is it an APK from giode build android?", in)
	}
	var names []string
	for abi := range abis {
		found, _ := filepath.Glob(filepath.Join(libs, abi, "*.so"))
		if len(found) == 0 {
			return 0, fmt.Errorf("no libraries for %s in %s; run go run ./tools/fetchsherpa android", abi, libs)
		}
		names = append(names, found...)
	}
	sort.Strings(names)
	added := 0
	for _, src := range names {
		name := "lib/" + filepath.Base(filepath.Dir(src)) + "/" + filepath.Base(src)
		if have[name] {
			continue
		}
		if err := addStored(w, name, src); err != nil {
			return 0, err
		}
		added++
	}
	if err := w.Close(); err != nil {
		return 0, err
	}
	return added, f.Close()
}

func addStored(w *zip.Writer, name, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	dst, err := w.CreateHeader(&zip.FileHeader{Name: name, Method: zip.Store})
	if err != nil {
		return err
	}
	_, err = io.Copy(dst, in)
	return err
}

// buildTools returns the newest build-tools folder of the Android SDK.
func buildTools() (string, error) {
	sdk := os.Getenv("ANDROID_HOME")
	if sdk == "" {
		sdk = os.Getenv("ANDROID_SDK_ROOT")
	}
	if sdk == "" {
		return "", fmt.Errorf("set ANDROID_HOME to the Android SDK")
	}
	dirs, _ := filepath.Glob(filepath.Join(sdk, "build-tools", "*"))
	if len(dirs) == 0 {
		return "", fmt.Errorf("no build-tools in %s; install them with sdkmanager", sdk)
	}
	sort.Strings(dirs)
	return dirs[len(dirs)-1], nil
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", filepath.Base(name), err)
	}
	return nil
}
