package core

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/ui/icons"
	"github.com/mehran9675/sonora/ui/tr"
)

// ShareClip hands the file at path to another app through the share
// sheet on phones, or shows it in the file manager on desktops.
func (c *Ctx) ShareClip(path string) {
	c.Background(func() error {
		return c.App.Share(context.Background(), giode.ShareFile{Path: path, MIME: mimeType(path)})
	}, nil)
}

// ShareTextFile writes text to a file called name and shares it, for
// subtitles that another app should get as a file.
func (c *Ctx) ShareTextFile(name, text string) {
	path := filepath.Join(CacheDir(), name)
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		c.Error(err)
		return
	}
	c.ShareClip(path)
}

// mimeType names the types Sonora shares; phones don't all know Opus
// or SubRip by extension.
func mimeType(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".wav":
		return "audio/wav"
	case ".opus", ".ogg":
		return "audio/ogg"
	case ".mp3":
		return "audio/mpeg"
	case ".flac":
		return "audio/flac"
	case ".srt":
		return "application/x-subrip"
	case ".vtt":
		return "text/vtt"
	case ".txt":
		return "text/plain"
	}
	return ""
}

// ShareIcon is the icon for ShareClip: a share arrow where it opens
// the share sheet, a folder where it opens the file manager.
func ShareIcon() string {
	if Mobile {
		return icons.Share
	}
	return icons.Folder
}

// ShareText hands text to another app through the share sheet on
// phones, or copies it to the clipboard on desktops.
func (c *Ctx) ShareText(text, title string) {
	c.Background(func() error {
		return c.App.Share(context.Background(), giode.ShareText{Text: text, Subject: title, Title: title})
	}, func(err error) {
		if err != nil {
			c.Error(err)
		} else if !Mobile {
			c.Success(tr.T("Copied to the clipboard"))
		}
	})
}
