package engines

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/styles/properties"

	"github.com/mehran9675/sonora/internal/manager"
	"github.com/mehran9675/sonora/ui/components/progressbar"
	"github.com/mehran9675/sonora/ui/format"
	"github.com/mehran9675/sonora/ui/tr"
)

// status is the one-line state area: install progress, a problem, or
// whether the model is in memory.
func (cd *card) status() giode.Element {
	rt := cd.rt
	m := rt.Manifest()
	switch {
	case rt.State == manager.StateInstalling:
		return giode.Stack(giode.Styles{Width: -1, Gap: 8},
			progressbar.New(rt.Progress),
			giode.Text(tr.T(rt.Status), giode.Styles{Color: "theme.muted", FontSize: 12, TextOverflow: giode.TextOverflowEllipsis}),
		)
	case rt.Err != "":
		return giode.Text(tr.T(format.FirstLine(rt.Err)), giode.Styles{Color: "theme.danger", FontSize: 12, LineHeight: 16, MaxLines: 2, TextOverflow: giode.TextOverflowEllipsis})
	case rt.State == manager.StateUnavailable && m.Hint != "":
		return giode.Text(tr.T(m.Hint), giode.Styles{Color: "theme.warning", FontSize: 12, LineHeight: 16, MaxLines: 2, TextOverflow: giode.TextOverflowEllipsis})
	case rt.Engine.Loaded():
		return giode.Row(giode.Styles{Gap: 8, Align: giode.AlignCenter},
			giode.Box(giode.Styles{Width: 8, Height: 8, BorderRadius: 99, Background: "theme.success"}),
			giode.Text(tr.T("Model loaded in memory"), giode.Styles{Color: "theme.secondary", FontSize: 12}))
	case cd.downloadedModels() != "":
		text := giode.Text(cd.downloadedModels(), giode.Styles{Color: "theme.secondary", FontSize: 12, LineHeight: 16,
			MaxLines: 2, TextOverflow: giode.TextOverflowEllipsis})
		if rt.State != manager.StateReady || cd.allDownloaded() {
			return text
		}
		// Download another model.
		return giode.Row(giode.Styles{Width: -1, Gap: 8, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			giode.Box(giode.Styles{FlexGrow: 1}, text), cd.addModel)
	}
	return giode.Box(giode.Styles{})
}

// fixed reserves height dp for el and clips anything taller.
func fixed(height int, el giode.Element) giode.Element {
	return giode.Box(giode.Styles{Width: -1, Height: properties.Height(height), Overflow: giode.OverflowHidden}, el)
}
