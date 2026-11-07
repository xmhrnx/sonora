// Package shell is the window frame: title bar, sidebar navigation and
// the page area.
package shell

import (
	"fmt"

	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/ui/components/navitem"
	"github.com/mehran9675/sonora/ui/core"
	"github.com/mehran9675/sonora/ui/icons"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

// Sidebar is the left navigation column.
type Sidebar struct {
	c     *core.Ctx
	items []*navitem.Item
	theme *giode.PressableElement
}

// NewSidebar builds the navigation for c.
func NewSidebar(c *core.Ctx, onTheme func()) *Sidebar {
	s := &Sidebar{c: c}
	nav := func(icon, label, tab string, p core.Page, badge func() string) {
		s.items = append(s.items, navitem.New(icon, label,
			func() bool { return c.Page() == p }, badge, func() { c.Go(p) }).TabLabel(tab))
	}
	// Each page has a short label for the phone's bottom bar, where six
	// share the width.
	nav(icons.Waves, tr.Msg("Studio"), tr.Msg("Studio (tab)"), core.PageStudio, nil)
	nav(icons.Captions, tr.Msg("Transcribe"), tr.Msg("Transcribe (tab)"), core.PageTranscribe, func() string {
		if n := len(c.Transcripts.Entries()); n > 0 {
			return fmt.Sprint(n)
		}
		return ""
	})
	nav(icons.Library, tr.Msg("Library"), tr.Msg("Library (tab)"), core.PageLibrary, func() string {
		if n := len(c.Library.Entries()); n > 0 {
			return fmt.Sprint(n)
		}
		return ""
	})
	nav(icons.UserVoice, tr.Msg("Profiles"), tr.Msg("Profiles (tab)"), core.PageProfiles, func() string {
		if n := len(c.Profiles.List()); n > 0 {
			return fmt.Sprint(n)
		}
		return ""
	})
	nav(icons.Chip, tr.Msg("Engines"), tr.Msg("Engines (tab)"), core.PageEngines, func() string {
		return fmt.Sprintf("%d/%d", len(c.Manager.Ready()), len(c.Manager.Engines()))
	})
	nav(icons.Sliders, tr.Msg("Settings"), tr.Msg("Settings (tab)"), core.PageSettings, nil)
	s.theme = giode.Pressable(func(ps giode.PressState) giode.Element {
		icon, label := icons.Sun, tr.T("Light mode")
		if !theme.IsDark() {
			icon, label = icons.Moon, tr.T("Dark mode")
		}
		fg := "theme.muted"
		bg := ""
		if ps.Hovered {
			fg, bg = "theme.text", "theme.text/0.05"
		}
		return giode.Row(giode.Styles{Width: -1, Gap: 12, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap,
			Padding: giode.SymmetricInset(12, 10), BorderRadius: theme.RadiusSm + 2, Background: bg},
			icons.Icon(icon, 18, fg), giode.Text(label, giode.Styles{Color: fg, FontSize: 13, FontWeight: giode.FontWeightMedium}))
	}).OnClick(func(giode.ClickEvent) { onTheme() })
	tr.Bind(func(l string) { s.theme.Label(l) }, "Toggle theme")
	return s
}
