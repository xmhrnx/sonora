package core

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/ui/tr"
)

// SaveClip asks where to save the audio file at path (the system
// picker on Android, a native dialog on desktops) and copies it there
// as it is - WAV, Opus or whatever it was recorded in. A desktop
// without a dialog service (a Linux system without the XDG portal)
// gets the file in its Downloads folder instead.
func (c *Ctx) SaveClip(path, title string) {
	c.save(clipName(title, filepath.Ext(path)), func() (io.ReadCloser, error) { return os.Open(path) })
}

// SaveText saves text as a file, suggesting name, the same way.
func (c *Ctx) SaveText(name, text string) {
	c.save(name, func() (io.ReadCloser, error) { return io.NopCloser(strings.NewReader(text)), nil })
}

func (c *Ctx) save(name string, open func() (io.ReadCloser, error)) {
	if Mobile {
		c.saveToDownloads(name, open)
		return
	}
	c.saveAs(name, open)
}

// saveToDownloads puts the file in Downloads/Sonora on a phone, no
// questions asked; where that isn't possible (Android before 10) the
// system asks where, as on desktops.
func (c *Ctx) saveToDownloads(name string, open func() (io.ReadCloser, error)) {
	tmp := filepath.Join(CacheDir(), name)
	c.Background(func() error {
		defer os.Remove(tmp)
		if err := copyTo(tmp, open); err != nil {
			return err
		}
		_, err := c.App.SaveToDownloads(context.Background(), tmp, "Sonora")
		return err
	}, func(err error) {
		switch {
		case errors.Is(err, giode.ErrNotSupported):
			c.saveAs(name, open)
		case err != nil:
			c.Error(err)
		default:
			c.Success(tr.T("Saved to Downloads"))
		}
	})
}

// copyTo writes what open returns to a new file at path.
func copyTo(path string, open func() (io.ReadCloser, error)) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	in, err := open()
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

// saveAs asks where to save with the system's dialog; a desktop
// without one saves into Downloads.
func (c *Ctx) saveAs(name string, open func() (io.ReadCloser, error)) {
	var fallback string
	c.Background(func() error {
		w, err := c.App.SaveFile(context.Background(), name)
		if err != nil && !errors.Is(err, giode.ErrCanceled) && !Mobile {
			fallback = filepath.Join(downloadsDir(), name)
			w, err = os.Create(fallback)
		}
		if err != nil {
			return err
		}
		in, err := open()
		if err != nil {
			w.Close()
			return err
		}
		defer in.Close()
		if _, err := io.Copy(w, in); err != nil {
			w.Close()
			return err
		}
		return w.Close()
	}, func(err error) {
		switch {
		case err != nil:
			c.Error(err)
		case fallback != "":
			c.Success(tr.T("Saved to {path}", "path", fallback))
		default:
			c.Success(tr.T("Saved"))
		}
	})
}

// downloadsDir is the user's Downloads folder, or home.
func downloadsDir() string {
	home, _ := os.UserHomeDir()
	if st, err := os.Stat(filepath.Join(home, "Downloads")); err == nil && st.IsDir() {
		return filepath.Join(home, "Downloads")
	}
	return home
}
