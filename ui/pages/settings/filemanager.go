package settings

import (
	"runtime"

	"github.com/xmhrnx/giode"

	st "github.com/mehran9675/sonora/internal/settings"
	"github.com/mehran9675/sonora/ui/components/toggle"
	"github.com/mehran9675/sonora/ui/core"
	"github.com/mehran9675/sonora/ui/tr"
)

// fileManager is the switch that puts "Transcribe with Sonora" and
// "Read aloud with Sonora" in the file manager's right-click menu. macOS
// takes them from the app itself, so it has none.
type fileManager struct {
	c      *core.Ctx
	toggle *toggle.Toggle
}

func (f *fileManager) init(c *core.Ctx) {
	f.c = c
	f.toggle = toggle.New(tr.Msg("Show in the file manager's menu"), c.Settings.Get().ContextMenu, func(on bool) {
		c.Settings.Update(func(s *st.Settings) { s.ContextMenu = on })
		c.Background(func() error { return c.ContextMenu(on) }, nil)
	})
}

// element is the setting row, or nil where it doesn't apply.
func (f *fileManager) element(p *Page) giode.Element {
	if core.Mobile || runtime.GOOS == "darwin" {
		return nil
	}
	return p.setting(tr.T("File manager"),
		tr.T("Right-click audio or video to transcribe it, or a document to read it aloud, with Sonora."), f.toggle.Element())
}
