package onnx

import (
	"archive/tar"
	"compress/bzip2"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mehran9675/sonora/internal/engines"
)

// extraPath is where a model's extra file from url is kept.
func extraPath(dir, url string) string { return filepath.Join(dir, path.Base(url)) }

// fetchExtra downloads a single file unless it is already there.
func fetchExtra(ctx context.Context, dir, url string, log engines.Logger) error {
	dst := extraPath(dir, url)
	if _, err := os.Stat(dst); err == nil {
		return nil
	}
	resp, err := get(ctx, url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	tmp := dst + ".partial"
	body := &progressReader{r: resp.Body, total: resp.ContentLength, label: "Downloading " + path.Base(url), log: log}
	if err := writeFile(tmp, body); err != nil {
		os.Remove(tmp)
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}
	return os.Rename(tmp, dst)
}

// get starts a download, failing on anything but 200 OK.
func get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("download %s: %s", path.Base(url), resp.Status)
	}
	return resp, nil
}

// fetch downloads and unpacks a model package into dir/<id>, unless it
// is already there, along with the extra files it needs. Unpacking goes
// to a temporary folder that is renamed when complete, so an
// interrupted download never looks installed.
func fetch(ctx context.Context, dir string, m model, log engines.Logger) (string, error) {
	for _, url := range m.Extra {
		if err := fetchExtra(ctx, dir, url, log); err != nil {
			return "", err
		}
	}
	dst := filepath.Join(dir, m.ID)
	if _, err := os.Stat(dst); err == nil {
		return dst, nil
	}
	if m.Repo != "" {
		return fetchFiles(ctx, dst, m, log)
	}
	tmp := dst + ".partial"
	os.RemoveAll(tmp)
	resp, err := get(ctx, m.url())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body := &progressReader{r: resp.Body, total: resp.ContentLength, label: "Downloading " + m.Label, log: log}
	// Download and unpack in one pass; the bzip2 stream is decoded as it
	// arrives.
	tr := tar.NewReader(bzip2.NewReader(body))
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			os.RemoveAll(tmp)
			if ctx.Err() != nil {
				return "", ctx.Err()
			}
			return "", fmt.Errorf("unpack %s: %w", m.ID, err)
		}
		// Archives hold <id>/...; drop that top folder.
		_, rest, ok := strings.Cut(filepath.ToSlash(h.Name), "/")
		if !ok || rest == "" || !m.keep(rest) {
			continue
		}
		out := filepath.Join(tmp, filepath.FromSlash(rest))
		if !strings.HasPrefix(out, tmp+string(os.PathSeparator)) {
			continue // outside the folder
		}
		switch h.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(out, 0o755)
		case tar.TypeReg:
			err = writeFile(out, tr)
		}
		if err != nil {
			os.RemoveAll(tmp)
			return "", err
		}
	}
	if err := os.Rename(tmp, dst); err != nil {
		os.RemoveAll(tmp)
		return "", err
	}
	log(m.Label + " is ready.")
	return dst, nil
}

func writeFile(path string, r io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, r); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// fetchFiles downloads a package published as separate files (see
// model.Repo) into dst, through a temporary folder renamed when every
// file is there.
func fetchFiles(ctx context.Context, dst string, m model, log engines.Logger) (string, error) {
	tmp := dst + ".partial"
	os.RemoveAll(tmp)
	for i, file := range m.Files {
		resp, err := get(ctx, m.fileURL(file))
		if err != nil {
			os.RemoveAll(tmp)
			return "", err
		}
		label := fmt.Sprintf("Downloading %s (%d of %d)", m.Label, i+1, len(m.Files))
		err = writeFile(filepath.Join(tmp, file), &progressReader{r: resp.Body, total: resp.ContentLength, label: label, log: log})
		resp.Body.Close()
		if err != nil {
			os.RemoveAll(tmp)
			if ctx.Err() != nil {
				return "", ctx.Err()
			}
			return "", err
		}
	}
	if err := os.Rename(tmp, dst); err != nil {
		os.RemoveAll(tmp)
		return "", err
	}
	log(m.Label + " is ready.")
	return dst, nil
}
