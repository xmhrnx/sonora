// Package group is a titled group in a side panel: a small overline
// title over its content, as the Studio and Transcribe panels use for
// Engine, Voice and Audio.
package group

import (
	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/ui/theme"
)

// New returns body under title.
func New(title string, body giode.Element) giode.Element {
	return giode.Stack(giode.Styles{Width: -1, Gap: 8}, giode.Text(title, theme.Overline), body)
}
