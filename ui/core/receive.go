package core

import (
	"errors"
	"mime"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/internal/audio/decode"
	"github.com/mehran9675/sonora/internal/textfile"
	"github.com/mehran9675/sonora/ui/tr"
)

// Receive takes what another app shared with Sonora: audio or video
// opens on the Transcribe page, a document or text becomes the
// Studio's script. Files are read off the UI goroutine and closed.
func (c *Ctx) Receive(s giode.SharedContent) {
	for _, f := range s.Files {
		c.receive(f, func(dir string) (string, error) { return SavePicked(f, dir, "shared") })
	}
	if len(s.Files) == 0 && strings.TrimSpace(s.Text) != "" && c.Speak != nil {
		c.Speak(s.Text)
	}
}

// ReceivePaths takes files dropped on the window, as Receive does.
func (c *Ctx) ReceivePaths(paths []string) {
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			c.Error(err)
			continue
		}
		pf := giode.PickedFile{ReadCloser: f, Name: filepath.Base(p), MIME: mime.TypeByExtension(filepath.Ext(p)), Path: p}
		// A dropped file is already on disk: transcribe it where it is.
		c.receive(pf, func(string) (string, error) { return p, nil })
	}
}

// receive routes one file by its kind; keep stores audio where the
// Transcribe page can read it and returns the path.
func (c *Ctx) receive(f giode.PickedFile, keep func(dir string) (string, error)) {
	switch kindOf(f) {
	case "audio":
		var path string
		c.Background(func() error {
			defer f.Close()
			var err error
			path, err = keep(PickedDir())
			return err
		}, func(err error) {
			if err != nil {
				c.Error(err)
			} else if c.Transcribe != nil {
				c.Transcribe(path, strings.TrimSuffix(f.Name, filepath.Ext(f.Name)))
			}
		})
	case "text":
		var text string
		c.Background(func() error {
			defer f.Close()
			var err error
			text, err = textfile.Read(f.Name, f)
			return err
		}, func(err error) {
			if err != nil {
				c.Error(err)
			} else if c.Speak != nil {
				c.Speak(text)
				c.Success(tr.N("Loaded {n} words from {file}", len(strings.Fields(text)), "file", f.Name))
			}
		})
	default:
		f.Close()
		c.Error(errors.New(tr.T("Sonora can't open {file}: it reads audio, video and documents", "file", f.Name)))
	}
}

// kindOf says whether a file is "audio" (or video) to transcribe,
// "text" to speak, or neither (""), by its extension or MIME type.
func kindOf(f giode.PickedFile) string {
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(f.Name)), ".")
	switch {
	case slices.Contains(decode.Extensions, ext),
		strings.HasPrefix(f.MIME, "audio/"), strings.HasPrefix(f.MIME, "video/"):
		return "audio"
	case slices.Contains(textfile.Extensions, ext),
		strings.HasPrefix(f.MIME, "text/"), f.MIME == "application/pdf":
		return "text"
	}
	return ""
}
