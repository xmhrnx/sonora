package engines

import (
	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/ui/theme"
)

func group(title string, cards []giode.Element) giode.Element {
	return giode.Stack(giode.Styles{Width: -1, Gap: 14},
		giode.Text(title, theme.Overline),
		giode.Row(giode.Styles{Width: -1, Gap: 16}, cards...),
	)
}
