package library

import (
	"github.com/xmhrnx/giode"
)

// column stacks list rows.
func column(els []giode.Element) giode.Element {
	// Room for the scrollbar, which phones don't show.
	right := 12
	if giode.Compact() {
		right = 0
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 10, Padding: giode.Inset{Right: right, Bottom: 12}}, els...)
}
