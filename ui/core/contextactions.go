package core

import (
	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/internal/audio/decode"
	"github.com/mehran9675/sonora/internal/textfile"
	"github.com/mehran9675/sonora/ui/tr"
)

// ContextActions are Sonora's entries in the file manager's right-click
// menu, in the current language: transcribe audio and video, read a
// document aloud. giode.toml declares the same for the Windows
// installer and the macOS app.
func ContextActions() []giode.ContextAction {
	dotted := func(exts []string) []string {
		out := make([]string, len(exts))
		for i, e := range exts {
			out[i] = "." + e
		}
		return out
	}
	return []giode.ContextAction{
		{ID: "transcribe", Label: tr.T("Transcribe with Sonora"), Extensions: dotted(decode.Extensions)},
		{ID: "speak", Label: tr.T("Read aloud with Sonora"), Extensions: dotted(textfile.Extensions)},
	}
}

// ContextMenu adds Sonora's actions to the file manager's menu, or
// removes them, for the current user. Call it off the UI goroutine.
func (c *Ctx) ContextMenu(on bool) error {
	if on {
		return c.App.InstallContextActions()
	}
	return c.App.UninstallContextActions()
}
