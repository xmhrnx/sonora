package engines

import (
	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/internal/manager"
	"github.com/mehran9675/sonora/ui/components/badge"
	"github.com/mehran9675/sonora/ui/components/chip"
	"github.com/mehran9675/sonora/ui/components/picker"
	"github.com/mehran9675/sonora/ui/components/statuspill"
	"github.com/mehran9675/sonora/ui/core"
	"github.com/mehran9675/sonora/ui/format"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

// card is one engine tile.
type card struct {
	c  *core.Ctx
	rt *manager.Runtime

	install   *giode.ButtonElement
	cancel    *giode.ButtonElement
	use       *giode.ButtonElement
	logBtn    *giode.ButtonElement
	remove    *giode.ButtonElement
	configure *giode.ButtonElement
	site      *giode.ButtonElement
	eject     *giode.ButtonElement

	// Choosing a model to download, for engines with a model choice:
	// choosing shows the chooser in place of the card's details.
	choosing bool
	choice   string
	models   *picker.Picker
	model    *picker.Trigger
	download *giode.ButtonElement
	back     *giode.ButtonElement
	addModel *giode.ButtonElement
}

// Fixed heights for each part of a card, so every card is the same
// size whatever its state.
const (
	cardHeight   = 356
	descHeight   = 64 // three lines
	metaHeight   = 32 // two lines
	statusHeight = 38
)

func (cd *card) element() giode.Element {
	rt := cd.rt
	m := rt.Manifest()
	st := theme.Card()
	st.Width, st.Height, st.Gap, st.Overflow = 318, cardHeight, 12, giode.OverflowHidden
	if giode.Compact() {
		st.Width = -1 // one card per row on a phone
	}
	header := giode.Row(giode.Styles{Width: -1, Height: 48, Gap: 12, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
		badge.Engine(m, 44),
		giode.Stack(giode.Styles{FlexGrow: 1, Gap: 4},
			giode.Text(m.Name, giode.Styles{Color: "theme.text", FontSize: 16, FontWeight: giode.FontWeightBold, TextOverflow: giode.TextOverflowEllipsis}),
			giode.Row(giode.Styles{Gap: 8, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
				statuspill.For(rt),
				giode.Text(format.Parts(m.Vendor), giode.Styles{Color: "theme.muted", FontSize: 12, TextOverflow: giode.TextOverflowEllipsis}),
			),
		),
	)
	var tags []giode.Element
	for _, t := range m.Tags {
		tags = append(tags, chip.Chip(tr.T(t)))
	}
	if cd.choosing && cd.rt.State != manager.StateInstalling {
		// The chooser takes the place of the description, tags, details
		// and status.
		return giode.Stack(st,
			header,
			fixed(descHeight+22+metaHeight+statusHeight+3*12, cd.chooser()),
			giode.Box(giode.Styles{FlexGrow: 1}),
			giode.Row(giode.Styles{Width: -1, Align: giode.AlignCenter, Gap: 6, Wrap: giode.WrapNoWrap},
				cd.download, cd.back),
		)
	}
	return giode.Stack(st,
		header,
		fixed(descHeight, giode.Text(tr.T(m.Description), giode.Styles{Color: "theme.secondary", FontSize: 13, LineHeight: 19, MaxLines: 3, TextOverflow: giode.TextOverflowEllipsis})),
		giode.Row(giode.Styles{Width: -1, Height: 22, Gap: 6, Wrap: giode.WrapNoWrap, Overflow: giode.OverflowHidden}, tags...),
		fixed(metaHeight, giode.Text(metaLine(m), giode.Styles{Color: "theme.muted", FontSize: 12, LineHeight: 16, MaxLines: 2, TextOverflow: giode.TextOverflowEllipsis})),
		fixed(statusHeight, cd.status()),
		giode.Box(giode.Styles{FlexGrow: 1}),
		cd.actions(),
	)
}
