// Package playbutton is the round play/pause button of list rows.
package playbutton

import (
	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/ui/core"
	"github.com/mehran9675/sonora/ui/icons"
	"github.com/mehran9675/sonora/ui/tr"
)

// New plays and pauses the audio at path in the player bar, under
// title and subtitle. Create one per row.
func New(c *core.Ctx, path, title, subtitle string) *giode.PressableElement {
	b := giode.Pressable(func(s giode.PressState) giode.Element {
		playing := c.IsPlaying(path)
		icon := icons.Solid(icons.Play, 14, "theme.primary")
		bg := "theme.primary/0.14"
		if playing {
			icon = icons.Solid(icons.Pause, 14, "#ffffff")
			bg = "theme.primary"
		} else if s.Hovered {
			bg = "theme.primary/0.24"
		}
		return giode.Box(giode.Styles{Width: 40, Height: 40, BorderRadius: 99, Background: bg,
			Align: giode.AlignCenter, Justify: giode.JustifyCenter}, icon)
	}).OnClick(func(giode.ClickEvent) { c.PlayPause(path, title, subtitle) })
	tr.Bind(func(s string) { b.Label(s) }, "Play")
	return b
}
