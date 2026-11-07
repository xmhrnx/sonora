package shell

import (
	"runtime"

	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/ui/core"
)

// Shell lays out the frame around the current page.
type Shell struct {
	c       *core.Ctx
	bar     giode.Element // nil on mobile, where the OS draws no window frame
	sidebar *Sidebar
	pages   map[core.Page]func() giode.Element
}

// compactBelow is the width, in dp, up to which the app uses the
// one-column phone layout with a bottom tab bar. Inside it,
// giode.Compact() reports true, which is how pages and components
// adapt.
const compactBelow = 760

// New builds the shell; pages maps each page to its view function.
func New(c *core.Ctx, pages map[core.Page]func() giode.Element, onTheme func()) *Shell {
	s := &Shell{c: c, sidebar: NewSidebar(c, onTheme), pages: pages}
	if !core.Mobile {
		bar := giode.TitleBar(c.App, "Sonora", giode.Styles{Background: "theme.surface", Color: "theme.muted"}).
			Icon(giode.Image(logoImage(), giode.Styles{Width: 16, Height: 16, Fit: giode.FitContain, BorderRadius: 4}))
		if runtime.GOOS == "darwin" {
			bar.Variant(giode.TitleBarMac)
		}
		s.bar = bar
	}
	return s
}
